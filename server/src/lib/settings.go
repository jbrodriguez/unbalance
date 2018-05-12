package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/namsral/flag"
)

// ReservedSpace -
const ReservedSpace uint64 = 512 * 1024 * 1024 // 512Mb

// Config -
type Config struct {
	DryRun         bool     `json:"dryRun"`
	NotifyPlan     int      `json:"notifyPlan"`
	NotifyTransfer int      `json:"notifyTransfer"`
	ReservedAmount uint64   `json:"reservedAmount"`
	ReservedUnit   string   `json:"reservedUnit"`
	RsyncArgs      []string `json:"rsyncArgs"`
	Version        string   `json:"version"`
	Verbosity      int      `json:"verbosity"`
	CheckForUpdate int      `json:"checkForUpdate"`
	RefreshRate    int      `json:"refreshRate"`
}

// NotifyPlan/NotifyTransfer possible values
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
	var notifyCalc, notifyMove, notifyPlan, notifyTransfer, verbosity, checkForUpdate, refreshRate int

	flagset := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)

	flagset.StringVar(&port, "port", "6237", "port to run the server")
	flagset.StringVar(&logDir, "logdir", "/boot/logs", "pathname where log file will be written to")
	flagset.StringVar(&folders, "folders", "", "deprecated - do not use")
	flagset.BoolVar(&dryRun, "dryRun", true, "perform a dry-run rather than actual work")
	flagset.IntVar(&notifyCalc, "notifyCalc", 0, "deprecated - do not use") // deprecated
	flagset.IntVar(&notifyMove, "notifyMove", 0, "deprecated - do not use") // deprecated
	flagset.IntVar(&notifyPlan, "notifyPlan", 0, "notify via email after plan operation has completed (unraid notifications must be set up first): 0 - No notifications; 1 - Simple notifications; 2 - Detailed notifications")
	flagset.IntVar(&notifyTransfer, "notifyTransfer", 0, "notify via email after transfer operation has completed (unraid notifications must be set up first): 0 - No notifications; 1 - Simple notifications; 2 - Detailed notifications")
	flagset.StringVar(&rsyncFlags, "rsyncFlags", "", "deprecated - do not use") // deprecated
	flagset.StringVar(&rsyncArgs, "rsyncArgs", "-X", "custom rsync arguments")
	flagset.StringVar(&apiFolders, "apiFolders", "/var/local/emhttp", "folders to look for api endpoints")
	flagset.IntVar(&verbosity, "verbosity", 0, "include rsync output in log files: 0 (default) - include; 1 - do not include")
	flagset.IntVar(&checkForUpdate, "checkForUpdate", 1, "checkForUpdate: 0 - dont' check; 1 (default) - check")
	flagset.IntVar(&refreshRate, "refreshRate", 250, "how often to refresh the ui while running a command (in milliseconds)")

	location := SearchFile(name, locations)
	if location != "" {
		flagset.String("config", filepath.Join(location, name), "config location")
	}

	if err := flagset.Parse(os.Args[1:]); err != nil {
		return nil, err
	}

	s := &Settings{}

	if rsyncArgs == "" {
		s.RsyncArgs = make([]string, 0)
	} else {
		s.RsyncArgs = strings.Split(rsyncArgs, "|")
	}

	s.DryRun = dryRun
	s.NotifyPlan = notifyPlan
	s.NotifyTransfer = notifyTransfer
	s.ReservedAmount = ReservedSpace / 1024 / 1024
	s.ReservedUnit = "Mb"
	s.Verbosity = verbosity
	s.CheckForUpdate = checkForUpdate
	s.RefreshRate = refreshRate
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

	if err = WriteLine(tmpFile, fmt.Sprintf("notifyPlan=%d", s.NotifyPlan)); err != nil {
		return err
	}

	if err = WriteLine(tmpFile, fmt.Sprintf("notifyTransfer=%d", s.NotifyTransfer)); err != nil {
		return err
	}

	rsyncArgs := strings.Join(s.RsyncArgs, "|")
	if err = WriteLine(tmpFile, fmt.Sprintf("rsyncArgs=%s", rsyncArgs)); err != nil {
		return err
	}

	if err = WriteLine(tmpFile, fmt.Sprintf("verbosity=%d", s.Verbosity)); err != nil {
		return err
	}

	if err = WriteLine(tmpFile, fmt.Sprintf("refreshRate=%d", s.RefreshRate)); err != nil {
		return err
	}

	err = os.Rename(tmpFile, confLocation)
	if err != nil {
		return err
	}

	return
}
