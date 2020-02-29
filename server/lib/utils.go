package lib

import (
	"fmt"
	"io"
	"io/ioutil"
	"math"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
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
func ByteSize(bytes int64) string {
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
func Max(x, y int64) int64 {
	if x > y {
		return x
	}
	return y
}

// Min -
func Min(x, y int64) int64 {
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

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}

// Sendmail -
func Sendmail(cmd string, notify int, subject, message string, dryRun bool) (err error) {
	if notify == 0 {
		return nil
	}

	dry := ""
	if dryRun {
		dry = "-------\nDRY RUN\n-------\n"
	}

	msg := dry + message

	send := exec.Command(cmd, "-e", "unBALANCE operation update", "-s", subject, "-m", msg)
	err = send.Run()

	return
}
