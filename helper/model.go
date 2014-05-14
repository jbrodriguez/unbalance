package helper

import (
	"fmt"
)

type Item struct {
	Name string
	Size uint64
	Path string
}

type Bin struct {
	Size  uint64
	Items []*Item
}

func (self *Bin) add(item *Item) {
	self.Items = append(self.Items, item)
	self.Size += item.Size
}

func (self *Bin) Print() {
	for _, item := range self.Items {
		fmt.Println(fmt.Sprintf("[%d] %s", item.Size, item.Name))
	}
}

type ByFilled []*Bin

func (s ByFilled) Len() int           { return len(s) }
func (s ByFilled) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFilled) Less(i, j int) bool { return s[i].Size > s[j].Size }

type BySize []*Item

func (s BySize) Len() int           { return len(s) }
func (s BySize) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s BySize) Less(i, j int) bool { return s[i].Size > s[j].Size }

type Disk struct {
	Path string
	Free uint64
	Bin  *Bin
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
		fmt.Println(fmt.Sprintf("[0/%d] 0% (%s)", self.Free, self.Path))
		fmt.Println("---------------------------------------------------------")
		fmt.Println("---------------------------------------------------------")
		fmt.Println("")
	}
}
