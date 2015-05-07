package helper

import (
	"bufio"
	// "errors"
	// "log"
	"io"
	"os/exec"
)

type Callback func(line string, arg interface{})

func Shell(command string, callback Callback, arg interface{}) error {
	cmd := exec.Command("/bin/sh", "-c", command)
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
		//		log.Fatalf("Unable to stdoutpipe %s: %s", command, err)
	}

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
		// log.Fatalf("Unable to stderrpipe %s: %s", command, err)
	}

	multi := io.MultiReader(stdout, stderr)

	scanner := bufio.NewScanner(multi)

	if err := cmd.Start(); err != nil {
		return err
		// log.Fatal("Unable to start command: ", err)
	}

	for scanner.Scan() {
		callback(scanner.Text(), arg)
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		return err
		// log.Fatal("Unable to wait for process to finish: ", err)
	}

	return nil
}
