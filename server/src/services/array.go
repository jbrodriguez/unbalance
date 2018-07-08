package services

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"syscall"
	"time"

	"jbrodriguez/unbalance/server/src/common"
	"jbrodriguez/unbalance/server/src/domain"
	"jbrodriguez/unbalance/server/src/dto"
	"jbrodriguez/unbalance/server/src/lib"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	ini "github.com/vaughan0/go-ini"
)

const certDir = "/boot/config/ssl/certs"

// Array -
type Array struct {
	bus      *pubsub.PubSub
	settings *lib.Settings
	actor    *actor.Actor
}

// NewArray -
func NewArray(bus *pubsub.PubSub, settings *lib.Settings) *Array {
	array := &Array{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
	}

	return array
}

// Start -
func (a *Array) Start() (err error) {
	mlog.Info("starting service Array ...")

	err = a.SanityCheck(a.settings.APIFolders)
	if err != nil {
		return err
	}

	a.actor.Register(common.IntGetArrayStatus, a.getStatus)
	a.actor.Register(common.APIGetTree, a.getTree)
	a.actor.Register(common.APIGetLog, a.getLog)

	go a.actor.React()

	return nil
}

// Stop -
func (a *Array) Stop() {
	mlog.Info("stopped service Array ...")
}

// SanityCheck -
func (a *Array) SanityCheck(locations []string) error {
	location := lib.SearchFile("var.ini", locations)
	if location == "" {
		return fmt.Errorf("Unable to find var.ini (%s)", strings.Join(locations, ", "))
	}

	location = lib.SearchFile("disks.ini", locations)
	if location == "" {
		return fmt.Errorf("Unable to find var.ini (%s)", strings.Join(locations, ", "))
	}

	return nil
}

// GetCertificate -
func (a *Array) GetCertificate() string {
	// get ssl settings
	ident, err := ini.LoadFile("/boot/config/ident.cfg")
	if err != nil {
		return ""
	}

	usessl, _ := ident.Get("", "USE_SSL")
	usessl = strings.Replace(usessl, "\"", "", -1)

	// get array status
	file, err := ini.LoadFile("/var/local/emhttp/var.ini")
	if err != nil {
		return ""
	}

	name, _ := file.Get("", "NAME")
	name = strings.Replace(name, "\"", "", -1)

	cert := getCertificateName(certDir, name)

	secure := cert != "" && !(usessl == "" || usessl == "no" || usessl == "auto")

	if secure {
		return cert
	}

	return ""
}

func (a *Array) getStatus(msg *pubsub.Message) {
	unraid, err := getArrayData()
	if err != nil {
		msg.Reply <- dto.Message{Data: nil, Error: err}
	}

	msg.Reply <- dto.Message{Data: unraid, Error: nil}
}

