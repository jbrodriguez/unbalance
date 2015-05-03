package model

import (
	"encoding/json"
	"github.com/jbrodriguez/mlog"
	"io/ioutil"
	"os"
	"path/filepath"
)

type Config struct {
	ReservedSpace uint64   `json:"reservedSpace"`
	Folders       []string `json:"folders"`
	DryRun        bool     `json:"dryRun"`

	ConfigDir string `json:"-"`
	LogDir    string `json:"-"`

	Version string `json:"version"`
}

func (c *Config) Init(version string, config string, log string) {
	c.Version = version

	c.ConfigDir = config
	c.LogDir = os.Getenv("UNBALANCE_LOGFILEPATH")

	if log != "" {
		c.LogDir = log
	}

	// os.Setenv("GIN_MODE", "release")

	if c.LogDir != "" {
		mlog.Start(mlog.LevelInfo, filepath.Join(c.LogDir, "unbalance.log"))
	} else {
		mlog.Start(mlog.LevelInfo, "")
	}

	//	mlog.Info("logfilePath: %s (%s)", c.LogDir, log)

	c.setupOperatingEnv()

	c.LoadConfig()
}

func (c *Config) setupOperatingEnv() {
	if _, err := os.Stat(c.ConfigDir); os.IsNotExist(err) {
		if err = os.MkdirAll(c.ConfigDir, 0755); err != nil {
			mlog.Fatalf("Unable to create folder %s: %s", c.ConfigDir, err)
		}
	}
}

func (c *Config) LoadConfig() {
	path := filepath.Join(c.ConfigDir, "config.json")
	file, err := os.Open(path)
	if err != nil {
		mlog.Warning("Config file %s doesn't exist. Creating one ...", path)

		c.ReservedSpace = 250000000
		c.Folders = make([]string, 0)
		c.DryRun = true

		c.Save()

		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		mlog.Fatalf("Unable to load configuration: %s", err)
	}

	c.ReservedSpace = config.ReservedSpace
	c.Folders = config.Folders
	c.DryRun = config.DryRun
}

func (c *Config) Save() {
	b, err := json.MarshalIndent(c, "", "   ")
	if err != nil {
		mlog.Info("couldn't marshal: %s", err)
		return
	}

	path := filepath.Join(c.ConfigDir, "config.json")
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		mlog.Info("WriteFileJson ERROR: %+v", err)
	}

	mlog.Info("saved as: %s", string(b))
}
