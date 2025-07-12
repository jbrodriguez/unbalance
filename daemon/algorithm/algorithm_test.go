package algorithm

import (
	"fmt"
	"testing"

	"unbalance/daemon/domain"
)

// createTestDisk creates a disk for testing purposes
func createTestDisk(name string, free, size uint64) *domain.Disk {
	return &domain.Disk{
		Name:       name,
		Path:       "/mnt/" + name,
		Free:       free,
		Size:       size,
		BlocksFree: free / 4096,  // Assume 4K blocks
		BlocksTotal: size / 4096,
	}
}

// createTestItem creates an item for testing
func createTestItem(name string, size uint64) *domain.Item {
	return &domain.Item{
		Name:       name,
		Size:       size,
		Path:       name,
		BlocksUsed: size / 4096, // Assume 4K blocks
	}
}

// createHardlinkedItem creates an item with hardlink information
func createHardlinkedItem(name string, sharedSize, expandedSize uint64, linkCount int) *domain.Item {
	return &domain.Item{
		Name:       name,
		Size:       sharedSize, // Will be updated based on hardlink policy
		Path:       name,
		BlocksUsed: sharedSize / 4096,
		HardlinkInfo: &domain.HardlinkInfo{
			HasHardlinks: linkCount > 1,
			SharedSize:   sharedSize,
			ExpandedSize: expandedSize,
			InodeCount:   linkCount,
			Inode:        12345, // Mock inode
		},
	}
}

func TestKnapsack_RegularFiles(t *testing.T) {
	// Create test disk with 10GB free space
	disk := createTestDisk("disk1", 10*1024*1024*1024, 20*1024*1024*1024)
	
	// Create test items
	items := []*domain.Item{
		createTestItem("file1.txt", 2*1024*1024*1024), // 2GB
		createTestItem("file2.txt", 3*1024*1024*1024), // 3GB
		createTestItem("file3.txt", 1*1024*1024*1024), // 1GB
		createTestItem("file4.txt", 4*1024*1024*1024), // 4GB
		createTestItem("file5.txt", 1*1024*1024*1024), // 1GB
	}
	
	reserved := uint64(1024 * 1024 * 1024) // 1GB reserved
	blockSize := uint64(4096)
	
	knapsack := NewKnapsack(disk, items, reserved, blockSize)
	bin := knapsack.BestFit()
	
	if bin == nil {
		t.Fatal("Expected bin to be created")
	}
	
	// Should fit items that don't exceed available space (10GB - 1GB reserved = 9GB available)
	if bin.Size > 9*1024*1024*1024 {
		t.Errorf("Bin size (%d) exceeds available space", bin.Size)
	}
	
	if len(bin.Items) == 0 {
		t.Error("Expected at least one item in bin")
	}
}

func TestKnapsack_HardlinkedFiles_Preserved(t *testing.T) {
	// Create test disk
	disk := createTestDisk("disk1", 10*1024*1024*1024, 20*1024*1024*1024)
	
	// Create items with hardlinks - when preserved, should use shared size
	items := []*domain.Item{
		createHardlinkedItem("hardlinked1.txt", 2*1024*1024*1024, 6*1024*1024*1024, 3), // 2GB shared, 6GB expanded
		createHardlinkedItem("hardlinked2.txt", 1*1024*1024*1024, 4*1024*1024*1024, 4), // 1GB shared, 4GB expanded  
		createTestItem("regular.txt", 3*1024*1024*1024), // 3GB regular file
	}
	
	// Simulate hardlinks being preserved (items should use shared size)
	for _, item := range items {
		if item.HardlinkInfo != nil && item.HardlinkInfo.HasHardlinks {
			item.Size = item.HardlinkInfo.SharedSize
			item.BlocksUsed = item.Size / 4096
		}
	}
	
	reserved := uint64(512 * 1024 * 1024) // 512MB reserved
	blockSize := uint64(4096)
	
	knapsack := NewKnapsack(disk, items, reserved, blockSize)
	bin := knapsack.BestFit()
	
	if bin == nil {
		t.Fatal("Expected bin to be created")
	}
	
	// With hardlinks preserved, total size should be: 2GB + 1GB + 3GB = 6GB
	// Available space: 10GB - 0.5GB = 9.5GB, so all should fit
	expectedSize := uint64(6 * 1024 * 1024 * 1024)
	if bin.Size != expectedSize {
		t.Errorf("Expected bin size %d, got %d", expectedSize, bin.Size)
	}
	
	if len(bin.Items) != 3 {
		t.Errorf("Expected 3 items in bin, got %d", len(bin.Items))
	}
}

