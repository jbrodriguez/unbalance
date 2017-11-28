package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/namsral/flag"
)

// ReservedSpace -
const ReservedSpace = 450000000 // 450Mb

// Config -
type Config struct {
	DryRun         bool     `json:"dryRun"`
	NotifyCalc     int      `json:"notifyCalc"`
	NotifyMove     int      `json:"notifyMove"`
	ReservedAmount int64    `json:"reservedAmount"`
	ReservedUnit   string   `json:"reservedUnit"`
	RsyncArgs      []string `json:"rsyncArgs"`
	Version        string   `json:"version"`
	Verbosity      int      `json:"verbosity"`
	CheckForUpdate int      `json:"checkForUpdate"`
}

// NotifyCalc/NotifyMove possible values
// 0 - no notification
// 1 - simple notification
// 2 - detailed notification

// Settings -
type Settings struct {
	Config

	Port       string
	LogDir     string
	APIFolders []string

	Location string
	confName string
}

const defaultConfLocation = "/boot/config/plugins/unbalance"

// NewSettings -
func NewSettings(name, version string, locations []string) (*Settings, error) {
	var port, logDir, folders, rsyncFlags, rsyncArgs, apiFolders string
	var dryRun bool
	var notifyCalc, notifyMove, verbosity, checkForUpdate int

	flagset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	flagset.StringVar(&port, "port", "6237", "port to run the server")
	flagset.StringVar(&logDir, "logdir", "/boot/logs", "pathname where log file will be written to")
	flagset.StringVar(&folders, "folders", "", "deprecated - do not use")
	flagset.BoolVar(&dryRun, "dryRun", true, "perform a dry-run rather than actual work")
	flagset.IntVar(&notifyCalc, "notifyCalc", 0, "notify via email after calculation operation has completed (unraid notifications must be set up first): 0 - No notifications; 1 - Simple notifications; 2 - Detailed notifications")
	flagset.IntVar(&notifyMove, "notifyMove", 0, "notify via email after move operation has completed (unraid notifications must be set up first): 0 - No notifications; 1 - Simple notifications; 2 - Detailed notifications")
	flagset.StringVar(&rsyncFlags, "rsyncFlags", "", "custom rsync flags") // to be deprecated
	flagset.StringVar(&rsyncArgs, "rsyncArgs", "", "custom rsync arguments")
	flagset.StringVar(&apiFolders, "apiFolders", "/var/local/emhttp", "folders to look for api endpoints")
	flagset.IntVar(&verbosity, "verbosity", 0, "include rsync output in log files: 0 (default) - include; 1 - do not include")
	flagset.IntVar(&checkForUpdate, "checkForUpdate", 1, "checkForUpdate: 0 - dont' check; 1 (default) - check")

	location := SearchFile(name, locations)
	if location != "" {
		flagset.String("config", filepath.Join(location, name), "config location")
	}

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	s := &Settings{}

	if rsyncArgs == "" {
		s.RsyncArgs = []string{""}
	} else {
		s.RsyncArgs = strings.Split(rsyncArgs, "|")
	}

	s.DryRun = dryRun
	s.NotifyCalc = notifyCalc
	s.NotifyMove = notifyMove
	s.ReservedAmount = ReservedSpace / 1000 / 1000
	s.ReservedUnit = "Mb"
	s.Verbosity = verbosity
	s.CheckForUpdate = checkForUpdate
	s.Version = version

	s.Port = port
	s.LogDir = logDir
	s.APIFolders = strings.Split(apiFolders, "|")
	s.Location = location
	s.confName = name

	return s, nil
}

// ToggleDryRun -
func (s *Settings) ToggleDryRun() {
	s.DryRun = !s.DryRun
}

// Save -
func (s *Settings) Save() (err error) {
	location := s.Location
	if location == "" {
		location = defaultConfLocation
	}

	confLocation := filepath.Join(location, s.confName)
	tmpFile := confLocation + ".tmp"

	if err = WriteLine(tmpFile, fmt.Sprintf("dryRun=%t", s.DryRun)); err != nil {
		return err
	}

	if err = WriteLine(tmpFile, fmt.Sprintf("notifyCalc=%d", s.NotifyCalc)); err != nil {
		return err
	}

	if err = WriteLine(tmpFile, fmt.Sprintf("notifyMove=%d", s.NotifyMove)); err != nil {
		return err
	}

	rsyncArgs := strings.Join(s.RsyncArgs, "|")
	if err = WriteLine(tmpFile, fmt.Sprintf("rsyncArgs=%s", rsyncArgs)); err != nil {
		return err
	}

	if err = WriteLine(tmpFile, fmt.Sprintf("verbosity=%d", s.Verbosity)); err != nil {
		return err
	}

	err = os.Rename(tmpFile, confLocation)
	if err != nil {
		return err
	}

	return
}
