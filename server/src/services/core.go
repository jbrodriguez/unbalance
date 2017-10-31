package services

import (
	"encoding/json"
	"regexp"
	"strings"

	"jbrodriguez/unbalance/server/src/common"
	"jbrodriguez/unbalance/server/src/domain"
	"jbrodriguez/unbalance/server/src/dto"
	"jbrodriguez/unbalance/server/src/lib"

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
	settings *lib.Settings
	actor    *actor.Actor

	// this holds the state the app
	state *domain.State

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
		actor:    actor.NewActor(bus),
		state:    &domain.State{},
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

	return core
}

// Start -
func (c *Core) Start() (err error) {
	mlog.Info("starting service Core ...")

	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	c.bus.Pub(msg, common.GET_ARRAY_STATUS)
	reply := <-msg.Reply
	message := reply.(dto.Message)
	if message.Error != nil {
		return message.Error
	}

	c.state.Status = common.StateIdle
	c.state.Unraid = message.Data.(*domain.Unraid)
	c.state.Operations = make([]*domain.Operation, 0)
	c.state.Operation = resetOp(c.state.Unraid.Disks)

	c.actor.Register(common.API_GET_CONFIG, c.getConfig)
	c.actor.Register(common.API_GET_STATUS, c.getStatus)
	c.actor.Register(common.API_GET_STATE, c.getState)
	c.actor.Register(common.API_RESET_OP, c.resetOp)

	// c.actor.Register("/config/set/notifyCalc", c.setNotifyCalc)
	// c.actor.Register("/config/set/notifyMove", c.setNotifyMove)
	// c.actor.Register("/config/set/reservedSpace", c.setReservedSpace)
	// c.actor.Register("/config/set/verbosity", c.setVerbosity)
	// c.actor.Register("/config/set/checkUpdate", c.setCheckUpdate)
	// c.actor.Register("/get/storage", c.getStorage)
	// c.actor.Register("/get/update", c.getUpdate)
	// c.actor.Register("/config/toggle/dryRun", c.toggleDryRun)
	// c.actor.Register("/get/tree", c.getTree)
	// c.actor.Register("/disks/locate", c.locate)
	// c.actor.Register("/config/set/rsyncFlags", c.setRsyncFlags)

	// c.actor.Register("calculate", c.calc)
	// c.actor.Register("move", c.move)
	// c.actor.Register("copy", c.copy)
	// c.actor.Register("validate", c.validate)
	// c.actor.Register("getLog", c.getLog)
	// c.actor.Register("findTargets", c.findTargets)
	// c.actor.Register("gather", c.gather)

	go c.actor.React()

	return nil
}