func TestKnapsack_HardlinkedFiles_NotPreserved(t *testing.T) {
	// Create test disk
	disk := createTestDisk("disk1", 10*1024*1024*1024, 20*1024*1024*1024)
	
	// Create items with hardlinks - when not preserved, should use expanded size
	items := []*domain.Item{
		createHardlinkedItem("hardlinked1.txt", 2*1024*1024*1024, 6*1024*1024*1024, 3), // 2GB shared, 6GB expanded
		createHardlinkedItem("hardlinked2.txt", 1*1024*1024*1024, 4*1024*1024*1024, 4), // 1GB shared, 4GB expanded
		createTestItem("regular.txt", 2*1024*1024*1024), // 2GB regular file
	}
	
	// Simulate hardlinks NOT being preserved (items should use expanded size)
	for _, item := range items {
		if item.HardlinkInfo != nil && item.HardlinkInfo.HasHardlinks {
			item.Size = item.HardlinkInfo.ExpandedSize
			item.BlocksUsed = item.Size / 4096
		}
	}
	
	reserved := uint64(512 * 1024 * 1024) // 512MB reserved
	blockSize := uint64(4096)
	
	knapsack := NewKnapsack(disk, items, reserved, blockSize)
	bin := knapsack.BestFit()
	
	if bin == nil {
		t.Fatal("Expected bin to be created")
	}
	
	// With hardlinks expanded, total would be: 6GB + 4GB + 2GB = 12GB
	// Available space: 10GB - 0.5GB = 9.5GB
	// Algorithm should pick items that fit, likely excluding the largest
	if bin.Size > 9*1024*1024*1024+512*1024*1024 {
		t.Errorf("Bin size (%d) exceeds available space", bin.Size)
	}
	
	// Should have fewer items since they're larger when expanded
	if len(bin.Items) == 0 {
		t.Error("Expected at least one item in bin")
	}
}

func TestKnapsack_MixedHardlinkScenarios(t *testing.T) {
	// Create test disk
	disk := createTestDisk("disk1", 15*1024*1024*1024, 30*1024*1024*1024)
	
	// Mix of regular files and hardlinked files
	items := []*domain.Item{
		createTestItem("regular1.txt", 2*1024*1024*1024),                                  // 2GB regular
		createHardlinkedItem("hardlinked1.txt", 1*1024*1024*1024, 5*1024*1024*1024, 5),   // 1GB->5GB
		createTestItem("regular2.txt", 3*1024*1024*1024),                                  // 3GB regular
		createHardlinkedItem("hardlinked2.txt", 2*1024*1024*1024, 4*1024*1024*1024, 2),   // 2GB->4GB
		createTestItem("regular3.txt", 1*1024*1024*1024),                                  // 1GB regular
		createHardlinkedItem("hardlinked3.txt", 500*1024*1024, 1500*1024*1024, 3),        // 0.5GB->1.5GB
	}
	
	reserved := uint64(1024 * 1024 * 1024) // 1GB reserved
	blockSize := uint64(4096)
	
	// Test with hardlinks preserved
	preservedItems := make([]*domain.Item, len(items))
	for i, item := range items {
		preservedItems[i] = &domain.Item{
			Name:         item.Name,
			Size:         item.Size,
			Path:         item.Path,
			BlocksUsed:   item.BlocksUsed,
			HardlinkInfo: item.HardlinkInfo,
		}
		if preservedItems[i].HardlinkInfo != nil && preservedItems[i].HardlinkInfo.HasHardlinks {
			preservedItems[i].Size = preservedItems[i].HardlinkInfo.SharedSize
			preservedItems[i].BlocksUsed = preservedItems[i].Size / 4096
		}
	}
	
	knapsackPreserved := NewKnapsack(disk, preservedItems, reserved, blockSize)
	binPreserved := knapsackPreserved.BestFit()
	
	// Test with hardlinks expanded
	expandedItems := make([]*domain.Item, len(items))
	for i, item := range items {
		expandedItems[i] = &domain.Item{
			Name:         item.Name,
			Size:         item.Size,
			Path:         item.Path,
			BlocksUsed:   item.BlocksUsed,
			HardlinkInfo: item.HardlinkInfo,
		}
		if expandedItems[i].HardlinkInfo != nil && expandedItems[i].HardlinkInfo.HasHardlinks {
			expandedItems[i].Size = expandedItems[i].HardlinkInfo.ExpandedSize
			expandedItems[i].BlocksUsed = expandedItems[i].Size / 4096
		}
	}
	
	knapsackExpanded := NewKnapsack(disk, expandedItems, reserved, blockSize)
	binExpanded := knapsackExpanded.BestFit()
	
	if binPreserved == nil || binExpanded == nil {
		t.Fatal("Expected both bins to be created")
	}
	
	// Both algorithms should work and produce valid results
	// The exact sizes may vary due to algorithm choices, but both should be valid
	t.Logf("Preserved hardlinks bin size: %d, Expanded hardlinks bin size: %d", binPreserved.Size, binExpanded.Size)
	
	// Both should respect the available space limit
	availableSpace := disk.Free - reserved
	if binPreserved.Size > availableSpace {
		t.Errorf("Preserved bin size (%d) exceeds available space (%d)", binPreserved.Size, availableSpace)
	}
	if binExpanded.Size > availableSpace {
		t.Errorf("Expanded bin size (%d) exceeds available space (%d)", binExpanded.Size, availableSpace)
	}
}

