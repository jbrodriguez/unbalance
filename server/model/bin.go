package model

import (
	"github.com/jbrodriguez/mlog"
	"jbrodriguez/unbalance/server/lib"
)

type Bin struct {
	Size  int64
	Items []*Item
}

func (self *Bin) Add(item *Item) {
	self.Items = append(self.Items, item)
	self.Size += item.Size
}

func (self *Bin) Print() {
	for _, item := range self.Items {
		mlog.Info("[%s] %s\n", lib.ByteSize(item.Size), item.Name)
	}
}

type ByFilled []*Bin

func (s ByFilled) Len() int           { return len(s) }
func (s ByFilled) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s ByFilled) Less(i, j int) bool { return s[i].Size > s[j].Size }
