package domain

// Bin -
type Bin struct {
	Size       uint64  `json:"size"`
	Items      []*Item `json:"items"`
	BlocksUsed uint64  `json:"blocksUsed"`

	// Hardlink-aware fields for deduplicated size calculations
	ActualSize       uint64          `json:"actualSize"`       // Deduplicated size (each inode counted once)
	ActualBlocksUsed uint64          `json:"actualBlocksUsed"` // Deduplicated blocks
	CountedInodes    map[uint64]bool `json:"-"`                // Track which inodes have been counted
}

// NewBin creates a new Bin with initialized tracking maps
func NewBin() *Bin {
	return &Bin{
		Items:         make([]*Item, 0),
		CountedInodes: make(map[uint64]bool),
	}
}

// Add adds an item to the bin (apparent size only, for backward compatibility)
func (b *Bin) Add(item *Item) {
	b.Items = append(b.Items, item)
	b.Size += item.Size
	b.BlocksUsed += item.BlocksUsed
}

// AddWithHardlinkTracking adds an item, tracking both apparent and actual (deduplicated) sizes
func (b *Bin) AddWithHardlinkTracking(item *Item) {
	b.Items = append(b.Items, item)
	b.Size += item.Size
	b.BlocksUsed += item.BlocksUsed

	// Initialize map if needed (for bins not created with NewBin)
	if b.CountedInodes == nil {
		b.CountedInodes = make(map[uint64]bool)
	}

	// Only count actual size once per unique inode
	inodeKey := item.InodeKey()
	if !b.CountedInodes[inodeKey] {
		b.ActualSize += item.Size
		b.ActualBlocksUsed += item.BlocksUsed
		b.CountedInodes[inodeKey] = true
	}
}

// IsInodeCounted returns true if the given inode has already been counted in this bin
func (b *Bin) IsInodeCounted(inodeKey uint64) bool {
	if b.CountedInodes == nil {
		return false
	}
	return b.CountedInodes[inodeKey]
}

// // Print -
// func (b *Bin) Print() {
// 	for _, item := range b.Items {
// 		mlog.Info("[%s] %s\n", lib.ByteSize(item.Size), item.Name)
// 	}
// }

// // ByFilled -
// type ByFilled []*Bin

// func (s ByFilled) Len() int           { return len(s) }
// func (s ByFilled) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
// func (s ByFilled) Less(i, j int) bool { return s[i].Size > s[j].Size }
