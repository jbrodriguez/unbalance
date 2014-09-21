package services

import (
	"apertoire.net/unbalance/server/dto"
	"apertoire.net/unbalance/server/lib"
	"apertoire.net/unbalance/server/model"
	"bufio"
	"fmt"
	"github.com/apertoire/mlog"
	"github.com/apertoire/pubsub"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

type Core struct {
	bus     *pubsub.PubSub
	storage *model.Unraid

	chanStorageInfo      chan *pubsub.Message
	chanCalculateBestFit chan *pubsub.Message

	reFreeSpace *regexp.Regexp
	reItems     *regexp.Regexp
}

func NewCore(bus *pubsub.PubSub) *Core {
	core := &Core{bus: bus}

	re, _ := regexp.Compile(`(.*?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.*?)\s+(.*?)$`)
	core.reFreeSpace = re

	re, _ = regexp.Compile(`(.\d+)\s+(.*?)$`)
	core.reItems = re

	core.storage = &model.Unraid{}

	return core
}

func (self *Core) Start() {
	mlog.Info("starting service Core ...")

	self.chanStorageInfo = self.bus.Sub("cmd.getStorageInfo")
	self.chanCalculateBestFit = self.bus.Sub("cmd.calculateBestFit")

	go self.react()
}

func (self *Core) Stop() {
	mlog.Info("stopped service Core ...")
}

func (self *Core) react() {
	for {
		select {
		case msg := <-self.chanStorageInfo:
			go self.getStorageInfo(msg)
		case msg := <-self.chanCalculateBestFit:
			go self.calculateBestFit(msg)
		}
	}
}

func (self *Core) getStorageInfo(msg *pubsub.Message) {
	mlog.Info("La vita e bella")

	msg.Reply <- self.storage.Refresh()
}

func (self *Core) calculateBestFit(msg *pubsub.Message) {
	disks := make([]*model.Disk, len(self.storage.Disks))
	copy(disks, self.storage.Disks)

	dto := msg.Payload.(*dto.BestFit)

	var srcDisk *model.Disk
	for _, disk := range disks {
		if disk.Path == dto.SourceDisk {
			srcDisk = disk
		}
	}

	mlog.Info("srcDisk = %s", dto.SourceDisk)

	self.storage.SourceDiskName = srcDisk.Path

	sort.Sort(model.ByFree(disks))

	var folders []*model.Item
	paths := []string{"films/bluray", "films/blurip"}

	for _, path := range paths {
		list := self.getFolders(dto.SourceDisk, path)

		if list != nil {
			folders = append(folders, list...)
		}
	}

	// srcDiskSizeFreeFinal := srcDiskSizeFreeOriginal
	srcDisk.NewFree = srcDisk.Free

	for _, disk := range disks {
		disk.NewFree = disk.Free
		if disk.Path != srcDisk.Path {
			packer := lib.NewKnapsack(disk, folders)
			bin := packer.BestFit()
			if bin != nil {
				// srcDiskSizeFreeFinal += bin.Size
				srcDisk.NewFree += bin.Size
				disk.NewFree -= bin.Size
				self.storage.BytesToMove += bin.Size

				folders = self.removeFolders(folders, bin.Items)
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

	msg.Reply <- self.storage
}

func (self *Core) getFolders(src string, folder string) (items []*model.Item) {
	srcFolder := filepath.Join(src, folder)
	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
		mlog.Info("Folder does not exist ", srcFolder)
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

		result := self.reItems.FindStringSubmatch(line)
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
	mlog.Info("done")
	return items
}

func (self *Core) removeFolders(folders []*model.Item, list []*model.Item) []*model.Item {
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
