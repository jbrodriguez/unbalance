package model

import (
	"encoding/json"
	"fmt"
	"github.com/jbrodriguez/mlog"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	SsmtpConf string = "/usr/local/etc/ssmtp.conf"
)

type Config struct {
	ReservedSpace uint64   `json:"reservedSpace"`
	Folders       []string `json:"folders"`
	DryRun        bool     `json:"dryRun"`
	Notifications bool     `json:"notifications"`
	NotiFrom      string   `json:"notiFrom"`
	NotiTo        string   `json:"notiTo"`
	NotiHost      string   `json:"notiHost"`
	NotiPort      uint     `json:"notiPort"`
	NotiEncrypt   bool     `json:"notiEncrypt"`
	NotiUser      string   `json:"notiUser"`
	NotiPassword  string   `json:"notiPassword"`

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

		c.ReservedSpace = 350000000
		c.Folders = make([]string, 0)
		c.DryRun = true
		c.Notifications = false
		c.NotiFrom = "myaccount@gmail.com"
		c.NotiTo = "myaccount@gmail.com"
		c.NotiHost = "smtp.gmail.com"
		c.NotiPort = 465
		c.NotiEncrypt = true
		c.NotiUser = "myaccount"
		c.NotiPassword = "mypass"

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

	c.ReservedSpace = 350000000
	c.Folders = config.Folders
	c.DryRun = config.DryRun
	c.Notifications = config.Notifications
	c.NotiFrom = config.NotiFrom
	c.NotiTo = config.NotiTo
	c.NotiHost = config.NotiHost
	c.NotiPort = config.NotiPort
	c.NotiEncrypt = config.NotiEncrypt
	c.NotiUser = config.NotiUser
	c.NotiPassword = config.NotiPassword
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

	if c.Notifications && c.NotiPassword != "mypass" {
		c.saveSsmtpConf()
	}

	mlog.Info("saved as: \nReservedSpace=%d\nFolders=%+v\nDryRun=%+v\nNotifications=%+v\nFrom=%s\nTo=%s\nHost=%s\nPort=%d\nEncrypt=%+v\nUser=%s", c.ReservedSpace, c.Folders, c.DryRun, c.Notifications, c.NotiFrom, c.NotiTo, c.NotiHost, c.NotiPort, c.NotiEncrypt, c.NotiUser)
}

func (c *Config) saveSsmtpConf() {
	f, err := os.Create(SsmtpConf)
	if err != nil {
		mlog.Error(err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("Root=%s", c.NotiFrom)); err != nil {
		mlog.Error(err)
		return
	}

	idx := strings.Index(c.NotiFrom, "@")
	if idx == -1 {
		mlog.Warning("Unable to find @ in email From field")
		return
	}

	if _, err := f.WriteString(fmt.Sprintf("rewriteDomain=%s", c.NotiFrom[idx+1:])); err != nil {
		mlog.Error(err)
		return
	}

	if _, err := f.WriteString("FromLineOverride=YES"); err != nil {
		mlog.Error(err)
		return
	}

	if _, err := f.WriteString(fmt.Sprintf("Mailhub=%s:%d", c.NotiHost, c.NotiPort)); err != nil {
		mlog.Error(err)
		return
	}

	var tls string
	if c.NotiEncrypt {
		tls = "YES"
	} else {
		tls = "NO"
	}

	if _, err := f.WriteString(fmt.Sprintf("UseTLS=%s", tls)); err != nil {
		mlog.Error(err)
		return
	}

	if _, err := f.WriteString("UseSTARTTLS=NO"); err != nil {
		mlog.Error(err)
		return
	}

	if _, err := f.WriteString("AuthMethod=login"); err != nil {
		mlog.Error(err)
		return
	}

	if _, err := f.WriteString(fmt.Sprintf("AuthUser=%s", c.NotiUser)); err != nil {
		mlog.Error(err)
	}

	if _, err := f.WriteString(fmt.Sprintf("AuthPass=%s", c.NotiPassword)); err != nil {
		mlog.Error(err)
		return
	}

	f.Sync()
}
