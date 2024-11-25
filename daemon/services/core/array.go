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
	"unbalance/daemon/domain"
	"unbalance/daemon/lib"
	"unbalance/daemon/logger"

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

func (c *Core) refreshUnraid() *domain.Unraid {
	unraid := c.state.Unraid

	newunraid, err := getArrayData()
	if err != nil {
		logger.Yellow("Unable to get storage: %s", err)
		return unraid
	}

	return newunraid
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
	mounts := make(map[string]bool)

	// err = lib.Shell("df --block-size=1 /mnt/*", mlog.Warning, "Refresh error:", "", func(line string) {
	// 	data := strings.Fields(line)
	// 	size[data[5]], _ = strconv.ParseUint(data[1], 10, 64)
	// 	free[data[5]], _ = strconv.ParseUint(data[3], 10, 64)
	// })

	err = lib.Shell("df --block-size=1 /mnt/*", "", func(line string) {
		data := strings.Fields(line)
		size[data[5]], _ = strconv.ParseUint(data[1], 10, 64)
		free[data[5]], _ = strconv.ParseUint(data[3], 10, 64)
		mounts[data[5]] = true
	})
	if err != nil {
		return nil, fmt.Errorf("unable to get free/size data: %w", err)
	}

	var blockSize, totalSize, totalFree uint64
	var hasBlockSize bool
	disks := make([]*domain.Disk, 0)

	pools := make(map[string]bool)

	// identify cache pools among disks
	for _, section := range file.Sections() {
		// DEFAULT section
		if section.Key("name").String() == "" {
			continue
		}

		diskType := section.Key("type").String()
		diskName := section.Key("name").String()
		diskStatus := section.Key("status").String()

		if !(diskType == "Cache") {
			continue
		}

		if diskStatus == "DISK_NP" {
			continue
		}

		if _, ok := mounts["/mnt/"+diskName]; !ok {
			continue
		}

		pools[diskName] = true
	}

	for _, section := range file.Sections() {
		// DEFAULT section
		if section.Key("name").String() == "" {
			continue
		}

		diskType := section.Key("type").String()
		diskName := section.Key("name").String()
		diskStatus := section.Key("status").String()

		if diskType == "Parity" || diskType == "Flash" || (diskType == "Cache" && !pools[diskName]) || diskStatus == "DISK_NP" {
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
			if int64(blockSize) != stat.Bsize {
				if !hasBlockSize {
					blockSize = uint64(stat.Bsize)
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

	udevExists, err := lib.Exists("/var/state/unassigned.devices/unassigned.devices.json")
	if err != nil {
		return nil, fmt.Errorf("error while checking for udev presence", err)
	}

	if udevExists {
		udevs, err := lib.LoadUnassignedDevices("/var/state/unassigned.devices/unassigned.devices.json")
		if err != nil {
			return nil, fmt.Errorf("unable to get unassigned devices", err)
		}

		id := uint64(len(disks))

		for _, udev := range udevs {
			fmt.Print(" %+v", udev)

			disk := &domain.Disk{}

			id += 1

			disk.ID = id                                                       // 1
			disk.Name = strings.ReplaceAll(udev.Mountpoint, "/mnt/disks/", "") // disk1, cache
			disk.Path = udev.Mountpoint                                        // /mnt/disk1, /mnt/cache
			disk.Device = strings.ReplaceAll(udev.Disk, "/dev/", "")           // sdp
			disk.Type = "Data"                                                 // Flash, Parity, Data, Cache
			disk.FsType = udev.FileSystem                                      // xfs, reiserfs, btrfs
			disk.Free = udev.Avail
			disk.Size = udev.Size
			disk.Serial = udev.Serial // WDC_WD30EZRX-00DC0B0_WD-WMC9T204468
			disk.Status = "DISK_OK"

			var stat syscall.Statfs_t
			e := syscall.Statfs(disk.Path, &stat)
			if e == nil {
				disk.BlocksTotal = stat.Blocks
				disk.BlocksFree = stat.Bavail

				//
				if int64(blockSize) != stat.Bsize {
					if !hasBlockSize {
						blockSize = uint64(stat.Bsize)
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
	}

	unraid.Size = totalSize
	unraid.Free = totalFree
	unraid.BlockSize = blockSize

	sort.Slice(disks, func(i, j int) bool { return disks[i].ID < disks[j].ID })

	unraid.Disks = disks

	return unraid, nil
}

func (c *Core) GetTree(path, id string) domain.Branch {
	section := make(map[string]domain.Node)
	order := make([]string, 0)

	elements, _ := os.ReadDir(path)
	for _, element := range elements {
		var node domain.Node

		// default values
		node.ID = c.sid.MustGenerate()
		node.Parent = id
		node.Label = element.Name()

		if element.IsDir() {
			node.Dir = true
			// let's check if the folder is empty
			// we can still get an i/o error, if that's the case
			// we assume the folder's empty
			// otherwise we act accordingly
			empty, err := lib.IsEmpty(filepath.Join(path, element.Name()))
			if err != nil {
				// mlog.Warning("GetTree - Unable to determine if folder is empty: %s", folder)
				node.Leaf = true
			} else {
				node.Leaf = empty
			}
		} else {
			node.Leaf = true
		}

		section[node.ID] = node
		order = append(order, node.ID)
	}

	return domain.Branch{
		Nodes: section,
		Order: order,
	}
}

func (c *Core) Locate(path string) []string {
	logger.Olive("path %s", path)
	locations := make([]string, 0)

	for _, disk := range c.state.Unraid.Disks {
		name := strings.Replace(path, "/mnt/user", "", 1)
		entry := filepath.Join(disk.Path, name)

		logger.Olive("name %s", name)
		logger.Olive("entry %s", entry)

		exists := true
		if _, err := os.Stat(entry); err != nil {
			exists = !os.IsNotExist(err)
		}

		if !exists {
			continue
		}

		locations = append(locations, disk.Name)
	}

	return locations
}

func (c *Core) GetLog() []string {
	cmd := "tail -n 100 /var/log/unbalanced.log"

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
