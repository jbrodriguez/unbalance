package model

import (
	"bufio"
	"github.com/golang/glog"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type Unraid struct {
	Box   *Box    `json:"box"`
	Disks []*Disk `json:"disks"`
}

func delim(r rune) bool {
	return r == '.' || r == '='
}

func NewUnraid() (unraid *Unraid) {
	var box Box
	var _disks [25]*Disk
	var free map[string]uint64 = make(map[string]uint64)
	var size map[string]uint64 = make(map[string]uint64)
	var disks []*Disk

	txt := "/root/mdcmd status|strings"
	cmd := exec.Command("sh", "-c", txt)
	out, err := cmd.StdoutPipe()
	if err != nil {
		glog.Fatalf("Unable to stdoutpipe %s: %s", txt, err)
	}

	rd := bufio.NewReader(out)

	if err := cmd.Start(); err != nil {
		glog.Fatal("Unable to start mdcmd: ", err)
	}

	// Get unRaid array information
	for {

		// sbNumDisks=21
		// sbSynced=1400211269
		// sbSyncErrs=0
		// mdVersion=2.2.0
		// mdState=STARTED
		// mdNumProtected=21
		// diskNumber.0=0
		// diskName.0=
		// diskSize.0=1953514552
		// diskState.0=7
		// diskId.0=Hitachi_HDS5C3020ALA632_ML0220F30JR0ND
		// rdevNumber.0=0
		// rdevStatus.0=DISK_OK
		// rdevName.0=sdn
		// rdevSize.0=1953514552
		// rdevId.0=Hitachi_HDS5C3020ALA632_ML0220F30JR0ND
		// rdevNumErrors.0=0
		// rdevLastIO.0=0
		// rdevSpinupGroup.0=0

		line, err := rd.ReadString('\n')
		if err == io.EOF && len(line) == 0 {
			// Good end of file with no partial line
			break
		}
		if err == io.EOF {
			glog.Fatal("Last line not terminated: ", err)
		}
		line = line[:len(line)-1] // drop the '\n'
		if line[len(line)-1] == '\r' {
			line = line[:len(line)-1] // drop the '\r'
		}

		if strings.HasPrefix(line, "sbNumDisks") {
			nd := strings.Split(line, "=")
			box.NumDisks, _ = strconv.ParseUint(nd[1], 10, 64)
		}

		if strings.HasPrefix(line, "mdNumProtected") {
			np := strings.Split(line, "=")
			box.NumProtected, _ = strconv.ParseUint(np[1], 10, 64)
		}

		if strings.HasPrefix(line, "sbSynced") {
			sd := strings.Split(line, "=")
			ut, _ := strconv.ParseInt(sd[1], 10, 64)
			box.Synced = time.Unix(ut, 0)
		}

		if strings.HasPrefix(line, "sbSyncErrs") {
			sr := strings.Split(line, "=")
			box.SyncErrs, _ = strconv.ParseUint(sr[1], 10, 64)
		}

		if strings.HasPrefix(line, "mdResync") {
			rs := strings.Split(line, "=")
			box.Resync, _ = strconv.ParseUint(rs[1], 10, 64)
		}

		if strings.HasPrefix(line, "mdResyncPos") {
			rp := strings.Split(line, "=")
			box.ResyncPos, _ = strconv.ParseUint(rp[1], 10, 64)
		}

		if strings.HasPrefix(line, "mdState") {
			st := strings.Split(line, "=")
			box.State = st[1]
		}

		// Get Disks Information
		if strings.HasPrefix(line, "diskNumber") {
			dn := strings.FieldsFunc(line, delim)

			diskId, _ := strconv.Atoi(dn[2])
			if _disks[diskId] == nil {
				_disks[diskId] = &Disk{Id: diskId, Path: "/mnt/disk" + dn[2]}
			}
		}

		if strings.HasPrefix(line, "diskName") {
			dm := strings.FieldsFunc(line, delim)

			diskId, _ := strconv.Atoi(dm[1])
			glog.Info("diskName diskId ", diskId)
			if len(dm) > 2 {
				_disks[diskId].Name = dm[2]
			} else if diskId == 0 {
				_disks[diskId].Name = "Parity"
			}
		}

		if strings.HasPrefix(line, "diskId") {
			dm := strings.FieldsFunc(line, delim)

			diskId, _ := strconv.Atoi(dm[1])
			glog.Info("diskId diskId ", diskId)
			if len(dm) > 2 {
				_disks[diskId].Serial = dm[2]
			}
		}

		if strings.HasPrefix(line, "rdevStatus") {
			dm := strings.FieldsFunc(line, delim)

			diskId, _ := strconv.Atoi(dm[1])
			glog.Info("rdevStatus diskId ", diskId)
			_disks[diskId].Status = dm[2]
		}

		if strings.HasPrefix(line, "rdevName") {
			dm := strings.FieldsFunc(line, delim)

			diskId, _ := strconv.Atoi(dm[1])
			glog.Info("rdevName diskId ", diskId)
			if len(dm) > 2 {
				_disks[diskId].Device = dm[2]
			}
		}
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		glog.Fatal("Unable to wait for process to finish: ", err)
	}

	txt = "df --block-size=1 /mnt/disk*"
	cmd = exec.Command("sh", "-c", txt)
	out, err = cmd.StdoutPipe()
	if err != nil {
		glog.Fatalf("Unable to stdoutpipe %s: %s", txt, err)
	}

	rd = bufio.NewReader(out)

	if err := cmd.Start(); err != nil {
		glog.Fatal("Unable to start df: ", err)
	}

	// ignore first line since it's just headers
	line, err := rd.ReadString('\n')

	// Get unRaid array information
	for {

		// Filesystem           1B-blocks      Used Available Use% Mounted on
		// /dev/md1             2000337846272 1998411968512 1925877760 100% /mnt/disk1

		line, err = rd.ReadString('\n')
		if err == io.EOF && len(line) == 0 {
			// Good end of file with no partial line
			break
		}
		if err == io.EOF {
			glog.Fatal("Last line not terminated: ", err)
		}
		line = line[:len(line)-1] // drop the '\n'
		if line[len(line)-1] == '\r' {
			line = line[:len(line)-1] // drop the '\r'
		}

		data := strings.Fields(line)
		size[data[5]], _ = strconv.ParseUint(data[1], 10, 64)
		free[data[5]], _ = strconv.ParseUint(data[3], 0, 64)
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		glog.Fatal("Unable to wait for process to finish: ", err)
	}

	for _, disk := range _disks {
		if disk != nil && disk.Name != "Parity" && disk.Status == "DISK_OK" {
			disk.Size = size[disk.Path]
			disk.Free = free[disk.Path]
			disk.NewFree = disk.Free

			box.Size += disk.Size
			box.Free += disk.Free
			box.NewFree += disk.Free

			disks = append(disks, disk)
		}
	}

	// file, _ := os.Open("/var/local/emhttp/var.ini")

	return &Unraid{Box: &box, Disks: disks}
}

func (self *Unraid) Print() {
	glog.Infof("Unraid Box: %+v", self.Box)
	// glog.Info("NumDisks: ", self.Box.NumDisks)
	// glog.Info("NumProtected: ", self.Box.NumProtected)
	// glog.Info("Synced: ", self.Box.Synced)
	// glog.Info("SyncErrs: ", self.Box.SyncErrs)
	// glog.Info("Resync: ", self.Box.Resync)
	// glog.Info("ResyncPrcnt: ", self.Box.ResyncPrcnt)
	// glog.Info("ResyncPos: ", self.Box.ResyncPos)
	// glog.Info("State: ", self.Box.State)

	for _, disk := range self.Disks {
		glog.Infof("%+v", disk)
	}
}
