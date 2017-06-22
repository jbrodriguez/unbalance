package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"jbrodriguez/unbalance/server/src/algorithm"
	"jbrodriguez/unbalance/server/src/dto"
	"jbrodriguez/unbalance/server/src/lib"
	"jbrodriguez/unbalance/server/src/model"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
)

const (
	mailCmd    = "/usr/local/emhttp/webGui/scripts/notify"
	timeFormat = "Jan _2, 2006 15:04:05"
)

// Core service
type Core struct {
	bus      *pubsub.PubSub
	storage  *model.Unraid
	settings *lib.Settings

	// this holds the state of any operation
	operation model.Operation

	actor *actor.Actor

	reFreeSpace *regexp.Regexp
	reItems     *regexp.Regexp
	reRsync     *regexp.Regexp
	reStat      *regexp.Regexp
	reProgress  *regexp.Regexp

	rsyncErrors map[int]string
}

// NewCore -
func NewCore(bus *pubsub.PubSub, settings *lib.Settings) *Core {
	core := &Core{
		bus:      bus,
		settings: settings,
		// opState:  stateIdle,
		storage:   &model.Unraid{},
		actor:     actor.NewActor(bus),
		operation: model.Operation{OpState: model.StateIdle, PrevState: model.StateIdle},
	}

	core.reFreeSpace = regexp.MustCompile(`(.*?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.*?)\s+(.*?)$`)
	core.reItems = regexp.MustCompile(`(\d+)\s+(.*?)$`)
	core.reRsync = regexp.MustCompile(`exit status (\d+)`)
	core.reProgress = regexp.MustCompile(`(?s)^([\d,]+).*?\(.*?\)$|^([\d,]+).*?$`)
	core.reStat = regexp.MustCompile(`[-dclpsbD]([-rwxsS]{3})([-rwxsS]{3})([-rwxtT]{3})\|(.*?)\:(.*?)\|(.*?)\|(.*)`)

	core.rsyncErrors = map[int]string{
		0:  "Success",
		1:  "Syntax or usage error",
		2:  "Protocol incompatibility",
		3:  "Errors selecting input/output files, dirs",
		4:  "Requested action not supported: an attempt was made to manipulate 64-bit files on a platform that cannot support them, or an option was specified that is supported by the client and not by the server.",
		5:  "Error starting client-server protocol",
		6:  "Daemon unable to append to log-file",
		10: "Error in socket I/O",
		11: "Error in file I/O",
		12: "Error in rsync protocol data stream",
		13: "Errors with program diagnostics",
		14: "Error in IPC code",
		20: "Received SIGUSR1 or SIGINT",
		21: "Some error returned by waitpid()",
		22: "Error allocating core memory buffers",
		23: "Partial transfer due to error",
		24: "Partial transfer due to vanished source files",
		25: "The --max-delete limit stopped deletions",
		30: "Timeout in data send/receive",
		35: "Timeout waiting for daemon connection",
	}

	// core.ownerPerms = map[int]bool{
	// 	644
	// }

	return core
}

// Start -
func (c *Core) Start() (err error) {
	mlog.Info("starting service Core ...")

	c.actor.Register("/get/config", c.getConfig)
	c.actor.Register("/config/set/notifyCalc", c.setNotifyCalc)
	c.actor.Register("/config/set/notifyMove", c.setNotifyMove)
	c.actor.Register("/config/set/reservedSpace", c.setReservedSpace)
	c.actor.Register("/get/storage", c.getStorage)
	c.actor.Register("/config/toggle/dryRun", c.toggleDryRun)
	c.actor.Register("/get/tree", c.getTree)
	c.actor.Register("/config/set/rsyncFlags", c.setRsyncFlags)

	c.actor.Register("calculate", c.calc)
	c.actor.Register("move", c.move)
	c.actor.Register("copy", c.copy)
	c.actor.Register("validate", c.validate)
	c.actor.Register("getLog", c.getLog)

	err = c.storage.SanityCheck(c.settings.APIFolders)
	if err != nil {
		return err
	}

	go c.actor.React()

	return nil
}

// Stop -
func (c *Core) Stop() {
	mlog.Info("stopped service Core ...")
}

// SetStorage -
func (c *Core) SetStorage(unraid *model.Unraid) {
	c.storage = unraid
}

func (c *Core) getConfig(msg *pubsub.Message) {
	mlog.Info("Sending config")

	rsyncFlags := strings.Join(c.settings.RsyncFlags, " ")
	if rsyncFlags == "-avX --partial" || rsyncFlags == "-avRX --partial" {
		c.settings.RsyncFlags = []string{"-avPRX"}
		c.settings.Save()
	}

	msg.Reply <- &c.settings.Config
}

