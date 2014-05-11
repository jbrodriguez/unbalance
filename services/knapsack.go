package services

import (
	"apertoire.net/unbalance/bus"
	"apertoire.net/unbalance/helper"
	"apertoire.net/unbalance/message"
	"fmt"
	"log"
)

type Knapsack struct {
	Bus *bus.Bus
}

func (self *Knapsack) Start() {
	log.Printf("starting Knapsack service ...")

	go self.react()

	log.Printf("Knapsack service started")
}

func (self *Knapsack) Stop() {
	// nothing right now
	log.Printf("Knapsack service stopped")
}

func (self *Knapsack) react() {
	for {
		select {
		case msg := <-self.Bus.GetBestFit:
			go self.doGetBestFit(msg)
		}
	}
}

func (self *Knapsack) doGetBestFit(msg *message.FitData) {
	packer := helper.NewPacker(msg.SourceDisk, msg.TargetDisk)

	free, err := packer.GetFreeSpace()
	if err != nil {
		log.Println(fmt.Sprintf("Available Space on %s: %d", msg.TargetDisk, free))
	}

	packer.GetItems("films/bluray")
	packer.GetItems("films/blurip")

	// packer.print()

	packer.BestFit()

	packer.Print()

	// items := self.getItems(msg.SourceDisk, "films/blurip")
	// for item := range items {
	// 	packer.Add(item)
	// }

	// items := self.getItems(msg.SourceDisk, "films/xvid")
	// for item := range items {
	// 	packer.Add(item)
	// }

	// items := self.getItems(msg.SourceDisk, "films/dvd")
	// for item := range items {
	// 	packer.Add(item)
	// }

}
