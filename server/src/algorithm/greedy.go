package algorithm

import (
	"jbrodriguez/unbalance/server/src/domain"
)

// Greedy -
type Greedy struct {
	disk    *domain.Disk
	entries []*domain.Item
	toFit   int64
	buffer  int64

	Bins []*domain.Bin
}

// NewGreedy -
func NewGreedy(disk *domain.Disk, entries []*domain.Item, reserved int64) *Greedy {
	g := &Greedy{}
	g.disk = disk
	g.entries = entries
	g.buffer = reserved
	// g.toFit = total
	return g
}

// FitAll - Make sure that the disk has enough space to hold all given entries
func (g *Greedy) FitAll() *domain.Bin {
	// sizeToFit := g.toFit
	bin := &domain.Bin{}

	for _, entry := range g.entries {
		// entry doesnt exist in this disk, add it to the bin which also
		// accumulates the total size
		if entry.Location != g.disk.Path {
			bin.Add(entry)
		}

		// if entry.Location == g.disk.Path {
		// 	// entry exists in this disk, subtract its size from the total
		// 	// and don't add it to list of entries that will be transferred
		// 	sizeToFit -= entry.Size
		// } else {
		// 	bin.Add(entry)
		// }
	}

	if bin.Size+g.buffer > g.disk.Free {
		return nil
	}

	return bin
}
