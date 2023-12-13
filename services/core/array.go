package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"unbalance/domain"
	"unbalance/lib"
	"unbalance/logger"

	"gopkg.in/ini.v1"
)

func (c *Core) sanityCheck() error {
	locations := []string{"/var/local/emhttp"}

	location := lib.SearchFile("var.ini", locations)
	if location == "" {
		return fmt.Errorf("unable to find var.ini (%s)", strings.Join(locations, ", "))
	}

	location = lib.SearchFile("disks.ini", locations)
	if location == "" {
		return fmt.Errorf("unable to find var.ini (%s)", strings.Join(locations, ", "))
	}

	return nil
}

func (c *Core) getCertificate() string {
	// get ssl settings
	ident, err := ini.Load("/boot/config/ident.cfg")
	if err != nil {
		return ""
	}

	usessl := ident.Section("").Key("USE_SSL").String()

	// get array status
	file, err := ini.Load("/var/local/emhttp/var.ini")
	if err != nil {
		return ""
	}

	name := file.Section("").Key("NAME").String()

	cert := getCertificateName(certDir, name)
	logger.LightBlue("cert: %s", cert)

	secure := cert != "" && !(usessl == "" || usessl == "no")

	if secure {
		return cert
	}

	return ""
}

func (c *Core) getStatus() (*domain.Unraid, error) {
	return getArrayData()
}

