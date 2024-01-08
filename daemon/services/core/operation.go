package core

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"unbalance/daemon/common"
	"unbalance/daemon/domain"
	"unbalance/daemon/lib"
	"unbalance/daemon/logger"

	"github.com/teris-io/shortid"
)

func (c *Core) runOperation(opName string) {
	logger.Blue("running %s operation ...", opName)

	operation := c.state.Operation
	operation.Started = time.Now()

	if c.ctx.NotifyTransfer == 2 {
		c.notifyCommandsToRun(opName, operation)
	}

	operation.Line = "Waiting to collect stats ..."
	packet := &domain.Packet{Topic: common.EventTransferStarted, Payload: operation}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	commandsExecuted := make([]string, 0)

	for _, command := range operation.Commands {
		args := append(
			operation.RsyncArgs,
			command.Entry,
			command.Dst,
		)

		cmd := fmt.Sprintf(`rsync %s %s %s`, operation.RsyncStrArgs, strconv.Quote(command.Entry), strconv.Quote(command.Dst))
		logger.Blue("command started: (src: %s) %s ", command.Src, cmd)

		operation.Line = cmd
		packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
		c.ctx.Hub.Pub(packet, "socket:broadcast")

		cmdTransferred, err := c.runCommand(operation, command, args)
		if err != nil {
			c.commandInterrupted(opName, operation, command, cmd, err, cmdTransferred, commandsExecuted)
			return
		}

		c.commandCompleted(operation, command)

		commandsExecuted = append(commandsExecuted, cmd)
	}

	c.operationCompleted(opName, operation, commandsExecuted)
}

func (c *Core) runCommand(operation *domain.Operation, command *domain.Command, args []string) (uint64, error) {
	// make sure the command will run
	c.stopped = false

	// start rsync command
	cmd, err := lib.StartRsync(command.Src, args...)
	if err != nil {
		return 0, err
	}

	// give some time for /proc/pid to come alive
	time.Sleep(500 * time.Millisecond)

	// monitor rsync progress
	retcode, transferred := c.monitorRsync(operation, command, cmd.Process.Pid)

	// 99 is a magic number meaning the user clicked on the stop operation button
	if retcode == 99 {
		logger.Blue("command:received:StopCommand")
		err = lib.KillRsync(cmd)
		if err != nil {
			logger.Yellow("command:kill:error(%s)", err)
		}

		command.Status = common.CmdStopped
		return transferred, fmt.Errorf("exit status %d", retcode)
	}

	command.Status = common.CmdCompleted

	// end rsync process
	exitCode, err := lib.EndRsync(cmd)
	if err != nil {
		logger.Yellow("command:end:error(%s)", err)
		command.Status = common.CmdStopped
	}

	logger.Blue("command:retcode(%d):exitcode(%d)", retcode, exitCode)

	// 23: "Partial transfer due to error"
	// 13: "Errors with program diagnostics"
	if exitCode == 23 || exitCode == 13 {
		err = nil
		command.Status = common.CmdFlagged
	}

	return transferred, err
}

func (c *Core) monitorRsync(operation *domain.Operation, command *domain.Command, procPid int) (int, uint64) {
	var transferred uint64
	var current string
	var retcode int
	var zombie bool
	var err error

	start := time.Now()
	display := true

	pid := strconv.Itoa(procPid)

	procStat := "/proc/" + pid + "/stat"
	procIo := "/proc/" + pid + "/io"
	procFd := "/proc/" + pid + "/fd/3"

	for {
		// c.stopped will be true if the user has stopped the command via the gui
		if c.stopped {
			retcode = 99
			break
		}

		// isZombie
		zombie, retcode, err = isZombie(procStat)
		if err != nil {
			logger.Yellow("isZombie:err(%s)", err)
			break
		}

		// pid finished processing, exit the loop and Wait on the pid, to finish it properly
		if zombie {
			break
		}

		// throttle both /proc consumption and messaging to the front end
		time.Sleep(time.Duration(c.ctx.RefreshRate) * time.Millisecond)

		// getReadBytes
		transferred, err = getReadBytes(procIo)
		if err != nil {
			logger.Yellow("getReadBytes:err(%s):xfer(%d)", err, transferred)
			break
		}

		// getCurrentTransfer
		current, err = getCurrentTransfer(procFd, filepath.Join(command.Src, command.Entry))
		if err != nil {
			continue
		}

		// on the first loop (and every 10 min after that, arbitrarily), display the file currently being transferred
		if display {
			logger.Blue("monitor:transfer:(%s)", current)
			start = time.Now()
		}
		display = time.Since(start) >= time.Minute*10

		// update progress stats
		transferred = lib.Min(transferred, command.Size)

		percent, left, speed := progress(operation.BytesToTransfer, operation.BytesTransferred+transferred, time.Since(operation.Started))

		operation.Line = current
		operation.Completed = percent
		operation.Speed = speed
		operation.Remaining = left.String()
		operation.DeltaTransfer = transferred
		command.Transferred = transferred
		command.Status = common.CmdInProgress

		packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
		c.ctx.Hub.Pub(packet, "socket:broadcast")
	}

	return retcode, transferred
}

