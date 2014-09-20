package main

import (
	"apertoire.net/unbalance/server/services"
	"fmt"
	"github.com/apertoire/mlog"
	"github.com/apertoire/pubsub"
)

func main() {
	mlog.Start(mlog.LevelInfo, "./log/app.log")
	mlog.Info("starting up ...")

	bus := pubsub.New(1)

	server := services.NewServer(bus)
	core := services.NewCore(bus)

	server.Start()
	core.Start()

	mlog.Info("press enter to stop ...")
	var input string
	fmt.Scanln(&input)

	core.Stop()
	server.Stop()
}
