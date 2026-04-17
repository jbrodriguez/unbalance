package domain

// HardlinkGroup represents a set of paths that share the same inode
// (i.e., are hardlinked to the same underlying data)
type HardlinkGroup struct {
	Inode  uint64   `json:"inode"`  // Filesystem inode number
	Device uint64   `json:"device"` // Device ID (filesystem identifier)
	Size   uint64   `json:"size"`   // Size of the data (counted only once)
	Paths  []string `json:"paths"`  // All paths sharing this inode
	Items  []*Item  `json:"items"`  // All items sharing this inode
}

// HardlinkSummary provides statistics about hardlinks in a scan
type HardlinkSummary struct {
	TotalHardlinkedFiles  int    `json:"totalHardlinkedFiles"`  // Count of files with linkCount > 1
	TotalHardlinkGroups   int    `json:"totalHardlinkGroups"`   // Count of unique inode groups
	ApparentSize          uint64 `json:"apparentSize"`          // Sum of all file sizes (counting hardlinks multiple times)
	ActualSize            uint64 `json:"actualSize"`            // True disk usage (counting each inode once)
	PotentialSavings      uint64 `json:"potentialSavings"`      // ApparentSize - ActualSize
	HardlinkGroups        []*HardlinkGroup `json:"hardlinkGroups,omitempty"` // Detailed group info (optional)
}

// NewHardlinkGroup creates a new hardlink group from the first item
func NewHardlinkGroup(item *Item) *HardlinkGroup {
	return &HardlinkGroup{
		Inode:  item.Inode,
		Device: item.Device,
		Size:   item.Size,
		Paths:  []string{item.Path},
		Items:  []*Item{item},
	}
}

// AddItem adds an item to this hardlink group
func (hg *HardlinkGroup) AddItem(item *Item) {
	hg.Paths = append(hg.Paths, item.Path)
	hg.Items = append(hg.Items, item)
}

// PathCount returns the number of paths in this group
func (hg *HardlinkGroup) PathCount() int {
	return len(hg.Paths)
}

// HardlinkAnalyzer analyzes items for hardlink relationships
type HardlinkAnalyzer struct {
	// Map from inode key (device:inode) to hardlink group
	groups map[uint64]*HardlinkGroup
}

// NewHardlinkAnalyzer creates a new analyzer
func NewHardlinkAnalyzer() *HardlinkAnalyzer {
	return &HardlinkAnalyzer{
		groups: make(map[uint64]*HardlinkGroup),
	}
}

// AddItem adds an item to the analyzer, grouping by inode
func (ha *HardlinkAnalyzer) AddItem(item *Item) {
	// Skip items without inode info (shouldn't happen but be defensive)
	if item.Inode == 0 {
		return
	}

	key := item.InodeKey()
	if group, exists := ha.groups[key]; exists {
		group.AddItem(item)
	} else {
		ha.groups[key] = NewHardlinkGroup(item)
	}
}

// Analyze processes a list of items and returns a summary
func (ha *HardlinkAnalyzer) Analyze(items []*Item) *HardlinkSummary {
	// Reset state
	ha.groups = make(map[uint64]*HardlinkGroup)

	// Group all items by inode
	for _, item := range items {
		ha.AddItem(item)
	}

	return ha.GetSummary()
}

// GetSummary returns statistics about the analyzed hardlinks
func (ha *HardlinkAnalyzer) GetSummary() *HardlinkSummary {
	summary := &HardlinkSummary{
		HardlinkGroups: make([]*HardlinkGroup, 0),
	}

	for _, group := range ha.groups {
		// Every item contributes to apparent size
		for _, item := range group.Items {
			summary.ApparentSize += item.Size
		}

		// Actual size counts each unique inode only once
		summary.ActualSize += group.Size

		// Track groups with multiple paths (actual hardlinks)
		if group.PathCount() > 1 {
			summary.TotalHardlinkGroups++
			summary.TotalHardlinkedFiles += group.PathCount()
			summary.HardlinkGroups = append(summary.HardlinkGroups, group)
		}
	}

	summary.PotentialSavings = summary.ApparentSize - summary.ActualSize

	return summary
}