// Stop -
func (c *Core) Stop() {
	mlog.Info("stopped service Core ...")
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

func (c *Core) getStatus(msg *pubsub.Message) {
	mlog.Info("Sending status")

	msg.Reply <- c.state.Status
}

func (c *Core) getState(msg *pubsub.Message) {
	mlog.Info("Sending state")

	msg.Reply <- c.state
}

func (c *Core) resetOp(msg *pubsub.Message) {
	mlog.Info("resetting op")

	c.state.Operation = resetOp(c.state.Unraid.Disks)

	msg.Reply <- c.state.Operation
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

func (c *Core) setVerbosity(msg *pubsub.Message) {
	fverbosity := msg.Payload.(float64)
	verbosity := int(fverbosity)

	mlog.Info("Setting verbosity to (%d)", verbosity)

	c.settings.Verbosity = verbosity
	err := c.settings.Save()
	if err != nil {
		mlog.Warning("not right %s", err)
	}

	msg.Reply <- &c.settings.Config
}

func (c *Core) setCheckUpdate(msg *pubsub.Message) {
	fupdate := msg.Payload.(float64)
	update := int(fupdate)

	mlog.Info("Setting checkForUpdate to (%d)", update)

	c.settings.CheckForUpdate = update
	err := c.settings.Save()
	if err != nil {
		mlog.Warning("not right %s", err)
	}

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

func resetOp(disks []*domain.Disk) *domain.Operation {
	op := &domain.Operation{
		OpKind: common.OP_NEUTRAL,
		VDisks: make(map[string]*domain.VDisk, 0),
	}

	for _, disk := range disks {
		vdisk := &domain.VDisk{Path: disk.Path, PlannedFree: disk.Free, Src: false, Dst: false}
		op.VDisks[disk.Path] = vdisk
	}

	return op
}

// func (c *Core) getStorage(msg *pubsub.Message) {
// 	var stats string

// 	if c.operation.OpState == model.StateIdle {
// 		c.storage.Refresh()
// 	} else if c.operation.OpState == model.StateMove || c.operation.OpState == model.StateCopy || c.operation.OpState == model.StateGather {
// 		percent, left, speed := progress(c.operation.BytesToTransfer, c.operation.BytesTransferred, time.Since(c.operation.Started))
// 		stats = fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)
// 	}

// 	c.storage.Stats = stats
// 	c.storage.OpState = c.operation.OpState
// 	c.storage.PrevState = c.operation.PrevState
// 	c.storage.BytesToTransfer = c.operation.BytesToTransfer

// 	msg.Reply <- c.storage
// }

// func (c *Core) getUpdate(msg *pubsub.Message) {
// 	var newest string

// 	if c.settings.CheckForUpdate == 1 {
// 		latest, err := lib.GetLatestVersion("https://raw.githubusercontent.com/jbrodriguez/unbalance/master/VERSION")
// 		if err != nil {
// 			return
// 		}

// 		latest = strings.TrimSuffix(latest, "\n")
// 		if version.Compare(latest, c.settings.Version, ">") {
// 			newest = latest
// 		}
// 	}

// 	msg.Reply <- newest
// }

// func (c *Core) toggleDryRun(msg *pubsub.Message) {
// 	mlog.Info("Toggling dryRun from (%t)", c.settings.DryRun)

// 	c.settings.ToggleDryRun()
// 	c.settings.Save()

// 	msg.Reply <- &c.settings.Config
// }

// func (c *Core) getTree(msg *pubsub.Message) {
// 	path := msg.Payload.(string)

// 	msg.Reply <- c.storage.GetTree(path)
// }

// func (c *Core) locate(msg *pubsub.Message) {
// 	chosen := msg.Payload.([]string)
// 	msg.Reply <- c.storage.Locate(chosen)
// }

// func (c *Core) setRsyncFlags(msg *pubsub.Message) {
// 	// mlog.Warning("payload: %+v", msg.Payload)
// 	payload, ok := msg.Payload.(string)
// 	if !ok {
// 		mlog.Warning("Unable to convert Rsync Flags parameters")
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert Rsync Flags parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		msg.Reply <- &c.settings.Config

// 		return
// 	}

// 	var rsync dto.Rsync
// 	err := json.Unmarshal([]byte(payload), &rsync)
// 	if err != nil {
// 		mlog.Warning("Unable to bind rsyncFlags parameters: %s", err)
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind rsyncFlags parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 		// mlog.Fatalf(err.Error())
// 	}

// 	mlog.Info("Setting rsyncFlags to (%s)", strings.Join(rsync.Flags, " "))

// 	c.settings.RsyncFlags = rsync.Flags
// 	c.settings.Save()

// 	msg.Reply <- &c.settings.Config
// }

// func (c *Core) calc(msg *pubsub.Message) {
// 	c.operation = model.Operation{OpState: model.StateCalc, PrevState: model.StateIdle}

// 	go c._calc(msg)
// }

// func (c *Core) _calc(msg *pubsub.Message) {
// 	defer func() { c.operation.OpState = model.StateIdle }()

// 	payload, ok := msg.Payload.(string)
// 	if !ok {
// 		mlog.Warning("Unable to convert calculate parameters")
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert calculate parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 	}

// 	var dtoCalc dto.Calculate
// 	err := json.Unmarshal([]byte(payload), &dtoCalc)
// 	if err != nil {
// 		mlog.Warning("Unable to bind calculate parameters: %s", err)
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind calculate parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 		// mlog.Fatalf(err.Error())
// 	}

// 	mlog.Info("Running calculate operation ...")
// 	c.operation.Started = time.Now()

// 	outbound := &dto.Packet{Topic: "calcStarted", Payload: "Operation started"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	disks := make([]*model.Disk, 0)

// 	// create array of destination disks
// 	var srcDisk *model.Disk
// 	for _, disk := range c.storage.Disks {
// 		// reset disk
// 		disk.NewFree = disk.Free
// 		disk.Bin = nil
// 		disk.Src = false
// 		disk.Dst = dtoCalc.DestDisks[disk.Path]

// 		if disk.Path == dtoCalc.SourceDisk {
// 			disk.Src = true
// 			srcDisk = disk
// 		} else {
// 			// add it to the target disk list, only if the user selected it
// 			if val, ok := dtoCalc.DestDisks[disk.Path]; ok && val {
// 				// double check, if it's a cache disk, make sure it's the main cache disk
// 				if disk.Type == "Cache" && len(disk.Name) > 5 {
// 					continue
// 				}

// 				disks = append(disks, disk)
// 			}
// 		}
// 	}

// 	mlog.Info("_calc:Begin:srcDisk(%s); dstDisks(%d)", srcDisk.Path, len(disks))

// 	for _, disk := range disks {
// 		mlog.Info("_calc:elegibleDestDisk(%s)", disk.Path)
// 	}

// 	sort.Sort(model.ByFree(disks))

// 	srcDiskWithoutMnt := srcDisk.Path[5:]

// 	owner := ""
// 	lib.Shell("id -un", mlog.Warning, "owner", "", func(line string) {
// 		owner = line
// 	})

// 	group := ""
// 	lib.Shell("id -gn", mlog.Warning, "group", "", func(line string) {
// 		group = line
// 	})

// 	c.operation.OwnerIssue = 0
// 	c.operation.GroupIssue = 0
// 	c.operation.FolderIssue = 0
// 	c.operation.FileIssue = 0

// 	// Check permission and gather folders to be transferred from
// 	// source disk
// 	folders := make([]*model.Item, 0)
// 	for _, path := range dtoCalc.Folders {
// 		msg := fmt.Sprintf("Scanning %s on %s", path, srcDiskWithoutMnt)
// 		outbound = &dto.Packet{Topic: "calcProgress", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		c.checkOwnerAndPermissions(&c.operation, dtoCalc.SourceDisk, path, owner, group)

// 		msg = "Checked permissions ..."
// 		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		mlog.Info("_calc:%s", msg)

// 		list := c.getFolders(dtoCalc.SourceDisk, path)
// 		if list != nil {
// 			folders = append(folders, list...)
// 		}
// 	}

// 	mlog.Info("_calc:totalIssues=owner(%d),group(%d),folder(%d),file(%d)", c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)
// 	mlog.Info("_calc:foldersToBeTransferredTotal(%d)", len(folders))

// 	for _, v := range folders {
// 		mlog.Info("_calc:toBeTransferred:Path(%s); Size(%s)", v.Path, lib.ByteSize(v.Size))
// 	}

// 	willBeTransferred := make([]*model.Item, 0)
// 	if len(folders) > 0 {
// 		// Initialize fields
// 		// c.storage.BytesToTransfer = 0
// 		// c.storage.SourceDiskName = srcDisk.Path
// 		c.operation.BytesToTransfer = 0
// 		c.operation.SourceDiskName = srcDisk.Path

// 		for _, disk := range disks {
// 			diskWithoutMnt := disk.Path[5:]
// 			msg := fmt.Sprintf("Trying to allocate folders to %s ...", diskWithoutMnt)
// 			outbound = &dto.Packet{Topic: "calcProgress", Payload: msg}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			mlog.Info("_calc:%s", msg)
// 			// time.Sleep(2 * time.Second)

// 			if disk.Path != srcDisk.Path {
// 				// disk.NewFree = disk.Free

// 				var reserved int64
// 				switch c.settings.ReservedUnit {
// 				case "%":
// 					fcalc := disk.Size * c.settings.ReservedAmount / 100
// 					reserved = int64(fcalc)
// 					break
// 				case "Mb":
// 					reserved = c.settings.ReservedAmount * 1000 * 1000
// 					break
// 				case "Gb":
// 					reserved = c.settings.ReservedAmount * 1000 * 1000 * 1000
// 					break
// 				default:
// 					reserved = lib.ReservedSpace
// 				}

// 				ceil := lib.Max(lib.ReservedSpace, reserved)
// 				mlog.Info("_calc:FoldersLeft(%d):ReservedSpace(%d)", len(folders), ceil)

// 				packer := algorithm.NewKnapsack(disk, folders, ceil)
// 				bin := packer.BestFit()
// 				if bin != nil {
// 					srcDisk.NewFree += bin.Size
// 					disk.NewFree -= bin.Size
// 					c.operation.BytesToTransfer += bin.Size
// 					// c.storage.BytesToTransfer += bin.Size

// 					willBeTransferred = append(willBeTransferred, bin.Items...)
// 					folders = c.removeFolders(folders, bin.Items)

// 					mlog.Info("_calc:BinAllocated=[Disk(%s); Items(%d)];Freespace=[original(%s); final(%s)]", disk.Path, len(bin.Items), lib.ByteSize(srcDisk.Free), lib.ByteSize(srcDisk.NewFree))
// 				} else {
// 					mlog.Info("_calc:NoBinAllocated=Disk(%s)", disk.Path)
// 				}
// 			}
// 		}
// 	}

// 	c.operation.Finished = time.Now()
// 	elapsed := lib.Round(time.Since(c.operation.Started), time.Millisecond)

// 	fstarted := c.operation.Started.Format(timeFormat)
// 	ffinished := c.operation.Finished.Format(timeFormat)

// 	// Send to frontend console started/ended/elapsed times
// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	if len(willBeTransferred) == 0 {
// 		mlog.Info("_calc:No folders can be transferred.")
// 	} else {
// 		mlog.Info("_calc:%d folders will be transferred.", len(willBeTransferred))
// 		for _, folder := range willBeTransferred {
// 			mlog.Info("_calc:willBeTransferred(%s)", folder.Path)
// 		}
// 	}

// 	// send to frontend the folders that will not be transferred, if any
// 	// notTransferred holds a string representation of all the folders, separated by a '\n'
// 	c.operation.FoldersNotTransferred = make([]string, 0)
// 	notTransferred := ""
// 	if len(folders) > 0 {
// 		outbound = &dto.Packet{Topic: "calcProgress", Payload: "The following folders will not be transferred, because there's not enough space in the target disks:\n"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		mlog.Info("_calc:%d folders will NOT be transferred.", len(folders))
// 		for _, folder := range folders {
// 			c.operation.FoldersNotTransferred = append(c.operation.FoldersNotTransferred, folder.Path)

// 			notTransferred += folder.Path + "\n"

// 			outbound = &dto.Packet{Topic: "calcProgress", Payload: folder.Path}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			mlog.Info("_calc:notTransferred(%s)", folder.Path)
// 		}
// 	}

// 	// send mail according to user preferences
// 	subject := "unBALANCE - CALCULATE operation completed"
// 	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", fstarted, ffinished, elapsed)
// 	if notTransferred != "" {
// 		switch c.settings.NotifyCalc {
// 		case 1:
// 			message += "\n\nSome folders will not be transferred because there's not enough space for them in any of the destination disks."
// 		case 2:
// 			message += "\n\nThe following folders will not be transferred because there's not enough space for them in any of the destination disks:\n\n" + notTransferred
// 		}
// 	}

// 	if c.operation.OwnerIssue > 0 || c.operation.GroupIssue > 0 || c.operation.FolderIssue > 0 || c.operation.FileIssue > 0 {
// 		message += fmt.Sprintf(`
// 			\n\nThere are some permission issues:
// 			\n\n%d file(s)/folder(s) with an owner other than 'nobody'
// 			\n%d file(s)/folder(s) with a group other than 'users'
// 			\n%d folder(s) with a permission other than 'drwxrwxrwx'
// 			\n%d files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'
// 			\n\nCheck the log file (/boot/logs/unbalance.log) for additional information
// 			\n\nIt's strongly suggested to install the Fix Common Plugins and run the Docker Safe New Permissions command
// 		`, c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)
// 	}

// 	if sendErr := c.sendmail(c.settings.NotifyCalc, subject, message, false); sendErr != nil {
// 		mlog.Error(sendErr)
// 	}

// 	// some local logging
// 	mlog.Info("_calc:FoldersLeft(%d)", len(folders))
// 	mlog.Info("_calc:src(%s):Listing (%d) disks ...", srcDisk.Path, len(c.storage.Disks))
// 	for _, disk := range c.storage.Disks {
// 		// mlog.Info("the mystery of the year(%s)", disk.Path)
// 		disk.Print()
// 	}

// 	mlog.Info("=========================================================")
// 	mlog.Info("Results for %s", srcDisk.Path)
// 	mlog.Info("Original Free Space: %s", lib.ByteSize(srcDisk.Free))
// 	mlog.Info("Final Free Space: %s", lib.ByteSize(srcDisk.NewFree))
// 	mlog.Info("Gained Space: %s", lib.ByteSize(srcDisk.NewFree-srcDisk.Free))
// 	mlog.Info("Bytes To Move: %s", lib.ByteSize(c.operation.BytesToTransfer))
// 	mlog.Info("---------------------------------------------------------")

// 	c.storage.Print()
// 	// msg.Reply <- c.storage

// 	mlog.Info("_calc:End:srcDisk(%s)", srcDisk.Path)

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: "Operation Finished"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	c.storage.BytesToTransfer = c.operation.BytesToTransfer
// 	c.storage.OpState = c.operation.OpState
// 	c.storage.PrevState = c.operation.PrevState

// 	// send to front end the signal of operation finished
// 	outbound = &dto.Packet{Topic: "calcFinished", Payload: c.storage}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
// 	// and there actually perm issues
// 	if c.operation.BytesToTransfer > 0 && (c.operation.OwnerIssue+c.operation.GroupIssue+c.operation.FolderIssue+c.operation.FileIssue > 0) {
// 		outbound = &dto.Packet{Topic: "calcPermIssue", Payload: fmt.Sprintf("%d|%d|%d|%d", c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	}
// }

// func (c *Core) getFolders(src string, folder string) (items []*model.Item) {
// 	srcFolder := filepath.Join(src, folder)

// 	mlog.Info("getFolders:Scanning disk(%s):folder(%s)", src, folder)

// 	var fi os.FileInfo
// 	var err error
// 	if fi, err = os.Stat(srcFolder); os.IsNotExist(err) {
// 		mlog.Warning("getFolders:Folder does not exist: %s", srcFolder)
// 		return nil
// 	}

// 	if !fi.IsDir() {
// 		mlog.Info("getFolder-found(%s)-size(%d)", srcFolder, fi.Size())

// 		item := &model.Item{Name: folder, Size: fi.Size(), Path: folder, Location: src}
// 		items = append(items, item)

// 		msg := fmt.Sprintf("Found %s (%s)", item.Name, lib.ByteSize(item.Size))
// 		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		return
// 	}

// 	dirs, err := ioutil.ReadDir(srcFolder)
// 	if err != nil {
// 		mlog.Warning("getFolders:Unable to readdir: %s", err)
// 	}

// 	mlog.Info("getFolders:Readdir(%d)", len(dirs))

// 	if len(dirs) == 0 {
// 		mlog.Info("getFolders:No subdirectories under %s", srcFolder)
// 		return nil
// 	}

// 	scanFolder := srcFolder + "/."
// 	cmdText := fmt.Sprintf("find \"%s\" ! -name . -prune -exec du -bs {} +", scanFolder)

// 	mlog.Info("getFolders:Executing %s", cmdText)

// 	lib.Shell(cmdText, mlog.Warning, "getFolders:find/du:", "", func(line string) {
// 		mlog.Info("getFolders:find(%s): %s", scanFolder, line)

// 		result := c.reItems.FindStringSubmatch(line)
// 		// mlog.Info("[%s] %s", result[1], result[2])

// 		size, _ := strconv.ParseInt(result[1], 10, 64)

// 		item := &model.Item{Name: result[2], Size: size, Path: filepath.Join(folder, filepath.Base(result[2])), Location: src}
// 		items = append(items, item)

// 		msg := fmt.Sprintf("Found %s (%s)", filepath.Base(item.Name), lib.ByteSize(size))
// 		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	})

// 	return
// }

// func (c *Core) checkOwnerAndPermissions(operation *model.Operation, src, folder, ownerName, groupName string) {
// 	srcFolder := filepath.Join(src, folder)

// 	// outbound := &dto.Packet{Topic: "calcProgress", Payload: "Checking permissions ..."}
// 	// c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	ownerIssue, groupIssue, folderIssue, fileIssue := 0, 0, 0, 0

// 	mlog.Info("perms:Scanning disk(%s):folder(%s)", src, folder)

// 	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
// 		mlog.Warning("perms:Folder does not exist: %s", srcFolder)
// 		return
// 	}

// 	scanFolder := srcFolder + "/."
// 	cmdText := fmt.Sprintf(`find "%s" -exec stat --format "%%A|%%U:%%G|%%F|%%n" {} \;`, scanFolder)

// 	mlog.Info("perms:Executing %s", cmdText)

// 	lib.Shell(cmdText, mlog.Warning, "perms:find/stat:", "", func(line string) {
// 		result := c.reStat.FindStringSubmatch(line)
// 		if result == nil {
// 			mlog.Warning("perms:Unable to parse (%s)", line)
// 			return
// 		}

// 		u := result[1]
// 		g := result[2]
// 		o := result[3]
// 		user := result[4]
// 		group := result[5]
// 		kind := result[6]
// 		name := result[7]

// 		perms := u + g + o

// 		if user != "nobody" {
// 			if c.settings.Verbosity == 1 {
// 				mlog.Info("perms:User != nobody: [%s]: %s", user, name)
// 			}

// 			operation.OwnerIssue++
// 			ownerIssue++
// 		}

// 		if group != "users" {
// 			if c.settings.Verbosity == 1 {
// 				mlog.Info("perms:Group != users: [%s]: %s", group, name)
// 			}

// 			operation.GroupIssue++
// 			groupIssue++
// 		}

// 		if kind == "directory" {
// 			if perms != "rwxrwxrwx" {
// 				if c.settings.Verbosity == 1 {
// 					mlog.Info("perms:Folder perms != rwxrwxrwx: [%s]: %s", perms, name)
// 				}

// 				operation.FolderIssue++
// 				folderIssue++
// 			}
// 		} else {
// 			match := strings.Compare(perms, "r--r--r--") == 0 || strings.Compare(perms, "rw-rw-rw-") == 0
// 			if !match {
// 				if c.settings.Verbosity == 1 {
// 					mlog.Info("perms:File perms != rw-rw-rw- or r--r--r--: [%s]: %s", perms, name)
// 				}

// 				operation.FileIssue++
// 				fileIssue++
// 			}
// 		}
// 	})

// 	mlog.Info("perms:issues=owner(%d),group(%d),folder(%d),file(%d)", ownerIssue, groupIssue, folderIssue, fileIssue)

// 	return
// }

// func (c *Core) removeFolders(folders []*model.Item, list []*model.Item) []*model.Item {
// 	w := 0 // write index

// loop:
// 	for _, fld := range folders {
// 		for _, itm := range list {
// 			if itm.Name == fld.Name {
// 				continue loop
// 			}
// 		}
// 		folders[w] = fld
// 		w++
// 	}

// 	return folders[:w]
// }

// func (c *Core) move(msg *pubsub.Message) {
// 	c.operation.OpState = model.StateMove
// 	c.operation.PrevState = model.StateMove
// 	go c.transfer("Move", false, msg)
// }

// func (c *Core) copy(msg *pubsub.Message) {
// 	c.operation.OpState = model.StateCopy
// 	c.operation.PrevState = model.StateCopy
// 	go c.transfer("Copy", false, msg)
// }

// func (c *Core) validate(msg *pubsub.Message) {
// 	c.operation.OpState = model.StateValidate
// 	c.operation.PrevState = model.StateValidate
// 	go c.checksum(msg)
// }

// func (c *Core) transfer(opName string, multiSource bool, msg *pubsub.Message) {
// 	defer func() {
// 		c.operation.OpState = model.StateIdle
// 		c.operation.Started = time.Time{}
// 		c.operation.BytesTransferred = 0
// 		c.operation.Target = ""
// 	}()

// 	mlog.Info("Running %s operation ...", opName)
// 	c.operation.Started = time.Now()

// 	outbound := &dto.Packet{Topic: "transferStarted", Payload: "Operation started"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "progressStats", Payload: "Waiting to collect stats ..."}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// user may have changed rsync flags or dry-run setting, adjust for it
// 	c.operation.RsyncFlags = c.settings.RsyncFlags
// 	c.operation.DryRun = c.settings.DryRun
// 	if c.operation.DryRun {
// 		c.operation.RsyncFlags = append(c.operation.RsyncFlags, "--dry-run")
// 		mlog.Info("dry-run ON")
// 	} else {
// 		mlog.Info("dry-run OFF")
// 	}
// 	c.operation.RsyncStrFlags = strings.Join(c.operation.RsyncFlags, " ")

// 	workdir := c.operation.SourceDiskName

// 	c.operation.Commands = make([]model.Command, 0)
// 	for _, disk := range c.storage.Disks {
// 		if disk.Bin == nil || disk.Src {
// 			continue
// 		}

// 		for _, item := range disk.Bin.Items {
// 			var src, dst string
// 			if strings.Contains(c.operation.RsyncStrFlags, "R") {
// 				if item.Path[0] == filepath.Separator {
// 					src = item.Path[1:]
// 				} else {
// 					src = item.Path
// 				}

// 				dst = disk.Path + string(filepath.Separator)
// 			} else {
// 				src = filepath.Join(c.operation.SourceDiskName, item.Path)
// 				dst = filepath.Join(disk.Path, filepath.Dir(item.Path)) + string(filepath.Separator)
// 			}

// 			if multiSource {
// 				workdir = item.Location
// 			}

// 			c.operation.Commands = append(c.operation.Commands, model.Command{
// 				Src:     src,
// 				Dst:     dst,
// 				Path:    item.Path,
// 				Size:    item.Size,
// 				WorkDir: workdir,
// 			})
// 		}
// 	}

// 	if c.settings.NotifyMove == 2 {
// 		c.notifyCommandsToRun(opName)
// 	}

// 	// execute each rsync command created in the step above
// 	c.runOperation(opName, c.operation.RsyncFlags, c.operation.RsyncStrFlags, multiSource)
// }

// func (c *Core) checksum(msg *pubsub.Message) {
// 	defer func() {
// 		c.operation.OpState = model.StateIdle
// 		c.operation.PrevState = model.StateIdle
// 		c.operation.Started = time.Time{}
// 		c.operation.BytesTransferred = 0
// 	}()

// 	opName := "Validate"

// 	mlog.Info("Running %s operation ...", opName)
// 	c.operation.Started = time.Now()

// 	outbound := &dto.Packet{Topic: "transferStarted", Payload: "Operation started"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	multiSource := false

// 	if !strings.HasPrefix(c.operation.RsyncStrFlags, "-a") {
// 		finished := time.Now()
// 		elapsed := time.Since(c.operation.Started)

// 		subject := fmt.Sprintf("unBALANCE - %s operation INTERRUPTED", strings.ToUpper(opName))
// 		headline := fmt.Sprintf("For proper %s operation, rsync flags MUST begin with -a", opName)

// 		mlog.Warning(headline)
// 		outbound := &dto.Packet{Topic: "opError", Payload: fmt.Sprintf("%s operation was interrupted. Check log (/boot/logs/unbalance.log) for details.", opName)}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		_, _, speed := progress(c.operation.BytesToTransfer, 0, elapsed)
// 		c.finishTransferOperation(subject, headline, make([]string, 0), c.operation.Started, finished, elapsed, 0, speed, multiSource)

// 		return
// 	}

// 	// Initialize local variables
// 	// we use the rsync flags that were created by the transfer operation,
// 	// but replace -a with -rc, to perform the validation
// 	checkRsyncFlags := make([]string, 0)
// 	for _, flag := range c.operation.RsyncFlags {
// 		checkRsyncFlags = append(checkRsyncFlags, strings.Replace(flag, "-a", "-rc", -1))
// 	}

// 	checkRsyncStrFlags := strings.Join(checkRsyncFlags, " ")

// 	// execute each rsync command created in the transfer phase
// 	c.runOperation(opName, checkRsyncFlags, checkRsyncStrFlags, multiSource)
// }

// func (c *Core) runOperation(opName string, rsyncFlags []string, rsyncStrFlags string, multiSource bool) {
// 	// Initialize local variables
// 	var calls int64
// 	var callsPerDelta int64

// 	var finished time.Time
// 	var elapsed time.Duration

// 	commandsExecuted := make([]string, 0)

// 	c.operation.BytesTransferred = 0

// 	for _, command := range c.operation.Commands {
// 		args := append(
// 			rsyncFlags,
// 			command.Src,
// 			command.Dst,
// 		)
// 		cmd := fmt.Sprintf(`rsync %s %s %s`, rsyncStrFlags, strconv.Quote(command.Src), strconv.Quote(command.Dst))
// 		mlog.Info("Command Started: (src: %s) %s ", command.WorkDir, cmd)

// 		outbound := &dto.Packet{Topic: "transferProgress", Payload: cmd}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		bytesTransferred := c.operation.BytesTransferred

// 		var deltaMoved int64

// 		// actual shell execution
// 		err := lib.ShellEx(func(text string) {
// 			line := strings.TrimSpace(text)

// 			if len(line) <= 0 {
// 				return
// 			}

// 			if callsPerDelta <= 50 {
// 				calls++
// 			}

// 			delta := int64(time.Since(c.operation.Started) / time.Second)
// 			if delta == 0 {
// 				delta = 1
// 			}
// 			callsPerDelta = calls / delta

// 			match := c.reProgress.FindStringSubmatch(line)
// 			if match == nil {
// 				// this is a regular output line from rsync, print it
// 				// according to verbosity settings
// 				if c.settings.Verbosity == 1 {
// 					mlog.Info("%s", line)
// 				}

// 				if callsPerDelta <= 50 {
// 					outbound := &dto.Packet{Topic: "transferProgress", Payload: line}
// 					c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 				}

// 				return
// 			}

// 			// this is a file transfer progress output line
// 			if match[1] == "" {
// 				// this happens when the file hasn't finished transferring
// 				moved := strings.Replace(match[2], ",", "", -1)
// 				deltaMoved, _ = strconv.ParseInt(moved, 10, 64)
// 			} else {
// 				// the file has finished transferring
// 				moved := strings.Replace(match[1], ",", "", -1)
// 				deltaMoved, _ = strconv.ParseInt(moved, 10, 64)
// 				bytesTransferred += deltaMoved
// 			}

// 			percent, left, speed := progress(c.operation.BytesToTransfer, bytesTransferred+deltaMoved, time.Since(c.operation.Started))
// 			msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)

// 			if callsPerDelta <= 50 {
// 				outbound := &dto.Packet{Topic: "progressStats", Payload: msg}
// 				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			}

// 		}, mlog.Warning, command.WorkDir, "rsync", args...)

// 		finished = time.Now()
// 		elapsed = time.Since(c.operation.Started)

// 		if err != nil {
// 			subject := fmt.Sprintf("unBALANCE - %s operation INTERRUPTED", strings.ToUpper(opName))
// 			headline := fmt.Sprintf("Command Interrupted: %s (%s)", cmd, err.Error()+" : "+getError(err.Error(), c.reRsync, c.rsyncErrors))

// 			mlog.Warning(headline)
// 			outbound := &dto.Packet{Topic: "opError", Payload: fmt.Sprintf("%s operation was interrupted. Check log (/boot/logs/unbalance.log) for details.", opName)}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 			_, _, speed := progress(c.operation.BytesToTransfer, bytesTransferred+deltaMoved, elapsed)
// 			c.finishTransferOperation(subject, headline, commandsExecuted, c.operation.Started, finished, elapsed, bytesTransferred+deltaMoved, speed, multiSource)

// 			return
// 		}

// 		mlog.Info("Command Finished")

// 		c.operation.BytesTransferred = c.operation.BytesTransferred + command.Size
// 		percent, left, speed := progress(c.operation.BytesToTransfer, c.operation.BytesTransferred, elapsed)

// 		msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)
// 		mlog.Info("Current progress: %s", msg)

// 		outbound = &dto.Packet{Topic: "progressStats", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		commandsExecuted = append(commandsExecuted, cmd)

// 		if c.operation.DryRun && c.operation.OpState == model.StateGather {
// 			parent := filepath.Dir(command.Src)
// 			// mlog.Info("parent(%s)-workdir(%s)-src(%s)-dst(%s)-path(%s)", parent, command.WorkDir, command.Src, command.Dst, command.Path)
// 			if parent != "." {
// 				mlog.Info(`Would delete empty folders starting from (%s) - (find "%s" -type d -empty -prune -exec rm -rf {} \;) `, filepath.Join(command.WorkDir, parent), filepath.Join(command.WorkDir, parent))
// 			} else {
// 				mlog.Info(`WONT DELETE: find "%s" -type d -empty -prune -exec rm -rf {} \;`, filepath.Join(command.WorkDir, parent))
// 			}
// 		}

// 		// if it isn't a dry-run and the operation is Move or Gather, delete the source folder
// 		if !c.operation.DryRun && (c.operation.OpState == model.StateMove || c.operation.OpState == model.StateGather) {
// 			exists, _ := lib.Exists(filepath.Join(command.Dst, command.Src))
// 			if exists {
// 				rmrf := fmt.Sprintf("rm -rf \"%s\"", filepath.Join(command.WorkDir, command.Path))
// 				mlog.Info("Removing: %s", rmrf)
// 				err = lib.Shell(rmrf, mlog.Warning, "transferProgress:", "", func(line string) {
// 					mlog.Info(line)
// 				})

// 				if err != nil {
// 					msg := fmt.Sprintf("Unable to remove source folder (%s): %s", filepath.Join(command.WorkDir, command.Path), err)

// 					outbound := &dto.Packet{Topic: "transferProgress", Payload: msg}
// 					c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 					mlog.Warning(msg)
// 				}

// 				if c.operation.OpState == model.StateGather {
// 					parent := filepath.Dir(command.Src)
// 					if parent != "." {
// 						rmdir := fmt.Sprintf(`find "%s" -type d -empty -prune -exec rm -rf {} \;`, filepath.Join(command.WorkDir, parent))
// 						mlog.Info("Running %s", rmdir)

// 						err = lib.Shell(rmdir, mlog.Warning, "transferProgress:", "", func(line string) {
// 							mlog.Info(line)
// 						})

// 						if err != nil {
// 							msg := fmt.Sprintf("Unable to remove parent folder (%s): %s", filepath.Join(command.WorkDir, parent), err)

// 							outbound := &dto.Packet{Topic: "transferProgress", Payload: msg}
// 							c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 							mlog.Warning(msg)
// 						}
// 					}
// 				}
// 			} else {
// 				mlog.Warning("Skipping deletion (file/folder not present in destination): %s", filepath.Join(command.Dst, command.Src))
// 			}
// 		}
// 	}

// 	subject := fmt.Sprintf("unBALANCE - %s operation completed", strings.ToUpper(opName))
// 	headline := fmt.Sprintf("%s operation has finished", opName)

// 	_, _, speed := progress(c.operation.BytesToTransfer, c.operation.BytesTransferred, elapsed)
// 	c.finishTransferOperation(subject, headline, commandsExecuted, c.operation.Started, finished, elapsed, c.operation.BytesTransferred, speed, multiSource)
// }

// func (c *Core) finishTransferOperation(subject, headline string, commands []string, started, finished time.Time, elapsed time.Duration, transferred int64, speed float64, multiSource bool) {
// 	fstarted := started.Format(timeFormat)
// 	ffinished := finished.Format(timeFormat)
// 	elapsed = lib.Round(time.Since(started), time.Millisecond)

// 	outbound := &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Transferred %s at ~ %.2f MB/s", lib.ByteSize(transferred), speed)}
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

// 	finishMsg := "transferFinished"
// 	if multiSource {
// 		finishMsg = "gatherFinished"
// 	}

// 	outbound = &dto.Packet{Topic: finishMsg, Payload: ""}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s\n\n%s\n\nTransferred %s at ~ %.2f MB/s", fstarted, ffinished, elapsed, headline, lib.ByteSize(transferred), speed)
// 	switch c.settings.NotifyMove {
// 	case 1:
// 		message += fmt.Sprintf("\n\n%d commands were executed.", len(commands))
// 	case 2:
// 		message += "\n\nThese are the commands that were executed:\n\n" + printedCommands
// 	}

// 	go func() {
// 		if sendErr := c.sendmail(c.settings.NotifyMove, subject, message, c.settings.DryRun); sendErr != nil {
// 			mlog.Error(sendErr)
// 		}
// 	}()

// 	mlog.Info("\n%s\n%s", subject, message)
// }

// func (c *Core) findTargets(msg *pubsub.Message) {
// 	c.operation = model.Operation{OpState: model.StateFindTargets, PrevState: model.StateIdle}
// 	go c._findTargets(msg)
// }

// func (c *Core) _findTargets(msg *pubsub.Message) {
// 	defer func() { c.operation.OpState = model.StateIdle }()

// 	data, ok := msg.Payload.(string)
// 	if !ok {
// 		mlog.Warning("Unable to convert findTargets parameters")
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert findTargets parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 	}

// 	var chosen []string
// 	err := json.Unmarshal([]byte(data), &chosen)
// 	if err != nil {
// 		mlog.Warning("Unable to bind findTargets parameters: %s", err)
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind findTargets parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 		// mlog.Fatalf(err.Error())
// 	}

// 	mlog.Info("Running findTargets operation ...")
// 	c.operation.Started = time.Now()

// 	outbound := &dto.Packet{Topic: "calcStarted", Payload: "Operation started"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// disks := make([]*model.Disk, 0)
// 	c.storage.Refresh()

// 	owner := ""
// 	lib.Shell("id -un", mlog.Warning, "owner", "", func(line string) {
// 		owner = line
// 	})

// 	group := ""
// 	lib.Shell("id -gn", mlog.Warning, "group", "", func(line string) {
// 		group = line
// 	})

// 	c.operation.OwnerIssue = 0
// 	c.operation.GroupIssue = 0
// 	c.operation.FolderIssue = 0
// 	c.operation.FileIssue = 0

// 	entries := make([]*model.Item, 0)

// 	// Check permission and look for the chosen folders on every disk
// 	for _, disk := range c.storage.Disks {
// 		for _, path := range chosen {
// 			msg := fmt.Sprintf("Scanning %s on %s", path, disk.Path)
// 			outbound = &dto.Packet{Topic: "calcProgress", Payload: msg}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			mlog.Info("_find:%s", msg)

// 			c.checkOwnerAndPermissions(&c.operation, disk.Path, path, owner, group)

// 			msg = "Checked permissions ..."
// 			outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			mlog.Info("_find:%s", msg)

// 			list := c.getFolders(disk.Path, path)
// 			if list != nil {
// 				entries = append(entries, list...)
// 			}
// 		}
// 	}

// 	mlog.Info("_find:elegibleFolders(%d)", len(entries))

// 	var totalSize int64
// 	for _, entry := range entries {
// 		totalSize += entry.Size
// 		mlog.Info("_find:elegibleFolder:Location(%s); Size(%s)", filepath.Join(entry.Location, entry.Path), lib.ByteSize(entry.Size))
// 	}

// 	mlog.Info("_find:potentialSizeToBeTransferred(%s)", lib.ByteSize(totalSize))

// 	if len(entries) > 0 {
// 		// Initialize fields
// 		// c.storage.BytesToTransfer = 0
// 		// c.storage.SourceDiskName = srcDisk.Path
// 		c.operation.BytesToTransfer = 0
// 		// c.operation.SourceDiskName = mntUser

// 		for _, disk := range c.storage.Disks {
// 			diskWithoutMnt := disk.Path[5:]
// 			msg := fmt.Sprintf("Trying to allocate folders to %s ...", diskWithoutMnt)
// 			outbound = &dto.Packet{Topic: "calcProgress", Payload: msg}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			mlog.Info("_find:%s", msg)
// 			// time.Sleep(2 * time.Second)

// 			var reserved int64
// 			switch c.settings.ReservedUnit {
// 			case "%":
// 				fcalc := disk.Size * c.settings.ReservedAmount / 100
// 				reserved = int64(fcalc)
// 				break
// 			case "Mb":
// 				reserved = c.settings.ReservedAmount * 1000 * 1000
// 				break
// 			case "Gb":
// 				reserved = c.settings.ReservedAmount * 1000 * 1000 * 1000
// 				break
// 			default:
// 				reserved = lib.ReservedSpace
// 			}

// 			ceil := lib.Max(lib.ReservedSpace, reserved)
// 			mlog.Info("_find:FoldersLeft(%d):ReservedSpace(%d)", len(entries), ceil)

// 			packer := algorithm.NewGreedy(disk, entries, totalSize, ceil)
// 			bin := packer.FitAll()
// 			if bin != nil {
// 				disk.NewFree -= bin.Size
// 				disk.Src = false
// 				disk.Dst = false
// 				c.operation.BytesToTransfer += bin.Size
// 				mlog.Info("_find:BinAllocated=[Disk(%s); Items(%d)]", disk.Path, len(bin.Items))
// 			} else {
// 				mlog.Info("_find:NoBinAllocated=Disk(%s)", disk.Path)
// 			}
// 		}
// 	}

// 	c.operation.Finished = time.Now()
// 	elapsed := lib.Round(time.Since(c.operation.Started), time.Millisecond)

// 	fstarted := c.operation.Started.Format(timeFormat)
// 	ffinished := c.operation.Finished.Format(timeFormat)

// 	// Send to frontend console started/ended/elapsed times
// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// send to frontend the folders that will not be transferred, if any
// 	// notTransferred holds a string representation of all the folders, separated by a '\n'
// 	c.operation.FoldersNotTransferred = make([]string, 0)

// 	// send mail according to user preferences
// 	subject := "unBALANCE - CALCULATE operation completed"
// 	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", fstarted, ffinished, elapsed)

// 	if c.operation.OwnerIssue > 0 || c.operation.GroupIssue > 0 || c.operation.FolderIssue > 0 || c.operation.FileIssue > 0 {
// 		message += fmt.Sprintf(`
// 			\n\nThere are some permission issues:
// 			\n\n%d file(s)/folder(s) with an owner other than 'nobody'
// 			\n%d file(s)/folder(s) with a group other than 'users'
// 			\n%d folder(s) with a permission other than 'drwxrwxrwx'
// 			\n%d files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'
// 			\n\nCheck the log file (/boot/logs/unbalance.log) for additional information
// 			\n\nIt's strongly suggested to install the Fix Common Plugins and run the Docker Safe New Permissions command
// 		`, c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)
// 	}

// 	if sendErr := c.sendmail(c.settings.NotifyCalc, subject, message, false); sendErr != nil {
// 		mlog.Error(sendErr)
// 	}

// 	// some local logging
// 	mlog.Info("_find:Listing (%d) disks ...", len(c.storage.Disks))
// 	for _, disk := range c.storage.Disks {
// 		// mlog.Info("the mystery of the year(%s)", disk.Path)
// 		disk.Print()
// 	}

// 	c.storage.Print()

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: "Operation Finished"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	c.storage.BytesToTransfer = c.operation.BytesToTransfer
// 	c.storage.OpState = c.operation.OpState
// 	c.storage.PrevState = c.operation.PrevState

// 	// send to front end the signal of operation finished
// 	outbound = &dto.Packet{Topic: "findFinished", Payload: c.storage}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
// 	// and there actually perm issues
// 	if c.operation.BytesToTransfer > 0 && (c.operation.OwnerIssue+c.operation.GroupIssue+c.operation.FolderIssue+c.operation.FileIssue > 0) {
// 		outbound = &dto.Packet{Topic: "calcPermIssue", Payload: fmt.Sprintf("%d|%d|%d|%d", c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	}
// }

// func (c *Core) gather(msg *pubsub.Message) {
// 	mlog.Info("%+v", msg.Payload)
// 	data, ok := msg.Payload.(string)
// 	if !ok {
// 		mlog.Warning("Unable to convert gather parameters")
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert gather parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 	}

