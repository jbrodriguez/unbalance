package model

import (
	"apertoire.net/unbalance/server/helper"
	"github.com/jbrodriguez/mlog"
	"strconv"
	"strings"
	"time"
)

type Unraid struct {
	Condition *Condition `json:"condition"`
	Disks     []*Disk    `json:"disks"`

	SourceDiskName string
	BytesToMove    uint64

	disks [25]*Disk
}

func NewUnraid() (unraid *Unraid) {
	unraid = &Unraid{}
	return unraid
}

type DiskInfoDTO struct {
	Free map[string]uint64
	Size map[string]uint64
}

func (self *Unraid) Refresh() *Unraid {
	self.Disks = make([]*Disk, 0)
	self.Condition = &Condition{}

	helper.Shell("/root/mdcmd status|strings", self.readUnraidConfig, nil)

	di := &DiskInfoDTO{Free: make(map[string]uint64), Size: make(map[string]uint64)}
	helper.Shell("df --block-size=1 /mnt/disk*", self.getDiskInfo, di)

	for _, disk := range self.disks {
		if disk != nil && disk.Name != "Parity" && disk.Status == "DISK_OK" {
			disk.Size = di.Size[disk.Path]
			disk.Free = di.Free[disk.Path]
			disk.NewFree = disk.Free

			self.Condition.Size += disk.Size
			self.Condition.Free += disk.Free
			self.Condition.NewFree += disk.Free

			self.Disks = append(self.Disks, disk)
		}
	}

	self.Print()

	// condition := &Condition{
	// 	NumDisks:     14,
	// 	NumProtected: 14,
	// 	Synced:       time.Date(1969, 12, 31, 19, 0, 0, 0, time.UTC),
	// 	SyncErrs:     0,
	// 	Resync:       0,
	// 	ResyncPrcnt:  0,
	// 	ResyncPos:    0,
	// 	State:        "STARTED",
	// 	Size:         41006844604416,
	// 	Free:         768601899008,
	// 	NewFree:      768601899008,
	// }

	// disks := []*Disk{
	// 	(&Disk{Id: 1, Name: "md1", Path: "/mnt/disk1", Device: "sdn", Free: 2798870528, NewFree: 2798870528, Size: 3000501350400, Serial: "WDC_WD30EZRX-00DC0B0_WD-WMC1T0345089", Status: "DISK_OK"}),
	// 	(&Disk{Id: 2, Name: "md2", Path: "/mnt/disk2", Device: "sdm", Free: 12431601664, NewFree: 12431601664, Size: 3000501350400, Serial: "WDC_WD30EZRX-00DC0B0_WD-WMC1T0373550", Status: "DISK_OK"}),
	// 	(&Disk{Id: 3, Name: "md3", Path: "/mnt/disk3", Device: "sdk", Free: 8654426112, NewFree: 8654426112, Size: 3000501350400, Serial: "ST3000DM001-9YN166_W1F181AR", Status: "DISK_OK"}),
	// 	(&Disk{Id: 4, Name: "md4", Path: "/mnt/disk4", Device: "sdl", Free: 110264877056, NewFree: 110264877056, Size: 3000501350400, Serial: "ST3000DM001-9YN166_Z1F1546H", Status: "DISK_OK"}),
	// 	(&Disk{Id: 5, Name: "md5", Path: "/mnt/disk5", Device: "sdi", Free: 7675904, NewFree: 7675904, Size: 3000501350400, Serial: "TOSHIBA_DT01ACA300_23CEUGZWS", Status: "DISK_OK"}),
	// 	(&Disk{Id: 6, Name: "md6", Path: "/mnt/disk6", Device: "sdj", Free: 13362188288, NewFree: 13362188288, Size: 3000501350400, Serial: "TOSHIBA_DT01ACA300_23CENSPWS", Status: "DISK_OK"}),
	// 	(&Disk{Id: 7, Name: "md7", Path: "/mnt/disk7", Device: "sdh", Free: 10317832192, NewFree: 10317832192, Size: 3000501350400, Serial: "TOSHIBA_DT01ACA300_23DG6Z7WS", Status: "DISK_OK"}),
	// 	(&Disk{Id: 8, Name: "md8", Path: "/mnt/disk8", Device: "sdb", Free: 116319207424, NewFree: 116319207424, Size: 3000501350400, Serial: "ST3000DM001-1CH166_W1F45LE8", Status: "DISK_OK"}),
	// 	(&Disk{Id: 9, Name: "md9", Path: "/mnt/disk9", Device: "sdg", Free: 25462644736, NewFree: 25462644736, Size: 3000501350400, Serial: "TOSHIBA_DT01ACA300_Y3UEB7GGS", Status: "DISK_OK"}),
	// 	(&Disk{Id: 10, Name: "md10", Path: "/mnt/disk10", Device: "sdf", Free: 380406677504, NewFree: 380406677504, Size: 3000501350400, Serial: "TOSHIBA_DT01ACA300_X3V9V7TGS", Status: "DISK_OK"}),
	// 	(&Disk{Id: 11, Name: "md11", Path: "/mnt/disk11", Device: "sde", Free: 0, NewFree: 0, Size: 3000501350400, Serial: "WDC_WD30EFRX-68AX9N0_WD-WMC1T0571629", Status: "DISK_OK"}),
	// 	(&Disk{Id: 12, Name: "md12", Path: "/mnt/disk12", Device: "sdd", Free: 6960766976, NewFree: 6960766976, Size: 4000664875008, Serial: "ST4000DM000-1F2168_Z301LVKC", Status: "DISK_OK"}),
	// 	(&Disk{Id: 13, Name: "md13", Path: "/mnt/disk13", Device: "sdc", Free: 67401682944, NewFree: 67401682944, Size: 4000664875008, Serial: "WDC_WD40EZRX-00SPEB0_WD-WCC4EM0WN2RE", Status: "DISK_OK"}),
	// }

	// self.Disks = disks
	// self.Condition = condition

	return self
}

