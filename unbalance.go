package main

import (
	"log"
	"path/filepath"

	"github.com/alecthomas/kong"
	"github.com/cskr/pubsub"
	"gopkg.in/natefinch/lumberjack.v2"

	"unbalance/daemon/cmd"
	"unbalance/daemon/domain"
)

var Version string

// const ReservedSpace int64 = 512 * 1024 * 1024 // 512Mb

var cli struct {
	Port    string `default:"6237" help:"port to listen on"`
	LogsDir string `default:"/var/log" help:"directory to store logs"`

	// Config vars
	DryRun         bool     `env:"DRY_RUN" default:"true" help:"perform a dry-run rather than actual work"`
	NotifyPlan     int      `env:"NOTIFY_PLAN" default:"0" help:"notify via email after plan operation has completed (unraid notifications must be set up first): 0 - No notifications; 1 - Simple notifications; 2 - Detailed notifications"`
	NotifyTransfer int      `env:"NOTIFY_TRANSFER" default:"0" help:"notify via email after transfer operation has completed (unraid notifications must be set up first): 0 - No notifications; 1 - Simple notifications; 2 - Detailed notifications"`
	ReservedAmount uint64   `env:"RESERVED_AMOUNT" default:"512" help:"Minimun Amount of space to reserve"`
	ReservedUnit   string   `env:"RESERVED_UNIT" default:"Mb" help:"Reserved Amount unit: Mb, Gb or %"`
	RsyncArgs      []string `env:"RSYNC_ARGS" default:"-X" help:"custom rsync arguments"`
	Verbosity      int      `env:"VERBOSITY" default:"0" help:"include rsync output in log files: 0 (default) - include; 1 - do not include"`
	CheckForUpdate int      `env:"CHECK_FOR_UPDATE" default:"1" help:"checkForUpdate: 0 - dont' check; 1 (default) - check"`
	RefreshRate    int      `env:"REFRESH_RATE" default:"250" help:"how often to refresh the ui while running a command (in milliseconds)"`

	Boot cmd.Boot `cmd:"" default:"1" help:"start processing"`
}

func main() {
	// Users can set some value that falls below ReservedSpace, but during planning we force ReservedSpace if
	// reservation is less than that
	// Also, if they enter some unrecognized unit, we will used ReservedSpace (in planning as well)
	// ctx := kong.Parse(&cli, kong.Vars{
	// 	"reserved_amount": strconv.FormatUint(common.ReservedSpace/1024/1024, 10),
	// })
	ctx := kong.Parse(&cli)

	log.Printf("cli: %+v", cli)

	log.SetOutput(&lumberjack.Logger{
		Filename:   filepath.Join(cli.LogsDir, "unbalanced.log"),
		MaxSize:    10, // megabytes
		MaxBackups: 10,
		MaxAge:     28, //days
		// Compress:   true, // disabled by default
	})

	err := ctx.Run(&domain.Context{
		Port: cli.Port,
		Config: domain.Config{
			Version:        Version,
			DryRun:         cli.DryRun,
			NotifyPlan:     cli.NotifyPlan,
			NotifyTransfer: cli.NotifyTransfer,
			ReservedAmount: cli.ReservedAmount,
			ReservedUnit:   cli.ReservedUnit,
			RsyncArgs:      cli.RsyncArgs,
			Verbosity:      cli.Verbosity,
			CheckForUpdate: cli.CheckForUpdate,
			RefreshRate:    cli.RefreshRate,
		},
		Hub: pubsub.New(23),
	})
	ctx.FatalIfErrorf(err)
}
