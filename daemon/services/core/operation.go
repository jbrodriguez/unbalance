package core

import (
	"fmt"
	"math"
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

const (
	defaultSpeedWindow = 90 * time.Second
	maxSpeedWindow     = 10 * time.Minute
)

func (c *Core) runOperation(opName string) {
	logger.Blue("running %s operation ...", opName)

	operation := c.state.Operation
	if err := validateRsyncArgs(operation.RsyncArgs); err != nil {
		logger.Yellow("operation:rsync-args:blocked:%s", err)
		c.publishOperationError("unable to start %s operation: %s", opName, err)
		c.state.Status = common.OpNeutral
		c.state.Operation = nil
		return
	}

	operation.Started = time.Now()

	if c.ctx.NotifyTransfer == 2 {
		c.notifyCommandsToRun(opName, operation)
	}

	operation.Line = "Waiting to collect stats ..."
	packet := &domain.Packet{Topic: common.EventTransferStarted, Payload: operation}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	commandsExecuted := make([]string, 0)

	for _, command := range operation.Commands {
		paths, err := c.validateCommandForExecution(command)
		if err != nil {
			cmd := fmt.Sprintf(`rsync %s %s %s`, operation.RsyncStrArgs, strconv.Quote(command.Entry), strconv.Quote(command.Dst))
			c.commandInterrupted(opName, operation, command, cmd, fmt.Errorf("unsafe command path: %w", err), 0, commandsExecuted)
			return
		}

		command.Src = paths.SrcRoot
		command.Dst = paths.DstRoot
		command.Entry = paths.Entry

		cmd := fmt.Sprintf(`rsync %s %s %s`, operation.RsyncStrArgs, strconv.Quote(paths.Entry), strconv.Quote(paths.DstRoot))
		logger.Blue("command started: (src: %s) %s ", paths.SrcRoot, cmd)

		operation.Line = cmd
		packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
		c.ctx.Hub.Pub(packet, "socket:broadcast")

		cmdTransferred, err := c.runCommand(operation, command)
		if err != nil {
			c.commandInterrupted(opName, operation, command, cmd, err, cmdTransferred, commandsExecuted)
			return
		}

		c.commandCompleted(operation, command)

		commandsExecuted = append(commandsExecuted, cmd)
	}

	c.operationCompleted(opName, operation, commandsExecuted)
}

func (c *Core) runCommand(operation *domain.Operation, command *domain.Command) (uint64, error) {
	paths, err := c.validateCommandForExecution(command)
	if err != nil {
		return 0, fmt.Errorf("unsafe command path: %w", err)
	}

	command.Src = paths.SrcRoot
	command.Dst = paths.DstRoot
	command.Entry = paths.Entry

	args := append(
		operation.RsyncArgs,
		paths.Entry,
		paths.DstRoot,
	)

	// make sure the command will run
	c.stopped = false

	// start rsync command
	cmd, err := lib.StartRsync(paths.SrcRoot, args...)
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

		bytesTransferred := operation.BytesTransferred + transferred
		percent, _, _ := progress(operation.BytesToTransfer, bytesTransferred, time.Since(operation.Started))

		c.updateSamples(operation, bytesTransferred)
		speed := c.calculateSpeed(operation)

		operation.Line = current
		operation.Completed = percent
		operation.Speed = speed
		operation.Remaining = remainingAtSpeed(operation.BytesToTransfer, bytesTransferred, speed)
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

	subject := fmt.Sprintf("unbalanced - %s operation INTERRUPTED", strings.ToUpper(opName))
	headline := fmt.Sprintf("Command Interrupted: %s (%s)", cmd, err.Error()+" : "+getError(err.Error(), reRsync, rsyncErrors))

	logger.Yellow("%s", headline)
	packet := &domain.Packet{Topic: common.EventOperationError, Payload: fmt.Sprintf("%s operation was interrupted. Check log (/var/log/unbalanced.log) for additional details.", opName)}
	c.ctx.Hub.Pub(packet, "socket:broadcast")

	operation.BytesTransferred += cmdTransferred
	percent, _, _ := progress(operation.BytesToTransfer, operation.BytesTransferred, elapsed)

	speed := c.calculateSpeed(operation)

	operation.Completed = percent
	operation.Speed = speed
	operation.Remaining = remainingAtSpeed(operation.BytesToTransfer, operation.BytesTransferred, speed)
	operation.DeltaTransfer = cmdTransferred
	command.Transferred = cmdTransferred

	c.endOperation(subject, headline, commandsExecuted, operation)
}

func (c *Core) commandCompleted(operation *domain.Operation, command *domain.Command) {
	text := "Command Finished"
	logger.Blue("%s", text)

	operation.BytesTransferred += command.Size
	c.updateSamples(operation, operation.BytesTransferred)

	percent, _, _ := progress(operation.BytesToTransfer, operation.BytesTransferred, time.Since(operation.Started))

	speed := c.calculateSpeed(operation)

	operation.Completed = percent
	operation.Speed = speed
	operation.Remaining = remainingAtSpeed(operation.BytesToTransfer, operation.BytesTransferred, speed)
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
			logger.Yellow("%s", msg)

			return
		}

		operation.Line = fmt.Sprintf("Removing source %s", filepath.Join(command.Src, command.Entry))

		packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
		c.ctx.Hub.Pub(packet, "socket:broadcast")

		removed, pruned, err := removeTransferredSource(command, operation.OpKind == common.OpGatherMove)
		if err != nil {
			msg := fmt.Sprintf("Unable to remove source folder (%s): %s", filepath.Join(command.Src, command.Entry), err)
			operation.Line = msg

			packet := &domain.Packet{Topic: common.EventTransferProgress, Payload: operation}
			c.ctx.Hub.Pub(packet, "socket:broadcast")

			logger.Yellow("%s", msg)
			return
		}

		logger.Blue("removed:(%s)", removed)
		for _, parent := range pruned {
			logger.Blue("pruned:(%s)", parent)
		}
	}
}

