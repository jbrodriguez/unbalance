package services

import (
	"fmt"
	"io/ioutil"
	"jbrodriguez/unbalance/server/src/algorithm"
	"jbrodriguez/unbalance/server/src/common"
	"jbrodriguez/unbalance/server/src/domain"
	"jbrodriguez/unbalance/server/src/dto"
	"jbrodriguez/unbalance/server/src/lib"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
)

type Calc struct {
	bus      *pubsub.PubSub
	settings *lib.Settings
	actor    *actor.Actor

	reItems *regexp.Regexp
	reStat  *regexp.Regexp
}

// NewCore -
func NewCalc(bus *pubsub.PubSub, settings *lib.Settings) *Calc {
	calc := &Calc{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
	}

	calc.reItems = regexp.MustCompile(`(\d+)\s+(.*?)$`)
	calc.reStat = regexp.MustCompile(`[-dclpsbD]([-rwxsS]{3})([-rwxsS]{3})([-rwxtT]{3})\|(.*?)\:(.*?)\|(.*?)\|(.*)`)

	return calc
}

// Start -
func (c *Calc) Start() (err error) {
	mlog.Info("starting service Calc ...")

	c.actor.Register(common.INT_SCATTER_CALCULATE, c.scatter)
	c.actor.Register(common.INT_GATHER_CALCULATE, c.gather)

	go c.actor.React()

	return nil
}

// Stop -
func (c *Calc) Stop() {
	mlog.Info("stopped service Calc ...")
}

func getSourceAndDestinationDisks(state *domain.State) ([]*domain.Disk, []*domain.Disk) {
	srcDisks := make([]*domain.Disk, 0)
	dstDisks := make([]*domain.Disk, 0)

	for _, disk := range state.Unraid.Disks {
		if state.Operation.VDisks[disk.Path].Src {
			srcDisks = append(srcDisks, disk)
		}

		if state.Operation.VDisks[disk.Path].Dst {
			dstDisks = append(dstDisks, disk)
		}
	}

	return srcDisks, dstDisks
}

func (c *Calc) getIssues(disk *domain.Disk, path string) (int64, int64, int64, int64) {
	var ownerIssue, groupIssue, folderIssue, fileIssue int64

	// Check owner and permission issues
	outbound := &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: fmt.Sprintf("Scanning %s on %s", path, disk.Path)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	folder := filepath.Join(disk.Path, path)

	mlog.Info("getIssues:Scanning disk(%s):folder(%s)", disk.Path, path)

	if _, err := os.Stat(folder); os.IsNotExist(err) {
		mlog.Warning("getIssues:Folder does not exist:(%s)", folder)
		return ownerIssue, groupIssue, folderIssue, fileIssue
	}

	scanFolder := folder + "/."
	cmdText := fmt.Sprintf(`find "%s" -exec stat --format "%%A|%%U:%%G|%%F|%%n" {} \;`, scanFolder)

	mlog.Info("getIssues:Executing:(%s)", cmdText)

	lib.Shell(cmdText, mlog.Warning, "getIssues:find/stat:", "", func(line string) {
		result := c.reStat.FindStringSubmatch(line)
		if result == nil {
			mlog.Warning("getIssues:Unable to parse:(%s)", line)
			return
		}

		u := result[1]
		g := result[2]
		o := result[3]
		user := result[4]
		group := result[5]
		kind := result[6]
		name := result[7]

		perms := u + g + o

		if user != "nobody" {
			if c.settings.Verbosity == 1 {
				mlog.Info("getIssues:User != nobody:[%s]:(%s)", user, name)
			}

			ownerIssue++
		}

		if group != "users" {
			if c.settings.Verbosity == 1 {
				mlog.Info("getIssues:Group != users:[%s]:(%s)", group, name)
			}

			groupIssue++
		}

		if kind == "directory" {
			if perms != "rwxrwxrwx" {
				if c.settings.Verbosity == 1 {
					mlog.Info("getIssues:Folder perms != rwxrwxrwx:[%s]:(%s)", perms, name)
				}

				folderIssue++
			}
		} else {
			match := strings.Compare(perms, "r--r--r--") == 0 || strings.Compare(perms, "rw-rw-rw-") == 0
			if !match {
				if c.settings.Verbosity == 1 {
					mlog.Info("getIssues:File perms != rw-rw-rw- or r--r--r--:[%s]:(%s)", perms, name)
				}

				fileIssue++
			}
		}
	})

	mlog.Info("getIssues:owner(%d),group(%d),folder(%d),file(%d)", ownerIssue, groupIssue, folderIssue, fileIssue)

	outbound = &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: "Checked permissions ..."}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	return ownerIssue, groupIssue, folderIssue, fileIssue
}

