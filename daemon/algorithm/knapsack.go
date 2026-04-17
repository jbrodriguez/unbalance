package algorithm

import (
	"sort"

	"unbalance/daemon/domain"
)

// Knapsack -
type Knapsack struct {
	disk *domain.Disk

	Bins []*domain.Bin
	list []*domain.Item
	over []*domain.Item

	buffer    uint64
	blockSize uint64
}

// NewKnapsack -
func NewKnapsack(disk *domain.Disk, items []*domain.Item, reserved, blockSize uint64) *Knapsack {
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

	for _, item := range k.list {
		if item.Size > (k.disk.Free - k.buffer) {
			k.over = append(k.over, item)
		} else {
			targetBin := -1
			remainingSpace := k.disk.Free

			for i, bin := range k.Bins {
				// Use ActualSize for space calculations (deduplicated)
				binSpaceUsed := bin.ActualSize

				// Check if this item's inode is already counted in this bin
				itemActualSize := item.Size
				if bin.IsInodeCounted(item.InodeKey()) {
					itemActualSize = 0 // Already counted, no additional space needed
				}

				binSpaceLeft := k.disk.Free - binSpaceUsed - itemActualSize

				if binSpaceLeft < remainingSpace && binSpaceLeft >= k.buffer {
					remainingSpace = binSpaceLeft
					targetBin = i
				}
			}

			if targetBin >= 0 {
				k.Bins[targetBin].AddWithHardlinkTracking(item)
			} else {
				newbin := domain.NewBin()
				newbin.AddWithHardlinkTracking(item)
				k.Bins = append(k.Bins, newbin)
			}
		}
	}

	if len(k.Bins) > 0 {
		// Sort by ActualSize for consistent selection
		sort.Slice(k.Bins, func(i, j int) bool { return k.Bins[i].ActualSize > k.Bins[j].ActualSize })
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
				// Use ActualBlocksUsed for space calculations (deduplicated)
				binSpaceUsed := bin.ActualBlocksUsed

				// Check if this item's inode is already counted in this bin
				itemActualBlocks := item.BlocksUsed
				if bin.IsInodeCounted(item.InodeKey()) {
					itemActualBlocks = 0 // Already counted, no additional space needed
				}

				binSpaceLeft := k.disk.BlocksFree - binSpaceUsed - itemActualBlocks

				if binSpaceLeft < remainingSpace && binSpaceLeft >= buffer {
					remainingSpace = binSpaceLeft
					targetBin = i
				}
			}

			if targetBin >= 0 {
				k.Bins[targetBin].AddWithHardlinkTracking(item)
			} else {
				newbin := domain.NewBin()
				newbin.AddWithHardlinkTracking(item)
				k.Bins = append(k.Bins, newbin)
			}
		}
	}

	if len(k.Bins) > 0 {
		// Sort by ActualBlocksUsed for consistent selection
		sort.Slice(k.Bins, func(i, j int) bool { return k.Bins[i].ActualBlocksUsed > k.Bins[j].ActualBlocksUsed })
		bin = k.Bins[0]
	}

	return bin
}
