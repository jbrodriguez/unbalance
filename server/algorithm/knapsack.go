package algorithm

import (
	"sort"

	"unbalance/domain"
)

// Knapsack -
type Knapsack struct {
	disk *domain.Disk

	Bins []*domain.Bin
	list []*domain.Item
	over []*domain.Item

	buffer    int64
	blockSize int64
}

// NewKnapsack -
func NewKnapsack(disk *domain.Disk, items []*domain.Item, reserved, blockSize int64) *Knapsack {
	p := &Knapsack{}

	p.disk = disk
	p.list = items
	p.buffer = reserved
	p.blockSize = blockSize

	return p
}

// BestFit -
func (k *Knapsack) BestFit() *domain.Bin {
	if k.blockSize > 0 {
		return k.fitBlocks()
	}

	return k.fitBytes()
}

func (k *Knapsack) fitBytes() (bin *domain.Bin) {
	sort.Slice(k.list, func(i, j int) bool { return k.list[i].Size > k.list[j].Size })

	// for _, itm := range k.list {
	// 	fmt.Printf("disk(%s):item(%s):size(%d)\n", k.disk.Path, itm.Path, itm.Size)
	// }

	for _, item := range k.list {
		// fmt.Printf("loop:size(%d):free(%d)-buffer(%d):resta(%d)\n", item.Size, k.disk.Free, k.buffer, k.disk.Free-k.buffer)

		if item.Size > (k.disk.Free - k.buffer) {
			// fmt.Printf("size: %d, disk: %s, free: %d\n", item.Size, k.disk.Path, k.disk.Free)
			k.over = append(k.over, item)
		} else {
			targetBin := -1
			remainingSpace := k.disk.Free

			// fmt.Printf("else:disk(%s)-bins(%d); item(%s)-size(%d); remainingSpace(%d)\n", k.disk.Name, len(k.Bins), item.Name, item.Size, remainingSpace)

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
				newbin := &domain.Bin{}
				newbin.Add(item)
				k.Bins = append(k.Bins, newbin)
			}
		}
	}

	if len(k.Bins) > 0 {
		sort.Slice(k.Bins, func(i, j int) bool { return k.Bins[i].Size > k.Bins[j].Size })
		bin = k.Bins[0]
	}

	return bin
}

func (k *Knapsack) fitBlocks() (bin *domain.Bin) {
	sort.Slice(k.list, func(i, j int) bool { return k.list[i].BlocksUsed > k.list[j].BlocksUsed })

	// how many blocks used by k.buffer bytes
	buffer := k.buffer / k.blockSize

	for _, item := range k.list {
		if item.BlocksUsed > (k.disk.BlocksFree - buffer) {
			k.over = append(k.over, item)
		} else {
			targetBin := -1
			remainingSpace := k.disk.BlocksFree

			for i, bin := range k.Bins {
				binSpaceUsed := bin.BlocksUsed
				binSpaceLeft := k.disk.BlocksFree - binSpaceUsed - item.BlocksUsed

				if binSpaceLeft < remainingSpace && binSpaceLeft >= buffer {
					remainingSpace = binSpaceLeft
					targetBin = i
				}
			}

			if targetBin >= 0 {
				k.Bins[targetBin].Add(item)
			} else {
				newbin := &domain.Bin{}
				newbin.Add(item)
				k.Bins = append(k.Bins, newbin)
			}
		}
	}

	if len(k.Bins) > 0 {
		sort.Slice(k.Bins, func(i, j int) bool { return k.Bins[i].BlocksUsed > k.Bins[j].BlocksUsed })
		bin = k.Bins[0]
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

// func (k *Knapsack) printList() {
// 	for _, item := range k.list {
// 		mlog.Info(fmt.Sprintf("Item (%s): %d", item.Name, item.Size))
// 	}
// }

// func (k *Knapsack) sortBins() {
// 	sort.Sort(model.ByFilled(k.Bins))
// }

// // Print -
// func (k *Knapsack) Print() {
// 	for i, bin := range k.Bins {
// 		mlog.Info("=========================================================")
// 		mlog.Info(fmt.Sprintf("%0d [%d/%d] %2.2f%% (%s)", i, bin.Size, k.disk.Free, (float64(bin.Size)/float64(k.disk.Free))*100, k.disk.Path))
// 		mlog.Info("---------------------------------------------------------")
// 		bin.Print()
// 		mlog.Info("---------------------------------------------------------")
// 		mlog.Info("")
// 	}
// }
