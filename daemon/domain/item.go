package domain

// Item -
type Item struct {
	Name       string `json:"name"`
	Size       uint64 `json:"size"`
	Path       string `json:"path"`
	Location   string `json:"location"`
	BlocksUsed uint64 `json:"blocksUsed"`

	// Hardlink detection fields
	Inode     uint64 `json:"inode"`     // Filesystem inode number
	Device    uint64 `json:"device"`    // Device ID (filesystem identifier)
	LinkCount uint64 `json:"linkCount"` // Number of hardlinks to this inode
}

// IsHardlinked returns true if this item has multiple hardlinks
func (i *Item) IsHardlinked() bool {
	return i.LinkCount > 1
}

// InodeKey returns a unique key for identifying hardlink groups
// combining device and inode ensures uniqueness across filesystems
func (i *Item) InodeKey() uint64 {
	// Use device ID in upper 32 bits and inode in lower 32 bits
	// This works for typical inode values; for very large inodes
	// a string key would be safer
	return (i.Device << 32) | (i.Inode & 0xFFFFFFFF)
}

// // BySize -
// type BySize []*Item

// func (s BySize) Len() int           { return len(s) }
// func (s BySize) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
// func (s BySize) Less(i, j int) bool { return s[i].Size > s[j].Size }
