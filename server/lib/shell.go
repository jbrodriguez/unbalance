package lib

import (
	"bufio"
	"bytes"
	// "errors"
	// "log"
	"io"
	"os/exec"
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

	for {
		var line string
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

func Shell(command string, writer StderrWriter, prefix string, callback Callback) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
		//		log.Fatalf("Unable to stdoutpipe %s: %s", command, err)
	}

	cmd.Stderr = NewStreamer(writer, prefix)

	// stderr, err := cmd.StderrPipe()
	// if err != nil {
	// 	return err
	// 	// log.Fatalf("Unable to stderrpipe %s: %s", command, err)
	// }

	// multi := io.MultiReader(stdout, stderr)

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
		return err
		// log.Fatal("Unable to wait for process to finish: ", err)
	}

	return nil
}

func Pipeline(cmds ...*exec.Cmd) (pipeLineOutput, collectedStandardError []byte, pipeLineError error) {
	// Require at least one command
	if len(cmds) < 1 {
		return nil, nil, nil
	}

	// Collect the output from the command(s)
	var output bytes.Buffer
	var stderr bytes.Buffer

	last := len(cmds) - 1
	for i, cmd := range cmds[:last] {
		var err error
		// Connect each command's stdin to the previous command's stdout
		if cmds[i+1].Stdin, err = cmd.StdoutPipe(); err != nil {
			return nil, nil, err
		}
		// Connect each command's stderr to a buffer
		cmd.Stderr = &stderr
	}

	// Connect the output and error for the last command
	cmds[last].Stdout, cmds[last].Stderr = &output, &stderr

	// Start each command
	for _, cmd := range cmds {
		if err := cmd.Start(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Wait for each command to complete
	for _, cmd := range cmds {
		if err := cmd.Wait(); err != nil {
			return output.Bytes(), stderr.Bytes(), err
		}
	}

	// Return the pipeline output and the collected standard error
	return output.Bytes(), stderr.Bytes(), nil
}