func (c *Core) commandInterrupted(opName string, operation *domain.Operation, command *domain.Command, cmd string, err error, cmdTransferred uint64, commandsExecuted []string) {
	operation.Finished = time.Now()
	elapsed := time.Since(operation.Started)

	subject := fmt.Sprintf("unBALANCE - %s operation INTERRUPTED", strings.ToUpper(opName))
	headline := fmt.Sprintf("Command Interrupted: %s (%s)", cmd, err.Error()+" : "+getError(err.Error(), reRsync, rsyncErrors))

	logger.Yellow(headline)
	packet := &domain.Packet{Topic: common.EventOperationError, Payload: fmt.Sprintf("%s operation was interrupted. Check log (/boot/logs/unbalance.log) for details.", opName)}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	operation.BytesTransferred += cmdTransferred
	percent, left, speed := progress(operation.BytesToTransfer, operation.BytesTransferred, elapsed)

	operation.Completed = percent
	operation.Speed = speed
	operation.Remaining = left.String()
	operation.DeltaTransfer = cmdTransferred
	command.Transferred = cmdTransferred

	c.endOperation(subject, headline, commandsExecuted, operation)
}

func (c *Core) commandCompleted(operation *domain.Operation, command *domain.Command) {
	text := "Command Finished"
	logger.Blue(text)

	operation.BytesTransferred += command.Size
	percent, left, speed := progress(operation.BytesToTransfer, operation.BytesTransferred, time.Since(operation.Started))

	operation.Completed = percent
	operation.Speed = speed
	operation.Remaining = left.String()
	operation.DeltaTransfer = 0
	operation.Line = text
	command.Transferred = command.Size

	msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, operation.Remaining, speed)
	logger.Blue("Current progress: %s", msg)

	packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	// this is just a heads up for the user, shows which folders would/wouldn't be pruned if run without dry-run
	showPotentiallyPrunedItems(operation, command)

	// if it isn't a dry-run and the operation is Move or Gather, delete the source folder
	c.handleItemDeletion(operation, command)
}

func (c *Core) handleItemDeletion(operation *domain.Operation, command *domain.Command) {
	if !operation.DryRun && (operation.OpKind == common.OpScatterMove || operation.OpKind == common.OpGatherMove) {
		// the command was flagged due to an error, don't delete the source file/folder in these cases
		if command.Status == common.CmdFlagged {
			msg := fmt.Sprintf("skipping:deletion:(rsync command was flagged):(%s)", filepath.Join(command.Dst, command.Entry))
			operation.Line = msg

			packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
			c.ctx.Hub.Pub(packet, "socket:broadcast")
			logger.Yellow(msg)

			return
		}

		exists, _ := lib.Exists(filepath.Join(command.Dst, command.Entry))
		if exists {
			rmrf := fmt.Sprintf("rm -rf \"%s\"", filepath.Join(command.Src, command.Entry))
			operation.Line = fmt.Sprintf("Removing source %s", filepath.Join(command.Src, command.Entry))

			packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
			c.ctx.Hub.Pub(packet, "socket:broadcast")
			logger.Blue("removing:(%s)", rmrf)

			err := lib.Shell(rmrf, "", func(line string) {
				logger.Blue(line)
			})

			if err != nil {
				msg := fmt.Sprintf("Unable to remove source folder (%s): %s", filepath.Join(command.Src, command.Entry), err)
				operation.Line = msg

				packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
				c.ctx.Hub.Pub(packet, "socket:broadcast")

				logger.Yellow(msg)
			}

			if operation.OpKind == common.OpGatherMove {
				parent := filepath.Dir(command.Entry)
				// if entry is a user share (tvshows), Dir returns ".", so we won't touch it
				// if entry is a top-level children of a user share (tvshows/Billions), Dir returns "tvshows"
				// if entry is a nested children (tvshows/Billions/Season 01), Dir returns "tvshows/Billions"
				// in the first 2 cases no "/" is present in parent, so I won't prune them
				if strings.Contains(parent, "/") {
					rmdir := fmt.Sprintf(`find "%s" -type d -empty -prune -exec rm -rf {} \;`, filepath.Join(command.Src, parent))
					operation.Line = fmt.Sprintf("Pruning parent %s", filepath.Join(command.Src, parent))

					packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
					c.ctx.Hub.Pub(packet, "socket:broadcast")

					logger.Blue("pruning:(%s)", rmdir)

					err = lib.Shell(rmdir, "", func(line string) {
						logger.Blue(line)
					})

					if err != nil {
						msg := fmt.Sprintf("Unable to remove parent folder (%s): %s", filepath.Join(command.Src, parent), err)
						operation.Line = msg

						packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
						c.ctx.Hub.Pub(packet, "socket:broadcast")

						logger.Yellow(msg)
					}
				} else {
					logger.Yellow("skipping:prune:(%s)", filepath.Join(command.Src, parent))
				}
			}
		} else {
			logger.Yellow("skipping:deletion:(file/folder not present in destination):(%s)", filepath.Join(command.Dst, command.Entry))
		}
	}
}

