package lib

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"os/exec"
	"syscall"
)

const cRsyncBin = "/usr/bin/rsync"

// maxLineSize bounds a single output line; anything beyond aborts the scan
// instead of deadlocking the pipe (see runCmd)
const maxLineSize = 1024 * 1024

// Callback -
type Callback func(line string)

// runCmd streams the command's stdout line by line into callback. If the
// scanner aborts (line too long, read error), the pipe keeps getting drained
// so the child can exit and Wait can never deadlock against a full pipe.
func runCmd(cmd *exec.Cmd, callback Callback) error {
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(stdout)
	scanner.Buffer(make([]byte, 0, 64*1024), maxLineSize)

	if err = cmd.Start(); err != nil {
		return err
	}

	for scanner.Scan() {
		callback(scanner.Text())
	}

	scanErr := scanner.Err()
	if scanErr != nil {
		go func() { _, _ = io.Copy(io.Discard, stdout) }()
	}

	return errors.Join(scanErr, cmd.Wait())
}

func Shell(command string, workDir string, callback Callback) error {
	return ShellContext(context.Background(), command, workDir, callback)
}

// ShellContext is Shell with cancellation: the child is killed when ctx is
// cancelled or times out.
func ShellContext(ctx context.Context, command string, workDir string, callback Callback) error {
	cmd := exec.CommandContext(ctx, "/bin/sh", "-c", command)
	if workDir != "" {
		cmd.Dir = workDir
	}
	cmd.Stderr = os.Stdout

	return runCmd(cmd, callback)
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

// Stream runs name with args directly (no shell), invoking callback for each
// stdout line. It is the no-shell replacement for the legacy Shell2 helper:
// arguments are passed as argv tokens, so caller-supplied paths cannot be
// interpreted as shell metacharacters.
func Stream(name string, args []string, callback Callback) error {
	return StreamContext(context.Background(), name, args, callback)
}

// StreamContext is Stream with cancellation: the child is killed when ctx is
// cancelled or times out.
func StreamContext(ctx context.Context, name string, args []string, callback Callback) error {
	return runCmd(exec.CommandContext(ctx, name, args...), callback)
}

// StartRsync -
func StartRsync(workDir string, args ...string) (*exec.Cmd, error) {
	cmd := exec.Command(cRsyncBin, args...)
	cmd.Dir = workDir
	cmd.Stderr = os.Stdout

	err := cmd.Start()
	return cmd, err
}

// EndRsync -
func EndRsync(cmd *exec.Cmd) (int, error) {
	exitCode := 0

	err := cmd.Wait()
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				exitCode = waitStatus.ExitStatus()
			}
		}
	}

	return exitCode, err
}

// KillRsync -
func KillRsync(cmd *exec.Cmd) error {
	return cmd.Process.Kill()
}
