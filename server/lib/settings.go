package lib

import (
	// "errors"
	"fmt"
	"github.com/namsral/flag"
	"io"
	"os"
	"strings"
)

const RESERVED_SPACE = 450000000 // 450Mb

type Config struct {
	Folders    []string `json:"folders"`
	DryRun     bool     `json:"dryRun"`
	NotifyCalc int      `json:"notifyCalc"`
	NotifyMove int      `json:"notifyMove"`
	Version    string   `json:"version"`
}

// NotifyCalc/NotifyMove possible values
// 0 - no notification
// 1 - simple notification
// 2 - detailed notification

type Settings struct {
	Config

	Conf            string
	Log             string
	RunningInDocker bool
	ReservedSpace   int64
}

func NewSettings(version string) (*Settings, error) {
	var config, log, folders string
	var dryRun, runningInDocker bool
	var notifyCalc, notifyMove int

	// location := SearchFile(name, locations)
	// if location != "" {
	// 	flag.Set("config", filepath.Join(location, name))
	// }

	// /boot/config/plugins/unbalance/
	flag.StringVar(&config, "config", "/boot/config/plugins/unbalance/unbalance.conf", "config location")
	flag.StringVar(&log, "log", "", "folder where log file will be written to")
	flag.StringVar(&folders, "folders", "", "folders that will be scanned for media")
	flag.BoolVar(&dryRun, "dryRun", true, "perform a dry-run rather than actual work")
	flag.IntVar(&notifyCalc, "notifyCalc", 0, "notify via email after calculation operation has completed (unraid notifications must be set up first): 0 - No notifications; 1 - Simple notifications; 2 - Detailed notifications")
	flag.IntVar(&notifyMove, "notifyMove", 0, "notify via email after move operation has completed (unraid notifications must be set up first): 0 - No notifications; 1 - Simple notifications; 2 - Detailed notifications")
	flag.BoolVar(&runningInDocker, "docker", false, "notify via email after move operation has completed (unraid notifications must be set up first)")

	flag.Parse()

	// fmt.Printf("mediaFolders: %s\n", mediaFolders)

	s := &Settings{}
	if folders == "" {
		s.Folders = make([]string, 0)
	} else {
		s.Folders = strings.Split(folders, "|")
	}
	s.DryRun = dryRun
	s.NotifyCalc = notifyCalc
	s.NotifyMove = notifyMove
	s.Version = version

	s.Conf = config
	s.Log = log
	s.RunningInDocker = runningInDocker
	s.ReservedSpace = RESERVED_SPACE

	return s, nil
}

func (s *Settings) Save() error {
	file, err := os.Create(s.Conf)
	defer file.Close()

	if err != nil {
		return err
	}

	// if err = writeLine(file, "datadir", s.DataDir); err != nil {
	// 	return err
	// }

	// if err = writeLine(file, "webdir", s.WebDir); err != nil {
	// 	return err
	// }

	// if err = writeLine(file, "logdir", s.LogDir); err != nil {
	// 	return err
	// }

	// mediaFolders := strings.Join(s.MediaFolders, "|")
	// if err = writeLine(file, "mediafolders", mediaFolders); err != nil {
	// 	return err
	// }

	return nil
}

func writeLine(file *os.File, key, value string) error {
	_, err := io.WriteString(file, fmt.Sprintf("%s=%s\n", key, value))
	if err != nil {
		return err
	}

	return nil
}