func (c *Core) setNotifyCalc(msg *pubsub.Message) {
	fnotify := msg.Payload.(float64)
	notify := int(fnotify)

	mlog.Info("Setting notifyCalc to (%d)", notify)

	c.settings.NotifyCalc = notify
	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) setNotifyMove(msg *pubsub.Message) {
	fnotify := msg.Payload.(float64)
	notify := int(fnotify)

	mlog.Info("Setting notifyMove to (%d)", notify)

	c.settings.NotifyMove = notify
	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) setReservedSpace(msg *pubsub.Message) {
	mlog.Warning("payload: %+v", msg.Payload)
	payload, ok := msg.Payload.(string)
	if !ok {
		mlog.Warning("Unable to convert Reserved Space parameters")
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert Reserved Space parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		msg.Reply <- &c.settings.Config

		return
	}

	var reserved dto.Reserved
	err := json.Unmarshal([]byte(payload), &reserved)
	if err != nil {
		mlog.Warning("Unable to bind reservedSpace parameters: %s", err)
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind reservedSpace parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		return
	}

	amount := int64(reserved.Amount)
	unit := reserved.Unit

	mlog.Info("Setting reservedAmount to (%d)", amount)
	mlog.Info("Setting reservedUnit to (%s)", unit)

	c.settings.ReservedAmount = amount
	c.settings.ReservedUnit = unit
	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) getStorage(msg *pubsub.Message) {
	var stats string

	if c.operation.OpState == model.StateIdle {
		c.storage.Refresh()
	} else if c.operation.OpState == model.StateMove || c.operation.OpState == model.StateCopy {
		percent, left, speed := progress(c.operation.BytesToTransfer, c.operation.BytesTransferred, c.operation.Started)
		stats = fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)
	}

	c.storage.Stats = stats
	c.storage.OpState = c.operation.OpState
	c.storage.PrevState = c.operation.PrevState
	c.storage.BytesToTransfer = c.operation.BytesToTransfer

	msg.Reply <- c.storage
}

func (c *Core) toggleDryRun(msg *pubsub.Message) {
	mlog.Info("Toggling dryRun from (%t)", c.settings.DryRun)

	c.settings.ToggleDryRun()

	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) getTree(msg *pubsub.Message) {
	path := msg.Payload.(string)

	msg.Reply <- c.storage.GetTree(path)
}

func (c *Core) setRsyncFlags(msg *pubsub.Message) {
	// mlog.Warning("payload: %+v", msg.Payload)
	payload, ok := msg.Payload.(string)
	if !ok {
		mlog.Warning("Unable to convert Rsync Flags parameters")
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert Rsync Flags parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		msg.Reply <- &c.settings.Config

		return
	}

	var rsync dto.Rsync
	err := json.Unmarshal([]byte(payload), &rsync)
	if err != nil {
		mlog.Warning("Unable to bind rsyncFlags parameters: %s", err)
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind rsyncFlags parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		return
		// mlog.Fatalf(err.Error())
	}

	mlog.Info("Setting rsyncFlags to (%s)", strings.Join(rsync.Flags, " "))

	c.settings.RsyncFlags = rsync.Flags
	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) calc(msg *pubsub.Message) {
	c.operation = model.Operation{OpState: model.StateCalc, PrevState: model.StateIdle}

	go c._calc(msg)
}

