package services

import (
	"apertoire.net/unbalance/bus"
	"apertoire.net/unbalance/message"
	"bufio"
	"fmt"
	"io"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
)

type Knapsack struct {
	Bus         *bus.Bus
	reFreeSpace *regexp.Regexp
	reItems     *regexp.Regexp
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
	freespace, _ := self.getFreespace(msg.TargetDisk)

	packer := &Packer{Name: msg.TargetDisk, Size: freespace}

	items := self.getItems(msg.SourceDisk, "films/bluray")
	for _, item := range items {
		packer.add(item)
	}

	// packer.print()

	packer.bestFit()

	packer.sortBins()

	packer.printBins()

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
