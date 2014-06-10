package bus

import (
	"apertoire.net/unbalance/message"
	"github.com/golang/glog"
)

type Bus struct {
	GetBestFit chan *message.FitData
	GetStatus  chan *message.Status
}

func (self *Bus) Start() {
	glog.Info("Bus starting up ...")

	self.GetBestFit = make(chan *message.FitData)
	self.GetStatus = make(chan *message.Status)
}