func (c *Calc) getEntries(src string, folder string) []*domain.Item {
	srcFolder := filepath.Join(src, folder)

	items := make([]*domain.Item, 0)

	mlog.Info("getEntries:Scanning disk(%s):folder(%s)", src, folder)

	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(srcFolder); os.IsNotExist(err) {
		mlog.Warning("getEntries:Folder does not exist: %s", srcFolder)
		return items
	}

	if !fi.IsDir() {
		mlog.Info("getEntries:found(%s):size(%d)", srcFolder, fi.Size())

		item := &domain.Item{Name: folder, Size: fi.Size(), Path: folder, Location: src}
		items = append(items, item)

		msg := fmt.Sprintf("Found %s (%s)", item.Name, lib.ByteSize(item.Size))
		outbound := &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: msg}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return items
	}

	entries, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		mlog.Warning("getEntries:Unable to readdir:(%s)", err)
	}

	mlog.Info("getEntries:Readdir(%d)", len(entries))

	if len(entries) == 0 {
		mlog.Info("getEntries:No entries under %s", srcFolder)
		return nil
	}

	scanFolder := srcFolder + "/."
	cmdText := fmt.Sprintf("find \"%s\" ! -name . -prune -exec du -bs {} +", scanFolder)

	mlog.Info("getEntries:Executing:(%s)", cmdText)

	lib.Shell(cmdText, mlog.Warning, "getEntries:find/du:", "", func(line string) {
		mlog.Info("getEntries:find(%s): %s", scanFolder, line)

		result := c.reItems.FindStringSubmatch(line)

		size, _ := strconv.ParseInt(result[1], 10, 64)

		item := &domain.Item{Name: result[2], Size: size, Path: filepath.Join(folder, filepath.Base(result[2])), Location: src}
		items = append(items, item)

		msg := fmt.Sprintf("Found %s (%s)", filepath.Base(item.Name), lib.ByteSize(size))
		outbound := &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: msg}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	})

	return items
}

func (c *Calc) sendTimeFeedbackToFrontend(fstarted, ffinished string, elapsed time.Duration) {
	outbound := &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: fmt.Sprintf("Ended: %s", ffinished)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
}

func (c *Calc) sendMailFeeback(fstarted, ffinished string, elapsed time.Duration, operation *domain.Operation, notTransferred string) {
	subject := "unBALANCE - CALCULATE operation completed"
	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", fstarted, ffinished, elapsed)
	if notTransferred != "" {
		switch c.settings.NotifyCalc {
		case 1:
			message += "\n\nSome folders will not be transferred because there's not enough space for them in any of the destination disks."
		case 2:
			message += "\n\nThe following folders will not be transferred because there's not enough space for them in any of the destination disks:\n\n" + notTransferred
		}
	}

	if operation.OwnerIssue > 0 || operation.GroupIssue > 0 || operation.FolderIssue > 0 || operation.FileIssue > 0 {
		message += fmt.Sprintf(`
			\n\nThere are some permission issues:
			\n\n%d file(s)/folder(s) with an owner other than 'nobody'
			\n%d file(s)/folder(s) with a group other than 'users'
			\n%d folder(s) with a permission other than 'drwxrwxrwx'
			\n%d files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'
			\n\nCheck the log file (/boot/logs/unbalance.log) for additional information
			\n\nIt's strongly suggested to install the Fix Common Plugins and run the Docker Safe New Permissions command
		`, operation.OwnerIssue, operation.GroupIssue, operation.FolderIssue, operation.FileIssue)
	}

	if sendErr := lib.Sendmail(common.MailCmd, c.settings.NotifyCalc, subject, message, false); sendErr != nil {
		mlog.Error(sendErr)
	}
}

func (c *Calc) getReservedAmount(size int64) int64 {
	var reserved int64

	switch c.settings.ReservedUnit {
	case "%":
		fcalc := size * c.settings.ReservedAmount / 100
		reserved = int64(fcalc)
		break
	case "Mb":
		reserved = c.settings.ReservedAmount * 1000 * 1000
		break
	case "Gb":
		reserved = c.settings.ReservedAmount * 1000 * 1000 * 1000
		break
	default:
		reserved = lib.ReservedSpace
	}

	return reserved
}

