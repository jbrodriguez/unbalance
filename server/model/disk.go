package model

import (
	"fmt"
	"github.com/jbrodriguez/mlog"
	"jbrodriguez/unbalance/server/lib"
)

type Disk struct {
	Id      int64  `json:"id"`
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

func (self *Disk) Print() {
	// this disk was not assigned to a bin
	if self.Bin != nil {
		mlog.Info("=========================================================")
		mlog.Info("Disk(%s):ALLOCATED %d folders:[%s/%s] %2.2f%%\n", self.Path, len(self.Bin.Items), lib.ByteSize(self.Bin.Size), lib.ByteSize(self.Free), (float64(self.Bin.Size)/float64(self.Free))*100)
		mlog.Info("---------------------------------------------------------")
		self.Bin.Print()
		mlog.Info("---------------------------------------------------------")
		mlog.Info("")
	} else {
		mlog.Info("=========================================================")
		mlog.Info("Disk(%s):NO ALLOCATION:[0/%s] 0%%\n", self.Path, lib.ByteSize(self.Free))
		mlog.Info("---------------------------------------------------------")
		mlog.Info("---------------------------------------------------------")
		mlog.Info("")
	}
}

func (self *Disk) toString() string {
	return fmt.Sprintf("Id(%d); Name(%s); Path(%s); Device(%s); Type(%s); FsType(%s); Free(%s); NewFree(%s); Size(%s); Serial(%s); Status(%s); Bin(%v)",
		self.Id,
		self.Name,
		self.Path,
		self.Device,
		self.Type,
		self.FsType,
		lib.ByteSize(self.Free),
		lib.ByteSize(self.NewFree),
		lib.ByteSize(self.Size),
		self.Serial,
		self.Status, self.Bin)
}

type ByFree []*Disk

func (s ByFree) Len() int           { return len(s) }
func (s ByFree) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFree) Less(i, j int) bool { return s[i].Free > s[j].Free }

type ById []*Disk

func (s ById) Len() int           { return len(s) }
func (s ById) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ById) Less(i, j int) bool { return s[i].Id < s[j].Id }
