package algorithm

import (
	"jbrodriguez/unbalance/server/src/domain"
)

// Greedy -
type Greedy struct {
	disk      *domain.Disk
	entries   []*domain.Item
	buffer    int64
	blockSize int64

	Bins []*domain.Bin
}

// NewGreedy -
func NewGreedy(disk *domain.Disk, entries []*domain.Item, reserved, blockSize int64) *Greedy {
	g := &Greedy{}

	g.disk = disk
	g.entries = entries
	g.buffer = reserved
	g.blockSize = blockSize

	return g
}

// FitAll - Make sure that the disk has enough space to hold all given entries
func (g *Greedy) FitAll() *domain.Bin {
	if g.blockSize > 0 {
		return g.fitBlocks()
	}

	return g.fitBytes()
}

func (g *Greedy) fitBytes() *domain.Bin {
	bin := &domain.Bin{}

	for _, entry := range g.entries {
		// entry doesnt exist in this disk, add it to the bin which also
		// accumulates the total size
		if entry.Location != g.disk.Path {
			bin.Add(entry)
		}
	}

	if bin.Size+g.buffer > g.disk.Free {
		return nil
	}

	return bin
}

func (g *Greedy) fitBlocks() *domain.Bin {
	bin := &domain.Bin{}

	for _, entry := range g.entries {
		// entry doesnt exist in this disk, add it to the bin which also
		// accumulates the total size
		if entry.Location != g.disk.Path {
			bin.Add(entry)
		}
	}

	// how many blocks are used in g.buffer bytes
	buffer := g.buffer / g.blockSize

	if bin.BlocksUsed+buffer > g.disk.BlocksFree {
		return nil
	}

	return bin
}
