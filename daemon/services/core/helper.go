package core

import (
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/cskr/pubsub"

	"unbalance/daemon/common"
	"unbalance/daemon/domain"
	"unbalance/daemon/lib"
	"unbalance/daemon/logger"
)

// COMMON PLANNER
func getSourceAndDestinationDisks(disks []*domain.Disk, plan *domain.Plan) (*domain.Disk, []*domain.Disk) {
	var srcDisk *domain.Disk
	dstDisks := make([]*domain.Disk, 0)

	for _, disk := range disks {
		if plan.VDisks[disk.Path].Src {
			srcDisk = disk
		}

		if plan.VDisks[disk.Path].Dst {
			dstDisks = append(dstDisks, disk)
		}
	}

	return srcDisk, dstDisks
}

func getIssues(re *regexp.Regexp, disk *domain.Disk, path string) (int64, int64, int64, int64, error) {
	var ownerIssue, groupIssue, folderIssue, fileIssue int64

	folder := filepath.Join(disk.Path, path)

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		return ownerIssue, groupIssue, folderIssue, fileIssue, err
	}

	scanFolder := folder + "/."
	cmd := fmt.Sprintf(`find "%s" -exec stat --format "%%A|%%U:%%G|%%F|%%n" {} \;`, scanFolder)

	err := lib.Shell2(cmd, func(line string) {
		result := re.FindStringSubmatch(line)
		if result == nil {
			return
		}

		u := result[1]
		g := result[2]
		o := result[3]
		user := result[4]
		group := result[5]
		kind := result[6]

		perms := u + g + o

		if user != "nobody" {
			ownerIssue++
		}

		if group != "users" {
			groupIssue++
		}

		if kind == "directory" {
			if perms != "rwxrwxrwx" {
				folderIssue++
			}
		} else {
			match := strings.Compare(perms, "r--r--r--") == 0 || strings.Compare(perms, "rw-rw-rw-") == 0
			if !match {
				fileIssue++
			}
		}
	})

	return ownerIssue, groupIssue, folderIssue, fileIssue, err
}

func getItems(blockSize uint64, re *regexp.Regexp, src, folder string) ([]*domain.Item, uint64, error) {
	var total, blocks uint64
	fBlockSize := float64(blockSize)
	srcFolder := filepath.Join(src, folder)

	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(srcFolder); os.IsNotExist(err) {
		return nil, total, err
	}

	if !fi.IsDir() {
		item := &domain.Item{Name: folder, Size: uint64(fi.Size()), Path: folder, Location: src}
		// Get hardlink info for single file
		populateHardlinkInfo(item, srcFolder)
		return []*domain.Item{item}, uint64(fi.Size()), nil
	}

	entries, err := os.ReadDir(srcFolder)
	if err != nil {
		return nil, total, err
	}

	if len(entries) == 0 {
		// Size: 1 is a trick to allow natural processing of this empty folder: if set to zero, many comparison
		// would misinterpret this as a pending transfer and so on
		return []*domain.Item{&domain.Item{Name: srcFolder, Size: 1, Path: folder, Location: src}}, 0, nil
	}

	items := make([]*domain.Item, 0)

	cmd := fmt.Sprintf(`find "%s" ! -name . -prune -exec du -bs {} +`, srcFolder+"/.")

	err = lib.Shell2(cmd, func(line string) {
		result := re.FindStringSubmatch(line)

		size, _ := strconv.ParseInt(result[1], 10, 64)
		total += uint64(size)

		if blockSize > 0 {
			blocks = uint64(math.Ceil(float64(size) / fBlockSize))
		} else {
			blocks = 0
		}

		itemPath := result[2]
		item := &domain.Item{Name: itemPath, Size: uint64(size), Path: filepath.Join(folder, filepath.Base(itemPath)), Location: src, BlocksUsed: uint64(blocks)}

		// Get hardlink info for this item
		populateHardlinkInfo(item, itemPath)

		items = append(items, item)
	})

	if err != nil {
		return nil, total, err
	}

	return items, total, err
}

