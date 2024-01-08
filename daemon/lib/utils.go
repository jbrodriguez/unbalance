package lib

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"gopkg.in/ini.v1"

	"unbalance/daemon/domain"
)

// Exists - Check if File / Directory Exists
func Exists(path string) (bool, error) {
	_, err := os.Stat(path)

	if err == nil {
		return true, nil
	}

	if os.IsNotExist(err) {
		return false, nil
	}

	return false, err
}

// IsEmpty - checks if a folder is empty
func IsEmpty(folder string) (bool, error) {
	f, err := os.Open(folder)
	if err != nil {
		return false, err
	}
	defer func() { _ = f.Close() }()

	_, err = f.Readdirnames(1) // Or f.Readdir(1)
	if err == io.EOF {
		return true, nil
	}

	return false, err // Either not empty or error, suits both cases
}

// SearchFile -
func SearchFile(name string, locations []string) string {
	for _, location := range locations {
		if b, _ := Exists(filepath.Join(location, name)); b {
			return location
		}
	}

	return ""
}

var sizes = []string{"B", "KB", "MB", "GB", "TB", "PB", "EB", "ZB", "YB"}

// ByteSize -
func ByteSize(bytes uint64) string {
	if bytes == 0 {
		return "0B"
	}

	k := float64(1000)
	i := math.Floor(math.Log(float64(bytes)) / math.Log(k))

	return fmt.Sprintf("%.2f %s", float64(bytes)/math.Pow(k, i), sizes[int64(i)])
}

// WriteLine -
func WriteLine(fullpath, line string) error {
	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(line + "\n")

	return err
}

// WriteLines -
func WriteLines(fullpath string, lines []string) error {
	f, err := os.OpenFile(fullpath, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, line := range lines {
		_, err = f.WriteString(line + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Round -
func Round(d, r time.Duration) time.Duration {
	if r <= 0 {
		return d
	}
	neg := d < 0
	if neg {
		d = -d
	}
	if m := d % r; m+m < r {
		d -= d
	} else {
		d = d + r - m
	}
	if neg {
		return -d
	}
	return d
}

// Max -
func Max(x, y uint64) uint64 {
	if x > y {
		return x
	}
	return y
}

// Min -
func Min(x, y uint64) uint64 {
	if x < y {
		return x
	}
	return y
}

// GetLatestVersion  -
func GetLatestVersion(url string) (dst string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

func Bind(content any, data any) error {
	m := content.(map[string]interface{})
	s, err := json.Marshal(m)
	if err != nil {
		return err
	}

	err = json.Unmarshal(s, &data)
	if err != nil {
		return err
	}

	return nil
}

func LoadEnv(location string, config *domain.Config) error {
	// load file
	file, err := ini.Load(location)
	if err != nil {
		return err
	}

	// fill data
	config.DryRun, _ = file.Section("").Key("DRY_RUN").Bool()
	config.NotifyPlan, _ = file.Section("").Key("NOTIFY_PLAN").Int()
	config.NotifyTransfer, _ = file.Section("").Key("NOTIFY_TRANSFER").Int()
	config.ReservedAmount, _ = file.Section("").Key("RESERVED_AMOUNT").Uint64()
	config.ReservedUnit = file.Section("").Key("RESERVED_UNIT").String()
	config.RsyncArgs = file.Section("").Key("RSYNC_ARGS").Strings(",")
	config.Verbosity, _ = file.Section("").Key("VERBOSITY").Int()
	config.CheckForUpdate, _ = file.Section("").Key("CHECK_FOR_UPDATE").Int()
	config.RefreshRate, _ = file.Section("").Key("REFRESH_RATE").Int()

	return nil
}

func SaveEnv(location string, config domain.Config) error {
	// load file
	file, err := ini.Load(location)
	if err != nil {
		return err
	}

	ini.PrettyFormat = false

	// fill data
	file.Section("").Key("DRY_RUN").SetValue(strconv.FormatBool(config.DryRun))
	file.Section("").Key("NOTIFY_PLAN").SetValue(strconv.Itoa(config.NotifyPlan))
	file.Section("").Key("NOTIFY_TRANSFER").SetValue(strconv.Itoa(config.NotifyTransfer))
	file.Section("").Key("RESERVED_AMOUNT").SetValue(strconv.FormatUint(config.ReservedAmount, 10))
	file.Section("").Key("RESERVED_UNIT").SetValue(config.ReservedUnit)
	file.Section("").Key("RSYNC_ARGS").SetValue(strings.Join(config.RsyncArgs, ","))
	file.Section("").Key("VERBOSITY").SetValue(strconv.Itoa(config.Verbosity))
	file.Section("").Key("CHECK_FOR_UPDATE").SetValue(strconv.Itoa(config.CheckForUpdate))
	file.Section("").Key("REFRESH_RATE").SetValue(strconv.Itoa(config.RefreshRate))

	file.SaveTo(location)

	return nil
}
