package algorithm

import (
	"fmt"
	"github.com/jbrodriguez/mlog"
	"jbrodriguez/unbalance/server/src/model"
	"sort"
)

// Knapsack -
type Knapsack struct {
	disk *model.Disk

	Bins []*model.Bin
	list []*model.Item
	over []*model.Item

	buffer int64
}

// NewKnapsack -
func NewKnapsack(disk *model.Disk, items []*model.Item, reserved int64) *Knapsack {
	p := new(Knapsack)
	p.disk = disk
	p.list = items
	p.buffer = reserved
	return p
}

// BestFit -
func (k *Knapsack) BestFit() (bin *model.Bin) {
	sort.Sort(model.BySize(k.list))

	// for _, itm := range k.list {
	// 	mlog.Info("disk (%s) > item: %s", k.disk.Path, itm.Path)
	// }

	for _, item := range k.list {
		if item.Size > (k.disk.Free - k.buffer) {
			// if item.Size > k.disk.Free {
			// mlog.Info("size: %d, disk: %s, free: %d", item.Size, k.disk.Path, k.disk.Free)
			k.over = append(k.over, item)
		} else {
			targetBin := -1
			remainingSpace := k.disk.Free

			// mlog.Info("disk(%s)-bins(%d); item(%s)-size(%d); remainingSpace(%d)", k.disk.Name, len(k.Bins), item.Name, item.Size, remainingSpace)

			for i, bin := range k.Bins {
				binSpaceUsed := bin.Size
				binSpaceLeft := k.disk.Free - binSpaceUsed - item.Size

				// mlog.Info("su(%d); sl(%d)", binSpaceUsed, binSpaceLeft)
				// if k.disk.Path == "/mnt/disk8" {
				// 	mlog.Info("[/mnt/disk/14] Bin: %d ", i)
				// }

				if binSpaceLeft < remainingSpace && binSpaceLeft >= k.buffer {
					// mlog.Info("[%s] Used: %d | Left: %d\n", k.disk.Path, binSpaceUsed, binSpaceLeft)
					// mlog.Info("Disk: %s Folder: %s Bin: %d Used: %d | Left: %d\n", k.disk.Path, item.Name, i, binSpaceUsed, binSpaceLeft)

					remainingSpace = binSpaceLeft
					targetBin = i
				}
			}

			if targetBin >= 0 {
				k.Bins[targetBin].Add(item)
			} else {
				newbin := &model.Bin{}
				newbin.Add(item)
				k.Bins = append(k.Bins, newbin)
			}
		}
	}

	if len(k.Bins) > 0 {
		sort.Sort(model.ByFilled(k.Bins))
		k.disk.Bin = k.Bins[0]
		bin = k.disk.Bin
	}

	return bin
}

// func (k *Knapsack) add(item *model.Item) {
// 	if item.Size > k.disk.Free {
// 		k.over = append(k.over, item)
// 	} else {
// 		k.list = append(k.list, item)
// 	}
// }

func (k *Knapsack) printList() {
	for _, item := range k.list {
		mlog.Info(fmt.Sprintf("Item (%s): %d", item.Name, item.Size))
	}
}

func (k *Knapsack) sortBins() {
	sort.Sort(model.ByFilled(k.Bins))
}

// Print -
func (k *Knapsack) Print() {
	for i, bin := range k.Bins {
		mlog.Info("=========================================================")
		mlog.Info(fmt.Sprintf("%0d [%d/%d] %2.2f%% (%s)", i, bin.Size, k.disk.Free, (float64(bin.Size)/float64(k.disk.Free))*100, k.disk.Path))
		mlog.Info("---------------------------------------------------------")
		bin.Print()
		mlog.Info("---------------------------------------------------------")
		mlog.Info("")
	}
}