func (c *Core) _calc(msg *pubsub.Message) {
	defer func() { c.operation.OpState = model.StateIdle }()

	payload, ok := msg.Payload.(string)
	if !ok {
		mlog.Warning("Unable to convert calculate parameters")
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert calculate parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		return
	}

	var dtoCalc dto.Calculate
	err := json.Unmarshal([]byte(payload), &dtoCalc)
	if err != nil {
		mlog.Warning("Unable to bind calculate parameters: %s", err)
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind calculate parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		return
		// mlog.Fatalf(err.Error())
	}

	mlog.Info("Running calculate operation ...")
	c.operation.Started = time.Now()

	outbound := &dto.Packet{Topic: "calcStarted", Payload: "Operation started"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	disks := make([]*model.Disk, 0)

	// create array of destination disks
	var srcDisk *model.Disk
	for _, disk := range c.storage.Disks {
		// reset disk
		disk.NewFree = disk.Free
		disk.Bin = nil
		disk.Src = false
		disk.Dst = dtoCalc.DestDisks[disk.Path]

		if disk.Path == dtoCalc.SourceDisk {
			disk.Src = true
			srcDisk = disk
		} else {
			// add it to the target disk list, only if the user selected it
			if val, ok := dtoCalc.DestDisks[disk.Path]; ok && val {
				// double check, if it's a cache disk, make sure it's the main cache disk
				if disk.Type == "Cache" && len(disk.Name) > 5 {
					continue
				}

				disks = append(disks, disk)
			}
		}
	}

	mlog.Info("_calc:Begin:srcDisk(%s); dstDisks(%d)", srcDisk.Path, len(disks))

	for _, disk := range disks {
		mlog.Info("_calc:elegibleDestDisk(%s)", disk.Path)
	}

	sort.Sort(model.ByFree(disks))

	srcDiskWithoutMnt := srcDisk.Path[5:]

	owner := ""
	lib.Shell("id -un", mlog.Warning, "owner", "", func(line string) {
		owner = line
	})

	group := ""
	lib.Shell("id -gn", mlog.Warning, "group", "", func(line string) {
		group = line
	})

	c.operation.OwnerIssue = 0
	c.operation.GroupIssue = 0
	c.operation.FolderIssue = 0
	c.operation.FileIssue = 0

	// Check permission and gather folders to be transferred from
	// source disk
	folders := make([]*model.Item, 0)
	for _, path := range dtoCalc.Folders {
		msg := fmt.Sprintf("Scanning %s on %s", path, srcDiskWithoutMnt)
		outbound = &dto.Packet{Topic: "calcProgress", Payload: msg}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		mlog.Info("_calc:%s", msg)

		c.checkOwnerAndPermissions(c.operation, dtoCalc.SourceDisk, path, owner, group)

		msg = "Checked permissions ..."
		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		mlog.Info("_calc:%s", msg)

		list := c.getFolders(dtoCalc.SourceDisk, path)
		if list != nil {
			folders = append(folders, list...)
		}
	}

	mlog.Info("_calc:foldersToBeTransferredTotal(%d)", len(folders))

	for _, v := range folders {
		mlog.Info("_calc:toBeTransferred:Path(%s); Size(%s)", v.Path, lib.ByteSize(v.Size))
	}

	willBeTransferred := make([]*model.Item, 0)
	if len(folders) > 0 {
		// Initialize fields
		// c.storage.BytesToTransfer = 0
		// c.storage.SourceDiskName = srcDisk.Path
		c.operation.BytesToTransfer = 0
		c.operation.SourceDiskName = srcDisk.Path
		c.operation.RsyncFlags = c.settings.RsyncFlags
		if c.settings.DryRun {
			c.operation.RsyncFlags = append(c.operation.RsyncFlags, "--dry-run")
		}
		c.operation.RsyncStrFlags = strings.Join(c.operation.RsyncFlags, " ")
		c.operation.Commands = make([]model.Command, 0)

		for _, disk := range disks {
			diskWithoutMnt := disk.Path[5:]
			msg := fmt.Sprintf("Trying to allocate folders to %s ...", diskWithoutMnt)
			outbound = &dto.Packet{Topic: "calcProgress", Payload: msg}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			mlog.Info("_calc:%s", msg)
			// time.Sleep(2 * time.Second)

			if disk.Path != srcDisk.Path {
				// disk.NewFree = disk.Free

				var reserved int64
				switch c.settings.ReservedUnit {
				case "%":
					fcalc := disk.Size * c.settings.ReservedAmount / 100
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

				ceil := lib.Max(lib.ReservedSpace, reserved)
				mlog.Info("_calc:FoldersLeft(%d):ReservedSpace(%d)", len(folders), ceil)

				packer := algorithm.NewKnapsack(disk, folders, ceil)
				bin := packer.BestFit()
				if bin != nil {
					srcDisk.NewFree += bin.Size
					disk.NewFree -= bin.Size
					c.operation.BytesToTransfer += bin.Size
					// c.storage.BytesToTransfer += bin.Size

					willBeTransferred = append(willBeTransferred, bin.Items...)
					folders = c.removeFolders(folders, bin.Items)

					for _, item := range bin.Items {
						var src, dst string
						if strings.Contains(c.operation.RsyncStrFlags, "R") {
							if item.Path[0] == filepath.Separator {
								src = item.Path[1:]
							} else {
								src = item.Path
							}

							dst = disk.Path + string(filepath.Separator)
						} else {
							src = filepath.Join(c.operation.SourceDiskName, item.Path)
							dst = filepath.Join(disk.Path, filepath.Dir(item.Path)) + string(filepath.Separator)
						}

						c.operation.Commands = append(c.operation.Commands, model.Command{
							Src:  src,
							Dst:  dst,
							Path: item.Path,
							Size: item.Size,
						})

						// args := append(
						// 	rsyncArgs,
						// 	src,
						// 	dst,
						// )
						// // cmd := exec.Command("rsync", args)

						// // cmd := fmt.Sprintf("rsync %s \"%s\" \"%s/\"", strings.Join(rsyncArgs, " "), filepath.Join(c.storage.SourceDiskName, item.Path), filepath.Join(disk.Path, filepath.Dir(item.Path)))
						// cmd := fmt.Sprintf(`rsync %s %s %s`, rsyncStrArgs, strconv.Quote(src), strconv.Quote(dst))
						// mlog.Info("cmd(%s)", cmd)
					}

					mlog.Info("_calc:BinAllocated=[Disk(%s); Items(%d)];Freespace=[original(%s); final(%s)]", disk.Path, len(bin.Items), lib.ByteSize(srcDisk.Free), lib.ByteSize(srcDisk.NewFree))
				} else {
					mlog.Info("_calc:NoBinAllocated=Disk(%s)", disk.Path)
				}
			}
		}
	}

	c.operation.Finished = time.Now()
	elapsed := lib.Round(time.Since(c.operation.Started), time.Millisecond)

	fstarted := c.operation.Started.Format(timeFormat)
	ffinished := c.operation.Finished.Format(timeFormat)

	// Send to frontend console started/ended/elapsed times
	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	if len(willBeTransferred) == 0 {
		mlog.Info("_calc:No folders can be transferred.")
	} else {
		mlog.Info("_calc:%d folders will be transferred.", len(willBeTransferred))
		for _, folder := range willBeTransferred {
			mlog.Info("_calc:willBeTransferred(%s)", folder.Path)
		}
	}

	// send to frontend the folders that will not be transferred, if any
	// notTransferred holds a string representation of all the folders, separated by a '\n'
	c.operation.FoldersNotTransferred = make([]string, 0)
	notTransferred := ""
	if len(folders) > 0 {
		outbound = &dto.Packet{Topic: "calcProgress", Payload: "The following folders will not be transferred, because there's not enough space in the target disks:\n"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		mlog.Info("_calc:%d folders will NOT be transferred.", len(folders))
		for _, folder := range folders {
			c.operation.FoldersNotTransferred = append(c.operation.FoldersNotTransferred, folder.Path)

			notTransferred += folder.Path + "\n"

			outbound = &dto.Packet{Topic: "calcProgress", Payload: folder.Path}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			mlog.Info("_calc:notTransferred(%s)", folder.Path)
		}
	}

	// send mail according to user preferences
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

	if c.operation.OwnerIssue > 0 || c.operation.GroupIssue > 0 || c.operation.FolderIssue > 0 || c.operation.FileIssue > 0 {
		message += fmt.Sprintf(`
			\n\nThere are some permission issues:
			\n\n%d file(s)/folder(s) with an owner other than 'nobody'
			\n%d file(s)/folder(s) with a group other than 'users'
			\n%d folder(s) with a permission other than 'drwxrwxrwx'
			\n%d files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'
			\n\nCheck the log file (/boot/logs/unbalance.log) for additional information
			\n\nIt's strongly suggested to install the Fix Common Plugins and run the Docker Safe New Permissions command
		`, c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)
	}

	if sendErr := c.sendmail(c.settings.NotifyCalc, subject, message, false); sendErr != nil {
		mlog.Error(sendErr)
	}

	// some local logging
	mlog.Info("_calc:FoldersLeft(%d)", len(folders))
	mlog.Info("_calc:src(%s):Listing (%d) disks ...", srcDisk.Path, len(c.storage.Disks))
	for _, disk := range c.storage.Disks {
		// mlog.Info("the mystery of the year(%s)", disk.Path)
		disk.Print()
	}

	mlog.Info("=========================================================")
	mlog.Info("Results for %s", srcDisk.Path)
	mlog.Info("Original Free Space: %s", lib.ByteSize(srcDisk.Free))
	mlog.Info("Final Free Space: %s", lib.ByteSize(srcDisk.NewFree))
	mlog.Info("Gained Space: %s", lib.ByteSize(srcDisk.NewFree-srcDisk.Free))
	mlog.Info("Bytes To Move: %s", lib.ByteSize(c.operation.BytesToTransfer))
	mlog.Info("---------------------------------------------------------")

	c.storage.Print()
	// msg.Reply <- c.storage

	mlog.Info("_calc:End:srcDisk(%s)", srcDisk.Path)

	outbound = &dto.Packet{Topic: "calcProgress", Payload: "Operation Finished"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	c.storage.BytesToTransfer = c.operation.BytesToTransfer
	c.storage.OpState = c.operation.OpState
	c.storage.PrevState = c.operation.PrevState

	mlog.Warning("prev state (%d)", c.storage.PrevState)

	// send to front end the signal of operation finished
	outbound = &dto.Packet{Topic: "calcFinished", Payload: c.storage}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
	// and there actually perm issues
	if c.operation.BytesToTransfer > 0 && (c.operation.OwnerIssue+c.operation.GroupIssue+c.operation.FolderIssue+c.operation.FileIssue > 0) {
		outbound = &dto.Packet{Topic: "calcPermIssue", Payload: fmt.Sprintf("%d|%d|%d|%d", c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}
}

func (c *Core) getFolders(src string, folder string) (items []*model.Item) {
	srcFolder := filepath.Join(src, folder)

	mlog.Info("getFolders:Scanning source-disk(%s):folder(%s)", src, folder)

	var fi os.FileInfo
	var err error
	if fi, err = os.Stat(srcFolder); os.IsNotExist(err) {
		mlog.Warning("getFolders:Folder does not exist: %s", srcFolder)
		return nil
	}

	if !fi.IsDir() {
		mlog.Info("getFolder-found(%s)-size(%d)", srcFolder, fi.Size())

		item := &model.Item{Name: folder, Size: fi.Size(), Path: folder}
		items = append(items, item)

		msg := fmt.Sprintf("Found %s (%s)", item.Name, lib.ByteSize(item.Size))
		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	dirs, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		mlog.Warning("getFolders:Unable to readdir: %s", err)
	}

	mlog.Info("getFolders:Readdir(%d)", len(dirs))

	if len(dirs) == 0 {
		mlog.Info("getFolders:No subdirectories under %s", srcFolder)
		return nil
	}

	scanFolder := srcFolder + "/."
	cmdText := fmt.Sprintf("find \"%s\" ! -name . -prune -exec du -bs {} +", scanFolder)

	mlog.Info("getFolders:Executing %s", cmdText)

	lib.Shell(cmdText, mlog.Warning, "getFolders:find/du:", "", func(line string) {
		mlog.Info("getFolders:find(%s): %s", scanFolder, line)

		result := c.reItems.FindStringSubmatch(line)
		// mlog.Info("[%s] %s", result[1], result[2])

		size, _ := strconv.ParseInt(result[1], 10, 64)

		item := &model.Item{Name: result[2], Size: size, Path: filepath.Join(folder, filepath.Base(result[2]))}
		items = append(items, item)

		msg := fmt.Sprintf("Found %s (%s)", filepath.Base(item.Name), lib.ByteSize(size))
		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	})

	return
}

func (c *Core) checkOwnerAndPermissions(operation model.Operation, src, folder, ownerName, groupName string) {
	srcFolder := filepath.Join(src, folder)

	// outbound := &dto.Packet{Topic: "calcProgress", Payload: "Checking permissions ..."}
	// c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	mlog.Info("perms:Scanning disk(%s):folder(%s)", src, folder)

	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
		mlog.Warning("perms:Folder does not exist: %s", srcFolder)
		return
	}

	scanFolder := srcFolder + "/."
	cmdText := fmt.Sprintf(`find "%s" -exec stat --format "%%A|%%U:%%G|%%F|%%n" {} \;`, scanFolder)

	mlog.Info("perms:Executing %s", cmdText)

	lib.Shell(cmdText, mlog.Warning, "perms:find/stat:", "", func(line string) {
		result := c.reStat.FindStringSubmatch(line)
		if result == nil {
			mlog.Warning("perms:Unable to parse (%s)", line)
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
			mlog.Info("perms:User != nobody: [%s]: %s", user, name)
			operation.OwnerIssue++
		}

		if group != "users" {
			mlog.Info("perms:Group != users: [%s]: %s", group, name)
			operation.GroupIssue++
		}

		if kind == "directory" {
			if perms != "rwxrwxrwx" {
				mlog.Info("perms:Folder perms != rwxrwxrwx: [%s]: %s", perms, name)
				operation.FolderIssue++
			}
		} else {
			match := strings.Compare(perms, "r--r--r--") == 0 || strings.Compare(perms, "rw-rw-rw-") == 0
			if !match {
				mlog.Info("perms:File perms != rw-rw-rw- or r--r--r--: [%s]: %s", perms, name)
				operation.FileIssue++
			}
		}
	})

	return
}

func (c *Core) removeFolders(folders []*model.Item, list []*model.Item) []*model.Item {
	w := 0 // write index

loop:
	for _, fld := range folders {
		for _, itm := range list {
			if itm.Name == fld.Name {
				continue loop
			}
		}
		folders[w] = fld
		w++
	}

	return folders[:w]
}

func (c *Core) move(msg *pubsub.Message) {
	c.operation.OpState = model.StateMove
	c.operation.PrevState = model.StateMove
	go c.transfer(msg)
}

func (c *Core) copy(msg *pubsub.Message) {
	c.operation.OpState = model.StateCopy
	c.operation.PrevState = model.StateCopy
	go c.transfer(msg)
}

func (c *Core) validate(msg *pubsub.Message) {
	c.operation.OpState = model.StateValidate
	c.operation.PrevState = model.StateValidate
	// go c.checksum(msg)
}

func (c *Core) transfer(msg *pubsub.Message) {
	defer func() {
		c.operation.OpState = model.StateIdle
		c.operation.Started = time.Time{}
		c.operation.BytesTransferred = 0
	}()

	mlog.Info("Running transfer operation ...")
	c.operation.Started = time.Now()

	outbound := &dto.Packet{Topic: "transferStarted", Payload: "Operation started"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// rsyncArgs := c.settings.RsyncFlags

	// if c.settings.DryRun {
	// 	rsyncArgs = append(rsyncArgs, "--dry-run")
	// }

	commandsExecuted := make([]string, 0)

	outbound = &dto.Packet{Topic: "progressStats", Payload: "Waiting to collect stats ..."}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	var calls int64
	var callsPerDelta int64

	// execute each rsync command created during the calculate phase
	for _, command := range c.operation.Commands {
		args := append(
			c.operation.RsyncFlags,
			command.Src,
			command.Dst,
		)
		// cmd := exec.Command("rsync", args)

		// cmd := fmt.Sprintf("rsync %s \"%s\" \"%s/\"", strings.Join(rsyncArgs, " "), filepath.Join(c.storage.SourceDiskName, item.Path), filepath.Join(disk.Path, filepath.Dir(item.Path)))
		cmd := fmt.Sprintf(`rsync %s %s %s`, c.operation.RsyncStrFlags, strconv.Quote(command.Src), strconv.Quote(command.Dst))
		mlog.Info("cmd(%s)", cmd)

		outbound = &dto.Packet{Topic: "transferProgress", Payload: cmd}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		bytesTransferred := c.operation.BytesTransferred

		var deltaMoved int64

		// actual shell execution
		err := lib.ShellEx(func(text string) {
			line := strings.TrimSpace(text)

			if len(line) <= 0 {
				return
			}

			if callsPerDelta <= 50 {
				calls++
			}

			delta := int64(time.Since(c.operation.Started) / time.Second)
			if delta == 0 {
				delta = 1
			}
			// mlog.Info("calls(%d)-seconds(%d)", calls, delta)
			callsPerDelta = calls / delta

			match := c.reProgress.FindStringSubmatch(line)
			if match == nil {
				// this is a regular output line from rsync, print it
				mlog.Info("%s", line)

				if callsPerDelta <= 50 {
					outbound := &dto.Packet{Topic: "transferProgress", Payload: line}
					c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
				}

				return
			}

			// this is a file transfer progress output line
			if match[1] == "" {
				// this happens when the file hasn't finished transferring
				moved := strings.Replace(match[2], ",", "", -1)
				deltaMoved, _ = strconv.ParseInt(moved, 10, 64)
			} else {
				// the file has finished transferring
				moved := strings.Replace(match[1], ",", "", -1)
				deltaMoved, _ = strconv.ParseInt(moved, 10, 64)
				bytesTransferred += deltaMoved
			}

			percent, left, speed := progress(c.operation.BytesToTransfer, bytesTransferred+deltaMoved, c.operation.Started)

			msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)

			if callsPerDelta <= 50 {
				outbound := &dto.Packet{Topic: "progressStats", Payload: msg}
				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			}

		}, c.operation.SourceDiskName, "rsync", args...)

		if err != nil {
			finished := time.Now()
			elapsed := time.Since(c.operation.Started)

			subject := "unBALANCE - TRANSFER operation INTERRUPTED"

			headline := fmt.Sprintf("Transfer command (%s) was interrupted: %s", cmd, err.Error()+" : "+getError(err.Error(), c.reRsync, c.rsyncErrors))

			mlog.Warning(headline)
			outbound := &dto.Packet{Topic: "opError", Payload: "Transfer operation was interrupted. Check log (/boot/logs/unbalance.log) for details."}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			c.finishTransferOperation(subject, headline, commandsExecuted, c.operation.Started, finished, elapsed)

			return
		}

		c.operation.BytesTransferred = c.operation.BytesTransferred + command.Size

		percent, left, speed := progress(c.operation.BytesToTransfer, c.operation.BytesTransferred, c.operation.Started)

		msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)
		mlog.Info("Current progress: %s", msg)

		outbound := &dto.Packet{Topic: "progressStats", Payload: msg}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		commandsExecuted = append(commandsExecuted, cmd)

		// if the operation is Move, delete the source folder
		if !c.settings.DryRun && c.operation.OpState == model.StateMove {
			rmrf := fmt.Sprintf("rm -rf \"%s\"", filepath.Join(c.operation.SourceDiskName, command.Path))
			mlog.Info("Removing: (%s)", rmrf)
			err = lib.Shell(rmrf, mlog.Warning, "transferProgress:", "", func(line string) {
				mlog.Info(line)
			})

			if err != nil {
				msg := fmt.Sprintf("Unable to remove source folder:(%s)", filepath.Join(c.operation.SourceDiskName, command.Path))

				outbound := &dto.Packet{Topic: "transferProgress", Payload: msg}
				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

				mlog.Warning(msg)
			}
		}
	}

	finished := time.Now()
	elapsed := time.Since(c.operation.Started)

	subject := "unBALANCE - TRANSFER operation completed"
	headline := "Transfer operation has finished"

	c.finishTransferOperation(subject, headline, commandsExecuted, c.operation.Started, finished, elapsed)
}

func (c *Core) finishTransferOperation(subject, headline string, commands []string, started, finished time.Time, elapsed time.Duration) {
	fstarted := started.Format(timeFormat)
	ffinished := finished.Format(timeFormat)
	elapsed = lib.Round(time.Since(started), time.Millisecond)

	outbound := &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "transferProgress", Payload: headline}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "transferProgress", Payload: "These are the commands that were executed:"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	printedCommands := ""
	for _, command := range commands {
		printedCommands += command + "\n"
		outbound = &dto.Packet{Topic: "transferProgress", Payload: command}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	outbound = &dto.Packet{Topic: "transferProgress", Payload: "Operation Finished"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// send to front end the signal of operation finished
	if c.settings.DryRun {
		outbound = &dto.Packet{Topic: "transferProgress", Payload: "--- IT WAS A DRY RUN ---"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	outbound = &dto.Packet{Topic: "transferFinished", Payload: ""}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s\n\n%s", fstarted, ffinished, elapsed, headline)
	switch c.settings.NotifyMove {
	case 1:
		message += fmt.Sprintf("\n\n%d commands were executed.", len(commands))
	case 2:
		message += "\n\nThese are the commands that were executed:\n\n" + printedCommands
	}

	if sendErr := c.sendmail(c.settings.NotifyCalc, subject, message, c.settings.DryRun); sendErr != nil {
		mlog.Error(sendErr)
	}

	mlog.Info(subject)
	mlog.Info(message)
}

func (c *Core) checksum(msg *pubsub.Message) {
	defer func() {
		c.operation.OpState = model.StateIdle
		c.operation.PrevState = model.StateIdle
		c.operation.Started = time.Time{}
		c.operation.BytesTransferred = 0
	}()

	mlog.Info("Running validate operation ...")
	c.operation.Started = time.Now()

	outbound := &dto.Packet{Topic: "transferStarted", Payload: "Operation started"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	commandsExecuted := make([]string, 0)

	checkRsyncFlags := make([]string, 0)
	for _, flag := range c.operation.RsyncFlags {
		checkRsyncFlags = append(checkRsyncFlags, strings.Replace(flag, "-a", "-rc", -1))
	}

	checkRsyncStrFlags := strings.Join(checkRsyncFlags, " ")

	var toTransfer int64
	for _, command := range c.operation.Commands {
		toTransfer += command.Size
	}

	for _, command := range c.operation.Commands {
		args := append(
			checkRsyncFlags,
			command.Src,
			command.Dst,
		)
		// cmd := exec.Command("rsync", args)

		// cmd := fmt.Sprintf("rsync %s \"%s\" \"%s/\"", strings.Join(rsyncArgs, " "), filepath.Join(c.storage.SourceDiskName, item.Path), filepath.Join(disk.Path, filepath.Dir(item.Path)))
		cmd := fmt.Sprintf(`rsync %s %s %s`, checkRsyncStrFlags, strconv.Quote(command.Src), strconv.Quote(command.Dst))
		mlog.Info("cmd(%s)", cmd)

		outbound = &dto.Packet{Topic: "transferProgress", Payload: cmd}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		err := lib.ShellEx(func(text string) {
			line := strings.TrimSpace(text)

			if len(line) <= 0 {
				return
			}

			// match := c.reProgress.FindStringSubmatch(line)
			// if match == nil {
			// 	// this is a regular output line from rsync, print it
			// 	mlog.Info("%s", line)

			// 	outbound := &dto.Packet{Topic: "transferProgress", Payload: line}
			// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			// 	return
			// }

			// // this is a file transfer progress output line
			// if match[1] == "" {
			// 	// this happens when the file hasn't finished transferring
			// 	moved := strings.Replace(match[2], ",", "", -1)
			// 	deltaMoved, _ = strconv.ParseInt(moved, 10, 64)
			// } else {
			// 	// the file has finished transferring
			// 	moved := strings.Replace(match[1], ",", "", -1)
			// 	deltaMoved, _ = strconv.ParseInt(moved, 10, 64)
			// 	bytesTransferred += deltaMoved
			// }

			// percent, left, speed := progress(c.operation.BytesToTransfer, bytesTransferred+deltaMoved, c.operation.Started)

			// msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)

			// outbound := &dto.Packet{Topic: "progressStats", Payload: msg}
			// c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		}, c.operation.SourceDiskName, "rsync", args...)

		if err != nil {
			finished := time.Now()
			elapsed := time.Since(c.operation.Started)

			subject := "unBALANCE - VALIDATE operation INTERRUPTED"

			headline := fmt.Sprintf("Validate command (%s) was interrupted: %s", cmd, err.Error()+" : "+getError(err.Error(), c.reRsync, c.rsyncErrors))

			mlog.Warning(headline)
			outbound := &dto.Packet{Topic: "opError", Payload: "Validate operation was interrupted. Check log (/boot/logs/unbalance.log) for details."}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			c.finishTransferOperation(subject, headline, commandsExecuted, c.operation.Started, finished, elapsed)

			return
		}

		c.operation.BytesTransferred = c.operation.BytesTransferred + command.Size

		percent, left, speed := progress(c.operation.BytesToTransfer, c.operation.BytesTransferred, c.operation.Started)

		msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)
		mlog.Info("Current progress: %s", msg)

		outbound := &dto.Packet{Topic: "progressStats", Payload: msg}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		commandsExecuted = append(commandsExecuted, cmd)
	}

	finished := time.Now()
	elapsed := time.Since(c.operation.Started)

	// finish validate operation
	subject := "unBALANCE - VALIDATE operation completed"
	headline := "Validate operation has finished"

	c.finishTransferOperation(subject, headline, commandsExecuted, c.operation.Started, finished, elapsed)

}

// func (c *Core) finishValidateOperation(subject, headline string, commands []string, started, finished time.Time, elapsed time.Duration) {
// 	fstarted := started.Format(timeFormat)
// 	ffinished := finished.Format(timeFormat)
// 	elapsed = lib.Round(time.Since(started), time.Millisecond)

// 	outbound := &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: headline}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: "These are the commands that were executed:"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	printedCommands := ""
// 	for _, command := range commands {
// 		printedCommands += command + "\n"
// 		outbound = &dto.Packet{Topic: "transferProgress", Payload: command}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	}

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: "Operation Finished"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// send to front end the signal of operation finished
// 	if c.settings.DryRun {
// 		outbound = &dto.Packet{Topic: "transferProgress", Payload: "--- IT WAS A DRY RUN ---"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	}

// 	outbound = &dto.Packet{Topic: "transferFinished", Payload: ""}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s\n\n%s", fstarted, ffinished, elapsed, headline)
// 	switch c.settings.NotifyMove {
// 	case 1:
// 		message += fmt.Sprintf("\n\n%d commands were executed.", len(commands))
// 	case 2:
// 		message += "\n\nThese are the commands that were executed:\n\n" + printedCommands
// 	}

// 	if sendErr := c.sendmail(c.settings.NotifyCalc, subject, message, c.settings.DryRun); sendErr != nil {
// 		mlog.Error(sendErr)
// 	}

// 	mlog.Info(subject)
// 	mlog.Info(message)
// }

func (c *Core) getLog(msg *pubsub.Message) {
	log := c.storage.GetLog()

	outbound := &dto.Packet{Topic: "gotLog", Payload: log}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	return
}

func (c *Core) sendmail(notify int, subject, message string, dryRun bool) (err error) {
	if notify == 0 {
		return nil
	}

	dry := ""
	if dryRun {
		dry = "-------\nDRY RUN\n-------\n"
	}

	msg := dry + message

	// strCmd := fmt.Sprintf("-s \"%s\" -m \"%s\"", mailCmd, subject, msg)
	cmd := exec.Command(mailCmd, "-e", "unBALANCE operation update", "-s", subject, "-m", msg)
	err = cmd.Run()

	return
}

func progress(bytesToTransfer, bytesTransferred int64, started time.Time) (percent float64, left time.Duration, speed float64) {
	delta := time.Since(started)

	bytesPerSec := float64(bytesTransferred) / delta.Seconds()
	speed = bytesPerSec / 1024 / 1024 // MB/s

	percent = (float64(bytesTransferred) / float64(bytesToTransfer)) * 100 // %

	left = time.Duration(float64(bytesToTransfer-bytesTransferred)/bytesPerSec) * time.Second

	return
}

func getError(line string, re *regexp.Regexp, errors map[int]string) string {
	result := re.FindStringSubmatch(line)
	status, _ := strconv.Atoi(result[1])
	msg, ok := errors[status]
	if !ok {
		msg = "unknown error"
	}

	return msg
}