func removeItems(items []*domain.Item, list []*domain.Item) []*domain.Item {
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

func (c *Calc) scatter(msg *pubsub.Message) {
	state := msg.Payload.(*domain.State)

	mlog.Info("Running calculate for scatter operation ...")
	state.Operation.Started = time.Now()

	outbound := &dto.Packet{Topic: common.WS_CALC_STARTED, Payload: "Operation started"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// create two slices
	// one of source disks, the other of destinations disks
	// in scatter srcDisks will contain only one element
	srcDisks, dstDisks := getSourceAndDestinationDisks(state)

	// get dest disks with more free space to the top
	sort.Slice(dstDisks, func(i, j int) bool { return dstDisks[i].Free < dstDisks[j].Free })

	// some logging
	for _, disk := range srcDisks {
		mlog.Info("scatterCalc:sourceDisk:(%s)", disk.Path)
	}

	for _, disk := range dstDisks {
		mlog.Info("scatterCalc:destDisk:(%s)", disk.Path)
	}

	var ownerIssue, groupIssue, folderIssue, fileIssue int64
	items := make([]*domain.Item, 0)

	// Get owner/permission issues
	// Get items to be transferred
	for _, disk := range srcDisks {
		for _, path := range state.Operation.ChosenFolders {
			// check owner and permissions issues for this folder/disk combination
			ownIssue, grpIssue, fldIssue, filIssue := c.getIssues(disk, path)
			ownerIssue += ownIssue
			groupIssue += grpIssue
			folderIssue += fldIssue
			fileIssue += filIssue

			// get children files/folders to be transferred
			entries := c.getEntries(disk.Path, path)
			if entries != nil {
				items = append(items, entries...)
			}
		}
	}

	// this is true for scatter calculation, where only one source disk is allowed
	srcDisk := srcDisks[0]

	state.Operation.OwnerIssue = ownerIssue
	state.Operation.GroupIssue = groupIssue
	state.Operation.FolderIssue = folderIssue
	state.Operation.FileIssue = fileIssue

	mlog.Info("scatterCalc:issues:owner(%d),group(%d),folder(%d),file(%d)", state.Operation.OwnerIssue, state.Operation.GroupIssue, state.Operation.FolderIssue, state.Operation.FileIssue)
	mlog.Info("scatterCalc:totalItemsToBeTransferred(%d)", len(items))

	for _, item := range items {
		mlog.Info("scatterCalc:toBeTransferred:Path(%s):Size(%s)", item.Path, lib.ByteSize(item.Size))
	}

	willBeTransferred := make([]*domain.Item, 0)

	if len(items) > 0 {
		// Initialize fields
		state.Operation.BytesToTransfer = 0

		for _, disk := range dstDisks {
			diskWithoutMnt := disk.Path[5:]
			msg := fmt.Sprintf("Trying to allocate folders to %s ...", diskWithoutMnt)
			outbound = &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: msg}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			mlog.Info("scatterCalc:%s", msg)
			// time.Sleep(2 * time.Second)

			reserved := c.getReservedAmount(disk.Size)

			ceil := lib.Max(lib.ReservedSpace, reserved)
			mlog.Info("scatterCalc:ItemsLeft(%d):ReservedSpace(%d)", len(items), ceil)

			packer := algorithm.NewKnapsack(disk, items, ceil)
			bin := packer.BestFit()
			if bin != nil {
				state.Operation.VDisks[disk.Path].Bin = bin
				state.Operation.VDisks[disk.Path].PlannedFree -= bin.Size
				state.Operation.VDisks[srcDisk.Path].PlannedFree += bin.Size

				state.Operation.BytesToTransfer += bin.Size

				willBeTransferred = append(willBeTransferred, bin.Items...)
				items = removeItems(items, bin.Items)

				mlog.Info("scatterCalc:BinAllocated=[Disk(%s); Items(%d)];Freespace=[original(%s); final(%s)]", disk.Path, len(bin.Items), lib.ByteSize(srcDisk.Free), lib.ByteSize(state.Operation.VDisks[srcDisk.Path].PlannedFree))
			} else {
				mlog.Info("scatterCalc:NoBinAllocated=Disk(%s)", disk.Path)
			}
		}
	}

	state.Operation.Finished = time.Now()
	elapsed := lib.Round(time.Since(state.Operation.Started), time.Millisecond)

	fstarted := state.Operation.Started.Format(timeFormat)
	ffinished := state.Operation.Finished.Format(timeFormat)

	// Send to frontend console started/ended/elapsed times
	c.sendTimeFeedbackToFrontend(fstarted, ffinished, elapsed)

	// some logging
	if len(willBeTransferred) == 0 {
		mlog.Info("scatterCalc:No items can be transferred.")
	} else {
		mlog.Info("scatterCalc:%d items will be transferred.", len(willBeTransferred))
		for _, folder := range willBeTransferred {
			mlog.Info("scatterCalc:willBeTransferred(%s)", folder.Path)
		}
	}

	// send to frontend the items that will not be transferred, if any
	// notTransferred holds a string representation of all the items, separated by a '\n'
	notTransferred := ""

	if len(items) > 0 {
		outbound = &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: "The following items will not be transferred, because there's not enough space in the target disks:\n"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		mlog.Info("scatterCalc:%d items will NOT be transferred.", len(items))
		for _, item := range items {
			notTransferred += item.Path + "\n"

			outbound = &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: item.Path}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			mlog.Info("scatterCalc:notTransferred(%s)", item.Path)
		}
	}

	// send mail according to user preferences
	c.sendMailFeeback(fstarted, ffinished, elapsed, state.Operation, notTransferred)

	// some local logging
	mlog.Info("scatterCalc:ItemsLeft(%d)", len(items))
	mlog.Info("scatterCalc:src(%s):Listing (%d) disks ...", srcDisk.Path, len(state.Unraid.Disks))
	for _, disk := range state.Unraid.Disks {
		if state.Operation.VDisks[disk.Path].Bin != nil {
			mlog.Info("=========================================================")
			mlog.Info("Disk(%s):ALLOCATED %d items (%s): %s planned free space remaining", disk.Path, len(state.Operation.VDisks[disk.Path].Bin.Items), lib.ByteSize(state.Operation.VDisks[disk.Path].Bin.Size), lib.ByteSize(state.Operation.VDisks[disk.Path].PlannedFree))
			mlog.Info("---------------------------------------------------------")

			for _, item := range state.Operation.VDisks[disk.Path].Bin.Items {
				mlog.Info("[%s] %s", lib.ByteSize(item.Size), item.Name)
			}

			mlog.Info("---------------------------------------------------------")
			mlog.Info("")
		} else {
			mlog.Info("=========================================================")
			mlog.Info("Disk(%s):NO ALLOCATION: %s free", disk.Path, lib.ByteSize(disk.Free))
			mlog.Info("---------------------------------------------------------")
			mlog.Info("---------------------------------------------------------")
			mlog.Info("")
		}
	}

	mlog.Info("=========================================================")
	mlog.Info("Results for %s", srcDisk.Path)
	mlog.Info("Original Free Space: %s", lib.ByteSize(srcDisk.Free))
	mlog.Info("Final Free Space: %s", lib.ByteSize(state.Operation.VDisks[srcDisk.Path].PlannedFree))
	mlog.Info("Bytes To Transfer: %s", lib.ByteSize(state.Operation.BytesToTransfer))
	mlog.Info("---------------------------------------------------------")

	outbound = &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: "Operation Finished"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	c.bus.Pub(&pubsub.Message{Payload: state.Operation}, common.INT_SCATTER_CALCULATE_FINISHED)
}