// 	var target model.Disk
// 	err := json.Unmarshal([]byte(data), &target)
// 	if err != nil {
// 		mlog.Warning("Unable to bind gather parameters: %s", err)
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind gather parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 		// mlog.Fatalf(err.Error())
// 	}

// 	// user chose a target disk, adjust bytestotransfer to the size of its bin, since
// 	// that's the amount of data we need to transfer. Also remove bin from all other disks,
// 	// since only the target will have work to do
// 	for _, disk := range c.storage.Disks {
// 		if disk.Path == target.Path {
// 			c.operation.BytesToTransfer = disk.Bin.Size
// 		} else {
// 			disk.Bin = nil
// 		}
// 	}

// 	c.operation.OpState = model.StateGather
// 	c.operation.PrevState = model.StateGather

// 	go c.transfer("Move", true, msg)
// }

// func (c *Core) getLog(msg *pubsub.Message) {
// 	log := c.storage.GetLog()

// 	outbound := &dto.Packet{Topic: "gotLog", Payload: log}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	return
// }

// func (c *Core) sendmail(notify int, subject, message string, dryRun bool) (err error) {
// 	if notify == 0 {
// 		return nil
// 	}

// 	dry := ""
// 	if dryRun {
// 		dry = "-------\nDRY RUN\n-------\n"
// 	}

