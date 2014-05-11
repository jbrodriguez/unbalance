package helper

import (
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

type Item struct {
	Name string
	Size uint64
	Path string
}

type Bin struct {
	Size uint64
	list []Item
}

func (self *Bin) add(item Item) {
	self.list = append(self.list, item)
	self.Size += item.Size
}

func (self *Bin) print() {
	for _, item := range self.list {
		fmt.Println(fmt.Sprintf("[%d] %s", item.Size, item.Name))
	}
}

type Packer struct {
	SourceDisk string
	TargetDisk string
	MaxSize    uint64

	Bins []Bin
	list []Item
	over []Item

	reFreeSpace *regexp.Regexp
	reItems     *regexp.Regexp
}

func NewPacker(src string, dst string) *Packer {
	p := new(Packer)
	p.SourceDisk = src
	p.TargetDisk = dst
	p.MaxSize = 0

	re, _ := regexp.Compile(`(.*?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.*?)\s+(.*?)$`)
	p.reFreeSpace = re

	re, _ = regexp.Compile(`(.\d+)\s+(.*?)$`)
	p.reItems = re

	return p
}

func (self *Packer) GetFreeSpace() (size uint64, err error) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("df --block-size=1 %s", self.TargetDisk))
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

	line, err = rd.ReadString('\n')
	// if err == io.EOF && len(line) == 0 {
	// 	// Good end of file with no partial line
	// 	return 0, nil
	// }
	if err == io.EOF {
		log.Fatal("Last line not terminated: ", err)
	}
	line = line[:len(line)-1] // drop the '\n'
	if line[len(line)-1] == '\r' {
		line = line[:len(line)-1] // drop the '\r'
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Unable to wait for process to finish: ", err)
	}

	result := self.reFreeSpace.FindStringSubmatch(line)
	free, err := strconv.ParseUint(result[4], 10, 64)
	self.MaxSize = free
	return free, err
}

func (self *Packer) GetItems(folder string) {
	cmd := exec.Command("sh", "-c", fmt.Sprintf("du -bs %s", filepath.Join(self.SourceDisk, folder, "*")))
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

		self.add(Item{Name: result[2], Size: size, Path: folder})

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
}

func (self *Packer) BestFit() {
	sort.Sort(BySize(self.list))

	for _, item := range self.list {
		if item.Size > self.MaxSize {
			self.over = append(self.over, item)
		} else {
			targetBin := -1
			remainingSpace := self.MaxSize

			for i, bin := range self.Bins {
				binSpaceUsed := bin.Size
				binSpaceLeft := self.MaxSize - binSpaceUsed - item.Size

				if binSpaceLeft < remainingSpace && binSpaceLeft >= 0 {
					remainingSpace = binSpaceLeft
					targetBin = i
				}
			}

			if targetBin >= 0 {
				self.Bins[targetBin].add(item)
			} else {
				newbin := Bin{}
				newbin.add(item)
				self.Bins = append(self.Bins, newbin)
			}
		}
	}

	sort.Sort(ByFilled(self.Bins))
}

func (self *Packer) add(item Item) {
	if item.Size > self.MaxSize {
		self.over = append(self.over, item)
	} else {
		self.list = append(self.list, item)
	}
}

func (self *Packer) printList() {
	for _, item := range self.list {
		log.Println(fmt.Sprintf("Item (%s): %d", item.Name, item.Size))
	}
}

func (self *Packer) sortBins() {
	sort.Sort(ByFilled(self.Bins))
}

func (self *Packer) Print() {
	for i, bin := range self.Bins {
		fmt.Println("=========================================================")
		fmt.Println(fmt.Sprintf("%0d [%d/%d] %2.2f%% (%s)", i, bin.Size, self.MaxSize, (float64(bin.Size)/float64(self.MaxSize))*100, self.TargetDisk))
		fmt.Println("---------------------------------------------------------")
		bin.print()
		fmt.Println("---------------------------------------------------------")
		fmt.Println("")
	}
}

type ByFilled []Bin

func (s ByFilled) Len() int           { return len(s) }
func (s ByFilled) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFilled) Less(i, j int) bool { return s[i].Size > s[j].Size }

type BySize []Item

func (s BySize) Len() int           { return len(s) }
func (s BySize) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BySize) Less(i, j int) bool { return s[i].Size > s[j].Size }
