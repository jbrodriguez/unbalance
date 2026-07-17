package core

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"unbalance/daemon/domain"
	"unbalance/daemon/logger"
)

// Size reports, for each disk where path exists, the storage the entry
// occupies on that disk, plus the total across all disks.
func (c *Core) Size(path string) domain.Sizes {
	sizes := domain.Sizes{Disks: make(map[string]uint64)}

	for _, disk := range c.state.Unraid.Disks {
		name := strings.Replace(path, "/mnt/user", "", 1)
		entry := filepath.Join(disk.Path, name)

		fi, err := os.Lstat(entry)
		if err != nil {
			continue
		}

		size, err := entrySize(entry, fi)
		if err != nil {
			logger.Yellow("unable to size %s: %s", entry, err)
			continue
		}

		sizes.Disks[disk.Name] = size
		sizes.Total += size
	}

	return sizes
}

// entrySize returns the apparent size in bytes of the file or directory tree
// rooted at entry. Symlinks are not followed.
func entrySize(entry string, fi os.FileInfo) (uint64, error) {
	if fi.Mode()&os.ModeSymlink != 0 {
		return 0, nil
	}

	if !fi.IsDir() {
		return uint64(fi.Size()), nil
	}

	var total uint64
	err := filepath.WalkDir(entry, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			if path == entry {
				return err
			}
			// keep sizing the rest of the tree past unreadable entries
			return nil
		}

		if d.Type().IsRegular() {
			info, ierr := d.Info()
			if ierr != nil {
				return nil
			}
			total += uint64(info.Size())
		}

		return nil
	})

	return total, err
}
