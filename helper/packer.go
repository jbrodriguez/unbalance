package helper

import (
	"apertoire.net/unbalance/model"
	"fmt"
	"github.com/golang/glog"
	"sort"
)

type Packer struct {
	// SourceDisk string
	// TargetDisk string
	// MaxSize    uint64

	disk *model.Disk

	Bins []*model.Bin
	list []*model.Item
	over []*model.Item
}

func NewPacker(disk *model.Disk, items []*model.Item) *Packer {
	p := new(Packer)
	p.disk = disk
	p.list = items
	return p
}

func (self *Packer) BestFit() (bin *model.Bin) {
	sort.Sort(model.BySize(self.list))

	for _, item := range self.list {
		if item.Size > self.disk.Free {
			// glog.Info(fmt.Sprintf("size: %d, disk: %s, free: %d", item.Size, self.disk.Path, self.disk.Free))
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
				self.Bins[targetBin].Add(item)
			} else {
				newbin := &model.Bin{}
				newbin.Add(item)
				self.Bins = append(self.Bins, newbin)
			}
		}
	}

	if len(self.Bins) > 0 {
		sort.Sort(model.ByFilled(self.Bins))
		self.disk.Bin = self.Bins[0]
		bin = self.disk.Bin
	}

	return bin
}

func (self *Packer) add(item *model.Item) {
	if item.Size > self.disk.Free {
		self.over = append(self.over, item)
	} else {
		self.list = append(self.list, item)
	}
}

func (self *Packer) printList() {
	for _, item := range self.list {
		glog.Info(fmt.Sprintf("Item (%s): %d", item.Name, item.Size))
	}
}

func (self *Packer) sortBins() {
	sort.Sort(model.ByFilled(self.Bins))
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
