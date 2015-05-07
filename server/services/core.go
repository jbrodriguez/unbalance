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
)

const dockerEnv = "UNBALANCE_DOCKER"
const diskmvCmd = "./diskmv"
const diskmvDockerCmd = "/usr/bin/diskmv"

type Core struct {
	bus     *pubsub.PubSub
	storage *model.Unraid
	config  *model.Config

	chanConfigInfo       chan *pubsub.Message
	chanSaveConfig       chan *pubsub.Message
	chanStorageInfo      chan *pubsub.Message
	chanCalculateBestFit chan *pubsub.Message
	chanMove             chan *pubsub.Message
	storageMove          chan *pubsub.Message
	storageUpdate        chan *pubsub.Message

	reFreeSpace *regexp.Regexp
	reItems     *regexp.Regexp
}

func NewCore(bus *pubsub.PubSub, config *model.Config) *Core {
	core := &Core{bus: bus, config: config}

	re, _ := regexp.Compile(`(.*?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.*?)\s+(.*?)$`)
	core.reFreeSpace = re

	re, _ = regexp.Compile(`(.\d+)\s+(.*?)$`)
	core.reItems = re

	core.storage = &model.Unraid{}

	core.chanConfigInfo = core.bus.Sub("cmd.getConfig")
	core.chanSaveConfig = core.bus.Sub("cmd.saveConfig")
	core.chanStorageInfo = core.bus.Sub("cmd.getStorageInfo")
	core.chanCalculateBestFit = core.bus.Sub("cmd.calculateBestFit")
	core.chanMove = core.bus.Sub("cmd.move")
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
		case msg := <-c.chanMove:
			go c.move(msg)
		case msg := <-c.storageMove:
			go c.doStorageMove(msg)
		case msg := <-c.storageUpdate:
			go c.doStorageUpdate(msg)
		}
	}
}

func (c *Core) getConfigInfo(msg *pubsub.Message) {
	mlog.Info("Sending config")

	msg.Reply <- c.config
}

func (c *Core) saveConfig(msg *pubsub.Message) {
	mlog.Info("Saving config")

	config := msg.Payload.(*model.Config)
	c.config.Folders = config.Folders
	c.config.DryRun = config.DryRun

	c.config.Save()

	msg.Reply <- c.config
}

func (c *Core) getStorageInfo(msg *pubsub.Message) {
	msg.Reply <- c.storage.Refresh()
}

func (c *Core) calculateBestFit(msg *pubsub.Message) {
	disks := make([]*model.Disk, len(c.storage.Disks))
	copy(disks, c.storage.Disks)

	dto := msg.Payload.(*dto.BestFit)

	var srcDisk *model.Disk
	for _, disk := range disks {
		if disk.Path == dto.SourceDisk {
			srcDisk = disk
		}
	}

	// Initializae fields
	c.storage.BytesToMove = 0
	c.storage.SourceDiskName = srcDisk.Path

	sort.Sort(model.ByFree(disks))

	var folders []*model.Item
	for _, path := range c.config.Folders {
		list := c.getFolders(dto.SourceDisk, path)

		if list != nil {
			folders = append(folders, list...)
		}
	}

	srcDisk.NewFree = srcDisk.Free

	for _, disk := range disks {
		//		disk.NewFree = disk.Free
		if disk.Path != srcDisk.Path {
			disk.NewFree = disk.Free

			packer := lib.NewKnapsack(disk, folders, c.config.ReservedSpace)
			bin := packer.BestFit()
			if bin != nil {
				srcDisk.NewFree += bin.Size
				disk.NewFree -= bin.Size
				c.storage.BytesToMove += bin.Size

				mlog.Info("Original Free Space: %s", lib.ByteSize(srcDisk.Free))
				mlog.Info("Final Free Space: %s", lib.ByteSize(srcDisk.NewFree))

				folders = c.removeFolders(folders, bin.Items)
			}
		}
	}

	for _, disk := range disks {
		disk.Print()
	}

	mlog.Info("=========================================================")
	mlog.Info("Results for %s", srcDisk.Path)
	mlog.Info("Original Free Space: %s", lib.ByteSize(srcDisk.Free))
	mlog.Info("Final Free Space: %s", lib.ByteSize(srcDisk.NewFree))
	mlog.Info("Gained Space: %s", lib.ByteSize(srcDisk.NewFree-srcDisk.Free))
	mlog.Info("---------------------------------------------------------")

	c.storage.Print()

	msg.Reply <- c.storage
}

