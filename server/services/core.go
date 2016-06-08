package services

import (
	"encoding/json"
	"fmt"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"io/ioutil"
	"jbrodriguez/unbalance/server/algorithm"
	"jbrodriguez/unbalance/server/dto"
	"jbrodriguez/unbalance/server/lib"
	"jbrodriguez/unbalance/server/model"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

const (
	MAIL_CMD    = "/usr/local/emhttp/webGui/scripts/notify"
	TIME_FORMAT = "Jan _2, 2006 15:04:05"
)

const (
	IDLE = 0
	CALC = 1
	MOVE = 2
)

type Core struct {
	Service

	bus      *pubsub.PubSub
	storage  *model.Unraid
	settings *lib.Settings

	opState uint64

	foldersNotMoved []string

	mailbox chan *pubsub.Mailbox

	reFreeSpace *regexp.Regexp
	reItems     *regexp.Regexp

	// diskmvLocation string

	bytesMoved int64
	started    time.Time
}

func NewCore(bus *pubsub.PubSub, settings *lib.Settings) *Core {
	core := &Core{
		bus:      bus,
		settings: settings,
		opState:  IDLE,
		storage:  &model.Unraid{},
	}
	core.init()

	re, _ := regexp.Compile(`(.*?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.*?)\s+(.*?)$`)
	core.reFreeSpace = re

	re, _ = regexp.Compile(`(\d+)\s+(.*?)$`)
	core.reItems = re

	return core
}

func (c *Core) Start() (err error) {
	mlog.Info("starting service Core ...")

	c.mailbox = c.register(c.bus, "/get/config", c.getConfig)
	c.registerAdditional(c.bus, "/config/set/notifyCalc", c.setNotifyCalc, c.mailbox)
	c.registerAdditional(c.bus, "/config/set/notifyMove", c.setNotifyMove, c.mailbox)
	c.registerAdditional(c.bus, "/config/set/reservedSpace", c.setReservedSpace, c.mailbox)
	c.registerAdditional(c.bus, "/config/add/folder", c.addFolder, c.mailbox)
	c.registerAdditional(c.bus, "/config/delete/folder", c.deleteFolder, c.mailbox)
	c.registerAdditional(c.bus, "/get/storage", c.getStorage, c.mailbox)
	c.registerAdditional(c.bus, "/config/toggle/dryRun", c.toggleDryRun, c.mailbox)
	c.registerAdditional(c.bus, "/get/tree", c.getTree, c.mailbox)
	c.registerAdditional(c.bus, "/config/set/rsyncFlags", c.setRsyncFlags, c.mailbox)

	c.registerAdditional(c.bus, "calculate", c.calc, c.mailbox)
	c.registerAdditional(c.bus, "move", c.move, c.mailbox)
	// c.registerAdditional(c.bus, "/set/config", c.setConfig)

	err = c.storage.SanityCheck(c.settings.ApiFolders)
	if err != nil {
		return err
	}

	// locations := []string{
	// 	"/usr/local/emhttp/plugins/unbalance",
	// 	".",
	// }

	// c.diskmvLocation = lib.SearchFile("diskmv", locations)
	// if c.diskmvLocation == "" {
	// 	msg := ""
	// 	for _, loc := range locations {
	// 		msg += fmt.Sprintf("%s, ", loc)
	// 	}
	// 	mlog.Fatalf("Unable to find diskmv. Exiting now. (searched in %s)", msg)
	// }

	go c.react()

	return nil
}

func (c *Core) Stop() {
	mlog.Info("stopped service Core ...")
}

func (c *Core) SetStorage(unraid *model.Unraid) {
	c.storage = unraid
}

func (c *Core) react() {
	for mbox := range c.mailbox {
		// mlog.Info("Core:Topic: %s", mbox.Topic)
		c.dispatch(mbox.Topic, mbox.Content)
	}
}

func (c *Core) getConfig(msg *pubsub.Message) {
	mlog.Info("Sending config")

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
		// mlog.Fatalf(err.Error())
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

func (c *Core) addFolder(msg *pubsub.Message) {

	folder := msg.Payload.(string)

	mlog.Info("Adding folder (%s)", folder)

	c.settings.AddFolder(folder)
	// c.config.Folders = config.Folders
	// c.config.DryRun = config.DryRun
	// c.config.Notifications = config.Notifications
	// c.Notifications = config.Notifications
	// c.NotiFrom = config.NotiFrom
	// c.NotiTo = config.NotiTo
	// c.NotiHost = config.NotiHost
	// c.NotiPort = config.NotiPort
	// c.NotiEncrypt = config.NotiEncrypt
	// c.NotiUser = config.NotiUser
	// c.NotiPassword = config.NotiPassword

	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) deleteFolder(msg *pubsub.Message) {

	folder := msg.Payload.(string)

	mlog.Info("Deleting folder (%s)", folder)

	c.settings.DeleteFolder(folder)

	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) getStorage(msg *pubsub.Message) {
	var stats string

	if c.opState == IDLE {
		c.storage.Refresh()
	} else if c.opState == MOVE {
		percent, left, speed := progress(c.storage.BytesToMove, c.bytesMoved, c.started)
		stats = fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)
	}

	c.storage.Stats = stats
	c.storage.OpState = c.opState
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
	c.opState = CALC
	go c._calc(msg)
}

func (c *Core) _calc(msg *pubsub.Message) {
	defer func() { c.opState = IDLE }()

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

	// dtoCalc, ok := msg.Payload.(*dto.Calculate)
	// if !ok {
	// 	mlog.Warning("Unable to convert calculate parameters")
	// 	outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert calculate parameters"}
	// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	// 	return
	// }

	mlog.Info("Running calculate operation ...")
	started := time.Now()

	// c.storage.Condition.State = "CALCULATING"

	outbound := &dto.Packet{Topic: "calcStarted", Payload: "Operation started"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// mlog.Info("payload received is: (%+v)", msg.Payload)

	disks := make([]*model.Disk, 0)

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

	c.foldersNotMoved = make([]string, 0)

	mlog.Info("_calc:Begin:srcDisk(%s); dstDisks(%d)", srcDisk.Path, len(disks))

	for _, disk := range disks {
		mlog.Info("_calc:elegibleDestDisk(%s)", disk.Path)
	}

	// Initialize fields
	c.storage.BytesToMove = 0
	c.storage.SourceDiskName = srcDisk.Path

	sort.Sort(model.ByFree(disks))

	srcDiskWithoutMnt := srcDisk.Path[5:]

	var folders []*model.Item
	for _, path := range c.settings.Folders {
		msg := fmt.Sprintf("Scanning folder %s on %s", path, srcDiskWithoutMnt)
		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		list := c.getFolders(dtoCalc.SourceDisk, path)

		if list != nil {
			folders = append(folders, list...)
		}
	}

	mlog.Info("_calc:foldersToBeMovedTotal(%d)", len(folders))
	var lsal string
	for _, v := range folders {
		lib.Shell(fmt.Sprintf("stat -c \"%%A %%y %%U %%G\" \"%s\"", v.Name), mlog.Warning, "_calc:lsal", "", func(line string) {
			lsal = line
		})

		mlog.Info("_calc:toBeMoved:Path(%s); Size(%s); linux(%s)", v.Path, lib.ByteSize(v.Size), lsal)
	}

	// srcDisk.NewFree = srcDisk.Free
	willBeMoved := make([]*model.Item, 0)

	if len(folders) > 0 {
		for _, disk := range disks {
			diskWithoutMnt := disk.Path[5:]
			msg := fmt.Sprintf("Trying to allocate folders to %s ...", diskWithoutMnt)
			outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
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
					reserved = lib.RESERVED_SPACE
				}

				ceil := lib.Max(lib.RESERVED_SPACE, reserved)
				mlog.Info("_calc:FoldersLeft(%d):ReservedSpace(%d)", len(folders), ceil)

				packer := algorithm.NewKnapsack(disk, folders, ceil)
				bin := packer.BestFit()
				if bin != nil {
					srcDisk.NewFree += bin.Size
					disk.NewFree -= bin.Size
					c.storage.BytesToMove += bin.Size

					willBeMoved = append(willBeMoved, bin.Items...)
					folders = c.removeFolders(folders, bin.Items)

					mlog.Info("_calc:BinAllocated=[Disk(%s); Items(%d)];Freespace=[original(%s); final(%s)]", disk.Path, len(bin.Items), lib.ByteSize(srcDisk.Free), lib.ByteSize(srcDisk.NewFree))
				} else {
					mlog.Info("_calc:NoBinAllocated=Disk(%s)", disk.Path)
				}
			}
		}
	}

	finished := time.Now()
	elapsed := lib.Round(time.Since(started), time.Millisecond)

	fstarted := started.Format(TIME_FORMAT)
	ffinished := finished.Format(TIME_FORMAT)

	// Send to frontend console started/ended/elapsed times
	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// send to frontend the folders that will not be moved, if any
	// notMoved holds a string representation of all the folders, separated by a '\n'

	if len(willBeMoved) == 0 {
		mlog.Info("_calc:No folders can be moved.")
	} else {
		mlog.Info("_calc:%d folders will be moved.", len(willBeMoved))
		for _, folder := range willBeMoved {
			mlog.Info("_calc:willBeMoved(%s)", folder.Path)
		}
	}

	notMoved := ""
	if len(folders) > 0 {
		// c.foldersNotMoved = append(make([]*model.Item, 0), folders...)

		outbound := &dto.Packet{Topic: "calcProgress", Payload: "The following folders will not be moved, because there's not enough space in the target disks:\n"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		mlog.Info("_calc:%d folders will NOT be moved.", len(folders))
		c.foldersNotMoved = make([]string, 0)
		for _, folder := range folders {
			c.foldersNotMoved = append(c.foldersNotMoved, folder.Path)

			notMoved += folder.Path + "\n"

			outbound = &dto.Packet{Topic: "calcProgress", Payload: folder.Path}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			mlog.Info("_calc:notMoved(%s)", folder.Path)
		}

		// for _, folder := range c.foldersNotMoved {
		// 	notMoved += folder + "\n"
		// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: folder}
		// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		// }
	}

	// send mail according to user preferences
	subject := "unBALANCE - CALCULATE operation completed"
	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", fstarted, ffinished, elapsed)
	if notMoved != "" {
		switch c.settings.NotifyCalc {
		case 1:
			message += "\n\nSome folders will not be moved because there's not enough space for them in any of the destination disks."
		case 2:
			message += "\n\nThe following folders will not be moved because there's not enough space for them in any of the destination disks:\n\n" + notMoved
		}
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
	mlog.Info("Bytes To Move: %s", lib.ByteSize(c.storage.BytesToMove))
	mlog.Info("---------------------------------------------------------")

	c.storage.Print()
	// msg.Reply <- c.storage

	mlog.Info("_calc:End:srcDisk(%s)", srcDisk.Path)

	outbound = &dto.Packet{Topic: "calcProgress", Payload: "Operation Finished"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// send to front end the signal of operation finished
	outbound = &dto.Packet{Topic: "calcFinished", Payload: c.storage}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
}

func (c *Core) getFolders(src string, folder string) (items []*model.Item) {
	srcFolder := filepath.Join(src, folder)

	mlog.Info("getFolders:Scanning source-disk(%s):folder(%s)", src, folder)

	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
		mlog.Warning("getFolders:Folder does not exist: %s", srcFolder)
		return nil
	}

	dirs, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		mlog.Fatalf("getFolders:Unable to readdir: %s", err)
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

// func (c *Core) getFolders(src string, folder string) (items []*model.Item) {
// 	srcFolder := filepath.Join(src, folder)

// 	mlog.Info("getFolders:Scanning source-disk(%s):folder(%s)", src, folder)
// 	// _, err := os.Stat(filepath.Join("/mnt/disk13/films", "*"))
// 	// mlog.Info("Error: %s", err)

// 	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
// 		mlog.Warning("getFolders:Folder does not exist: %s", srcFolder)
// 		return nil
// 	}

// 	dirs, err := ioutil.ReadDir(srcFolder)
// 	if err != nil {
// 		mlog.Fatalf("getFolders:Unable to readdir: %s", err)
// 	}

// 	mlog.Info("getFolders:Readdir(%d)", len(dirs))

// 	// mlog.Info("Dirs: %+v", dirs)

// 	if len(dirs) == 0 {
// 		mlog.Info("getFolders:No subdirectories under %s", srcFolder)
// 		return nil
// 	}

// 	// scanFolder := filepath.Join(fmt.Sprintf("\"%s\"", srcFolder), "*")
// 	// cmd := exec.Command("sh", "-c", fmt.Sprintf("du -bs %s", scanFolder))

// 	scanFolder := srcFolder + "/."
// 	cmdText := fmt.Sprintf("find \"%s\" ! -name . -prune -exec du -bs {} +", scanFolder)

// 	mlog.Info("getFolders:Executing %s", cmdText)

// 	cmd := exec.Command("sh", "-c", cmdText)
// 	out, err := cmd.StdoutPipe()
// 	if err != nil {
// 		mlog.Fatalf("getFolders:Unable to stdoutpipe cmd(%s): %s", cmdText, err)
// 	}

// 	stderr, err := cmd.StderrPipe()
// 	if err != nil {
// 		mlog.Fatalf("getFolders:Unable to stdoutpipe cmd(%s): %s", cmdText, err)
// 	}

// 	rd := bufio.NewReader(out)

// 	if err := cmd.Start(); err != nil {
// 		mlog.Fatalf("getFolders:Unable to start du: %s", err)
// 	}

// 	go func() {
// 		errbuf := bufio.NewScanner(stderr)
// 		for errbuf.Scan() {
// 			mlog.Warning("getFolders:find/du:stderr: %s", errbuf.Text())
// 		}
// 	}()

// 	for {
// 		line, err := rd.ReadString('\n')
// 		if err == io.EOF && len(line) == 0 {
// 			// Good end of file with no partial line
// 			break
// 		}
// 		if err == io.EOF {
// 			mlog.Fatalf("getFolders:Last line not terminated: %s", err)
// 		}

// 		if err != nil {
// 			mlog.Fatalf("getFolders:Unable to ReadString: %s", err)
// 		}

// 		line = line[:len(line)-1] // drop the '\n'
// 		if line[len(line)-1] == '\r' {
// 			line = line[:len(line)-1] // drop the '\r'
// 		}

// 		mlog.Info("getFolders:find(%s): %s", scanFolder, line)

// 		result := c.reItems.FindStringSubmatch(line)
// 		// mlog.Info("[%s] %s", result[1], result[2])

// 		size, _ := strconv.ParseInt(result[1], 10, 64)

// 		item := &model.Item{Name: result[2], Size: size, Path: filepath.Join(folder, filepath.Base(result[2]))}
// 		items = append(items, item)

// 		msg := fmt.Sprintf("Found %s (%s)", filepath.Base(item.Name), lib.ByteSize(size))
// 		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		// fmt.Println(line)
// 		// mlog.Info("getFolders:item: %+v", item)
// 	}

// 	// Wait for the result of the command; also closes our end of the pipe
// 	err = cmd.Wait()
// 	if err != nil {
// 		mlog.Fatalf("getFolders:Unable to wait for process to finish: %s", err)
// 	}

// 	// out, err := lib.Shell(fmt.Sprintf("du -sh %s", filepath.Join(disk, folder, "*")))
// 	// if err != nil {
// 	// 	glog.Fatal(err)
// 	// }

// 	// glog.Info(string(out))
// 	// mlog.Info("done")
// 	return items
// }

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
	c.opState = MOVE
	go c._move(msg)
}

func (c *Core) _move(msg *pubsub.Message) {
	defer func() {
		c.opState = IDLE
		c.started = time.Time{}
		c.bytesMoved = 0
	}()

	mlog.Info("Running move operation ...")
	c.started = time.Now()

	outbound := &dto.Packet{Topic: "moveStarted", Payload: "Operation started"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// rsyncArgs := []string{
	// 	"-avX",
	// 	"--partial",
	// }

	rsyncArgs := c.settings.RsyncFlags

	if c.settings.DryRun {
		rsyncArgs = append(rsyncArgs, "--dry-run")
	}

	commands := make([]string, 0)

	outbound = &dto.Packet{Topic: "progressStats", Payload: "Waiting to collect stats ..."}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	for _, disk := range c.storage.Disks {
		if disk.Bin == nil || disk.Path == c.storage.SourceDiskName {
			continue
		}

		for _, item := range disk.Bin.Items {
			workDir := c.storage.SourceDiskName
			// src := strconv.Quote(filepath.Join(c.storage.SourceDiskName, item.Path))
			// src := strconv.Quote(c.storage.SourceDiskName+string(filepath.Separator)+"."+string(filepath.Separator)+ item.Path)
			// src := strconv.Quote(c.storage.SourceDiskName + string(filepath.Separator) + "." + string(filepath.Separator) + item.Path)
			src := item.Path
			// dst := strconv.Quote(filepath.Join(disk.Path, filepath.Dir(item.Path)))
			// dst := strconv.Quote(filepath.Join(disk.Path, filepath.Dir(item.Path)) + string(filepath.Separator))
			dst := disk.Path + string(filepath.Separator)

			// args := append(make([]string, 0), rsyncArgs...)
			// args = append(
			// 	args,
			// 	filepath.Join(c.storage.SourceDiskName, item.Path),
			// 	filepath.Join(disk.Path, filepath.Dir(item.Path)),
			// )

			args := append(
				rsyncArgs,
				src,
				dst,
			)
			// cmd := exec.Command("rsync", args)

			// cmd := fmt.Sprintf("rsync %s \"%s\" \"%s/\"", strings.Join(rsyncArgs, " "), filepath.Join(c.storage.SourceDiskName, item.Path), filepath.Join(disk.Path, filepath.Dir(item.Path)))
			cmd := fmt.Sprintf(`rsync %s %s %s`, strings.Join(rsyncArgs, " "), strconv.Quote(src), dst)
			mlog.Info("cmd(%s)", cmd)

			outbound = &dto.Packet{Topic: "moveProgress", Payload: cmd}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			// err := lib.ShellEx(mlog.Warning, "moveProgress:", func(line string) {
			// 	outbound := &dto.Packet{Topic: "moveProgress", Payload: line}
			// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			// 	mlog.Info(line)
			// }, "cd", args...)

			err := lib.ShellEx(mlog.Warning, "moveProgress:", workDir, func(line string) {
				outbound := &dto.Packet{Topic: "moveProgress", Payload: line}
				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

				mlog.Info(line)
			}, "rsync", args...)

			if err != nil {
				finished := time.Now()
				elapsed := time.Since(c.started)

				subject := "unBALANCE - MOVE operation INTERRUPTED"
				headline := fmt.Sprintf("Move command (%s) was interrupted: %s", cmd, err.Error())

				mlog.Warning(headline)
				outbound := &dto.Packet{Topic: "opError", Payload: "Move operation was interrupted. Check logs for details."}
				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

				c.finishMoveOperation(subject, headline, commands, c.started, finished, elapsed)

				return
			}

			c.bytesMoved = c.bytesMoved + item.Size

			percent, left, speed := progress(c.storage.BytesToMove, c.bytesMoved, c.started)

			msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)
			mlog.Info("Current progress: %s", msg)

			outbound := &dto.Packet{Topic: "progressStats", Payload: msg}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			commands = append(commands, cmd)

			if !c.settings.DryRun {
				rmrf := fmt.Sprintf("rm -rf \"%s\"", filepath.Join(c.storage.SourceDiskName, item.Path))
				mlog.Info("Removing: (%s)", rmrf)
				err = lib.Shell(rmrf, mlog.Warning, "moveProgress:", "", func(line string) {
					mlog.Info(line)
				})

				if err != nil {
					msg := fmt.Sprintf("Unable to remove source folder:(%s)", filepath.Join(c.storage.SourceDiskName, item.Path))

					outbound := &dto.Packet{Topic: "moveProgress", Payload: msg}
					c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

					mlog.Warning(msg)
				}
			}
		}
	}

	finished := time.Now()
	elapsed := time.Since(c.started)

	subject := "unBALANCE - MOVE operation completed"
	headline := "Move operation has finished"

	c.finishMoveOperation(subject, headline, commands, c.started, finished, elapsed)
}

// //			dst := filepath.Join(disk.Path, item.Path)

// // mlog.Info("disk.Path = %s | item.Path = %s | dst = %s", disk.Path, item.Path, c.storage.SourceDiskName)
// // mlog.Info("disk.Path = %s | item.Name = %s | item.Path = %s | dst = %s", disk.Path, item.Name, item.Path, dst)
// // mlog.Info("mv %s %s", strconv.Quote(item.Name), strconv.Quote(dst))
// //			command := &dto.Move{Command: fmt.Sprintf("mv %s %s", strconv.Quote(item.Name), strconv.Quote(dst))}
// //			commands = append(commands, command)

// // sanePath := item.Path
// // if item.Path[0] == '/' {
// // 	sanePath = sanePath[1:]
// // }

// // execute shell command inline
// // command := fmt.Sprintf("%s/diskmv %s \"%s\" %s %s", c.diskmvLocation, dry, sanePath, c.storage.SourceDiskName, disk.Path)
// // mlog.Info("cmd(%s)", command)

// // cmd := exec.Command("/bin/sh", "-c", command)
// // cmd := exec.Command(command)
// cmd := exec.Command("rsync", rsyncArgs, filepath.Join(c.storage.SourceDiskName, item.Path), filepath.Join(disk.Path, item.Path))

// command := fmt.Sprintf("%s %s \"%s\" \"%s\"", cmd.Path, rsyncArgs, filepath.Join(c.storage.SourceDiskName, item.Path), filepath.Join(disk.Path, item.Path))
// mlog.Info("cmd(%s)", command)

// outbound = &dto.Packet{Topic: "moveProgress", Payload: command}
// c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// stdout, err := cmd.StdoutPipe()
// if err != nil {
// 	finished := time.Now()
// 	elapsed := time.Since(c.started)

// 	subject := "unBALANCE - MOVE operation INTERRUPTED"
// 	headline := fmt.Sprintf("Move command (%s) was interrupted: %s", command, err)

// 	mlog.Warning(headline)
// 	outbound := &dto.Packet{Topic: "opError", Payload: "Move operation was interrupted. Check logs for details."}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	c.finishMoveOperation(subject, headline, commands, c.started, finished, elapsed)
// 	return
// 	//		log.Fatalf("Unable to stdoutpipe %s: %s", command, err)
// }

// cmd.Stderr = lib.NewStreamer(mlog.Warning, "moveProgress")

// // stderr, err := cmd.StderrPipe()
// // if err != nil {
// // 	return err
// // 	// log.Fatalf("Unable to stderrpipe %s: %s", command, err)
// // }

// // multi := io.MultiReader(stdout, stderr)

// rd := bufio.NewReader(stdout)

// if err := cmd.Start(); err != nil {
// 	mlog.Warning("Path(%s):Args(%s):Dir(%s):String(%s)", cmd.Path, cmd.Args, cmd.Dir, cmd.ProcessState.String())
// 	finished := time.Now()
// 	elapsed := time.Since(c.started)

// 	subject := "unBALANCE - MOVE operation INTERRUPTED"
// 	headline := fmt.Sprintf("Move command (%s) was interrupted: %s", command, err)

// 	mlog.Warning(headline)
// 	outbound := &dto.Packet{Topic: "opError", Payload: "Move operation was interrupted. Check logs for details."}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	c.finishMoveOperation(subject, headline, commands, c.started, finished, elapsed)
// 	return
// }

// for {
// 	line, err := rd.ReadString('\n')

// 	if err == io.EOF && len(line) == 0 {
// 		// Good end of file with no partial line
// 		mlog.Warning("_move:ExitOk")
// 		break
// 	}
// 	if err == io.EOF {
// 		mlog.Warning("_move:lineNotTerminated(%s):(%s)", err, line)
// 		break
// 	}

// 	// mlog.Info("thisline:(%s)", line)
// 	line = line[:len(line)-1] // drop the '\n'
// 	if len(line) > 0 && line[len(line)-1] == '\r' {
// 		line = line[:len(line)-1] // drop the '\r'
// 	}

// 	outbound := &dto.Packet{Topic: "moveProgress", Payload: line}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	mlog.Info(line)
// }

// // Wait for the result of the command; also closes our end of the pipe
// err = cmd.Wait()
// if err != nil {
// 	var waitStatus syscall.WaitStatus
// 	if exiterr, ok := err.(*exec.ExitError); ok {
// 		waitStatus = exiterr.Sys().(syscall.WaitStatus)
// 		mlog.Warning("_move:waitError:Status(%d):Err(%s):ExitErr(%s)", waitStatus.ExitStatus(), err, exiterr)
// 	} else {
// 		mlog.Warning("_move:waitError:(%s)", err)
// 	}

// 	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)
// 	mlog.Warning("_move:waitStatus:(%d)", waitStatus.ExitStatus())

// 	finished := time.Now()
// 	elapsed := time.Since(c.started)

// 	subject := "unBALANCE - MOVE operation INTERRUPTED"
// 	headline := fmt.Sprintf("Move command (%s) was interrupted: %s", command, err)

// 	// mlog.Warning(headline)
// 	outbound := &dto.Packet{Topic: "opError", Payload: "Move operation was interrupted. Check logs for details."}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	c.finishMoveOperation(subject, headline, commands, c.started, finished, elapsed)
// 	// log.Fatal("Unable to wait for process to finish: ", err)
// 	return
// }

// c.bytesMoved = c.bytesMoved + item.Size

// percent, left, speed := progress(c.storage.BytesToMove, c.bytesMoved, c.started)

// msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)
// mlog.Info("Current progress: %s", msg)

// outbound := &dto.Packet{Topic: "progressStats", Payload: msg}
// c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// commands = append(commands, command)
// 	}
// }

// finished := time.Now()
// elapsed := time.Since(c.started)

// subject := "unBALANCE - MOVE operation completed"
// headline := "Move operation has finished"

// c.finishMoveOperation(subject, headline, commands, c.started, finished, elapsed)

// }

func (c *Core) finishMoveOperation(subject, headline string, commands []string, started, finished time.Time, elapsed time.Duration) {
	fstarted := started.Format(TIME_FORMAT)
	ffinished := finished.Format(TIME_FORMAT)
	elapsed = lib.Round(time.Since(started), time.Millisecond)

	outbound := &dto.Packet{Topic: "moveProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "moveProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "moveProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "moveProgress", Payload: headline}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "moveProgress", Payload: "These are the commands that were executed:"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	printedCommands := ""
	for _, command := range commands {
		printedCommands += command + "\n"
		outbound = &dto.Packet{Topic: "moveProgress", Payload: command}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	outbound = &dto.Packet{Topic: "moveProgress", Payload: "Operation Finished"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// send to front end the signal of operation finished
	if !c.settings.DryRun {
		c.storage.Refresh()
	} else {
		outbound = &dto.Packet{Topic: "moveProgress", Payload: "--- IT WAS A DRY RUN ---"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	outbound = &dto.Packet{Topic: "moveFinished", Payload: c.storage}
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

// func (c *Core) doStorageUpdate(msg *pubsub.Message) {
// 	outbound := &dto.Packet{Topic: "storage:update:completed", Payload: c.storage.Refresh()}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// }

func (c *Core) sendmail(notify int, subject, message string, dryRun bool) (err error) {
	if notify == 0 {
		return nil
	}

	dry := ""
	if dryRun {
		dry = "-------\nDRY RUN\n-------\n"
	}

	msg := dry + message

	// strCmd := fmt.Sprintf("-s \"%s\" -m \"%s\"", MAIL_CMD, subject, msg)
	cmd := exec.Command(MAIL_CMD, "-e", "unBALANCE operation update", "-s", subject, "-m", msg)
	err = cmd.Run()

	// from := "From: " + c.settings.NotiFrom
	// to := "To: " + c.settings.NotiTo
	// subject := "Subject: unBALANCE Notification"

	// dry := ""
	// if !dryRun {
	// 	dry = "-------\nDRY RUN\n-------\n"
	// }

	// mail := "\"" + from + "\n" + to + "\n" + subject + "\n\n" + dry + msg + "\""
	// mail := from + "\n" + to + "\n" + subject + "\n\n" + dry + msg

	// err := ioutil.WriteFile(msgLocation, []byte(mail), 0644)
	// if err != nil {
	// 	return err
	// }

	// echo := exec.Command("echo", "-e", mail)
	// ssmtp := exec.Command("ssmtp", c.settings.NotiTo)

	// // mlog.Info("sendmail:echo: %s %s (%s)", echo.Path, echo.Args, echo.Dir)
	// // mlog.Info("sendmail:ssmtp: %s %s (%s)", ssmtp.Path, ssmtp.Args, ssmtp.Dir)

	// _, _, err := lib.Pipeline(echo, ssmtp)
	// if err != nil {
	// 	return err
	// }

	return
}

func progress(bytesToMove, bytesMoved int64, started time.Time) (percent float64, left time.Duration, speed float64) {
	delta := time.Since(started)

	bytesPerSec := float64(bytesMoved) / delta.Seconds()
	speed = bytesPerSec / 1024 / 1024 // MB/s

	percent = (float64(bytesMoved) / float64(bytesToMove)) * 100 // %

	left = time.Duration(float64(bytesToMove-bytesMoved)/bytesPerSec) * time.Second

	return
}

// func (c *Core) printCommands(list []string) string {
// 	var str string
// 	for _, value := range list {
// 		str += value
// 	}
// 	return str
// }
