package main

import (
	"apertoire.net/unbalance/server/model"
	"apertoire.net/unbalance/server/services"
	"flag"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	"os"
	"os/signal"
	"path/filepath"
)

var Version string

const (
	defaultCfgLocation   = "/boot/config/plugins/unbalance/"
	defaultCfgUsage      = "Path to the config file"
	defaultLogLocation   = ""
	defaultLogUsage      = "Path to the log file"
	unbalanceLogFilePath = "UNBALANCE_LOGFILEPATH"
)

func main() {
	var cfg, log string

	flag.StringVar(&cfg, "config", defaultCfgLocation, defaultCfgUsage)
	flag.StringVar(&cfg, "c", defaultCfgLocation, defaultCfgUsage+" (shorthand)")

	flag.StringVar(&log, "log", defaultLogLocation, defaultLogUsage)
	flag.StringVar(&log, "l", defaultLogLocation, defaultLogUsage+" (shorthand)")

	flag.Parse()

	finalLog := os.Getenv(unbalanceLogFilePath)

	// command line param takes precedence over environment var
	if log != "" {
		finalLog = log
	}

	if finalLog != "" {
		mlog.Start(mlog.LevelInfo, filepath.Join(finalLog, "unbalance.log"))
	} else {
		mlog.Start(mlog.LevelInfo, "")
	}

	mlog.Info("unBALANCE v%s starting up ...", Version)

	settings := model.Settings{}
	settings.Init(Version, cfg, log)

	bus := pubsub.New(1)

	socket := services.NewSocket(bus, &settings)
	server := services.NewServer(bus, &settings, socket)
	core := services.NewCore(bus, &settings)

	socket.Start()
	server.Start()
	core.Start()

	mlog.Info("Press Ctrl+C to stop ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	for _ = range c {
		mlog.Info("Received an interrupt, shutting the app down ...")

		core.Stop()
		server.Stop()
		socket.Stop()

		break
	}

}
