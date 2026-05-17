package core

import (
	"os/exec"

	"unbalance/daemon/domain"
	"unbalance/daemon/lib"
)

type rsyncProcess interface {
	PID() int
	Kill() error
	Wait() (int, error)
}

type transferExecutor interface {
	PrepareCommand(command *domain.Command, allowedRoots []string) (safeTransferPaths, error)
	StartRsync(workDir string, args ...string) (rsyncProcess, error)
	RemoveTransferredSource(command *domain.Command, pruneParents bool) (string, []string, error)
}

type inProcessExecutor struct{}

func newInProcessExecutor() transferExecutor {
	return inProcessExecutor{}
}

func (inProcessExecutor) PrepareCommand(command *domain.Command, allowedRoots []string) (safeTransferPaths, error) {
	return validateCommandForExecution(command, allowedRoots)
}

func (inProcessExecutor) StartRsync(workDir string, args ...string) (rsyncProcess, error) {
	cmd, err := lib.StartRsync(workDir, args...)
	if err != nil {
		return nil, err
	}

	return execRsyncProcess{cmd: cmd}, nil
}

func (inProcessExecutor) RemoveTransferredSource(command *domain.Command, pruneParents bool) (string, []string, error) {
	return removeTransferredSource(command, pruneParents)
}

type execRsyncProcess struct {
	cmd *exec.Cmd
}

func (p execRsyncProcess) PID() int {
	return p.cmd.Process.Pid
}

func (p execRsyncProcess) Kill() error {
	return lib.KillRsync(p.cmd)
}

func (p execRsyncProcess) Wait() (int, error) {
	return lib.EndRsync(p.cmd)
}
