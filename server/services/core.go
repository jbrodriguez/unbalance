package services

import (
	"apertoire.net/unbalance/server/dto"
	"apertoire.net/unbalance/server/helper"
	"apertoire.net/unbalance/server/lib"
	"apertoire.net/unbalance/server/model"
	"bufio"
	"fmt"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"time"
)

const dockerEnv = "UNBALANCE_DOCKER"
const diskmvCmd = "./diskmv"
const diskmvDockerCmd = "/usr/bin/diskmv"
const msgLocation string = "/home/msg.txt"

type Core struct {
	bus      *pubsub.PubSub
	storage  *model.Unraid
	settings *model.Settings

	chanConfigInfo       chan *pubsub.Message
	chanSaveConfig       chan *pubsub.Message
	chanStorageInfo      chan *pubsub.Message
	chanCalculateBestFit chan *pubsub.Message
	storageMove          chan *pubsub.Message
	storageUpdate        chan *pubsub.Message

	reFreeSpace *regexp.Regexp
	reItems     *regexp.Regexp
}

func NewCore(bus *pubsub.PubSub, settings *model.Settings) *Core {
	core := &Core{bus: bus, settings: settings}

	re, _ := regexp.Compile(`(.*?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.*?)\s+(.*?)$`)
	core.reFreeSpace = re

	re, _ = regexp.Compile(`(\d+)\s+(.*?)$`)
	core.reItems = re

	core.storage = &model.Unraid{}

	core.chanConfigInfo = core.bus.Sub("cmd.getConfig")
	core.chanSaveConfig = core.bus.Sub("cmd.saveConfig")
	core.chanStorageInfo = core.bus.Sub("cmd.getStorageInfo")
	core.chanCalculateBestFit = core.bus.Sub("cmd.calculateBestFit")
	core.storageMove = core.bus.Sub("storage:move")
	core.storageUpdate = core.bus.Sub("storage:update")

	return core
}

func (c *Core) Start() {
	mlog.Info("starting service Core ...")
	go c.react()
}

func (c *Core) Stop() {
	mlog.Info("stopped service Core ...")
}

func (c *Core) react() {
	for {
		select {
		case msg := <-c.chanConfigInfo:
			go c.getConfigInfo(msg)
		case msg := <-c.chanSaveConfig:
			go c.saveConfig(msg)
		case msg := <-c.chanStorageInfo:
			go c.getStorageInfo(msg)
		case msg := <-c.chanCalculateBestFit:
			go c.calculateBestFit(msg)
		case msg := <-c.storageMove:
			go c.doStorageMove(msg)
		case msg := <-c.storageUpdate:
			go c.doStorageUpdate(msg)
		}
	}
}

func (c *Core) getConfigInfo(msg *pubsub.Message) {
	mlog.Info("Sending config")

	msg.Reply <- &c.settings.Config
}

