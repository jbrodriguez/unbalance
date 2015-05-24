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
	SsmtpConf     string = "/etc/ssmtp/ssmtp.conf"
	dockerEnvS    string = "UNBALANCE_DOCKER"
	reservedSpace int64  = 450000000
)

type Config struct {
	ReservedSpace int64    `json:"reservedSpace"`
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
	Version       string   `json:"version"`
}

type Settings struct {
	ConfigDir      string
	LogDir         string
	CurrentVersion string

	Config
}

func (s *Settings) Init(version string, config string, log string) {
	s.Version = version
	s.CurrentVersion = version

	s.ConfigDir = config
	s.LogDir = log

	s.setup()
	s.load()

	s.Save()

	mlog.Info("Config file loaded from (%s) as:\n%s", s.ConfigDir, s.toString())
}

func (s *Settings) setup() {
	if _, err := os.Stat(s.ConfigDir); os.IsNotExist(err) {
		if err = os.MkdirAll(s.ConfigDir, 0755); err != nil {
			mlog.Fatalf("Unable to create folder %s: %s", s.ConfigDir, err)
		}
	}
}

func (s *Settings) load() {
	path := filepath.Join(s.ConfigDir, "config.json")
	file, err := os.Open(path)
	if err != nil {
		mlog.Warning("Config file %s doesn't exist. Creating one ...", path)

		s.Save()

		return
	}
	defer file.Close()

	decoder := json.NewDecoder(file)

	config := Config{}
	err = decoder.Decode(&config)
	if err != nil {
		mlog.Fatalf("Unable to load configuration: %s", err)
	}

	s.Config = config
}

func (s *Settings) sanitize() {
	s.ReservedSpace = reservedSpace
	s.Folders = getArray(s.Folders, make([]string, 0))
	s.DryRun = s.DryRun || false
	s.Notifications = s.Notifications || false
	s.NotiFrom = getString(s.NotiFrom, "myaccount@gmail.com")
	s.NotiTo = getString(s.NotiTo, "myaccount@gmail.com")
	s.NotiHost = getString(s.NotiHost, "smtp.gmail.com")
	s.NotiPort = getUint(s.NotiPort, 465)
	s.NotiEncrypt = s.NotiEncrypt && true
	s.NotiUser = getString(s.NotiUser, "myaccount")
	s.NotiPassword = getString(s.NotiPassword, "mypass")
	s.Version = s.CurrentVersion
}

func (s *Settings) toString() string {
	return fmt.Sprintf("ReservedSpace=%d\nFolders=%+v\nDryRun=%+v\nNotifications=%+v\nFrom=%s\nTo=%s\nHost=%s\nPort=%d\nEncrypt=%+v\nUser=%s", s.ReservedSpace, s.Folders, s.DryRun, s.Notifications, s.NotiFrom, s.NotiTo, s.NotiHost, s.NotiPort, s.NotiEncrypt, s.NotiUser)
}

func (s *Settings) Save() {
	s.sanitize()

	b, err := json.MarshalIndent(s, "", "   ")
	if err != nil {
		mlog.Info("couldn't marshal: %s", err)
		return
	}

	path := filepath.Join(s.ConfigDir, "config.json")
	err = ioutil.WriteFile(path, b, 0644)
	if err != nil {
		mlog.Info("WriteFileJson ERROR: %+v", err)
	}

	s.saveSsmtpConf()

	mlog.Info("Config file saved as: \n%s", s.toString())
}

func (s *Settings) saveSsmtpConf() {
	if os.Getenv(dockerEnvS) != "y" {
		return
	}

	f, err := os.Create(SsmtpConf)
	if err != nil {
		mlog.Error(err)
		return
	}
	defer f.Close()

	if _, err := f.WriteString(fmt.Sprintf("Root=%s\n", s.NotiFrom)); err != nil {
		mlog.Error(err)
		return
	}

	idx := strings.Index(s.NotiFrom, "@")
	if idx == -1 {
		mlog.Warning("Unable to find @ in email From field")
		return
	}

	if _, err := f.WriteString(fmt.Sprintf("rewriteDomain=%s\n", s.NotiFrom[idx+1:])); err != nil {
		mlog.Error(err)
		return
	}

	if _, err := f.WriteString("FromLineOverride=YES\n"); err != nil {
		mlog.Error(err)
		return
	}

	if _, err := f.WriteString(fmt.Sprintf("Mailhub=%s:%d\n", s.NotiHost, s.NotiPort)); err != nil {
		mlog.Error(err)
		return
	}

	var tls string
	if s.NotiEncrypt {
		tls = "YES"
	} else {
		tls = "NO"
	}

	if _, err := f.WriteString(fmt.Sprintf("UseTLS=%s\n", tls)); err != nil {
		mlog.Error(err)
		return
	}

	if _, err := f.WriteString("UseSTARTTLS=NO\n"); err != nil {
		mlog.Error(err)
		return
	}

	if _, err := f.WriteString("AuthMethod=login\n"); err != nil {
		mlog.Error(err)
		return
	}

	if _, err := f.WriteString(fmt.Sprintf("AuthUser=%s\n", s.NotiUser)); err != nil {
		mlog.Error(err)
	}

	if _, err := f.WriteString(fmt.Sprintf("AuthPass=%s\n", s.NotiPassword)); err != nil {
		mlog.Error(err)
		return
	}

	f.Sync()
}

func getString(val, def string) string {
	if val != "" {
		return val
	}

	return def
}

func getUint(val, def uint) uint {
	if val != 0 {
		return val
	}

	return def
}

func getArray(val, def []string) []string {
	if val != nil {
		return val
	}

	return def
}