func getArrayData() (*domain.Unraid, error) {
	unraid := &domain.Unraid{}

	// get array status
	file, err := ini.LoadFile("/var/local/emhttp/var.ini")
	if err != nil {
		return nil, err
	}

	tmp, _ := file.Get("", "mdNumDisks")
	numDisks := strings.Replace(tmp, "\"", "", -1)
	unraid.NumDisks, _ = strconv.ParseInt(numDisks, 10, 64)

	tmp, _ = file.Get("", "mdNumProtected")
	numProtected := strings.Replace(tmp, "\"", "", -1)
	unraid.NumProtected, _ = strconv.ParseInt(numProtected, 10, 64)

	tmp, _ = file.Get("", "sbSynced")
	synced := strings.Replace(tmp, "\"", "", -1)
	ut, _ := strconv.ParseInt(synced, 10, 64)
	unraid.Synced = time.Unix(ut, 0)

	tmp, _ = file.Get("", "sbSyncErrs")
	syncErrs := strings.Replace(tmp, "\"", "", -1)
	unraid.SyncErrs, _ = strconv.ParseInt(syncErrs, 10, 64)

	tmp, _ = file.Get("", "mdResync")
	resync := strings.Replace(tmp, "\"", "", -1)
	unraid.Resync, _ = strconv.ParseInt(resync, 10, 64)

	tmp, _ = file.Get("", "mdResyncPos")
	resyncPos := strings.Replace(tmp, "\"", "", -1)
	unraid.ResyncPos, _ = strconv.ParseInt(resyncPos, 10, 64)

	tmp, _ = file.Get("", "mdState")
	unraid.State = strings.Replace(tmp, "\"", "", -1)

	// get disks
	file, err = ini.LoadFile("/var/local/emhttp/disks.ini")
	if err != nil {
		return nil, err
	}

	// get free/size data
	free := make(map[string]int64)
	size := make(map[string]int64)

	err = lib.Shell("df --block-size=1 /mnt/*", mlog.Warning, "Refresh error:", "", func(line string) {
		data := strings.Fields(line)
		size[data[5]], _ = strconv.ParseInt(data[1], 10, 64)
		free[data[5]], _ = strconv.ParseInt(data[3], 0, 64)
	})

	if err != nil {
		return nil, err
	}

	var totalSize, totalFree, blockSize int64
	var hasBlockSize bool
	disks := make([]*domain.Disk, 0)

	for _, section := range file {
		diskType := strings.Replace(section["type"], "\"", "", -1)
		diskName := strings.Replace(section["name"], "\"", "", -1)
		diskStatus := strings.Replace(section["status"], "\"", "", -1)

		if diskType == "Parity" || diskType == "Flash" || (diskType == "Cache" && len(diskName) > 5 || diskStatus == "DISK_NP") {
			continue
		}

		disk := &domain.Disk{}

		disk.ID, _ = strconv.ParseInt(strings.Replace(section["idx"], "\"", "", -1), 10, 64) // 1
		disk.Name = diskName                                                                 // disk1, cache
		disk.Path = "/mnt/" + disk.Name                                                      // /mnt/disk1, /mnt/cache
		disk.Device = strings.Replace(section["device"], "\"", "", -1)                       // sdp
		disk.Type = diskType                                                                 // Flash, Parity, Data, Cache
		disk.FsType = strings.Replace(section["fsType"], "\"", "", -1)                       // xfs, reiserfs, btrfs
		disk.Free = free[disk.Path]
		disk.Size = size[disk.Path]
		disk.Serial = strings.Replace(section["id"], "\"", "", -1) // WDC_WD30EZRX-00DC0B0_WD-WMC9T204468
		disk.Status = diskStatus                                   // DISK_OK

		var stat syscall.Statfs_t
		e := syscall.Statfs(disk.Path, &stat)
		if e == nil {
			disk.BlocksTotal = int64(stat.Blocks)
			disk.BlocksFree = int64(stat.Bavail)

			if blockSize != int64(stat.Bsize) {
				if !hasBlockSize {
					blockSize = int64(stat.Bsize)
				} else {
					blockSize = 0
				}

				hasBlockSize = true
			}
		}
		// mlog.Info("name(%s),blocks(%d),size(%d),avail(%d)", disk.Path, disk.BlocksTotal, blockSize, disk.BlocksFree)

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

// GetTree -
func (a *Array) getTree(msg *pubsub.Message) {
	path := msg.Payload.(string)

	entry := &dto.Entry{Path: path}

	items := make([]dto.Node, 0)

	elements, _ := ioutil.ReadDir(path)
	for _, element := range elements {
		var node dto.Node

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
				mlog.Warning("GetTree - Unable to determine if folder is empty: %s", folder)
				node.Children = nil
			} else {
				if empty {
					node.Children = nil
				} else {
					node.Children = []dto.Node{dto.Node{Label: "Loading ...", Collapsed: true, Checkbox: false, Children: nil}}
				}
			}
		} else {
			node.Children = nil
		}

		items = append(items, node)
	}

	entry.Nodes = items

	msg.Reply <- entry
}

func (a *Array) getLog(msg *pubsub.Message) {
	cmd := "tail -n 100 /boot/logs/unbalance.log"

	log := make([]string, 0)

	err := lib.Shell(cmd, mlog.Warning, "Get Log error:", "", func(line string) {
		log = append(log, line)
	})

	if err != nil {
		mlog.Warning("Unable to get log: %s", err)
	}

	outbound := &dto.Packet{Topic: "gotLog", Payload: log}
	a.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
}

func getCertificateName(certDir, name string) string {
	cert := filepath.Join(certDir, "certificate_bundle.pem")

	exists, err := lib.Exists(cert)
	if err != nil {
		mlog.Warning("unable to check for %s presence:(%s)", cert, err)
		return ""
	}

	if exists {
		return cert
	}

	cert = filepath.Join(certDir, name+"_unraid_bundle.pem")

	exists, err = lib.Exists(cert)
	if err != nil {
		mlog.Warning("unable to check for %s presence:(%s)", cert, err)
		return ""
	}

	if exists {
		return cert
	}

	return ""
}
