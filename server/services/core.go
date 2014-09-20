package services

import (
	"apertoire.net/unbalance/server/model"
	"github.com/apertoire/mlog"
	"github.com/apertoire/pubsub"
)

type Core struct {
	bus     *pubsub.PubSub
	storage *model.Unraid

	chanStorageInfo chan *pubsub.Message
}

func NewCore(bus *pubsub.PubSub) *Core {
	core := &Core{bus: bus}
	core.storage = &model.Unraid{}
	return core
}

func (self *Core) Start() {
	mlog.Info("starting service Core ...")

	self.chanStorageInfo = self.bus.Sub("cmd.getStorageInfo")

	go self.react()
}

func (self *Core) Stop() {
	mlog.Info("stopped service Core ...")
}

func (self *Core) react() {
	for {
		select {
		case msg := <-self.chanStorageInfo:
			go self.getStorageInfo(msg)
		}
	}
}

func (self *Core) getStorageInfo(msg *pubsub.Message) {
	mlog.Info("La vita e bella")

	msg.Reply <- self.storage.Refresh()
}
