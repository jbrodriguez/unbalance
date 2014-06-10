package model

import (
	"fmt"
)

type Disk struct {
	Id     int    `json:"id"`
	Name   string `json:"name"`
	Path   string `json:"path"`
	Device string `json:"device"`
	Free   uint64 `json:"free"`
	Size   uint64 `json:"size"`
	Serial string `json:"serial"`
	Status string `json:"status"`
	Bin    *Bin   `json:"-"`
}

func (self *Disk) Print() {
	if self.Bin != nil {
		fmt.Println("=========================================================")
		fmt.Println(fmt.Sprintf("[%d/%d] %2.2f%% (%s)", self.Bin.Size, self.Free, (float64(self.Bin.Size)/float64(self.Free))*100, self.Path))
		fmt.Println("---------------------------------------------------------")
		self.Bin.Print()
		fmt.Println("---------------------------------------------------------")
		fmt.Println("")
	} else {
		fmt.Println("=========================================================")
		fmt.Println(fmt.Sprintf("[0/%d] 0%% (%s)", self.Free, self.Path))
		fmt.Println("---------------------------------------------------------")
		fmt.Println("---------------------------------------------------------")
		fmt.Println("")
	}
}

type ByFree []*Disk

func (s ByFree) Len() int           { return len(s) }
func (s ByFree) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFree) Less(i, j int) bool { return s[i].Free > s[j].Free }
