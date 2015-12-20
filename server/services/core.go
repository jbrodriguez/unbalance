package services

import (
	"bufio"
	"fmt"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"io"
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
	"time"
)

// const dockerEnv = "UNBALANCE_DOCKER"
const (
	DISKMV_CMD        = "./diskmv"
	DISKMV_DOCKER_CMD = "/usr/bin/diskmv"
)

type Core struct {
	Service

	bus      *pubsub.PubSub
	storage  *model.Unraid
	settings *lib.Settings

	foldersNotMoved []string

	mailbox chan *pubsub.Mailbox

	// chanConfigInfo       chan *pubsub.Message
	// chanSaveConfig       chan *pubsub.Message
	// chanStorageInfo      chan *pubsub.Message
	// chanCalculateBestFit chan *pubsub.Message
	// storageMove          chan *pubsub.Message
	// storageUpdate        chan *pubsub.Message

	reFreeSpace *regexp.Regexp
	reItems     *regexp.Regexp
}

func NewCore(bus *pubsub.PubSub, settings *lib.Settings) *Core {
	core := &Core{
		bus:      bus,
		settings: settings,
	}
	core.init()

	re, _ := regexp.Compile(`(.*?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.*?)\s+(.*?)$`)
	core.reFreeSpace = re

	re, _ = regexp.Compile(`(\d+)\s+(.*?)$`)
	core.reItems = re

	core.storage = &model.Unraid{}

	return core
}

func (c *Core) Start() {
	mlog.Info("starting service Core ...")

	// core.chanConfigInfo = core.bus.Sub("cmd.getConfig")
	// core.chanSaveConfig = core.bus.Sub("cmd.saveConfig")
	// core.chanStorageInfo = core.bus.Sub("cmd.getStorageInfo")
	// core.chanCalculateBestFit = core.bus.Sub("cmd.calculateBestFit")
	// core.storageMove = core.bus.Sub("storage:move")
	// core.storageUpdate = core.bus.Sub("storage:update")

	c.mailbox = c.register(c.bus, "/get/config", c.getConfig)
	c.registerAdditional(c.bus, "/set/config", c.setConfig, c.mailbox)
	c.registerAdditional(c.bus, "/get/storage", c.getStorage, c.mailbox)
	c.registerAdditional(c.bus, "/calculate", c.calc, c.mailbox)
	c.registerAdditional(c.bus, "/move", c.move, c.mailbox)
	// c.registerAdditional(c.bus, "/set/config", c.setConfig)

	go c.react()
}

func (c *Core) Stop() {
	mlog.Info("stopped service Core ...")
}

func (c *Core) react() {
	for mbox := range c.mailbox {
		// mlog.Info("Core:Topic: %s", mbox.Topic)
		c.dispatch(mbox.Topic, mbox.Content)
	}
}

// func (c *Core) SetStorage(storage *model.Unraid) {
// 	c.storage = storage
// }

// func (c *Core) react() {
// 	for {
// 		select {
// 		case msg := <-c.chanConfigInfo:
// 			go c.getConfigInfo(msg)
// 		case msg := <-c.chanSaveConfig:
// 			go c.saveConfig(msg)
// 		case msg := <-c.chanStorageInfo:
// 			go c.getStorageInfo(msg)
// 		case msg := <-c.chanCalculateBestFit:
// 			go c.calculateBestFit(msg)
// 		case msg := <-c.storageMove:
// 			go c.doStorageMove(msg)
// 		case msg := <-c.storageUpdate:
// 			go c.doStorageUpdate(msg)
// 		}
// 	}
// }

func (c *Core) getConfig(msg *pubsub.Message) {
	mlog.Info("Sending config")

	msg.Reply <- &c.settings.Config
}

