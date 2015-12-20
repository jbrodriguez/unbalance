package main

import (
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	// "jbrodriguez/unbalance/server/model"
	"jbrodriguez/unbalance/server/lib"
	"jbrodriguez/unbalance/server/services"
	"os"
	"os/signal"
	// "path/filepath"
	"fmt"
	"log"
	"syscall"
)

var Version string

func main() {
	settings, err := lib.NewSettings(Version)
	if err != nil {
		log.Printf("Unable to load settings: %s", err.Error())
		os.Exit(1)
	}

	mlog.Start(mlog.LevelInfo, settings.Log)

	mlog.Info("unBALANCE v%s starting up ...", Version)

	// mlog.Info("%+v", settings)

	var msg string
	if exists, _ := lib.Exists(settings.Conf); exists {
		msg = fmt.Sprintf("Using config file %s ...", settings.Conf)
	} else {
		msg = "No config file exists yet. Using app defaults ..."
	}
	mlog.Info(msg)

	bus := pubsub.New(623)

	socket := services.NewSocket(bus, settings)
	server := services.NewServer(bus, settings)
	core := services.NewCore(bus, settings)

	socket.Start()
	server.Start()
	core.Start()

	mlog.Info("Press Ctrl+C to stop ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL)
	mlog.Info("Received signal: (%s) ... shutting down the app now ...", <-c)

	core.Stop()
	server.Stop()
	socket.Stop()
}
