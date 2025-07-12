package core

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"

	"unbalance/daemon/domain"
	"unbalance/daemon/lib"
	"unbalance/daemon/logger"
)

// detectHardlinks analyzes a file/directory for hardlinks and returns hardlink metadata
func detectHardlinks(itemPath string) (*domain.HardlinkInfo, error) {
	fi, err := os.Stat(itemPath)
	if err != nil {
		return nil, err
	}

	// Get file system stat
	stat, ok := fi.Sys().(*syscall.Stat_t)
	if !ok {
		return nil, fmt.Errorf("unable to get system stat for %s", itemPath)
	}

	// Check if it's a directory or file with multiple hardlinks
	hardlinkInfo := &domain.HardlinkInfo{
		Inode:        stat.Ino,
		InodeCount:   int(stat.Nlink),
		HasHardlinks: stat.Nlink > 1,
		SharedSize:   uint64(fi.Size()),
		ExpandedSize: uint64(fi.Size()),
	}

	// For files with hardlinks, calculate expanded size
	if hardlinkInfo.HasHardlinks && !fi.IsDir() {
		// If preserving hardlinks is disabled, each hardlink will become a separate file
		hardlinkInfo.ExpandedSize = uint64(fi.Size()) * uint64(stat.Nlink)
	}

	// For directories, we need to scan for hardlinked files within
	if fi.IsDir() {
		dirHardlinkInfo, err := scanDirectoryForHardlinks(itemPath)
		if err != nil {
			logger.Yellow("hardlink detection warning for %s: %s", itemPath, err)
			// Continue with file-level info even if directory scan fails
		} else {
			// Use directory scan results which are more comprehensive
			hardlinkInfo = dirHardlinkInfo
			hardlinkInfo.Inode = stat.Ino // Keep directory inode for tracking
		}
	}

	return hardlinkInfo, nil
}

// scanDirectoryForHardlinks scans a directory to find all hardlinked files within it
func scanDirectoryForHardlinks(dirPath string) (*domain.HardlinkInfo, error) {
	var totalSharedSize, totalExpandedSize uint64
	var hasAnyHardlinks bool
	inodeMap := make(map[uint64]int) // Track inodes and their link counts

	// Use cross-platform approach with filepath.Walk
	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // Continue on errors
		}
		
		// Only process regular files
		if !info.Mode().IsRegular() {
			return nil
		}
		
		// Get system stat information
		stat, ok := info.Sys().(*syscall.Stat_t)
		if !ok {
			return nil // Continue if we can't get stat info
		}
		
		inode := stat.Ino
		size := uint64(info.Size())
		
		// Track unique inodes
		if _, exists := inodeMap[inode]; !exists {
			inodeMap[inode] = int(stat.Nlink)
			
			// Add to shared size (count each unique inode once)
			totalSharedSize += size
			
			// Add to expanded size (count each hardlink separately if they would be duplicated)
			if stat.Nlink > 1 {
				hasAnyHardlinks = true
				totalExpandedSize += size * uint64(stat.Nlink)
			} else {
				totalExpandedSize += size
			}
		}
		
		return nil
	})

	if err != nil {
		return nil, err
	}

	// Calculate total hardlink count
	totalHardlinks := 0
	for _, count := range inodeMap {
		if count > 1 {
			totalHardlinks += count
		}
	}

	return &domain.HardlinkInfo{
		InodeCount:   totalHardlinks,
		SharedSize:   totalSharedSize,
		ExpandedSize: totalExpandedSize,
		HasHardlinks: hasAnyHardlinks,
		Inode:        0, // Not applicable for directory aggregates
	}, nil
}

// getEffectiveSize returns the size to use for space calculations based on hardlink configuration
func getEffectiveSize(item *domain.Item, preserveHardlinks bool) uint64 {
	if item.HardlinkInfo == nil || !item.HardlinkInfo.HasHardlinks {
		return item.Size
	}

	if preserveHardlinks {
		// If preserving hardlinks, use shared size for planning
		return item.HardlinkInfo.SharedSize
	} else {
		// If not preserving hardlinks, use expanded size (worst case)
		return item.HardlinkInfo.ExpandedSize
	}
}

// updateItemWithHardlinkInfo updates an item's size field based on hardlink detection
func updateItemWithHardlinkInfo(item *domain.Item, preserveHardlinks bool) {
	if item.HardlinkInfo != nil && item.HardlinkInfo.HasHardlinks {
		// Update the main Size field to reflect what will actually be transferred
		item.Size = getEffectiveSize(item, preserveHardlinks)
		
		if item.HardlinkInfo.HasHardlinks {
			logger.Blue("hardlink detected: %s (links: %d, shared: %s, expanded: %s, effective: %s)", 
				item.Name, 
				item.HardlinkInfo.InodeCount, 
				lib.ByteSize(item.HardlinkInfo.SharedSize), 
				lib.ByteSize(item.HardlinkInfo.ExpandedSize),
				lib.ByteSize(item.Size))
		}
	}
}