func (c *Core) setConfig(msg *pubsub.Message) {
	mlog.Info("Saving config")

	config := msg.Payload.(*lib.Config)
	c.settings.Config = *config
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

func (c *Core) getStorage(msg *pubsub.Message) {
	if c.storage.Condition.State == "STARTED" {
		c.storage.Refresh(c.settings.RunningInDocker)
	}

	msg.Reply <- c.storage
}

func (c *Core) calc(msg *pubsub.Message) {
	mlog.Info("Running calculate operation ...")
	started := time.Now()

	c.storage.Condition.State = "CALCULATING"

	outbound := &dto.Packet{Topic: "storage:calc:begin", Payload: "Operation started"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	dtoCalc := msg.Payload.(*dto.Calculate)

	disks := make([]*model.Disk, 0)

	var srcDisk *model.Disk
	for _, disk := range c.storage.Disks {
		disk.NewFree = 0
		disk.Bin = nil

		if disk.Path == dtoCalc.SourceDisk {
			srcDisk = disk
		} else {
			if val, ok := dtoCalc.DestDisks[disk.Path]; ok && val {
				disks = append(disks, disk)
			} else {
				// if the disk is not elegible as a target, let newFree = Free, to prevent the UI to think there was some change in it
				disk.NewFree = disk.Free
			}
		}

	}

	c.foldersNotMoved = make([]string, 0)

	mlog.Info("calculateBestFit:Begin:srcDisk(%s); dstDisks(%d)", srcDisk.Path, len(disks))

	for _, disk := range disks {
		mlog.Info("calculateBestFit:elegibleDestDisk(%s)", disk.Path)
	}

	// Initialize fields
	c.storage.BytesToMove = 0
	c.storage.SourceDiskName = srcDisk.Path

	sort.Sort(model.ByFree(disks))

	var folders []*model.Item
	for _, path := range c.settings.Folders {
		list := c.getFolders(dtoCalc.SourceDisk, path)

		if list != nil {
			folders = append(folders, list...)
		}
	}

	for _, v := range folders {
		mlog.Info("calculateBestFit:total(%d):toBeMoved:Path(%s); Size(%s)", len(folders), v.Path, lib.ByteSize(v.Size))
	}

	srcDisk.NewFree = srcDisk.Free

	for _, disk := range disks {
		msg := fmt.Sprintf("Processing disk %s ...", disk.Path)
		outbound := &dto.Packet{Topic: "storage:calc:progress", Payload: msg}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		mlog.Info("calculateBestFit:%s", msg)

		if disk.Path != srcDisk.Path {
			disk.NewFree = disk.Free

			mlog.Info("calculateBestFit:FoldersLeft(%d)", len(folders))

			packer := algorithm.NewKnapsack(disk, folders, c.settings.ReservedSpace)
			bin := packer.BestFit()
			if bin != nil {
				srcDisk.NewFree += bin.Size
				disk.NewFree -= bin.Size
				c.storage.BytesToMove += bin.Size

				folders = c.removeFolders(folders, bin.Items)

				mlog.Info("calculateBestFit:BinAllocated=[Disk(%s); Items(%d)];Freespace=[original(%s); final(%s)]", disk.Path, len(bin.Items), lib.ByteSize(srcDisk.Free), lib.ByteSize(srcDisk.NewFree))
			} else {
				mlog.Info("calculateBestFit:NoBinAllocated=Disk(%s)", disk.Path)
			}
		}
	}

	finished := time.Now()
	elapsed := time.Since(started)

	// Send to frontend console started/ended/elapsed times
	outbound = &dto.Packet{Topic: "storage:calc:progress", Payload: fmt.Sprintf("Started: %s", started)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "storage:calc:progress", Payload: fmt.Sprintf("Ended: %s", finished)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "storage:calc:progress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// send to frontend the folders that will not be moved, if any
	// notMoved holds a string representation of all the folders, separated by a '\n'
	notMoved := ""
	if len(folders) > 0 {
		// c.foldersNotMoved = append(make([]*model.Item, 0), folders...)

		outbound := &dto.Packet{Topic: "storage:calc:progress", Payload: "The following folders will not be moved, because there's not enough space in the target disks:\n"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		c.foldersNotMoved = make([]string, 0)
		for _, folder := range folders {
			c.foldersNotMoved = append(c.foldersNotMoved, folder.Path)

			notMoved += folder.Path + "\n"

			outbound = &dto.Packet{Topic: "storage:calc:progress", Payload: folder.Path}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		}

		// for _, folder := range c.foldersNotMoved {
		// 	notMoved += folder + "\n"
		// 	outbound = &dto.Packet{Topic: "storage:calc:progress", Payload: folder}
		// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		// }
	}

	// send to front end the signal of operation finished
	outbound = &dto.Packet{Topic: "storage:calc:end", Payload: "Operation Finished"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// send mail according to user preferences
	subject := "unBALANCE - CALCULATE operation completed"
	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", started, finished, elapsed)
	if notMoved != "" {
		switch c.settings.NotifyCalc {
		case 1:
			message += "\n\nSome folders are not elegible for moving because there's not enough space for them in any of the target disks."
		case 2:
			message += "\n\nThe following folders are not elegible for moving because there's not enough space for them in any of the target disks:\n\n" + notMoved
		}
	}

	if sendErr := c.sendmail(c.settings.NotifyCalc, subject, message, false); sendErr != nil {
		mlog.Error(sendErr)
	}

	// some local logging
	mlog.Info("calculateBestFit:FoldersLeft(%d)", len(folders))
	mlog.Info("calculateBestFit:src(%s):Listing (%d) disks ...", srcDisk.Path, len(c.storage.Disks))
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

	mlog.Info("calculateBestFit:End:srcDisk(%s)", srcDisk.Path)
}

func (c *Core) getFolders(src string, folder string) (items []*model.Item) {
	srcFolder := filepath.Join(src, folder)

	mlog.Info("getFolders:Scanning source-disk(%s):folder(%s)", src, folder)
	// _, err := os.Stat(filepath.Join("/mnt/disk13/films", "*"))
	// mlog.Info("Error: %s", err)

	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
		mlog.Warning("getFolders:Folder does not exist: %s", srcFolder)
		return nil
	}

	dirs, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		mlog.Fatalf("getFolders:Unable to readdir: %s", err)
	}

	mlog.Info("getFolders:Readdir(%d)", len(dirs))

	// mlog.Info("Dirs: %+v", dirs)

	if len(dirs) == 0 {
		mlog.Info("getFolders:No subdirectories under %s", srcFolder)
		return nil
	}

	// scanFolder := filepath.Join(fmt.Sprintf("\"%s\"", srcFolder), "*")
	// cmd := exec.Command("sh", "-c", fmt.Sprintf("du -bs %s", scanFolder))

	scanFolder := srcFolder + "/."
	cmdText := fmt.Sprintf("find \"%s\" ! -name . -prune -exec du -bs {} +", scanFolder)

	mlog.Info("getFolders:Executing %s", cmdText)

	cmd := exec.Command("sh", "-c", cmdText)
	out, err := cmd.StdoutPipe()
	if err != nil {
		mlog.Fatalf("getFolders:Unable to stdoutpipe cmd(%s): %s", cmdText, err)
	}

	rd := bufio.NewReader(out)

	if err := cmd.Start(); err != nil {
		mlog.Fatalf("getFolders:Unable to start du: %s", err)
	}

	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF && len(line) == 0 {
			// Good end of file with no partial line
			break
		}
		if err == io.EOF {
			mlog.Fatalf("getFolders:Last line not terminated: %s", err)
		}
		line = line[:len(line)-1] // drop the '\n'
		if line[len(line)-1] == '\r' {
			line = line[:len(line)-1] // drop the '\r'
		}

		mlog.Info("getFolders:find(%s): %s", scanFolder, line)

		result := c.reItems.FindStringSubmatch(line)
		// mlog.Info("[%s] %s", result[1], result[2])

		size, _ := strconv.ParseInt(result[1], 10, 64)

		item := &model.Item{Name: result[2], Size: size, Path: filepath.Join(folder, filepath.Base(result[2]))}
		items = append(items, item)
		// fmt.Println(line)
		// mlog.Info("getFolders:item: %+v", item)
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		mlog.Fatalf("getFolders:Unable to wait for process to finish: %s", err)
	}

	// out, err := lib.Shell(fmt.Sprintf("du -sh %s", filepath.Join(disk, folder, "*")))
	// if err != nil {
	// 	glog.Fatal(err)
	// }

	// glog.Info(string(out))
	// mlog.Info("done")
	return items
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

// func (c *Core) processDiskMv(line string, arg interface{}) {
// 	outbound := &dto.Packet{Topic: "storage:move:progress", Payload: line}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	mlog.Info(line)
// }

func (c *Core) move(msg *pubsub.Message) {
	mlog.Info("Running move operation ...")
	started := time.Now()

	c.storage.Condition.State = "MOVING"

	var dry string
	if c.settings.DryRun {
		dry = "-t"
	} else {
		dry = "-f"
	}

	var diskmv string
	if c.settings.RunningInDocker {
		diskmv = DISKMV_DOCKER_CMD
	} else {
		diskmv = DISKMV_CMD
	}

	// c.storage.InProgress = true

	outbound := &dto.Packet{Topic: "storage:move:begin", Payload: "Operation started"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// if err := c.sendmail("Move Operation started"); err != nil {
	// 	mlog.Error(err)
	// }

	commands := make([]string, 0)

	// notMoved := ""
	// if len(c.foldersNotMoved) > 0 {
	// 	notMoved = "The following folders will not be moved, because there's not enough space in the target disks:\n"
	// 	for _, folder := range c.foldersNotMoved {
	// 		notMoved += folder + "\n"
	// 	}
	// }

	for _, disk := range c.storage.Disks {
		if disk.Bin == nil || disk.Path == c.storage.SourceDiskName {
			continue
		}

		for _, item := range disk.Bin.Items {
			//			dst := filepath.Join(disk.Path, item.Path)

			// mlog.Info("disk.Path = %s | item.Path = %s | dst = %s", disk.Path, item.Path, c.storage.SourceDiskName)
			// mlog.Info("disk.Path = %s | item.Name = %s | item.Path = %s | dst = %s", disk.Path, item.Name, item.Path, dst)
			// mlog.Info("mv %s %s", strconv.Quote(item.Name), strconv.Quote(dst))
			//			command := &dto.Move{Command: fmt.Sprintf("mv %s %s", strconv.Quote(item.Name), strconv.Quote(dst))}
			//			commands = append(commands, command)

			cmd := fmt.Sprintf("%s %s \"%s\" %s %s", diskmv, dry, item.Path, c.storage.SourceDiskName, disk.Path)
			mlog.Info("cmd(%s)", cmd)

			commands = append(commands, cmd)

			outbound = &dto.Packet{Topic: "storage:move:progress", Payload: cmd}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			err := lib.Shell(cmd, func(line string) {
				outbound := &dto.Packet{Topic: "storage:move:progress", Payload: line}
				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

				mlog.Info(line)
			})

			if err != nil {
				finished := time.Now()
				elapsed := time.Since(started)

				subject := "unBALANCE - MOVE operation INTERRUPTED"
				headline := fmt.Sprintf("Move command (%s) was interrupted: %s", cmd, err.Error())

				mlog.Warning(headline)

				c.finishMoveOperation(subject, headline, commands, started, finished, elapsed)

				return
			}

		}
	}

	finished := time.Now()
	elapsed := time.Since(started)

	subject := "unBALANCE - MOVE operation completed"
	headline := "Move operation has finished"

	c.finishMoveOperation(subject, headline, commands, started, finished, elapsed)

}

func (c *Core) finishMoveOperation(subject, headline string, commands []string, started, finished time.Time, elapsed time.Duration) {
	outbound := &dto.Packet{Topic: "storage:move:progress", Payload: fmt.Sprintf("Started: %s", started)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "storage:move:progress", Payload: fmt.Sprintf("Ended: %s", finished)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "storage:move:progress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "storage:move:progress", Payload: headline}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	outbound = &dto.Packet{Topic: "storage:move:progress", Payload: "These are the commands that were executed:"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	printedCommands := ""
	for _, command := range commands {
		printedCommands += command + "\n"
		outbound = &dto.Packet{Topic: "storage:move:progress", Payload: command}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	outbound = &dto.Packet{Topic: "storage:move:end", Payload: "Operation Finished"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s\n\n%s", started, finished, elapsed, headline)
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

func (c *Core) sendmail(notify int, subject, message string, dryRun bool) error {
	if notify == 0 {
		return nil
	}

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

	return nil

}

// func (c *Core) printCommands(list []string) string {
// 	var str string
// 	for _, value := range list {
// 		str += value
// 	}
// 	return str
// }