func (c *Core) operationCompleted(opName string, operation *domain.Operation, commandsExecuted []string) {
	operation.Finished = time.Now()
	elapsed := operation.Finished.Sub(operation.Started)

	subject := fmt.Sprintf("unbalance - %s operation completed", strings.ToUpper(opName))
	headline := fmt.Sprintf("%s operation has finished", opName)

	percent, left, speed := progress(operation.BytesToTransfer, operation.BytesTransferred, elapsed)
	operation.Completed = percent
	operation.Speed = speed
	operation.Remaining = left.String()

	c.endOperation(subject, headline, commandsExecuted, operation)
}

func (c *Core) endOperation(subject, headline string, commands []string, operation *domain.Operation) {
	fstarted := operation.Started.Format(timeFormat)
	ffinished := operation.Finished.Format(timeFormat)
	elapsed := lib.Round(operation.Finished.Sub(operation.Started), time.Millisecond)

	c.state.Status = common.OpNeutral
	c.state.Operation = nil
	c.state.Unraid = c.refreshUnraid()
	// TODO: update history
	c.updateHistory(c.state.History, operation)

	packet := &domain.Packet{Topic: common.EventTransferEnded, Payload: c.state}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s\n\n%s\n\nTransferred %s at ~ %.2f MB/s",
		fstarted, ffinished, elapsed, headline, lib.ByteSize(operation.BytesTransferred), operation.Speed,
	)

	switch c.ctx.NotifyTransfer {
	case 1:
		message += fmt.Sprintf("\n\n%d commands were executed.", len(commands))
	case 2:
		printedCommands := ""
		for _, command := range commands {
			printedCommands += command + "\n"
		}
		message += "\n\nThese are the commands that were executed:\n\n" + printedCommands
	}

	go func() {
		if sendErr := sendmail(c.ctx.NotifyTransfer, subject, message, c.ctx.DryRun); sendErr != nil {
			logger.Red("op-sendmail: %s", sendErr)
		}
	}()

	logger.Blue("\n%s\n%s", subject, message)
}

func (c *Core) removeSource(operation *domain.Operation, command *domain.Command) {
	c.state.Status = operation.OpKind
	c.state.Operation = operation
	c.state.Unraid = c.refreshUnraid()

	go c.performRemoveSource(operation, command)
}

func (c *Core) performRemoveSource(operation *domain.Operation, cmd *domain.Command) {
	operation.Line = fmt.Sprintf("Removing source %s", filepath.Join(cmd.Src, cmd.Entry))

	packet := &domain.Packet{Topic: common.EventTransferStarted, Payload: operation}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	for _, command := range operation.Commands {
		if command.ID != cmd.ID {
			continue
		}

		logger.Blue("Removing source:(%s)", filepath.Join(command.Src, command.Entry))

		// status is cmdFlagged currently, let's change this to cmdSourceRemoval, so that handleItemDeletion works
		// correctly and the UI can display a proper feedback for this command
		command.Status = common.CmdSourceRemoval

		c.handleItemDeletion(operation, command)

		command.Status = common.CmdCompleted

		text := "Removal completed"
		operation.Line = text
		logger.Blue(text)

		c.state.Unraid = c.refreshUnraid()
		c.state.History.Items[operation.ID] = operation

		state := &domain.State{
			Operation: operation,
			History:   c.state.History,
			Unraid:    c.state.Unraid,
		}

		packet := &domain.Packet{Topic: common.EventTransferEnded, Payload: state}
		c.ctx.Hub.Pub(packet, "socket:broadcast")

		// TODO: how to handle this
		// c.updateHistory(c.state.History, operation)
		err := c.historyWrite(c.state.History)
		if err != nil {
			logger.Yellow("Unable to write history: %s", err)
		}

		c.state.Status = common.OpNeutral
		c.state.Operation = nil

		break
	}
}

func (c *Core) replay(operation domain.Operation) {
	c.state.Status = operation.OpKind
	c.state.Operation = c.createReplayOperation(operation)
	c.state.Unraid = c.refreshUnraid()

	opName := "Copy"
	if c.state.Operation.OpKind == common.OpScatterMove || c.state.Operation.OpKind == common.OpGatherMove {
		opName = "Move"
	}

	go c.runOperation(opName)
}

func (c *Core) createReplayOperation(original domain.Operation) *domain.Operation {
	operation := &domain.Operation{
		ID:              shortid.MustGenerate(),
		OpKind:          original.OpKind,
		BytesToTransfer: original.BytesToTransfer,
		DryRun:          false,
	}

	operation.RsyncArgs = original.RsyncArgs
	operation.RsyncStrArgs = strings.Join(operation.RsyncArgs, " ")

	operation.Commands = original.Commands

	for _, command := range operation.Commands {
		command.ID = shortid.MustGenerate()
		command.Transferred = 0
		command.Status = common.CmdPending
	}

	return operation
}
