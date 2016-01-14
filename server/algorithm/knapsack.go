package algorithm

import (
	"fmt"
	"github.com/jbrodriguez/mlog"
	"jbrodriguez/unbalance/server/model"
	"sort"
)

type Knapsack struct {
	disk *model.Disk

	Bins []*model.Bin
	list []*model.Item
	over []*model.Item

	buffer int64
}

func max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

func NewKnapsack(disk *model.Disk, items []*model.Item, amount int64, unit string, floor int64) *Knapsack {
	p := new(Knapsack)
	p.disk = disk
	p.list = items

	var reserved int64
	switch unit {
	case "%":
		fcalc := disk.Size * amount / 100
		reserved = int64(fcalc)
		break
	case "Mb":
		reserved = amount * 1000 * 1000
		break
	case "Gb":
		reserved = amount * 1000 * 1000 * 1000
		break
	default:
		reserved = floor
	}

	p.buffer = max(floor, reserved)
	return p
}

func (self *Knapsack) BestFit() (bin *model.Bin) {
	sort.Sort(model.BySize(self.list))

	// for _, itm := range self.list {
	// 	mlog.Info("disk (%s) > item: %s", self.disk.Path, itm.Path)
	// }

	for _, item := range self.list {
		if item.Size > (self.disk.Free - self.buffer) {
			// if item.Size > self.disk.Free {
			// mlog.Info("size: %d, disk: %s, free: %d", item.Size, self.disk.Path, self.disk.Free)
			self.over = append(self.over, item)
		} else {
			targetBin := -1
			remainingSpace := self.disk.Free

			// mlog.Info("disk(%s)-bins(%d); item(%s)-size(%d); remainingSpace(%d)", self.disk.Name, len(self.Bins), item.Name, item.Size, remainingSpace)

			for i, bin := range self.Bins {
				binSpaceUsed := bin.Size
				binSpaceLeft := self.disk.Free - binSpaceUsed - item.Size

				// mlog.Info("su(%d); sl(%d)", binSpaceUsed, binSpaceLeft)
				// if self.disk.Path == "/mnt/disk8" {
				// 	mlog.Info("[/mnt/disk/14] Bin: %d ", i)
				// }

				if binSpaceLeft < remainingSpace && binSpaceLeft >= self.buffer {
					// mlog.Info("[%s] Used: %d | Left: %d\n", self.disk.Path, binSpaceUsed, binSpaceLeft)
					// mlog.Info("Disk: %s Folder: %s Bin: %d Used: %d | Left: %d\n", self.disk.Path, item.Name, i, binSpaceUsed, binSpaceLeft)

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
