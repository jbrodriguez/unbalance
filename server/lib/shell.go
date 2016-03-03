package lib

import (
	"bufio"
	"bytes"
	// "errors"
	// "log"
	"io"
	"os/exec"
	"syscall"
)

type Callback func(line string)
type StderrWriter func(format string, a ...interface{})

type Streamer struct {
	buf    *bytes.Buffer
	writer StderrWriter
	prefix string
}

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
		s.writer("%s: %s", s.prefix, line)
	}

	return
}

func ShellEx(writer StderrWriter, prefix string, callback Callback, name string, args ...string) error {
	return shell(writer, prefix, callback, name, args...)
}

func Shell(command string, writer StderrWriter, prefix string, callback Callback) error {
	args := []string{
		"-c",
	}
	args = append(args, command)

	return shell(writer, prefix, callback, "/bin/sh", args...)
}

func shell(writer StderrWriter, prefix string, callback Callback, name string, args ...string) error {
	cmd := exec.Command(name, args...)

	cmd.Stderr = NewStreamer(writer, prefix)

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
		//		log.Fatalf("Unable to stdoutpipe %s: %s", command, err)
	}
	scanner := bufio.NewScanner(stdout)

	if err := cmd.Start(); err != nil {
		return err
		// log.Fatal("Unable to start command: ", err)
	}

	for scanner.Scan() {
		callback(scanner.Text())
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		var waitStatus syscall.WaitStatus
		if exiterr, ok := err.(*exec.ExitError); ok {
			waitStatus = exiterr.Sys().(syscall.WaitStatus)
			writer("%s:waitError:Status(%d):Err(%s):ExitErr(%s)", prefix, waitStatus.ExitStatus(), err, exiterr)
		} else {
			writer("%s:waitError:(%s)", prefix, err)
		}

		return err
		// log.Fatal("Unable to wait for process to finish: ", err)
	}

	return nil
}