// populateHardlinkInfo fills in the inode, device, and link count for an item
func populateHardlinkInfo(item *domain.Item, fullPath string) {
	fi, err := os.Lstat(fullPath)
	if err != nil {
		// If we can't stat, leave hardlink fields at zero (will be treated as non-hardlinked)
		logger.Yellow("hardlink:stat-failed:path(%s):err(%s)", fullPath, err)
		return
	}

	// Get the underlying syscall stat info
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		// Not on a Unix-like system, leave fields at zero
		return
	}

	item.Inode = stat.Ino
	item.Device = uint64(stat.Dev)
	item.LinkCount = uint64(stat.Nlink)

	// For directories, Nlink includes . and .. plus subdirs, not actual hardlinks
	// Directories generally cannot be hardlinked on most filesystems
	// So we treat them as non-hardlinked (linkCount = 1) for our purposes
	if fi.IsDir() {
		item.LinkCount = 1
	}
}

// analyzeHardlinks examines items for hardlink relationships and returns a summary
func analyzeHardlinks(items []*domain.Item, topic string, hub *pubsub.PubSub) *domain.HardlinkSummary {
	analyzer := domain.NewHardlinkAnalyzer()
	summary := analyzer.Analyze(items)

	// Log hardlink information
	logger.Blue("hardlink:analysis:files(%d):groups(%d):apparent(%s):actual(%s):savings(%s)",
		summary.TotalHardlinkedFiles,
		summary.TotalHardlinkGroups,
		lib.ByteSize(summary.ApparentSize),
		lib.ByteSize(summary.ActualSize),
		lib.ByteSize(summary.PotentialSavings))

	// Send summary to frontend
	if summary.TotalHardlinkedFiles > 0 {
		msg := fmt.Sprintf("Hardlink analysis: %d files in %d groups (apparent: %s, actual: %s)",
			summary.TotalHardlinkedFiles,
			summary.TotalHardlinkGroups,
			lib.ByteSize(summary.ApparentSize),
			lib.ByteSize(summary.ActualSize))
		packet := &domain.Packet{Topic: topic, Payload: msg}
		hub.Pub(packet, "socket:broadcast")

		// Log details of each hardlink group
		for _, group := range summary.HardlinkGroups {
			logger.Blue("hardlink:group:inode(%d):size(%s):paths(%d):%v",
				group.Inode,
				lib.ByteSize(group.Size),
				len(group.Paths),
				group.Paths)
		}
	} else {
		logger.Blue("hardlink:analysis:no hardlinked files found")
	}

	// Don't include full group details in response to avoid large payloads
	// They are logged for debugging purposes
	summary.HardlinkGroups = nil

	return summary
}

// discoverHardlinkSiblings finds all paths on a disk that share the same inode.
// This uses `find <disk> -inum <inode>` to discover siblings not in the selected items.
func discoverHardlinkSiblings(diskPath string, inode uint64) ([]string, error) {
	siblings := make([]string, 0)

	cmd := fmt.Sprintf(`find "%s" -inum %d 2>/dev/null`, diskPath, inode)

	err := lib.Shell2(cmd, func(line string) {
		if line != "" {
			siblings = append(siblings, line)
		}
	})

	if err != nil {
		// Log but don't fail - find may return non-zero for permission issues
		logger.Yellow("hardlink:sibling-discovery:inode(%d):err(%s)", inode, err)
	}

	return siblings, nil
}

