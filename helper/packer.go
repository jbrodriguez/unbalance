package helper

import (
	"fmt"
	"log"
	"sort"
)

type Packer struct {
	// SourceDisk string
	// TargetDisk string
	// MaxSize    uint64

	disk *Disk

	Bins []*Bin
	list []*Item
	over []*Item
}

func NewPacker(disk *Disk, items []*Item) *Packer {
	p := new(Packer)
	p.disk = disk
	p.list = items
	return p
}

func (self *Packer) BestFit() (bin *Bin) {
	sort.Sort(BySize(self.list))

	for _, item := range self.list {
		if item.Size > self.disk.Free {
			log.Println(fmt.Sprintf("size: %d, disk: %s, free: %d", item.Size, self.disk.Path, self.disk.Free))
			self.over = append(self.over, item)
		} else {
			targetBin := -1
			remainingSpace := self.disk.Free

			for i, bin := range self.Bins {
				binSpaceUsed := bin.Size
				binSpaceLeft := self.disk.Free - binSpaceUsed - item.Size

				if binSpaceLeft < remainingSpace && binSpaceLeft >= 0 {
					remainingSpace = binSpaceLeft
					targetBin = i
				}
			}

			if targetBin >= 0 {
				self.Bins[targetBin].add(item)
			} else {
				newbin := &Bin{}
				newbin.add(item)
				self.Bins = append(self.Bins, newbin)
			}
		}
	}

	if len(self.Bins) > 0 {
		sort.Sort(ByFilled(self.Bins))
		self.disk.Bin = self.Bins[0]
		bin = self.disk.Bin
	}

	return bin
}

func (self *Packer) add(item *Item) {
	if item.Size > self.disk.Free {
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
		fmt.Println(fmt.Sprintf("%0d [%d/%d] %2.2f%% (%s)", i, bin.Size, self.disk.Free, (float64(bin.Size)/float64(self.disk.Free))*100, self.disk.Path))
		fmt.Println("---------------------------------------------------------")
		bin.Print()
		fmt.Println("---------------------------------------------------------")
		fmt.Println("")
	}
}