func getArrayData() (*domain.Unraid, error) {
	unraid := &domain.Unraid{}

	// get array status
	file, err := ini.Load("/var/local/emhttp/var.ini")
	if err != nil {
		return nil, fmt.Errorf("unable to load var.ini: %w", err)
	}

	numDisks := file.Section("").Key("mdNumDisks").String()
	unraid.NumDisks, _ = strconv.ParseUint(numDisks, 10, 64)

	numProtected := file.Section("").Key("mdNumProtected").String()
	unraid.NumProtected, _ = strconv.ParseUint(numProtected, 10, 64)

	synced := file.Section("").Key("sbSynced").String()
	ut, _ := strconv.ParseUint(synced, 10, 64)
	unraid.Synced = time.Unix(int64(ut), 0)

	syncErrs := file.Section("").Key("sbSyncErrs").String()
	unraid.SyncErrs, _ = strconv.ParseUint(syncErrs, 10, 64)

	resync := file.Section("").Key("mdResync").String()
	unraid.Resync, _ = strconv.ParseUint(resync, 10, 64)

	resyncPos := file.Section("").Key("mdResyncPos").String()
	unraid.ResyncPos, _ = strconv.ParseUint(resyncPos, 10, 64)

	mdState := file.Section("").Key("mdState").String()
	unraid.State = mdState

	// get disks
	file, err = ini.Load("/var/local/emhttp/disks.ini")
	if err != nil {
		return nil, fmt.Errorf("unable to load disks.ini: %w", err)
	}

	// get free/size data
	free := make(map[string]uint64)
	size := make(map[string]uint64)

	// err = lib.Shell("df --block-size=1 /mnt/*", mlog.Warning, "Refresh error:", "", func(line string) {
	// 	data := strings.Fields(line)
	// 	size[data[5]], _ = strconv.ParseUint(data[1], 10, 64)
	// 	free[data[5]], _ = strconv.ParseUint(data[3], 10, 64)
	// })

	err = lib.Shell("df --block-size=1 /mnt/*", "", func(line string) {
		data := strings.Fields(line)
		size[data[5]], _ = strconv.ParseUint(data[1], 10, 64)
		free[data[5]], _ = strconv.ParseUint(data[3], 10, 64)
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get free/size data: %w", err)
	}

	var blockSize int64
	var totalSize, totalFree uint64
	var hasBlockSize bool
	disks := make([]*domain.Disk, 0)

	for _, section := range file.Sections() {
		// DEFAULT section
		if section.Key("name").String() == "" {
			continue
		}

		diskType := section.Key("type").String()
		diskName := section.Key("name").String()
		diskStatus := section.Key("status").String()

		if diskType == "Parity" || diskType == "Flash" || (diskType == "Cache" && len(diskName) > 5 || diskStatus == "DISK_NP") {
			continue
		}

		disk := &domain.Disk{}

		disk.ID, _ = strconv.ParseUint(section.Key("idx").String(), 10, 64) // 1
		disk.Name = diskName                                                // disk1, cache
		disk.Path = "/mnt/" + disk.Name                                     // /mnt/disk1, /mnt/cache
		disk.Device = section.Key("device").String()                        // sdp
		disk.Type = diskType                                                // Flash, Parity, Data, Cache
		disk.FsType = section.Key("fsType").String()                        // xfs, reiserfs, btrfs
		disk.Free = free[disk.Path]
		disk.Size = size[disk.Path]
		disk.Serial = section.Key("id").String() // WDC_WD30EZRX-00DC0B0_WD-WMC9T204468
		disk.Status = diskStatus                 // DISK_OK

		var stat syscall.Statfs_t
		e := syscall.Statfs(disk.Path, &stat)
		if e == nil {
			disk.BlocksTotal = stat.Blocks
			disk.BlocksFree = stat.Bavail

			//
			if blockSize != stat.Bsize {
				if !hasBlockSize {
					blockSize = stat.Bsize
				} else {
					blockSize = 0
				}

				hasBlockSize = true
			}
		}

		totalSize += disk.Size
		totalFree += disk.Free

		disks = append(disks, disk)
	}

	unraid.Size = totalSize
	unraid.Free = totalFree
	unraid.BlockSize = blockSize

	sort.Slice(disks, func(i, j int) bool { return disks[i].ID < disks[j].ID })

	unraid.Disks = disks

	return unraid, nil
}

func (c *Core) getTree(path string) *domain.Entry {

	entry := &domain.Entry{Path: path}

	items := make([]domain.Node, 0)

	elements, _ := os.ReadDir(path)
	for _, element := range elements {
		var node domain.Node

		// default values
		node.Label = element.Name()
		node.Collapsed = true
		node.Checkbox = true
		node.Path = filepath.Join(path, element.Name())

		if element.IsDir() {
			// let's check if the folder is empty
			// we can still get an i/o error, if that's the case
			// we assume the folder's empty
			// otherwise we act accordingly
			folder := filepath.Join(path, element.Name())
			empty, err := lib.IsEmpty(folder)
			if err != nil {
				logger.Yellow("get-tree: unable to determine if folder is empty: %s", folder)
				node.Children = nil
			} else {
				if empty {
					node.Children = nil
				} else {
					node.Children = []domain.Node{domain.Node{Label: "Loading ...", Collapsed: true, Checkbox: false, Children: nil}}
				}
			}
		} else {
			node.Children = nil
		}

		items = append(items, node)
	}

	entry.Nodes = items

	return entry
}

func (c *Core) getLog() []string {
	cmd := "tail -n 100 /boot/logs/unbalance.log"

	log := make([]string, 0)

	// err := lib.Shell(cmd, mlog.Warning, "Get Log error:", "", func(line string) {
	// 	log = append(log, line)
	// })

	err := lib.Shell(cmd, "", func(line string) {
		log = append(log, line)
	})
	if err != nil {
		logger.Yellow("unable to get log: %s", err)
	}

	return log
}

func getCertificateName(certDir, name string) string {
	cert := filepath.Join(certDir, "certificate_bundle.pem")

	exists, err := lib.Exists(cert)
	if err != nil {
		logger.Yellow("unable to check for %s presence:(%s)", cert, err)
		return ""
	}

	if exists {
		return cert
	}

	cert = filepath.Join(certDir, name+"_unraid_bundle.pem")

	exists, err = lib.Exists(cert)
	if err != nil {
		logger.Yellow("unable to check for %s presence:(%s)", cert, err)
		return ""
	}

	if exists {
		return cert
	}

	return ""
}
