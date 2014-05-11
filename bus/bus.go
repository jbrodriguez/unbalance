package bus

import (
	"apertoire.net/unbalance/message"
	"log"
)

type Bus struct {
	GetBestFit chan *message.FitData
}

func (self *Bus) Start() {
	log.Println("Bus starting up ...")

	self.GetBestFit = make(chan *message.FitData)
}
