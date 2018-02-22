package lib

import (
	"bufio"
	"bytes"
	"io"
	"os/exec"
)

const cRsyncBin = "/usr/bin/rsync"

// Callback -
type Callback func(line string)

// StderrWriter -
type StderrWriter func(format string, a ...interface{})

// Streamer -
type Streamer struct {
	buf    *bytes.Buffer
	writer StderrWriter
	prefix string
}

// NewStreamer -
func NewStreamer(writer StderrWriter, prefix string) *Streamer {
	return &Streamer{
		buf:    bytes.NewBuffer([]byte("")),
		writer: writer,
		prefix: prefix,
	}
}

func (s *Streamer) Write(p []byte) (n int, err error) {
	if n, err = s.buf.Write(p); err != nil {
		return
	}

	var line string
	for {
		line, err = s.buf.ReadString('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			return
		}

		// l.readLines += line
		s.writer("%s:(%s)", s.prefix, line)
	}

	return
}

// Shell -
func Shell(command string, writer StderrWriter, prefix, workDir string, callback Callback) error {
	args := []string{
		"-c",
	}
	args = append(args, command)

	return shell(writer, prefix, workDir, callback, "/bin/sh", args...)
}

// writer: mlog.Writer
// prefix: prefix for each log line
// callback: invoked on each output line
// name: command name
// args: command arguments
func shell(writer StderrWriter, prefix, workDir string, callback Callback, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	// cmd.Env = os.Environ()
	if workDir != "" {
		cmd.Dir = workDir
	}
	cmd.Stderr = NewStreamer(writer, prefix)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
		//		log.Fatalf("Unable to stdoutpipe %s: %s", command, err)
	}
	scanner := bufio.NewScanner(stdout)

	if err = cmd.Start(); err != nil {
		return err
		// log.Fatal("Unable to start command: ", err)
	}

	for scanner.Scan() {
		callback(scanner.Text())
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		writer("%s: waitError: %s", prefix, err)
		return err
		// log.Fatal("Unable to wait for process to finish: ", err)
	}

	return nil
}

// dropCR drops a terminal \r from the data.
func dropCR(data []byte) []byte {
	if len(data) > 0 && data[len(data)-1] == '\r' {
		return data[0 : len(data)-1]
	}
	return data
}

// scanLinesEx is a split function for a Scanner that returns each line of
// text, stripped of any trailing end-of-line marker. The returned line may
// be empty. The end-of-line marker is one optional carriage return followed
// by one mandatory newline. In regular expression notation, it is `\r?\n`.
// The last non-empty line of input will be returned even if it has no
// newline.
func scanLinesEx(data []byte, atEOF bool) (advance int, token []byte, err error) {
	if atEOF && len(data) == 0 {
		return 0, nil, nil
	}

	if i := bytes.IndexByte(data, '\r'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}

	if i := bytes.IndexByte(data, '\n'); i >= 0 {
		// We have a full newline-terminated line.
		return i + 1, dropCR(data[0:i]), nil
	}

	// If we're at EOF, we have a final, non-terminated line. Return it.
	if atEOF {
		return len(data), dropCR(data), nil
	}

	// Request more data.
	return 0, nil, nil
}

// ShellEx -
func ShellEx(callback Callback, writer StderrWriter, workDir, name string, args ...string) error {
	cmd := exec.Command(name, args...)

	if workDir != "" {
		cmd.Dir = workDir
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
		//		log.Fatalf("Unable to stdoutpipe %s: %s", command, err)
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Split(scanLinesEx)

	if err = cmd.Start(); err != nil {
		return err
		// log.Fatal("Unable to start command: ", err)
	}

	for scanner.Scan() {
		callback(scanner.Text())
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		writer("rsync:exit:(%s)", cmd.ProcessState.String())
		return err
	}

	return nil
}

// Shell2 -
func Shell2(command string, callback Callback) error {
	args := append([]string{"-c"}, command)
	cmd := exec.Command("/bin/sh", args...)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)

	if err = cmd.Start(); err != nil {
		return err
	}

	for scanner.Scan() {
		callback(scanner.Text())
	}

	// Wait for the result of the command; also closes our end of the pipe
	return cmd.Wait()
}

// StartRsync -
func StartRsync(workDir string, writer StderrWriter, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(cRsyncBin, args...)
	cmd.Dir = workDir
	cmd.Stderr = NewStreamer(writer, "rsync")

	return cmd, cmd.Start()
}

// EndRsync -
func EndRsync(cmd *exec.Cmd) error {
	return cmd.Wait()
}
