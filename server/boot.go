package main

import (
	"apertoire.net/unbalance/server/model"
	"apertoire.net/unbalance/server/services"
	"fmt"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
)

var Version string

func main() {
	config := model.Config{}
	config.Init(Version)

	bus := pubsub.New(1)

	server := services.NewServer(bus, &config)
	core := services.NewCore(bus, &config)

	server.Start()
	core.Start()

	mlog.Info("press enter to stop ...")
	var input string
	fmt.Scanln(&input)

	core.Stop()
	server.Stop()
}
