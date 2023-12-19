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
				binSpaceUsed := bin.Size
				binSpaceLeft := k.disk.Free - binSpaceUsed - item.Size

				if binSpaceLeft < remainingSpace && binSpaceLeft >= k.buffer {
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

	// log.Printf("buffer %d\n", buffer)

	for _, item := range k.list {
		// log.Printf("item(%+v)\n", item)
		if item.BlocksUsed > (k.disk.BlocksFree - buffer) {
			// log.Println("if")
			k.over = append(k.over, item)
		} else {
			// log.Println("else")
			targetBin := -1
			remainingSpace := k.disk.BlocksFree

			// log.Printf("remspac %d\n", remainingSpace)

			for i, bin := range k.Bins {
				binSpaceUsed := bin.BlocksUsed
				binSpaceLeft := k.disk.BlocksFree - binSpaceUsed - item.BlocksUsed

				// log.Printf("bsu(%d)-bsl(%d)\n", binSpaceUsed, binSpaceLeft)
				// log.Printf("bsu(%d)-bsl(%d)-rs(%d)-buf(%d)\n", binSpaceUsed, binSpaceLeft, remainingSpace, buffer)

				if binSpaceLeft < remainingSpace && binSpaceLeft >= buffer {
					// log.Println("ifcabe")

					remainingSpace = binSpaceLeft
					targetBin = i

					// log.Printf("rs(%d)-tb(%d)\n", remainingSpace, i)
				}
			}

			if targetBin >= 0 {
				// log.Println("ifbin")
				k.Bins[targetBin].Add(item)
			} else {
				// log.Println("elsebin")
				newbin := &domain.Bin{}
				newbin.Add(item)
				// log.Printf("bin(%+v)", newbin)
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
