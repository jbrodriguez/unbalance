package helper

import (
	"bufio"
	"log"
	"os/exec"
)

type Callback func(line string, arg interface{})

func Shell(command string, callback Callback, arg interface{}) {
	cmd := exec.Command("/bin/sh", "-c", command)
	out, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatalf("Unable to stdoutpipe %s: %s", command, err)
	}

	scanner := bufio.NewScanner(out)

	if err := cmd.Start(); err != nil {
		log.Fatal("Unable to start command: ", err)
	}

	for scanner.Scan() {
		callback(scanner.Text(), arg)
	}

	// Wait for the result of the command; also closes our end of the pipe
	err = cmd.Wait()
	if err != nil {
		log.Fatal("Unable to wait for process to finish: ", err)
	}
}