func (self *Unraid) readUnraidConfig(line string, arg interface{}) {
	if strings.HasPrefix(line, "sbNumDisks") {
		nd := strings.Split(line, "=")
		self.Condition.NumDisks, _ = strconv.ParseUint(nd[1], 10, 64)
	}

	if strings.HasPrefix(line, "mdNumProtected") {
		np := strings.Split(line, "=")
		self.Condition.NumProtected, _ = strconv.ParseUint(np[1], 10, 64)
	}

	if strings.HasPrefix(line, "sbSynced") {
		sd := strings.Split(line, "=")
		ut, _ := strconv.ParseInt(sd[1], 10, 64)
		self.Condition.Synced = time.Unix(ut, 0)
	}

	if strings.HasPrefix(line, "sbSyncErrs") {
		sr := strings.Split(line, "=")
		self.Condition.SyncErrs, _ = strconv.ParseUint(sr[1], 10, 64)
	}

	if strings.HasPrefix(line, "mdResync") {
		rs := strings.Split(line, "=")
		self.Condition.Resync, _ = strconv.ParseUint(rs[1], 10, 64)
	}

	if strings.HasPrefix(line, "mdResyncPos") {
		rp := strings.Split(line, "=")
		self.Condition.ResyncPos, _ = strconv.ParseUint(rp[1], 10, 64)
	}

	if strings.HasPrefix(line, "mdState") {
		st := strings.Split(line, "=")
		self.Condition.State = st[1]
	}

	// Get Disks Information
	if strings.HasPrefix(line, "diskNumber") {
		dn := strings.FieldsFunc(line, delim)

		diskId, _ := strconv.Atoi(dn[2])
		// mlog.Info("diskId = %d", diskId)
		if self.disks[diskId] == nil {
			self.disks[diskId] = &Disk{Id: diskId, Path: "/mnt/disk" + dn[2]}
		}
	}

	if strings.HasPrefix(line, "diskName") {
		dm := strings.FieldsFunc(line, delim)

		diskId, _ := strconv.Atoi(dm[1])
		// mlog.Info("diskName %+v diskId %d", self.disks, diskId)
		if len(dm) > 2 {
			self.disks[diskId].Name = dm[2]
		} else if diskId == 0 {
			self.disks[diskId].Name = "Parity"
		}
	}

	if strings.HasPrefix(line, "diskId") {
		dm := strings.FieldsFunc(line, delim)

		diskId, _ := strconv.Atoi(dm[1])
		// mlog.Info("diskId diskId %d", diskId)
		if len(dm) > 2 {
			self.disks[diskId].Serial = dm[2]
		}
	}

	if strings.HasPrefix(line, "rdevStatus") {
		dm := strings.FieldsFunc(line, delim)

		diskId, _ := strconv.Atoi(dm[1])
		// mlog.Info("rdevStatus diskId %d", diskId)
		self.disks[diskId].Status = dm[2]
	}

	if strings.HasPrefix(line, "rdevName") {
		dm := strings.FieldsFunc(line, delim)

		diskId, _ := strconv.Atoi(dm[1])
		// mlog.Info("rdevName diskId %d", diskId)
		if len(dm) > 2 {
			self.disks[diskId].Device = dm[2]
		}
	}
}

func delim(r rune) bool {
	return r == '.' || r == '='
}

func (self *Unraid) getDiskInfo(line string, arg interface{}) {
	di := arg.(*DiskInfoDTO)

	data := strings.Fields(line)
	di.Size[data[5]], _ = strconv.ParseUint(data[1], 10, 64)
	di.Free[data[5]], _ = strconv.ParseUint(data[3], 0, 64)
}

func (self *Unraid) Print() {
	mlog.Info("Unraid Box: %+v", self.Condition)
	// glog.Info("NumDisks: ", self.Box.NumDisks)
	// glog.Info("NumProtected: ", self.Box.NumProtected)
	// glog.Info("Synced: ", self.Box.Synced)
	// glog.Info("SyncErrs: ", self.Box.SyncErrs)
	// glog.Info("Resync: ", self.Box.Resync)
	// glog.Info("ResyncPrcnt: ", self.Box.ResyncPrcnt)
	// glog.Info("ResyncPos: ", self.Box.ResyncPos)
	// glog.Info("State: ", self.Box.State)

	for _, disk := range self.Disks {
		mlog.Info("%+v", disk)
	}
}