// analyzeOrphanedHardlinks examines items for hardlink siblings that are NOT selected.
// Returns an OrphanSummary with details about orphaned hardlinks.
func analyzeOrphanedHardlinks(items []*domain.Item, srcDiskPath string, topic string, hub *pubsub.PubSub) *domain.OrphanSummary {
	summary := domain.NewOrphanSummary()

	// Build a set of selected paths for quick lookup
	selectedPaths := make(map[string]bool)
	for _, item := range items {
		fullPath := filepath.Join(item.Location, item.Path)
		selectedPaths[fullPath] = true
		// Also add the Name field which may contain the full path
		selectedPaths[item.Name] = true
	}

	// Track which inodes we've already analyzed
	analyzedInodes := make(map[uint64]bool)

	// Only analyze items that are hardlinked (LinkCount > 1)
	for _, item := range items {
		if !item.IsHardlinked() {
			continue
		}

		inodeKey := item.InodeKey()
		if analyzedInodes[inodeKey] {
			continue
		}
		analyzedInodes[inodeKey] = true

		// Send progress update
		msg := fmt.Sprintf("Checking hardlink siblings for %s...", item.Path)
		packet := &domain.Packet{Topic: topic, Payload: msg}
		hub.Pub(packet, "socket:broadcast")

		// Discover all siblings on the source disk
		allSiblings, err := discoverHardlinkSiblings(srcDiskPath, item.Inode)
		if err != nil {
			logger.Yellow("hardlink:orphan-analysis:skipped:inode(%d):err(%s)", item.Inode, err)
			continue
		}

		// Categorize siblings as selected or unselected
		selectedSiblings := make([]string, 0)
		unselectedSiblings := make([]string, 0)

		for _, siblingPath := range allSiblings {
			if selectedPaths[siblingPath] {
				selectedSiblings = append(selectedSiblings, siblingPath)
			} else {
				unselectedSiblings = append(unselectedSiblings, siblingPath)
			}
		}

		// If there are unselected siblings, this is an orphaned hardlink
		if len(unselectedSiblings) > 0 {
			orphan := &domain.OrphanedHardlink{
				Inode:           item.Inode,
				Device:          item.Device,
				Size:            item.Size,
				SelectedPaths:   selectedSiblings,
				UnselectedPaths: unselectedSiblings,
				TotalLinkCount:  item.LinkCount,
				SpaceImpact:     item.Size, // Source won't free this space
			}

			summary.AddOrphan(orphan)

			logger.Yellow("hardlink:orphan:inode(%d):selected(%d):unselected(%d):size(%s)",
				item.Inode,
				len(selectedSiblings),
				len(unselectedSiblings),
				lib.ByteSize(item.Size))
		}
	}

	// Send summary to frontend if orphans were found
	if summary.HasOrphans() {
		msg := fmt.Sprintf("WARNING: %d file(s) have hardlink siblings that will remain on source. Space impact: %s will NOT be freed.",
			summary.TotalOrphanedFiles,
			lib.ByteSize(summary.TotalSpaceImpact))
		packet := &domain.Packet{Topic: topic, Payload: msg}
		hub.Pub(packet, "socket:broadcast")

		logger.Yellow("hardlink:orphan-summary:groups(%d):files(%d):impact(%s)",
			summary.TotalOrphanedGroups,
			summary.TotalOrphanedFiles,
			lib.ByteSize(summary.TotalSpaceImpact))
	}

	return summary
}

// calculateFreeableSpace returns the actual space that will be freed on source,
// accounting for orphaned hardlinks whose inodes will remain referenced.
func calculateFreeableSpace(bin *domain.Bin, orphanSummary *domain.OrphanSummary) uint64 {
	if orphanSummary == nil || !orphanSummary.HasOrphans() {
		return bin.ActualSize
	}

	// Get the set of orphaned inode keys
	orphanedInodes := orphanSummary.GetOrphanedInodeKeys()

	// Calculate freeable space (exclude orphaned inodes)
	var freeableSpace uint64
	countedInodes := make(map[uint64]bool)

	for _, item := range bin.Items {
		inodeKey := item.InodeKey()

		// Skip if already counted
		if countedInodes[inodeKey] {
			continue
		}
		countedInodes[inodeKey] = true

		// Only count if NOT orphaned (will actually free space)
		if !orphanedInodes[inodeKey] {
			freeableSpace += item.Size
		}
	}

	return freeableSpace
}

