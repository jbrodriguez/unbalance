package domain

// HardlinkInfo contains metadata about hardlinks for an item
type HardlinkInfo struct {
	InodeCount     int    `json:"inodeCount"`     // Number of hardlinks to this inode
	SharedSize     uint64 `json:"sharedSize"`     // Size reported by du -bs (shared across hardlinks)
	ExpandedSize   uint64 `json:"expandedSize"`   // Size if all hardlinks become separate files
	HasHardlinks   bool   `json:"hasHardlinks"`   // True if this item has multiple hardlinks
	Inode          uint64 `json:"inode"`          // Inode number for hardlink tracking
}

// Item -
type Item struct {
	Name         string        `json:"name"`
	Size         uint64        `json:"size"`         // For compatibility, will use ExpandedSize when hardlinks present
	Path         string        `json:"path"`
	Location     string        `json:"location"`
	BlocksUsed   uint64        `json:"blocksUsed"`
	HardlinkInfo *HardlinkInfo `json:"hardlinkInfo,omitempty"` // Present when hardlinks are detected
}

// // BySize -
// type BySize []*Item

// func (s BySize) Len() int           { return len(s) }
// func (s BySize) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
// func (s BySize) Less(i, j int) bool { return s[i].Size > s[j].Size }
