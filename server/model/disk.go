package model

import (
	"apertoire.net/unbalance/server/helper"
	"fmt"
)

type Disk struct {
	Id      int    `json:"id"`
	Name    string `json:"name"`
	Path    string `json:"path"`
	Device  string `json:"device"`
	Free    int64  `json:"free"`
	NewFree int64  `json:"newFree"`
	Size    int64  `json:"size"`
	Serial  string `json:"serial"`
	Status  string `json:"status"`
	Bin     *Bin   `json:"-"`
}

func (self *Disk) Print() {
	// this disk was not assigned to a bin
	if self.Bin != nil {
		fmt.Println("=========================================================")
		fmt.Printf("Disk(%s):ALLOCATED %d folders:[%s/%s] %2.2f%%\n", self.Path, len(self.Bin.Items), helper.ByteSize(self.Bin.Size), helper.ByteSize(self.Free), (float64(self.Bin.Size)/float64(self.Free))*100)
		fmt.Println("---------------------------------------------------------")
		self.Bin.Print()
		fmt.Println("---------------------------------------------------------")
		fmt.Println("")
	} else {
		fmt.Println("=========================================================")
		fmt.Printf("Disk(%s):NO ALLOCATION:[0/%s] 0%%\n", self.Path, helper.ByteSize(self.Free))
		fmt.Println("---------------------------------------------------------")
		fmt.Println("---------------------------------------------------------")
		fmt.Println("")
	}
}

func (self *Disk) toString() string {
	return fmt.Sprintf("Id(%d); Name(%s); Path(%s); Device(%s), Free(%s); NewFree(%s); Size(%s); Serial(%s); Status(%s); Bin(%v)",
		self.Id,
		self.Name,
		self.Path,
		self.Device,
		helper.ByteSize(self.Free),
		helper.ByteSize(self.NewFree),
		helper.ByteSize(self.Size),
		self.Serial,
		self.Status, self.Bin)
}

type ByFree []*Disk

func (s ByFree) Len() int           { return len(s) }
func (s ByFree) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFree) Less(i, j int) bool { return s[i].Free > s[j].Free }