// expandSelectionWithSiblings adds unselected hardlink siblings to the item list.
// This allows moving all hardlinks together so space IS freed on source.
func expandSelectionWithSiblings(items []*domain.Item, orphanSummary *domain.OrphanSummary, blockSize uint64) []*domain.Item {
	if orphanSummary == nil || !orphanSummary.HasOrphans() {
		return items
	}

	expanded := make([]*domain.Item, len(items))
	copy(expanded, items)

	fBlockSize := float64(blockSize)

	for _, orphan := range orphanSummary.OrphanedHardlinks {
		for _, path := range orphan.UnselectedPaths {
			// Calculate blocks used
			var blocks uint64
			if blockSize > 0 {
				blocks = uint64(math.Ceil(float64(orphan.Size) / fBlockSize))
			}

			// Extract the disk path (location) and relative path
			// Path format is like /mnt/disk1/some/path
			// We need to extract /mnt/disk1 as location and some/path as relative path
			diskPath, relativePath := extractDiskAndRelativePath(path)

			newItem := &domain.Item{
				Name:       path,
				Size:       orphan.Size,
				Path:       relativePath,
				Location:   diskPath,
				BlocksUsed: blocks,
				Inode:      orphan.Inode,
				Device:     orphan.Device,
				LinkCount:  orphan.TotalLinkCount,
			}
			expanded = append(expanded, newItem)
		}
	}

	logger.Blue("hardlink:expand-selection:original(%d):expanded(%d)", len(items), len(expanded))

	return expanded
}

// extractDiskAndRelativePath splits a full path into disk path and relative path.
// For example, /mnt/disk1/media/movie.mkv becomes /mnt/disk1 and media/movie.mkv
func extractDiskAndRelativePath(fullPath string) (diskPath string, relativePath string) {
	// Paths are expected to be like /mnt/disk1/... or /mnt/user/...
	parts := strings.Split(fullPath, string(filepath.Separator))

	if len(parts) >= 4 && parts[1] == "mnt" {
		// /mnt/disk1/path/to/file -> diskPath=/mnt/disk1, relativePath=path/to/file
		diskPath = filepath.Join(string(filepath.Separator), parts[1], parts[2])
		relativePath = filepath.Join(parts[3:]...)
	} else {
		// Fallback: use parent directory as location
		diskPath = filepath.Dir(fullPath)
		relativePath = filepath.Base(fullPath)
	}

	return diskPath, relativePath
}

// analyzeOrphanedHardlinksMultiDisk analyzes orphaned hardlinks for items from multiple source disks.
// This is used for gather operations where items come from different disks.
func analyzeOrphanedHardlinksMultiDisk(items []*domain.Item, topic string, hub *pubsub.PubSub) *domain.OrphanSummary {
	summary := domain.NewOrphanSummary()

	// Group items by source disk
	itemsByDisk := make(map[string][]*domain.Item)
	for _, item := range items {
		itemsByDisk[item.Location] = append(itemsByDisk[item.Location], item)
	}

	// Analyze orphans for each source disk
	for diskPath, diskItems := range itemsByDisk {
		diskSummary := analyzeOrphanedHardlinks(diskItems, diskPath, topic, hub)
		if diskSummary.HasOrphans() {
			for _, orphan := range diskSummary.OrphanedHardlinks {
				summary.AddOrphan(orphan)
			}
		}
	}

	return summary
}

