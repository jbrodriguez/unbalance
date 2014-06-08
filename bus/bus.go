package bus

import (
	"apertoire.net/unbalance/message"
	"github.com/golang/glog"
)

type Bus struct {
	GetBestFit chan *message.FitData
	GetDisks   chan *message.Disks
}

func (self *Bus) Start() {
	glog.Info("Bus starting up ...")

	self.GetBestFit = make(chan *message.FitData)
	self.GetDisks = make(chan *message.Disks)
}