// 	msg := dry + message

// 	// strCmd := fmt.Sprintf("-s \"%s\" -m \"%s\"", mailCmd, subject, msg)
// 	cmd := exec.Command(mailCmd, "-e", "unBALANCE operation update", "-s", subject, "-m", msg)
// 	err = cmd.Run()

// 	return
// }

// func progress(bytesToTransfer, bytesTransferred int64, elapsed time.Duration) (percent float64, left time.Duration, speed float64) {
// 	bytesPerSec := float64(bytesTransferred) / elapsed.Seconds()
// 	speed = bytesPerSec / 1024 / 1024 // MB/s

// 	percent = (float64(bytesTransferred) / float64(bytesToTransfer)) * 100 // %

// 	left = time.Duration(float64(bytesToTransfer-bytesTransferred)/bytesPerSec) * time.Second

// 	return
// }

// func getError(line string, re *regexp.Regexp, errors map[int]string) string {
// 	result := re.FindStringSubmatch(line)
// 	status, _ := strconv.Atoi(result[1])
// 	msg, ok := errors[status]
// 	if !ok {
// 		msg = "unknown error"
// 	}

// 	return msg
// }

// func (c *Core) notifyCommandsToRun(opName string) {
// 	message := "\n\nThe following commands will be executed:\n\n"

// 	for _, command := range c.operation.Commands {
// 		cmd := fmt.Sprintf(`(src: %s) rsync %s %s %s`, command.WorkDir, c.operation.RsyncStrFlags, strconv.Quote(command.Src), strconv.Quote(command.Dst))
// 		message += cmd + "\n"
// 	}

// 	subject := fmt.Sprintf("unBALANCE - %s operation STARTED", strings.ToUpper(opName))

// 	go func() {
// 		if sendErr := c.sendmail(c.settings.NotifyMove, subject, message, c.settings.DryRun); sendErr != nil {
// 			mlog.Error(sendErr)
// 		}
// 	}()
// }