func (c *Core) saveConfig(msg *pubsub.Message) {
	mlog.Info("Saving config")

	config := msg.Payload.(*model.Config)
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

func (c *Core) getStorageInfo(msg *pubsub.Message) {
	msg.Reply <- c.storage.Refresh()
}

func (c *Core) calculateBestFit(msg *pubsub.Message) {
	dto := msg.Payload.(*dto.BestFit)

	disks := make([]*model.Disk, 0)

	var srcDisk *model.Disk
	for _, disk := range c.storage.Disks {
		disk.NewFree = 0
		disk.Bin = nil

		if disk.Path == dto.SourceDisk {
			srcDisk = disk
		} else {
			if val, ok := dto.DestDisks[disk.Path]; ok && val {
				disks = append(disks, disk)
			} else {
				// if the disk is not elegible as a target, let newFree = Free, to prevent the UI to think there was some change in it
				disk.NewFree = disk.Free
			}
		}

	}

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
		list := c.getFolders(dto.SourceDisk, path)

		if list != nil {
			folders = append(folders, list...)
		}
	}

	for _, v := range folders {
		mlog.Info("calculateBestFit:total(%d):toBeMoved:Path(%s); Size(%s)", len(folders), v.Path, helper.ByteSize(v.Size))
	}

	srcDisk.NewFree = srcDisk.Free

	for _, disk := range disks {
		//		disk.NewFree = disk.Free
		if disk.Path != srcDisk.Path {
			disk.NewFree = disk.Free

			mlog.Info("calculateBestFit:FoldersLeft(%d)", len(folders))

			packer := lib.NewKnapsack(disk, folders, c.settings.ReservedSpace)
			bin := packer.BestFit()
			if bin != nil {
				srcDisk.NewFree += bin.Size
				disk.NewFree -= bin.Size
				c.storage.BytesToMove += bin.Size

				folders = c.removeFolders(folders, bin.Items)

				mlog.Info("calculateBestFit:BinAllocated=[Disk(%s); Items(%d)];Freespace=[original(%s); final(%s)]", disk.Path, len(bin.Items), helper.ByteSize(srcDisk.Free), helper.ByteSize(srcDisk.NewFree))
			} else {
				mlog.Info("calculateBestFit:NoBinAllocated=Disk(%s)", disk.Path)
			}
		}
	}

	mlog.Info("calculateBestFit:FoldersLeft(%d)", len(folders))
	mlog.Info("calculateBestFit:src(%s):Listing (%d) disks ...", srcDisk.Path, len(c.storage.Disks))

	for _, disk := range c.storage.Disks {
		// mlog.Info("the mystery of the year(%s)", disk.Path)
		disk.Print()
	}

	mlog.Info("=========================================================")
	mlog.Info("Results for %s", srcDisk.Path)
	mlog.Info("Original Free Space: %s", helper.ByteSize(srcDisk.Free))
	mlog.Info("Final Free Space: %s", helper.ByteSize(srcDisk.NewFree))
	mlog.Info("Gained Space: %s", helper.ByteSize(srcDisk.NewFree-srcDisk.Free))
	mlog.Info("Bytes To Move: %s", helper.ByteSize(c.storage.BytesToMove))
	mlog.Info("---------------------------------------------------------")

	c.storage.Print()

	msg.Reply <- c.storage

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

func (c *Core) processDiskMv(line string, arg interface{}) {
	outbound := &dto.MessageOut{Topic: "storage:move:progress", Payload: line}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	mlog.Info(line)
}

func (c *Core) doStorageMove(msg *pubsub.Message) {
	// var commands []*dto.Move

	// commands = make([]*dto.Move, 0)

	var dry string
	if c.settings.DryRun {
		dry = "-t"
	} else {
		dry = "-f"
	}

	var diskmv string
	env := os.Getenv(dockerEnv)
	if env == "y" {
		diskmv = diskmvDockerCmd
	} else {
		diskmv = diskmvCmd
	}

	c.storage.InProgress = true

	outbound := &dto.MessageOut{Topic: "storage:move:begin", Payload: "Operation started"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// if err := c.sendmail("Move Operation started"); err != nil {
	// 	mlog.Error(err)
	// }

	started := time.Now()

	commands := make([]string, 0)

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

			commands = append(commands, cmd+"\n")

			outbound = &dto.MessageOut{Topic: "storage:move:progress", Payload: cmd}
			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

			err := helper.Shell(cmd, c.processDiskMv, nil)
			if err != nil {
				mlog.Info("error running the diskmv command: %s", err.Error())
				c.storage.InProgress = false

				txt := fmt.Sprintf("Move command was closed prematurely: %s", err.Error())
				outbound = &dto.MessageOut{Topic: "storage:move:progress", Payload: txt}
				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

				outbound = &dto.MessageOut{Topic: "storage:move:end", Payload: "Operation finished"}
				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

				finished := time.Now()
				elapsed := time.Since(started)

				message := fmt.Sprintf("There was an error when executing\n\n%s\n\nThese are the commands that were executed:\n\n%s\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", cmd, c.printCommands(commands), started, finished, elapsed)
				if sendErr := c.sendmail(message); sendErr != nil {
					mlog.Error(sendErr)
				}

				mlog.Error(err)
				return
			}

		}
	}

	c.storage.InProgress = false

	outbound = &dto.MessageOut{Topic: "storage:move:end", Payload: "Operation finished"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	finished := time.Now()
	elapsed := time.Since(started)

	message := fmt.Sprintf("Move operation completed.\n\nThese are the commands that were executed:\n\n%s\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", c.printCommands(commands), started, finished, elapsed)
	if sendErr := c.sendmail(message); sendErr != nil {
		mlog.Error(sendErr)
	}

}

func (c *Core) doStorageUpdate(msg *pubsub.Message) {
	outbound := &dto.MessageOut{Topic: "storage:update:completed", Payload: c.storage.Refresh()}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
}

func (c *Core) sendmail(msg string) error {
	if !c.settings.Notifications {
		return nil
	}

	from := "From: " + c.settings.NotiFrom
	to := "To: " + c.settings.NotiTo
	subject := "Subject: unBALANCE Notification"

	var dry string
	if c.settings.DryRun {
		dry = "-------\nDRY RUN\n-------\n"
	} else {
		dry = ""
	}

	// mail := "\"" + from + "\n" + to + "\n" + subject + "\n\n" + dry + msg + "\""
	mail := from + "\n" + to + "\n" + subject + "\n\n" + dry + msg

	// err := ioutil.WriteFile(msgLocation, []byte(mail), 0644)
	// if err != nil {
	// 	return err
	// }

	echo := exec.Command("echo", "-e", mail)
	ssmtp := exec.Command("ssmtp", c.settings.NotiTo)

	// mlog.Info("sendmail:echo: %s %s (%s)", echo.Path, echo.Args, echo.Dir)
	// mlog.Info("sendmail:ssmtp: %s %s (%s)", ssmtp.Path, ssmtp.Args, ssmtp.Dir)

	_, _, err := helper.Pipeline(echo, ssmtp)
	if err != nil {
		return err
	}

	return nil

}

func (c *Core) printCommands(list []string) string {
	var str string
	for _, value := range list {
		str += value
	}
	return str
}
