package model

import (
	"fmt"
	"github.com/jbrodriguez/mlog"
	"jbrodriguez/unbalance/server/src/lib"
)

// Disk -
type Disk struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Device  string `json:"device"`
	Type    string `json:"type"`
	FsType  string `json:"fsType"`
	Free    int64  `json:"free"`
	NewFree int64  `json:"newFree"`
	Size    int64  `json:"size"`
	Serial  string `json:"serial"`
	Status  string `json:"status"`
	Bin     *Bin   `json:"-"`
	Src     bool   `json:"src"`
	Dst     bool   `json:"dst"`
}

// Print -
func (d *Disk) Print() {
	// this disk was not assigned to a bin
	if d.Bin != nil {
		mlog.Info("=========================================================")
		mlog.Info("Disk(%s):ALLOCATED %d folders:[%s/%s] %2.2f%%\n", d.Path, len(d.Bin.Items), lib.ByteSize(d.Bin.Size), lib.ByteSize(d.Free), (float64(d.Bin.Size)/float64(d.Free))*100)
		mlog.Info("---------------------------------------------------------")
		d.Bin.Print()
		mlog.Info("---------------------------------------------------------")
		mlog.Info("")
	} else {
		mlog.Info("=========================================================")
		mlog.Info("Disk(%s):NO ALLOCATION:[0/%s] 0%%\n", d.Path, lib.ByteSize(d.Free))
		mlog.Info("---------------------------------------------------------")
		mlog.Info("---------------------------------------------------------")
		mlog.Info("")
	}
}

func (d *Disk) toString() string {
	return fmt.Sprintf("Id(%d); Name(%s); Path(%s); Device(%s); Type(%s); FsType(%s); Free(%s); NewFree(%s); Size(%s); Serial(%s); Status(%s); Bin(%v)",
		d.ID,
		d.Name,
		d.Path,
		d.Device,
		d.Type,
		d.FsType,
		lib.ByteSize(d.Free),
		lib.ByteSize(d.NewFree),
		lib.ByteSize(d.Size),
		d.Serial,
		d.Status, d.Bin)
}

// ByFree -
type ByFree []*Disk

func (s ByFree) Len() int           { return len(s) }
func (s ByFree) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFree) Less(i, j int) bool { return s[i].Free > s[j].Free }

// ByID -
type ByID []*Disk

func (s ByID) Len() int           { return len(s) }
func (s ByID) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByID) Less(i, j int) bool { return s[i].ID < s[j].ID }