func (c *Core) getItemsAndIssues(status, blockSize uint64, reItems, reStat *regexp.Regexp, disks []*domain.Disk, folders []string) ([]*domain.Item, int64, int64, int64, int64) {
	var ownerIssue, groupIssue, folderIssue, fileIssue int64
	items := make([]*domain.Item, 0)

	// Get owner/permission issues
	// Get items to be transferred
	for _, disk := range disks {
		for _, path := range folders {
			// logging
			logger.Blue("scanning:disk(%s):folder(%s)", disk.Path, path)

			packet := &domain.Packet{Topic: getTopic(status), Payload: fmt.Sprintf("Scanning %s on %s", path, disk.Path)}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			// check owner and permissions issues for this folder/disk combination
			packet = &domain.Packet{Topic: getTopic(status), Payload: "Checking issues ..."}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			ownIssue, grpIssue, fldIssue, filIssue, err := getIssues(reStat, disk, path)
			if err != nil {
				logger.Yellow("issues:not-available:(%s)", err)
			} else {
				ownerIssue += ownIssue
				groupIssue += grpIssue
				folderIssue += fldIssue
				fileIssue += filIssue

				logger.Blue("issues:owner(%d):group(%d):folder(%d):file(%d)", ownIssue, grpIssue, fldIssue, filIssue)
			}

			// get children files/folders to be transferred
			packet = &domain.Packet{Topic: getTopic(status), Payload: "Getting items ..."}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			list, total, err := getItems(blockSize, reItems, disk.Path, path)
			if err != nil {
				logger.Yellow("items:not-available:(%s)", err)
			} else {
				logger.Blue("items:count(%d):size(%s)", len(list), lib.ByteSize(total))
				items = append(items, list...)
			}
		}
	}

	return items, ownerIssue, groupIssue, folderIssue, fileIssue
}