func (c *Core) getFolders(src string, folder string) (items []*model.Item) {
	srcFolder := filepath.Join(src, folder)

	// mlog.Info("Folder is: %s", srcFolder)
	// _, err := os.Stat(filepath.Join("/mnt/disk13/films", "*"))
	// mlog.Info("Error: %s", err)

	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
		mlog.Warning("Folder does not exist: %s", srcFolder)
		return nil
	}

	dirs, err := ioutil.ReadDir(srcFolder)
	if err != nil {
		mlog.Fatalf("Unable to readdir: %s", err)
	}

	mlog.Info("Dirs: %+v", dirs)

	if len(dirs) == 0 {
		mlog.Info("No subdirectories under %s", srcFolder)
		return nil
	}

	cmd := exec.Command("sh", "-c", fmt.Sprintf("du -bs %s", filepath.Join(srcFolder, "*")))
	out, err := cmd.StdoutPipe()
	if err != nil {
		mlog.Fatalf("Unable to stdoutpipe du: %s", err)
	}

	rd := bufio.NewReader(out)

	if err := cmd.Start(); err != nil {
		mlog.Fatalf("Unable to start du: %s", err)
	}

	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF && len(line) == 0 {
			// Good end of file with no partial line
			break
		}
		if err == io.EOF {
			mlog.Fatalf("Last line not terminated: %s", err)
		}
		line = line[:len(line)-1] // drop the '\n'
		if line[len(line)-1] == '\r' {
			line = line[:len(line)-1] // drop the '\r'
		}

		result := c.reItems.FindStringSubmatch(line)
		// mlog.Info("[%s] %s", result[1], result[2])

		size, _ := strconv.ParseUint(result[1], 10, 64)

		item := &model.Item{Name: result[2], Size: size, Path: filepath.Join(folder, filepath.Base(result[2]))}
		items = append(items, item)
		// fmt.Println(line)
		mlog.Info("item: %+v", item)
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		mlog.Fatalf("Unable to wait for process to finish: %s", err)
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

func (c *Core) move(msg *pubsub.Message) {
	var commands []*dto.Move

	commands = make([]*dto.Move, 0)

	for _, disk := range c.storage.Disks {
		if disk.Bin == nil || disk.Path == c.storage.SourceDiskName {
			continue
		}

		for _, item := range disk.Bin.Items {
			dst := filepath.Join(disk.Path, item.Path)

			mlog.Info("disk.Path = %s | item.Path = %s | dst = %s", disk.Path, item.Path, c.storage.SourceDiskName)
			// mlog.Info("disk.Path = %s | item.Name = %s | item.Path = %s | dst = %s", disk.Path, item.Name, item.Path, dst)
			// mlog.Info("mv %s %s", strconv.Quote(item.Name), strconv.Quote(dst))
			command := &dto.Move{Command: fmt.Sprintf("mv %s %s", strconv.Quote(item.Name), strconv.Quote(dst))}
			commands = append(commands, command)

			cmd := fmt.Sprintf("./diskmv \"%s\" %s %s", item.Path, c.storage.SourceDiskName, disk.Path)
			mlog.Info("cmd = %s", cmd)

			helper.Shell(cmd, c.processDiskMv, nil)

			// mover.Src = item.Name
			// mover.Dst = dst
			// mover.Progress = progress

			// glog.Infof("mover: %+v", mover)

			// mover.Copy()
			// for {
			// 	select {
			// 	case msg := <-mover.ProgressCh:
			// 		glog.Infof("Progress: %+v", msg)
			// 	case <-mover.DoneCh:
			// 		return
			// 	}
			// }
		}
	}

	msg.Reply <- commands
}

func (c *Core) doStorageMove(msg *pubsub.Message) {
	// var commands []*dto.Move

	// commands = make([]*dto.Move, 0)

	var dry string
	if c.config.DryRun {
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

	for _, disk := range c.storage.Disks {
		if disk.Bin == nil || disk.Path == c.storage.SourceDiskName {
			continue
		}

		for _, item := range disk.Bin.Items {
			//			dst := filepath.Join(disk.Path, item.Path)

			mlog.Info("disk.Path = %s | item.Path = %s | dst = %s", disk.Path, item.Path, c.storage.SourceDiskName)
			// mlog.Info("disk.Path = %s | item.Name = %s | item.Path = %s | dst = %s", disk.Path, item.Name, item.Path, dst)
			// mlog.Info("mv %s %s", strconv.Quote(item.Name), strconv.Quote(dst))
			//			command := &dto.Move{Command: fmt.Sprintf("mv %s %s", strconv.Quote(item.Name), strconv.Quote(dst))}
			//			commands = append(commands, command)

			cmd := fmt.Sprintf("%s %s \"%s\" %s %s", diskmv, dry, item.Path, c.storage.SourceDiskName, disk.Path)
			mlog.Info("cmd = %s", cmd)

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

				mlog.Error(err)
				return
			}

		}
	}

	c.storage.InProgress = false

	outbound = &dto.MessageOut{Topic: "storage:move:end", Payload: "Operation finished"}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

}

func (c *Core) doStorageUpdate(msg *pubsub.Message) {
	outbound := &dto.MessageOut{Topic: "storage:update:completed", Payload: c.storage.Refresh()}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
}
