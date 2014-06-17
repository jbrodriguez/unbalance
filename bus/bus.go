package bus

import (
	"apertoire.net/unbalance/message"
	"github.com/golang/glog"
)

type Bus struct {
	GetBestFit chan *message.BestFit
	GetStatus  chan *message.StorageStatus
}

func (self *Bus) Start() {
	glog.Info("Bus starting up ...")

	self.GetBestFit = make(chan *message.BestFit)
	self.GetStatus = make(chan *message.StorageStatus)
}