func (c *Core) sendTimeFeedbackToFrontend(topic, fended string, elapsed time.Duration) {
	packet := &domain.Packet{Topic: topic, Payload: fmt.Sprintf("Ended: %s", fended)}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	packet = &domain.Packet{Topic: topic, Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
	c.ctx.Hub.Pub(packet, "socket:broadcast")
}

func (c *Core) sendMailFeedback(fstarted, ffinished string, elapsed time.Duration, plan *domain.Plan, notTransferred string) {
	subject := "unbalanced - PLANNING completed"
	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", fstarted, ffinished, elapsed)
	if notTransferred != "" {
		switch c.ctx.Config.NotifyPlan {
		case 1:
			message += "\n\nSome folders will not be transferred because there's not enough space for them in any of the destination disks."
		case 2:
			message += "\n\nThe following folders will not be transferred because there's not enough space for them in any of the destination disks:\n\n" + notTransferred
		}
	}

	if plan.OwnerIssue > 0 || plan.GroupIssue > 0 || plan.FolderIssue > 0 || plan.FileIssue > 0 {
		message += fmt.Sprintf(`
			\n\nThere are some permission issues:
			\n\n%d file(s)/folder(s) with an owner other than 'nobody'
			\n%d file(s)/folder(s) with a group other than 'users'
			\n%d folder(s) with a permission other than 'drwxrwxrwx'
			\n%d files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'
			\n\nCheck the log file (/var/log/unbalanced.log) for additional information
			\n\nIt's strongly suggested to install the Fix Common Plugins and run the Docker Safe New Permissions command
		`, plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)
	}

	if sendErr := sendmail(c.ctx.Config.NotifyPlan, subject, message, false); sendErr != nil {
		logger.Red("unable to send mail: %s", sendErr)
	}
}

func (c *Core) getReservedAmount(size uint64) uint64 {
	var reserved uint64

	switch c.ctx.Config.ReservedUnit {
	case "%":
		fcalc := size * c.ctx.Config.ReservedAmount / 100
		reserved = fcalc
	case "Mb":
		reserved = c.ctx.Config.ReservedAmount * 1024 * 1024
	case "Gb":
		reserved = c.ctx.Config.ReservedAmount * 1024 * 1024 * 1024
	default:
		reserved = common.ReservedSpace
	}

	return reserved
}

func (c *Core) endPlan(status uint64, plan *domain.Plan, disks []*domain.Disk, items, toBeTransferred []*domain.Item) {
	plan.Ended = time.Now()
	elapsed := lib.Round(time.Since(plan.Started), time.Millisecond)
	logger.Blue("%s", elapsed) // otherwise it won't send correct value to frontend 🤷‍♂️

	fstarted := plan.Started.Format(timeFormat)
	fended := plan.Ended.Format(timeFormat)

	// Send to frontend console started/ended/elapsed times
	c.sendTimeFeedbackToFrontend(getTopic(status), fended, time.Since(plan.Started))

	// send to frontend the items that will not be transferred, if any
	// notTransferred holds a string representation of all the items, separated by a '\n'
	notTransferred := ""

	if status == common.OpScatterPlan {
		// some logging
		if len(toBeTransferred) == 0 {
			logger.Blue("%s:No items can be transferred.", getName(status))
		} else {
			logger.Blue("%s:%d items will be transferred.", getName(status), len(toBeTransferred))
			for _, folder := range toBeTransferred {
				logger.Blue("%s:willBeTransferred(%s)", getName(status), folder.Path)
			}
		}

		if len(items) > 0 {
			packet := &domain.Packet{Topic: getTopic(status), Payload: "The following items will not be transferred, because there's not enough space in the target disks:\n"}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			logger.Blue("%s:%d items will NOT be transferred.", getName(status), len(items))
			for _, item := range items {
				notTransferred += item.Path + "\n"

				packet = &domain.Packet{Topic: getTopic(status), Payload: item.Path}
				c.ctx.Hub.Pub(packet, "socket:broadcast")
				logger.Blue("%s:notTransferred(%s)", getName(status), item.Path)
			}
		}
	}

	// send mail according to user preferences
	c.sendMailFeedback(fstarted, fended, elapsed, plan, notTransferred)

	// some local logging
	logger.Blue("%s:ItemsLeft(%d)", getName(status), len(items))
	logger.Blue("%s:Listing (%d) disks ...", getName(status), len(disks))
	for _, disk := range disks {
		if plan.VDisks[disk.Path].Bin != nil {
			logger.Blue("=========================================================")
			logger.Blue("disk(%s):items(%d)-(%s):currentFree(%s)-plannedFree(%s)", disk.Path, len(plan.VDisks[disk.Path].Bin.Items), lib.ByteSize(plan.VDisks[disk.Path].Bin.Size), lib.ByteSize(disk.Free), lib.ByteSize(plan.VDisks[disk.Path].PlannedFree))
			logger.Blue("---------------------------------------------------------")

			for _, item := range plan.VDisks[disk.Path].Bin.Items {
				logger.Blue("[%s] %s", lib.ByteSize(item.Size), item.Name)
			}

			logger.Blue("---------------------------------------------------------")
			logger.Blue("")
		} else {
			logger.Blue("=========================================================")
			logger.Blue("disk(%s):no-items:currentFree(%s)", disk.Path, lib.ByteSize(disk.Free))
			logger.Blue("---------------------------------------------------------")
			logger.Blue("---------------------------------------------------------")
			logger.Blue("")
		}
	}

	logger.Blue("=========================================================")
	logger.Blue("Bytes To Transfer: %s", lib.ByteSize(plan.BytesToTransfer))
	logger.Blue("---------------------------------------------------------")

	packet := &domain.Packet{Topic: getTopic(status), Payload: "Planning Ended"}
	c.ctx.Hub.Pub(packet, "socket:broadcast")
}

func (c *Core) printDisks(disks []*domain.Disk, blockSize uint64) {
	logger.Blue("planner:array(%d disks):blockSize(%d)", len(disks), blockSize)
	for _, disk := range disks {
		logger.Blue("disk(%s):fs(%s):size(%d):free(%d):blocksTotal(%d):blocksFree(%d)", disk.Path, disk.FsType, disk.Size, disk.Free, disk.BlocksTotal, disk.BlocksFree)
	}
}

// HELPER FUNCTIONS
func getName(status uint64) string {
	if status == common.OpScatterPlan {
		return "scatterPlan"
	}

	return "gatherPlan"
}

func getTopic(status uint64) string {
	if status == common.OpScatterPlan {
		return common.EventScatterPlanProgress
	}

	return common.EventGatherPlanProgress
}

func removeItems(items, list []*domain.Item) []*domain.Item {
	w := 0 // write index

loop:
	for _, item := range items {
		for _, itm := range list {
			if itm.Name == item.Name {
				continue loop
			}
		}
		items[w] = item
		w++
	}

	return items[:w]
}

func isZombie(proc string) (bool, int, error) {
	var zombie bool
	var retcode int

	b, e := os.ReadFile(proc)
	if e != nil {
		return false, 0, e
	}

	fields := strings.Split(string(b), " ")
	state := fields[2]
	zombie = state == "Z"
	if zombie {
		retcode, _ = strconv.Atoi(fields[51])
	}

	return zombie, retcode, nil
}

func getReadBytes(proc string) (uint64, error) {
	var sRead string

	b, e := os.ReadFile(proc)
	if e != nil {
		return 0, e
	}

	lines := strings.Split(string(b), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "rchar:") {
			sRead = line[7:]
			break
		}
	}

	read, _ := strconv.ParseUint(sRead, 10, 64)

	return read, nil
}