func (c *Calc) gather(msg *pubsub.Message) {
	state := msg.Payload.(*domain.State)

	mlog.Info("Running calculate for gather operation ...")
	state.Operation.Started = time.Now()

	outbound := &dto.Packet{Topic: common.WS_CALC_STARTED, Payload: "Operation started"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	var ownerIssue, groupIssue, folderIssue, fileIssue int64
	items := make([]*domain.Item, 0)

	// Get owner/permission issues
	// Get items to be transferred
	for _, disk := range state.Unraid.Disks {
		for _, path := range state.Operation.ChosenFolders {
			// check owner and permissions issues for this folder/disk combination
			ownIssue, grpIssue, fldIssue, filIssue := c.getIssues(disk, path)
			ownerIssue += ownIssue
			groupIssue += grpIssue
			folderIssue += fldIssue
			fileIssue += filIssue

			// get children files/folders to be transferred
			entries := c.getEntries(disk.Path, path)
			if entries != nil {
				items = append(items, entries...)
			}
		}
	}

	state.Operation.OwnerIssue = ownerIssue
	state.Operation.GroupIssue = groupIssue
	state.Operation.FolderIssue = folderIssue
	state.Operation.FileIssue = fileIssue

	mlog.Info("gatherCalc:issues:owner(%d),group(%d),folder(%d),file(%d)", state.Operation.OwnerIssue, state.Operation.GroupIssue, state.Operation.FolderIssue, state.Operation.FileIssue)
	mlog.Info("gatherCalc:totalItemsToBeTransferred(%d)", len(items))

	var totalSize int64
	for _, item := range items {
		totalSize += item.Size
		mlog.Info("gatherCalc:toBeTransferred:Path(%s):Size(%s)", item.Path, lib.ByteSize(item.Size))
	}

	if len(items) > 0 {
		// Initialize fields
		state.Operation.BytesToTransfer = 0

		for _, disk := range state.Unraid.Disks {
			diskWithoutMnt := disk.Path[5:]
			msg := fmt.Sprintf("Trying to allocate folders to %s ...", diskWithoutMnt)
			outbound = &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: msg}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			mlog.Info("gatherCalc:%s", msg)
			// time.Sleep(2 * time.Second)

			reserved := c.getReservedAmount(disk.Size)

			ceil := lib.Max(lib.ReservedSpace, reserved)
			mlog.Info("gatherCalc:ItemsLeft(%d):ReservedSpace(%d)", len(items), ceil)

			packer := algorithm.NewGreedy(disk, items, totalSize, ceil)
			bin := packer.FitAll()
			if bin != nil {
				state.Operation.VDisks[disk.Path].Bin = bin
				state.Operation.VDisks[disk.Path].PlannedFree -= bin.Size

				state.Operation.BytesToTransfer += bin.Size

				mlog.Info("gatherCalc:BinAllocated=[Disk(%s); Items(%d)]", disk.Path, len(bin.Items))
			} else {
				mlog.Info("gatherCalc:NoBinAllocated=Disk(%s)", disk.Path)
			}
		}
	}

	state.Operation.Finished = time.Now()
	elapsed := lib.Round(time.Since(state.Operation.Started), time.Millisecond)

	fstarted := state.Operation.Started.Format(timeFormat)
	ffinished := state.Operation.Finished.Format(timeFormat)

	// Send to frontend console started/ended/elapsed times
	c.sendTimeFeedbackToFrontend(fstarted, ffinished, elapsed)

	// send to frontend the items that will not be transferred, if any
	// notTransferred holds a string representation of all the items, separated by a '\n'
	state.Operation.FoldersNotTransferred = make([]string, 0)
	notTransferred := ""

	// send mail according to user preferences
	c.sendMailFeeback(fstarted, ffinished, elapsed, state.Operation, notTransferred)

	// some local logging
	mlog.Info("gatherCalc:ItemsLeft(%d)", len(items))
	mlog.Info("gatherCalc:Listing (%d) disks ...", len(state.Unraid.Disks))
	for _, disk := range state.Unraid.Disks {
		if state.Operation.VDisks[disk.Path].Bin != nil {
			mlog.Info("=========================================================")
			mlog.Info("Disk(%s):ALLOCATED %d items (%s): %s planned free space remaining", disk.Path, len(state.Operation.VDisks[disk.Path].Bin.Items), lib.ByteSize(state.Operation.VDisks[disk.Path].Bin.Size), lib.ByteSize(state.Operation.VDisks[disk.Path].PlannedFree))
			mlog.Info("---------------------------------------------------------")

			for _, item := range state.Operation.VDisks[disk.Path].Bin.Items {
				mlog.Info("[%s] %s", lib.ByteSize(item.Size), item.Name)
			}

			mlog.Info("---------------------------------------------------------")
			mlog.Info("")
		} else {
			mlog.Info("=========================================================")
			mlog.Info("Disk(%s):NO ALLOCATION: %s free", disk.Path, lib.ByteSize(disk.Free))
			mlog.Info("---------------------------------------------------------")
			mlog.Info("---------------------------------------------------------")
			mlog.Info("")
		}
	}

	mlog.Info("=========================================================")
	mlog.Info("Bytes To Transfer: %s", lib.ByteSize(state.Operation.BytesToTransfer))
	mlog.Info("---------------------------------------------------------")

	outbound = &dto.Packet{Topic: common.WS_CALC_PROGRESS, Payload: "Operation Finished"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	c.bus.Pub(&pubsub.Message{Payload: state.Operation}, common.INT_GATHER_CALCULATE_FINISHED)
}
