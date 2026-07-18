package core

import (
	"context"
	"fmt"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

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

func getIssues(ctx context.Context, re *regexp.Regexp, disk *domain.Disk, path string, tick func(int)) (int64, int64, int64, int64, error) {
	var ownerIssue, groupIssue, folderIssue, fileIssue int64

	folder := filepath.Join(disk.Path, path)

	if _, err := os.Stat(folder); err != nil {
		return ownerIssue, groupIssue, folderIssue, fileIssue, err
	}

	scanFolder := folder + "/."
	// '+' batches many files per stat invocation; ';' would fork one stat
	// process per file, which dominates scan time on large trees
	findArgs := []string{scanFolder, "-exec", "stat", "--format=%A|%U:%G|%F|%n", "{}", "+"}

	err := lib.StreamContext(ctx, "find", findArgs, func(line string) {
		tick(1)

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

func getItems(ctx context.Context, blockSize uint64, re *regexp.Regexp, src, folder string, tick func(int)) ([]*domain.Item, uint64, error) {
	var total, blocks uint64
	fBlockSize := float64(blockSize)
	srcFolder := filepath.Join(src, folder)

	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(srcFolder); err != nil {
		return nil, total, err
	}

	if !fi.IsDir() {
		size := uint64(fi.Size())
		if blockSize > 0 {
			blocks = uint64(math.Ceil(float64(size) / fBlockSize))
		}
		return []*domain.Item{&domain.Item{Name: folder, Size: size, Path: folder, Location: src, BlocksUsed: blocks}}, size, nil
	}

	entries, err := os.ReadDir(srcFolder)
	if err != nil {
		return nil, total, err
	}

	if len(entries) == 0 {
		// Size: 1 is a trick to allow natural processing of this empty folder: if set to zero, many comparison
		// would misinterpret this as a pending transfer and so on
		return []*domain.Item{&domain.Item{Name: srcFolder, Size: 1, Path: folder, Location: src, BlocksUsed: 1}}, 0, nil
	}

	items := make([]*domain.Item, 0)

	findArgs := []string{srcFolder + "/.", "!", "-name", ".", "-prune", "-exec", "du", "-bs", "{}", "+"}

	err = lib.StreamContext(ctx, "find", findArgs, func(line string) {
		tick(1)

		result := re.FindStringSubmatch(line)
		if result == nil {
			// du output can carry lines that don't parse (e.g. names with
			// embedded newlines); skip them instead of crashing the daemon
			logger.Yellow("items:unparseable du output line, skipping: %q", line)
			return
		}

		size, _ := strconv.ParseInt(result[1], 10, 64)
		total += uint64(size)

		if blockSize > 0 {
			blocks = uint64(math.Ceil(float64(size) / fBlockSize))
		} else {
			blocks = 0
		}

		item := &domain.Item{Name: result[2], Size: uint64(size), Path: filepath.Join(folder, filepath.Base(result[2])), Location: src, BlocksUsed: uint64(blocks)}
		items = append(items, item)
	})

	if err != nil {
		return nil, total, err
	}

	return items, total, err
}

func (c *Core) getItemsAndIssues(ctx context.Context, status, blockSize uint64, reItems, reStat *regexp.Regexp, disks []*domain.Disk, folders []string) ([]*domain.Item, int64, int64, int64, int64) {
	var ownerIssue, groupIssue, folderIssue, fileIssue int64
	items := make([]*domain.Item, 0)

	// heartbeat: let the frontend know the scan is alive during long walks,
	// at most one packet every few seconds
	var scanned int64
	lastBeat := time.Now()
	tick := func(n int) {
		scanned += int64(n)
		if time.Since(lastBeat) < 2*time.Second {
			return
		}
		lastBeat = time.Now()
		packet := &domain.Packet{Topic: getTopic(status), Payload: fmt.Sprintf("Scanning ... %d entries so far", scanned)}
		c.ctx.Hub.Pub(packet, "socket:broadcast")
	}

	// Get owner/permission issues
	// Get items to be transferred
	for _, disk := range disks {
		for _, path := range folders {
			// the user cancelled the plan, no point in scanning any further
			if c.stopped.Load() || ctx.Err() != nil {
				logger.Blue("planner:cancelled:scan abandoned")
				return items, ownerIssue, groupIssue, folderIssue, fileIssue
			}

			// logging
			logger.Blue("scanning:disk(%s):folder(%s)", disk.Path, path)

			packet := &domain.Packet{Topic: getTopic(status), Payload: fmt.Sprintf("Scanning %s on %s", path, disk.Path)}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			// check owner and permissions issues for this folder/disk combination
			packet = &domain.Packet{Topic: getTopic(status), Payload: "Checking issues ..."}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			ownIssue, grpIssue, fldIssue, filIssue, err := getIssues(ctx, reStat, disk, path, tick)
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

			list, total, err := getItems(ctx, blockSize, reItems, disk.Path, path, tick)
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

// publishFoundItems streams the discovered items to the frontend in chunks,
// so a large plan doesn't flood the websocket with one packet per item.
func (c *Core) publishFoundItems(topic string, items []*domain.Item) {
	const chunkSize = 200

	lines := make([]string, 0, chunkSize)
	flush := func() {
		if len(lines) == 0 {
			return
		}
		packet := &domain.Packet{Topic: topic, Payload: strings.Join(lines, "\n")}
		c.ctx.Hub.Pub(packet, "socket:broadcast")
		lines = lines[:0]
	}

	for _, item := range items {
		logger.Blue("planner:found(%s):size(%d)", filepath.Join(item.Location, item.Path), item.Size)

		lines = append(lines, fmt.Sprintf("Found %s (%s)", filepath.Join(item.Location, item.Path), lib.ByteSize(item.Size)))
		if len(lines) == chunkSize {
			flush()
		}
	}
	flush()
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

			const chunkSize = 200
			var sb strings.Builder
			lines := make([]string, 0, chunkSize)
			for _, item := range items {
				sb.WriteString(item.Path)
				sb.WriteByte('\n')

				lines = append(lines, item.Path)
				if len(lines) == chunkSize {
					packet = &domain.Packet{Topic: getTopic(status), Payload: strings.Join(lines, "\n")}
					c.ctx.Hub.Pub(packet, "socket:broadcast")
					lines = lines[:0]
				}
				logger.Blue("%s:notTransferred(%s)", getName(status), item.Path)
			}
			if len(lines) > 0 {
				packet = &domain.Packet{Topic: getTopic(status), Payload: strings.Join(lines, "\n")}
				c.ctx.Hub.Pub(packet, "socket:broadcast")
			}

			notTransferred = sb.String()
		}
	}

	// the notification email can't hold an unbounded list: a single exec
	// argument is limited to 128KiB on Linux, so truncate what gets mailed
	const maxMailList = 64 * 1024
	if len(notTransferred) > maxMailList {
		kept := strings.LastIndexByte(notTransferred[:maxMailList], '\n')
		if kept < 0 {
			kept = maxMailList
		}
		omitted := strings.Count(notTransferred[kept:], "\n")
		notTransferred = notTransferred[:kept] + fmt.Sprintf("\n... and %d more items (see /var/log/unbalanced.log)\n", omitted)
	}

	// send mail according to user preferences, without holding up the plan:
	// a stuck notify script must not keep the app busy
	go c.sendMailFeedback(fstarted, fended, elapsed, plan, notTransferred)

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

func (c *Core) planCancelled(topic string) {
	c.state.Status = common.OpNeutral
	c.clearPlanContext()

	logger.Blue("planning cancelled by the user")

	packet := &domain.Packet{Topic: topic, Payload: "Planning cancelled by the user"}
	c.ctx.Hub.Pub(packet, "socket:broadcast")
}

func removeItems(items, list []*domain.Item) []*domain.Item {
	remove := make(map[string]struct{}, len(list))
	for _, itm := range list {
		remove[itm.Name] = struct{}{}
	}

	w := 0 // write index

	for _, item := range items {
		if _, ok := remove[item.Name]; ok {
			continue
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
	// guard against NaN/Inf: they poison the operation and history since the
	// JSON encoder refuses non-finite floats
	if bytesToTransfer == 0 || elapsed <= 0 {
		return 0, 0, 0
	}
	if bytesTransferred > bytesToTransfer {
		bytesTransferred = bytesToTransfer
	}

	bytesPerSec := float64(bytesTransferred) / elapsed.Seconds()
	speed = bytesPerSec / 1024 / 1024 // MB/s

	percent = (float64(bytesTransferred) / float64(bytesToTransfer)) * 100 // %

	if bytesPerSec > 0 {
		left = time.Duration(float64(bytesToTransfer-bytesTransferred)/bytesPerSec) * time.Second
	}

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

func rsyncExitReason(code int) string {
	msg, ok := rsyncErrors[code]
	if !ok {
		msg = "unknown error"
	}

	return fmt.Sprintf("%s (rsync exit code %d)", msg, code)
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
			logger.Blue("Would prune empty parent folders starting from (%s)", filepath.Join(command.Src, parent))
		} else {
			logger.Blue("WONT PRUNE: (%s)", filepath.Join(command.Src, parent))
		}
	}
}
