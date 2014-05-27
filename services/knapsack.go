package services

import (
	"apertoire.net/unbalance/bus"
	"apertoire.net/unbalance/helper"
	"apertoire.net/unbalance/message"
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

type Knapsack struct {
	Bus *bus.Bus

	reFreeSpace *regexp.Regexp
	reItems     *regexp.Regexp
}

func (self *Knapsack) Start() {
	log.Printf("starting Knapsack service ...")

	re, _ := regexp.Compile(`(.*?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.*?)\s+(.*?)$`)
	self.reFreeSpace = re

	re, _ = regexp.Compile(`(.\d+)\s+(.*?)$`)
	self.reItems = re

	go self.react()

	log.Printf("Knapsack service started")
}

func (self *Knapsack) Stop() {
	// nothing right now
	log.Printf("Knapsack service stopped")
}

func (self *Knapsack) react() {
	for {
		select {
		case msg := <-self.Bus.GetBestFit:
			go self.doGetBestFit(msg)
		}
	}
}

func (self *Knapsack) removeFolders(folders []*helper.Item, list []*helper.Item) []*helper.Item {
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

func (self *Knapsack) doGetBestFit(msg *message.FitData) {
	disks, srcDiskSizeFreeOriginal, _ := self.GetDisks(msg.SourceDisk, msg.TargetDisk)

	// folders := []*helper.Item{&helper.Item{Name: "/The Godfather (1974)", Size: 34, Path: "films/bluray"}, &helper.Item{Name: "/The Mist (2010)", Size: 423, Path: "films/bluray"}, &helper.Item{Name: "/Aventador (1974)", Size: 3524, Path: "films/bluray"}, &helper.Item{Name: "/Countach (1974)", Size: 3432, Path: "films/bluray"}, &helper.Item{Name: "/Iroc-Z (1974)", Size: 6433, Path: "films/bluray"}}
	// // items := []*helper.Item{&helper.Item{Name: "/The Godfather (1974)", Size: 34, Path: "films/bluray"}, &helper.Item{Name: "/Aventador (1974)", Size: 3524, Path: "films/bluray"}}
	// items := []*helper.Item{&helper.Item{Name: "/Aventador (1974)", Size: 3524, Path: "films/bluray"}}

	// folders = self.removeFolder(folders, items)

	// for _, itm := range folders {
	// 	log.Println("yes: ", itm.Name)
	// }

	sort.Sort(helper.ByFree(disks))

	var folders []*helper.Item
	paths := []string{"films/bluray", "films/blurip"}

	for _, path := range paths {
		list := self.GetFolders(msg.SourceDisk, path)
		folders = append(folders, list...)
	}

	srcDiskSizeFreeFinal := srcDiskSizeFreeOriginal

	for _, disk := range disks {
		packer := helper.NewPacker(disk, folders)
		bin := packer.BestFit()
		if bin != nil {
			srcDiskSizeFreeFinal += bin.Size
			self.removeFolders(folders, bin.Items)
		}
	}

	for _, disk := range disks {
		disk.Print()
	}

	fmt.Println("=========================================================")
	fmt.Println(fmt.Sprintf("Results for %s", msg.SourceDisk))
	fmt.Println(fmt.Sprintf("Original Free Space: %s", helper.ByteSize(srcDiskSizeFreeOriginal)))
	fmt.Println(fmt.Sprintf("Final Free Space: %s", helper.ByteSize(srcDiskSizeFreeFinal)))
	fmt.Println(fmt.Sprintf("Gained Space: %s", helper.ByteSize(srcDiskSizeFreeFinal-srcDiskSizeFreeOriginal)))
	fmt.Println("---------------------------------------------------------")

	// free, err := packer.GetFreeSpace()
	// if err != nil {
	// 	log.Println(fmt.Sprintf("Available Space on %s: %d", msg.TargetDisk, free))
	// }

	// packer.GetItems("films/bluray")
	// packer.GetItems("films/blurip")

	// // packer.print()

	// packer.BestFit()

	// packer.Print()

	// items := self.getItems(msg.SourceDisk, "films/blurip")
	// for item := range items {
	// 	packer.Add(item)
	// }

	// items := self.getItems(msg.SourceDisk, "films/xvid")
	// for item := range items {
	// 	packer.Add(item)
	// }

	// items := self.getItems(msg.SourceDisk, "films/dvd")
	// for item := range items {
	// 	packer.Add(item)
	// }
}

func (self *Knapsack) GetDisks(src string, dst string) (disks []*helper.Disk, srcDiskFree uint64, err error) {
	// var disks []Disk

	cmd := exec.Command("sh", "-c", "df --block-size=1 /mnt/disk*")
	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Unable to stdoutpipe df: ", err)
	}

	rd := bufio.NewReader(out)

	if err := cmd.Start(); err != nil {
		log.Fatal("Unable to start df: ", err)
	}

	// ignore first line since it's just headers
	line, err := rd.ReadString('\n')

	for {
		line, err = rd.ReadString('\n')
		if err == io.EOF && len(line) == 0 {
			// Good end of file with no partial line
			break
		}
		if err == io.EOF {
			log.Fatal("Last line not terminated: ", err)
		}
		line = line[:len(line)-1] // drop the '\n'
		if line[len(line)-1] == '\r' {
			line = line[:len(line)-1] // drop the '\r'
		}

		// Filesystem           1B-blocks      Used Available Use% Mounted on
		// /dev/md1             2000337846272 1998411968512 1925877760 100% /mnt/disk1

		result := self.reFreeSpace.FindStringSubmatch(line)
		free, _ := strconv.ParseUint(result[4], 10, 64)

		if result[6] == src {
			srcDiskFree = free
		}

		if dst != "" {
			if dst == result[6] {
				disks = append(disks, &helper.Disk{Path: result[6], Free: free})
				// break
			}
		} else {
			if src != result[6] {
				disks = append(disks, &helper.Disk{Path: result[6], Free: free})
			}
		}
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Unable to wait for process to finish: ", err)
	}

	return disks, srcDiskFree, nil
}

func (self *Knapsack) GetFolders(src string, folder string) (items []*helper.Item) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("du -bs %s", filepath.Join(src, folder, "*")))
	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Unable to stdoutpipe du: ", err)
	}

	rd := bufio.NewReader(out)

	if err := cmd.Start(); err != nil {
		log.Fatal("Unable to start du: ", err)
	}

	for {
		line, err := rd.ReadString('\n')
		if err == io.EOF && len(line) == 0 {
			// Good end of file with no partial line
			break
		}
		if err == io.EOF {
			log.Fatal("Last line not terminated: ", err)
		}
		line = line[:len(line)-1] // drop the '\n'
		if line[len(line)-1] == '\r' {
			line = line[:len(line)-1] // drop the '\r'
		}

		result := self.reItems.FindStringSubmatch(line)
		log.Printf("[%s] %s", result[1], result[2])

		size, _ := strconv.ParseUint(result[1], 10, 64)

		items = append(items, &helper.Item{Name: result[2], Size: size, Path: folder})
		// fmt.Println(line)
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Unable to wait for process to finish: ", err)
	}

	// out, err := helper.Shell(fmt.Sprintf("du -sh %s", filepath.Join(disk, folder, "*")))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// log.Println(string(out))
	log.Println("done")
	return items
}
