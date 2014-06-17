package bus

import (
	"apertoire.net/unbalance/lib"
	"github.com/golang/glog"
)

type Bus struct {
	GetBestFit chan *lib.BestFit
	GetStatus  chan *lib.Status
}

func (self *Bus) Start() {
	glog.Info("Bus starting up ...")

	self.GetBestFit = make(chan *lib.BestFit)
	self.GetStatus = make(chan *lib.Status)
}
