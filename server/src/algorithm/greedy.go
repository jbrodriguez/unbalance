package algorithm

import (
	"jbrodriguez/unbalance/server/src/domain"

	"github.com/jbrodriguez/mlog"
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
func NewGreedy(disk *domain.Disk, entries []*domain.Item, total, reserved int64) *Greedy {
	g := &Greedy{}
	g.disk = disk
	g.entries = entries
	g.buffer = reserved
	g.toFit = total
	return g
}

func (g *Greedy) FitAll() *domain.Bin {
	sizeToFit := g.toFit
	bin := &domain.Bin{}

	for _, entry := range g.entries {
		if entry.Location == g.disk.Path {
			// entry exists in this disk, subtract its size from the total
			// and don't add it to list of entries that will be transferred
			sizeToFit -= entry.Size
		} else {
			bin.Add(entry)
		}
	}

	mlog.Info("disk(%s)-sizeToFit+buffer(%d)-diskFree(%d)", g.disk.Path, sizeToFit+g.buffer, g.disk.Free)

	if sizeToFit+g.buffer > g.disk.Free {
		// g.disk.Bin = nil
		return nil
	}

	// g.disk.Bin = bin

	return bin
	// if g.disk.Free >= g.toFit + reserved
}
