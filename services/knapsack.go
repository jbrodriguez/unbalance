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

func (self *Knapsack) doGetBestFit(msg *message.FitData) {
	disks, _ := self.GetDisks(msg.SourceDisk, msg.TargetDisk)

	var folders []helper.Item
	paths := []string{"films/bluray", "films/blurip"}

	for _, path := range paths {
		list := self.GetFolders(msg.SourceDisk, path)
		folders = append(folders, list...)
	}

	for _, disk := range disks {
		packer := helper.NewPacker(disk, folders)
		packer.BestFit()
		// self.RemoveFolders(bin)
	}

	for _, disk := range disks {
		disk.Print()
	}

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

func (self *Knapsack) GetDisks(src string, dst string) (disks []helper.Disk, err error) {
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

		result := self.reFreeSpace.FindStringSubmatch(line)
		free, _ := strconv.ParseUint(result[4], 10, 64)

		if dst != "" {
			if dst == result[6] {
				disks = append(disks, helper.Disk{Path: result[6], Free: free})
				break
			}
		} else {
			if src != result[6] {
				disks = append(disks, helper.Disk{Path: result[6], Free: free})
			}
		}
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Unable to wait for process to finish: ", err)
	}

	return disks, nil
}

func (self *Knapsack) GetFolders(src string, folder string) (items []helper.Item) {
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

		items = append(items, helper.Item{Name: result[2], Size: size, Path: folder})
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