func TestGreedy_HardlinkedFiles(t *testing.T) {
	// Create test disk
	disk := createTestDisk("disk1", 8*1024*1024*1024, 10*1024*1024*1024) // 8GB free
	
	// Create items with hardlinks
	items := []*domain.Item{
		createHardlinkedItem("hardlinked1.txt", 2*1024*1024*1024, 6*1024*1024*1024, 3), // 2GB->6GB
		createHardlinkedItem("hardlinked2.txt", 1*1024*1024*1024, 4*1024*1024*1024, 4), // 1GB->4GB
		createTestItem("regular.txt", 3*1024*1024*1024), // 3GB regular
	}
	
	reserved := uint64(512 * 1024 * 1024) // 512MB reserved
	blockSize := uint64(4096)
	
	// Test with hardlinks preserved (use shared sizes)
	preservedItems := make([]*domain.Item, len(items))
	for i, item := range items {
		preservedItems[i] = &domain.Item{
			Name:         item.Name,
			Size:         item.Size,
			Path:         item.Path,
			Location:     "/mnt/source", // Set location for greedy algorithm
			BlocksUsed:   item.BlocksUsed,
			HardlinkInfo: item.HardlinkInfo,
		}
		if preservedItems[i].HardlinkInfo != nil && preservedItems[i].HardlinkInfo.HasHardlinks {
			preservedItems[i].Size = preservedItems[i].HardlinkInfo.SharedSize
			preservedItems[i].BlocksUsed = preservedItems[i].Size / 4096
		}
	}
	
	greedy := NewGreedy(disk, preservedItems, reserved, blockSize)
	bin := greedy.FitAll()
	
	if bin == nil {
		t.Fatal("Expected bin to be created")
	}
	
	// With preserved hardlinks: 2GB + 1GB + 3GB = 6GB total should fit in 8GB-0.5GB = 7.5GB available
	expectedSize := uint64(6 * 1024 * 1024 * 1024)
	if bin.Size != expectedSize {
		t.Errorf("Expected bin size %d, got %d", expectedSize, bin.Size)
	}
	
	if len(bin.Items) != 3 {
		t.Errorf("Expected 3 items in bin, got %d", len(bin.Items))
	}
}

func TestGreedy_ConservativeHardlinkPlanning(t *testing.T) {
	// Create test disk with limited space
	disk := createTestDisk("disk1", 6*1024*1024*1024, 10*1024*1024*1024) // 6GB free
	
	// Create items where expanded size matters
	items := []*domain.Item{
		createHardlinkedItem("hardlinked1.txt", 2*1024*1024*1024, 8*1024*1024*1024, 4), // 2GB->8GB
	}
	
	reserved := uint64(512 * 1024 * 1024) // 512MB reserved
	blockSize := uint64(4096)
	
	// Test with conservative planning (hardlinks not preserved)
	expandedItems := make([]*domain.Item, len(items))
	for i, item := range items {
		expandedItems[i] = &domain.Item{
			Name:         item.Name,
			Size:         item.HardlinkInfo.ExpandedSize, // Use expanded size (8GB)
			Path:         item.Path,
			Location:     "/mnt/source", // Set location for greedy algorithm
			BlocksUsed:   item.HardlinkInfo.ExpandedSize / 4096,
			HardlinkInfo: item.HardlinkInfo,
		}
	}
	
	greedy := NewGreedy(disk, expandedItems, reserved, blockSize)
	bin := greedy.FitAll()
	
	// With conservative planning, item needs 8GB but only 5.5GB available (6GB - 0.5GB reserved)
	// So it should not fit
	if bin != nil {
		t.Error("Conservative planning should reject oversized items")
	}
}

// Benchmark tests for algorithm performance with hardlinks
func BenchmarkKnapsack_ManyHardlinks(b *testing.B) {
	disk := createTestDisk("disk1", 50*1024*1024*1024, 100*1024*1024*1024)
	
	// Create many items with hardlinks
	items := make([]*domain.Item, 1000)
	for i := 0; i < 1000; i++ {
		if i%3 == 0 {
			// Every third item has hardlinks
			items[i] = createHardlinkedItem(
				fmt.Sprintf("hardlinked_%d.txt", i),
				50*1024*1024,     // 50MB shared
				200*1024*1024,    // 200MB expanded
				4,                // 4 hardlinks
			)
			items[i].Size = items[i].HardlinkInfo.SharedSize
			items[i].BlocksUsed = items[i].Size / 4096
		} else {
			items[i] = createTestItem(fmt.Sprintf("regular_%d.txt", i), 100*1024*1024) // 100MB
		}
	}
	
	reserved := uint64(1024 * 1024 * 1024) // 1GB
	blockSize := uint64(4096)
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		knapsack := NewKnapsack(disk, items, reserved, blockSize)
		knapsack.BestFit()
	}
}