func (c *Core) notifyCommandsToRun(opName string, operation *domain.Operation) {
	message := "\n\nThe following commands will be executed:\n\n"

	for _, command := range operation.Commands {
		cmd := fmt.Sprintf(`(src: %s) rsync %s %s %s`, command.Src, operation.RsyncStrArgs, strconv.Quote(command.Entry), strconv.Quote(command.Dst))
		message += cmd + "\n"
	}

	subject := fmt.Sprintf("unbalanced - %s operation STARTED", strings.ToUpper(opName))

	go func() {
		if sendErr := sendmail(c.ctx.NotifyTransfer, subject, message, c.ctx.DryRun); sendErr != nil {
			logger.Red("hp-sendmail %s", sendErr)
		}
	}()
}

func progress(bytesToTransfer, bytesTransferred uint64, elapsed time.Duration) (percent float64, left time.Duration, speed float64) {
	bytesPerSec := float64(bytesTransferred) / elapsed.Seconds()
	speed = bytesPerSec / 1024 / 1024 // MB/s

	percent = (float64(bytesTransferred) / float64(bytesToTransfer)) * 100 // %

	left = time.Duration(float64(bytesToTransfer-bytesTransferred)/bytesPerSec) * time.Second

	return
}

func getCurrentTransfer(proc, prefix string) (string, error) {
	var current string

	name, e := os.Readlink(proc)
	if e != nil {
		return "", e
	}

	if strings.HasPrefix(name, prefix) {
		current = name
	}

	return current, nil
}

func getError(line string, re *regexp.Regexp, ers map[int]string) string {
	result := re.FindStringSubmatch(line)
	if result == nil || len(result) < 1 {
		return "unknown error"
	}

	status, _ := strconv.Atoi(result[1])
	msg, ok := ers[status]
	if !ok {
		msg = "unknown error"
	}

	return msg
}

func sendmail(notify int, subject, message string, dryRun bool) (err error) {
	if notify == 0 {
		return nil
	}

	dry := ""
	if dryRun {
		dry = "-------\nDRY RUN\n-------\n"
	}

	msg := dry + message

	cmd := exec.Command(mailCmd, "-e", "unbalanced operation update", "-s", subject, "-m", msg)
	err = cmd.Run()

	return
}

func showPotentiallyPrunedItems(operation *domain.Operation, command *domain.Command) {
	if operation.DryRun && operation.OpKind == common.OpGatherMove {
		parent := filepath.Dir(command.Entry)
		if parent != "." {
			logger.Blue(`Would delete empty folders starting from (%s) - (find "%s" -type d -empty -prune -exec rm -rf {} \;) `, filepath.Join(command.Src, parent), filepath.Join(command.Src, parent))
		} else {
			logger.Blue(`WONT DELETE: find "%s" -type d -empty -prune -exec rm -rf {} \;`, filepath.Join(command.Src, parent))
		}
	}
}
