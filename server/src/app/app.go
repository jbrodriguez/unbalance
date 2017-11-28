package app

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"jbrodriguez/unbalance/server/src/lib"
	"jbrodriguez/unbalance/server/src/services"

	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
)

// App empty placeholder
type App struct {
}

// Setup app
func (a *App) Setup(version string) (*lib.Settings, error) {
	// look for unbalance.conf at the following places
	// /boot/config/plugins/unbalance/
	// <current dir>/unbalance.conf

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	locations := []string{
		"/boot/config/plugins/unbalance",
		cwd,
	}

	settings, err := lib.NewSettings("unbalance.conf", version, locations)

	return settings, err
}

// Run app
func (a *App) Run(settings *lib.Settings) {
	if settings.LogDir != "" {
		mlog.Start(mlog.LevelInfo, filepath.Join(settings.LogDir, "unbalance.log"))
	} else {
		mlog.Start(mlog.LevelInfo, "")
	}

	mlog.Info("unbalance v%s starting ...", settings.Version)

	var msg string
	if settings.Location == "" {
		msg = "No config file specified. Using app defaults ..."
	} else {
		msg = fmt.Sprintf("Using config file located at %s ...", settings.Location)
	}
	mlog.Info(msg)

	bus := pubsub.New(623)

	server := services.NewServer(bus, settings)
	array := services.NewArray(bus, settings)
	planner := services.NewPlanner(bus, settings)
	core := services.NewCore(bus, settings)

	server.Start()
	mlog.FatalIfError(array.Start())
	mlog.FatalIfError(planner.Start())
	mlog.FatalIfError(core.Start())

	mlog.Info("Press Ctrl+C to stop ...")

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, syscall.SIGINT)
	mlog.Info("Received signal: (%s) ... shutting down the app now ...", <-c)

	core.Stop()
	planner.Stop()
	array.Stop()
	server.Stop()

	err := mlog.Stop()
	if err != nil {
		log.Printf("Unable to stop mlog: %s", err)
	}
}
