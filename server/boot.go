package main

import (
	"apertoire.net/unbalance/server/model"
	"apertoire.net/unbalance/server/services"
	"flag"
	"fmt"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
)

var Version string

const (
	defaultCfgLocation = "/boot/config/plugins/unbalance/"
	defaultCfgUsage    = "Path to the config file"
	defaultLogLocation = ""
	defaultLogUsage    = "Path to the log file"
)

func main() {
	var cfg, log string

	flag.StringVar(&cfg, "config", defaultCfgLocation, defaultCfgUsage)
	flag.StringVar(&cfg, "c", defaultCfgLocation, defaultCfgUsage+" (shorthand)")

	flag.StringVar(&log, "log", defaultLogLocation, defaultLogUsage)
	flag.StringVar(&log, "l", defaultLogLocation, defaultLogUsage+" (shorthand)")

	flag.Parse()

	config := model.Config{}
	config.Init(Version, cfg, log)

	bus := pubsub.New(1)

	socket := services.NewSocket(bus, &config)
	server := services.NewServer(bus, &config, socket)
	core := services.NewCore(bus, &config)

	socket.Start()
	server.Start()
	core.Start()

	mlog.Info("press enter to stop ...")
	var input string
	fmt.Scanln(&input)

	core.Stop()
	server.Stop()
	socket.Stop()
}