func (c *Core) operationCompleted(opName string, operation *domain.Operation, commandsExecuted []string) {
	operation.Finished = time.Now()
	elapsed := operation.Finished.Sub(operation.Started)

	subject := fmt.Sprintf("unbalanced - %s operation completed", strings.ToUpper(opName))
	headline := fmt.Sprintf("%s operation has finished", opName)

	percent, _, _ := progress(operation.BytesToTransfer, operation.BytesTransferred, elapsed)

	speed := c.calculateSpeed(operation)

	operation.Completed = percent
	operation.Speed = speed
	operation.Remaining = remainingAtSpeed(operation.BytesToTransfer, operation.BytesTransferred, speed)

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

func (c *Core) removeSourceByID(operationID, commandID string) {
	operation, err := c.historyOperation(operationID)
	if err != nil {
		logger.Yellow("removeSource: %s", err)
		c.publishOperationError("unable to remove source: %s", err)
		return
	}

	command, err := findCommand(&operation, commandID)
	if err != nil {
		logger.Yellow("removeSource: %s", err)
		c.publishOperationError("unable to remove source: %s", err)
		return
	}

	c.state.Status = operation.OpKind
	c.state.Operation = &operation
	c.state.Unraid = c.refreshUnraid()

	go c.performRemoveSource(&operation, command)
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
		logger.Blue("%s", text)

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

func (c *Core) replay(operationID string) {
	operation, err := c.historyOperation(operationID)
	if err != nil {
		logger.Yellow("replay: %s", err)
		c.publishOperationError("unable to replay operation: %s", err)
		return
	}

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

	operation.RsyncArgs = append([]string(nil), original.RsyncArgs...)
	operation.RsyncStrArgs = strings.Join(operation.RsyncArgs, " ")

	operation.Commands = make([]*domain.Command, 0, len(original.Commands))
	for _, originalCommand := range original.Commands {
		if originalCommand == nil {
			continue
		}

		command := *originalCommand
		command.ID = shortid.MustGenerate()
		command.Transferred = 0
		command.Status = common.CmdPending
		operation.Commands = append(operation.Commands, &command)
	}

	return operation
}

func (c *Core) updateSamples(operation *domain.Operation, transferred uint64) {
	c.updateSamplesAt(operation, transferred, time.Now())
}

func (c *Core) updateSamplesAt(operation *domain.Operation, transferred uint64, sampledAt time.Time) {
	if transferred < operation.PrevSample {
		operation.PrevSample = transferred
		operation.PrevSampleAt = sampledAt
		return
	}

	elapsedFrom := operation.PrevSampleAt
	if elapsedFrom.IsZero() {
		elapsedFrom = operation.Started
	}
	if elapsedFrom.IsZero() {
		elapsedFrom = sampledAt
	}

	if !sampledAt.After(elapsedFrom) {
		operation.PrevSample = transferred
		operation.PrevSampleAt = sampledAt
		return
	}

	operation.Samples = append(operation.Samples, domain.SpeedSample{
		Bytes:     transferred,
		SampledAt: sampledAt,
	})
	operation.Samples = c.pruneSpeedSamples(operation.Samples, sampledAt)
	operation.SampleIndex = len(operation.Samples)
	operation.PrevSample = transferred
	operation.PrevSampleAt = sampledAt
}

func (c *Core) resetSamples(operation *domain.Operation) {
	if operation == nil {
		return
	}

	operation.Samples = nil
	operation.SampleIndex = 0
	operation.PrevSample = 0
	operation.PrevSampleAt = time.Time{}
}

func (c *Core) calculateSpeed(operation *domain.Operation) float64 {
	window := c.speedWindow()
	cutoff := time.Now().Add(-window)
	var samples []domain.SpeedSample

	for _, sample := range operation.Samples {
		if sample.SampledAt.IsZero() || sample.SampledAt.Before(cutoff) {
			continue
		}
		samples = append(samples, sample)
	}

	if len(samples) < 2 {
		return 0.0
	}

	oldest := samples[0]
	latest := samples[len(samples)-1]
	elapsed := latest.SampledAt.Sub(oldest.SampledAt)
	if elapsed <= 0 || latest.Bytes <= oldest.Bytes {
		return 0.0
	}

	speed := float64(latest.Bytes-oldest.Bytes) / elapsed.Seconds() / 1024 / 1024 // MB/s

	return speed
}

func (c *Core) pruneSpeedSamples(samples []domain.SpeedSample, now time.Time) []domain.SpeedSample {
	window := c.speedWindow()
	cutoff := now.Add(-window)
	keepFrom := 0
	for keepFrom < len(samples) && samples[keepFrom].SampledAt.Before(cutoff) {
		keepFrom++
	}
	if keepFrom == 0 {
		return samples
	}
	return samples[keepFrom:]
}

func (c *Core) speedWindow() time.Duration {
	if c == nil || c.ctx == nil || c.ctx.Config.SpeedWindow == "" {
		return defaultSpeedWindow
	}

	window, err := time.ParseDuration(strings.TrimSpace(c.ctx.Config.SpeedWindow))
	if err != nil || window <= 0 {
		return defaultSpeedWindow
	}
	if window > maxSpeedWindow {
		return maxSpeedWindow
	}
	return window
}

func remainingAtSpeed(bytesToTransfer, bytesTransferred uint64, speed float64) string {
	if bytesTransferred >= bytesToTransfer {
		return "0s"
	}
	if speed <= 0 {
		return "unknown"
	}

	bytesPerSec := speed * 1024 * 1024
	left := time.Duration(float64(bytesToTransfer-bytesTransferred) / bytesPerSec * float64(time.Second))
	return formatRemainingDuration(left)
}

func formatRemainingDuration(duration time.Duration) string {
	seconds := int64(math.Ceil(duration.Seconds()))
	if seconds <= 0 {
		return "0s"
	}

	hours := seconds / 3600
	minutes := (seconds - hours*3600) / 60
	seconds = seconds - hours*3600 - minutes*60

	var parts []string
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	return strings.Join(parts, " ")
}