// GetGroups returns all hardlink groups (including single-file groups)
func (ha *HardlinkAnalyzer) GetGroups() map[uint64]*HardlinkGroup {
	return ha.groups
}

// GetHardlinkedGroups returns only groups with multiple paths
func (ha *HardlinkAnalyzer) GetHardlinkedGroups() []*HardlinkGroup {
	groups := make([]*HardlinkGroup, 0)
	for _, group := range ha.groups {
		if group.PathCount() > 1 {
			groups = append(groups, group)
		}
	}
	return groups
}

// GetHardlinkedItems returns only items that are hardlinked
func (ha *HardlinkAnalyzer) GetHardlinkedItems() []*Item {
	items := make([]*Item, 0)
	for _, group := range ha.groups {
		if group.PathCount() > 1 {
			items = append(items, group.Items...)
		}
	}
	return items
}

// GetNonHardlinkedItems returns items that are not hardlinked
func (ha *HardlinkAnalyzer) GetNonHardlinkedItems() []*Item {
	items := make([]*Item, 0)
	for _, group := range ha.groups {
		if group.PathCount() == 1 {
			items = append(items, group.Items...)
		}
	}
	return items
}

// CalculateActualSize returns the deduplicated size for a list of items
func CalculateActualSize(items []*Item) uint64 {
	analyzer := NewHardlinkAnalyzer()
	summary := analyzer.Analyze(items)
	return summary.ActualSize
}

// OrphanedHardlink represents a hardlink where some siblings are not selected for transfer.
// This occurs when a user selects some but not all paths sharing the same inode.
// Moving only selected paths will NOT free space on source (siblings still reference the inode).
type OrphanedHardlink struct {
	Inode           uint64   `json:"inode"`           // Filesystem inode number
	Device          uint64   `json:"device"`          // Device ID (filesystem identifier)
	Size            uint64   `json:"size"`            // Size of the data
	SelectedPaths   []string `json:"selectedPaths"`   // Paths selected for transfer
	UnselectedPaths []string `json:"unselectedPaths"` // Sibling paths NOT selected (will remain on source)
	TotalLinkCount  uint64   `json:"totalLinkCount"`  // Total hardlinks to this inode (from filesystem)
	SpaceImpact     uint64   `json:"spaceImpact"`     // Space that will NOT be freed on source
}

// OrphanSummary provides an overview of orphaned hardlinks in the plan
type OrphanSummary struct {
	TotalOrphanedGroups int                 `json:"totalOrphanedGroups"` // Count of inode groups with orphans
	TotalOrphanedFiles  int                 `json:"totalOrphanedFiles"`  // Count of selected files that are orphaned
	TotalSpaceImpact    uint64              `json:"totalSpaceImpact"`    // Total space that won't be freed
	OrphanedHardlinks   []*OrphanedHardlink `json:"orphanedHardlinks"`   // Detailed orphan info
}

// NewOrphanSummary creates a new empty OrphanSummary
func NewOrphanSummary() *OrphanSummary {
	return &OrphanSummary{
		OrphanedHardlinks: make([]*OrphanedHardlink, 0),
	}
}

// AddOrphan adds an orphaned hardlink to the summary
func (os *OrphanSummary) AddOrphan(orphan *OrphanedHardlink) {
	os.OrphanedHardlinks = append(os.OrphanedHardlinks, orphan)
	os.TotalOrphanedGroups++
	os.TotalOrphanedFiles += len(orphan.SelectedPaths)
	os.TotalSpaceImpact += orphan.SpaceImpact
}

// HasOrphans returns true if there are any orphaned hardlinks
func (os *OrphanSummary) HasOrphans() bool {
	return os.TotalOrphanedGroups > 0
}

// GetOrphanedInodeKeys returns a set of inode keys that are orphaned
func (os *OrphanSummary) GetOrphanedInodeKeys() map[uint64]bool {
	keys := make(map[uint64]bool)
	for _, orphan := range os.OrphanedHardlinks {
		// Create inode key same way as Item.InodeKey()
		key := (orphan.Device << 32) | (orphan.Inode & 0xFFFFFFFF)
		keys[key] = true
	}
	return keys
}
