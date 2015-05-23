package lib

import (
	"apertoire.net/unbalance/server/model"
	"fmt"
	"github.com/jbrodriguez/mlog"
	"sort"
)

type Knapsack struct {
	// SourceDisk string
	// TargetDisk string
	// MaxSize    uint64

	disk *model.Disk

	Bins []*model.Bin
	list []*model.Item
	over []*model.Item

	buffer uint64
}

func NewKnapsack(disk *model.Disk, items []*model.Item, reserved uint64) *Knapsack {
	p := new(Knapsack)
	p.disk = disk
	p.list = items
	p.buffer = reserved
	return p
}

func (self *Knapsack) BestFit() (bin *model.Bin) {
	sort.Sort(model.BySize(self.list))

	// for _, itm := range self.list {
	// 	mlog.Info("disk (%s) > item: %s", self.disk.Path, itm.Path)
	// }

	for _, item := range self.list {
		// if item.Size > (self.disk.Free - self.buffer) {
		if item.Size > self.disk.Free {
			// glog.Info(fmt.Sprintf("size: %d, disk: %s, free: %d", item.Size, self.disk.Path, self.disk.Free))
			self.over = append(self.over, item)
		} else {
			targetBin := -1
			remainingSpace := self.disk.Free

			// log.Printf("Disk [%s]: remainingSpace: %d\n", self.disk.Name, remainingSpace)

			for i, bin := range self.Bins {
				binSpaceUsed := bin.Size
				binSpaceLeft := self.disk.Free - binSpaceUsed - item.Size

				// if self.disk.Path == "/mnt/disk8" {
				// 	log.Printf("[/mnt/disk/8] Bin: %d ", i)
				// }

				if binSpaceLeft < remainingSpace && binSpaceLeft >= self.buffer {
					// log.Printf("[%s] Used: %d | Left: %d\n", self.disk.Path, binSpaceUsed, binSpaceLeft)
					// log.Printf("Disk: %s Folder: %s Bin: %d Used: %d | Left: %d\n", self.disk.Path, item.Name, i, binSpaceUsed, binSpaceLeft)

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

// func (self *Knapsack) add(item *model.Item) {
// 	if item.Size > self.disk.Free {
// 		self.over = append(self.over, item)
// 	} else {
// 		self.list = append(self.list, item)
// 	}
// }

func (self *Knapsack) printList() {
	for _, item := range self.list {
		mlog.Info(fmt.Sprintf("Item (%s): %d", item.Name, item.Size))
	}
}

func (self *Knapsack) sortBins() {
	sort.Sort(model.ByFilled(self.Bins))
}

func (self *Knapsack) Print() {
	for i, bin := range self.Bins {
		mlog.Info("=========================================================")
		mlog.Info(fmt.Sprintf("%0d [%d/%d] %2.2f%% (%s)", i, bin.Size, self.disk.Free, (float64(bin.Size)/float64(self.disk.Free))*100, self.disk.Path))
		mlog.Info("---------------------------------------------------------")
		bin.Print()
		mlog.Info("---------------------------------------------------------")
		mlog.Info("")
	}
}
