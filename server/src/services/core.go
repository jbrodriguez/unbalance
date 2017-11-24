package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"jbrodriguez/unbalance/server/src/common"
	"jbrodriguez/unbalance/server/src/domain"
	"jbrodriguez/unbalance/server/src/dto"
	"jbrodriguez/unbalance/server/src/lib"

	"github.com/jbrodriguez/actor"
	"github.com/jbrodriguez/mlog"
	"github.com/jbrodriguez/pubsub"
	version "github.com/mcuadros/go-version"
	"github.com/teris-io/shortid"
)

const (
	mailCmd    = "/usr/local/emhttp/webGui/scripts/notify"
	timeFormat = "Jan _2, 2006 15:04:05"
)

// Core service
type Core struct {
	bus      *pubsub.PubSub
	settings *lib.Settings
	actor    *actor.Actor

	// this holds the state the app
	state *domain.State

	reFreeSpace *regexp.Regexp
	reRsync     *regexp.Regexp
	reProgress  *regexp.Regexp

	rsyncErrors map[int]string
}

// NewCore -
func NewCore(bus *pubsub.PubSub, settings *lib.Settings) *Core {
	core := &Core{
		bus:      bus,
		settings: settings,
		actor:    actor.NewActor(bus),
		state:    &domain.State{},
	}

	core.reFreeSpace = regexp.MustCompile(`(.*?)\s+(\d+)\s+(\d+)\s+(\d+)\s+(.*?)\s+(.*?)$`)
	core.reRsync = regexp.MustCompile(`exit status (\d+)`)
	core.reProgress = regexp.MustCompile(`(?s)^([\d,]+).*?\(.*?\)$|^([\d,]+).*?$`)

	core.rsyncErrors = map[int]string{
		0:  "Success",
		1:  "Syntax or usage error",
		2:  "Protocol incompatibility",
		3:  "Errors selecting input/output files, dirs",
		4:  "Requested action not supported: an attempt was made to manipulate 64-bit files on a platform that cannot support them, or an option was specified that is supported by the client and not by the server.",
		5:  "Error starting client-server protocol",
		6:  "Daemon unable to append to log-file",
		10: "Error in socket I/O",
		11: "Error in file I/O",
		12: "Error in rsync protocol data stream",
		13: "Errors with program diagnostics",
		14: "Error in IPC code",
		20: "Received SIGUSR1 or SIGINT",
		21: "Some error returned by waitpid()",
		22: "Error allocating core memory buffers",
		23: "Partial transfer due to error",
		24: "Partial transfer due to vanished source files",
		25: "The --max-delete limit stopped deletions",
		30: "Timeout in data send/receive",
		35: "Timeout waiting for daemon connection",
	}

	return core
}

// Start -
func (c *Core) Start() (err error) {
	mlog.Info("starting service Core ...")

	msg := &pubsub.Message{Reply: make(chan interface{}, capacity)}
	c.bus.Pub(msg, common.INT_GET_ARRAY_STATUS)
	reply := <-msg.Reply
	message := reply.(dto.Message)
	if message.Error != nil {
		return message.Error
	}

	c.state.Status = common.OP_NEUTRAL
	c.state.Unraid = message.Data.(*domain.Unraid)
	c.state.Operation = resetOp(c.state.Unraid.Disks)

	history, err := c.historyRead()
	if err != nil {
		mlog.Warning("Unable to read history: %s", err)
	}

	c.state.History = history
	// for _, op := range c.state.History {
	// 	mlog.Info(`op
	// 		%+v
	// 		`, op)
	// }

	c.actor.Register(common.API_GET_CONFIG, c.getConfig)
	c.actor.Register(common.API_GET_STATUS, c.getStatus)
	c.actor.Register(common.API_GET_STATE, c.getState)
	c.actor.Register(common.API_GET_HISTORY, c.getHistory)
	c.actor.Register(common.API_RESET_OP, c.resetOp)
	c.actor.Register(common.API_LOCATE_FOLDER, c.locate)

	c.actor.Register(common.API_SCATTER_CALCULATE, c.scatterCalculate)
	c.actor.Register(common.INT_SCATTER_CALCULATE_FINISHED, c.scatterCalculateFinished)
	c.actor.Register(common.API_SCATTER_MOVE, c.scatterMove)
	c.actor.Register(common.API_SCATTER_COPY, c.scatterCopy)

	c.actor.Register(common.INT_OPERATION_FINISHED, c.operationFinished)

	c.actor.Register(common.API_GATHER_CALCULATE, c.gatherCalculate)
	c.actor.Register(common.INT_GATHER_CALCULATE_FINISHED, c.gatherCalculateFinished)
	c.actor.Register(common.API_GATHER_MOVE, c.gatherMove)

	c.actor.Register(common.API_TOGGLE_DRYRUN, c.toggleDryRun)
	c.actor.Register(common.API_NOTIFY_CALC, c.setNotifyCalc)
	c.actor.Register(common.API_NOTIFY_MOVE, c.setNotifyMove)
	c.actor.Register(common.API_SET_RESERVED, c.setReservedSpace)
	c.actor.Register(common.API_SET_VERBOSITY, c.setVerbosity)
	c.actor.Register(common.API_SET_CHECKUPDATE, c.setCheckUpdate)
	c.actor.Register(common.API_GET_UPDATE, c.getUpdate)
	// c.actor.Register("/config/set/rsyncFlags", c.setRsyncFlags)
	// c.actor.Register("validate", c.validate)
	// c.actor.Register("getLog", c.getLog)

	go c.actor.React()

	return nil
}

// Stop -
func (c *Core) Stop() {
	mlog.Info("stopped service Core ...")
}

func (c *Core) getConfig(msg *pubsub.Message) {
	mlog.Info("Sending config")

	rsyncFlags := strings.Join(c.settings.RsyncFlags, " ")
	if rsyncFlags == "-avX --partial" || rsyncFlags == "-avRX --partial" {
		c.settings.RsyncFlags = []string{"-avPRX"}
		c.settings.Save()
	}

	msg.Reply <- &c.settings.Config
}

func (c *Core) getStatus(msg *pubsub.Message) {
	mlog.Info("Sending status")

	msg.Reply <- c.state.Status
}

func (c *Core) getState(msg *pubsub.Message) {
	mlog.Info("Sending state")

	msg.Reply <- c.state
}

func (c *Core) getHistory(msg *pubsub.Message) {
	mlog.Info("Sending history")
	msg.Reply <- c.state.History
}

func (c *Core) resetOp(msg *pubsub.Message) {
	mlog.Info("resetting op")

	c.state.Operation = resetOp(c.state.Unraid.Disks)

	msg.Reply <- c.state.Operation
}

func (c *Core) locate(msg *pubsub.Message) {
	chosen := msg.Payload.([]string)

	location := &dto.Location{
		Disks:    make(map[string]*domain.Disk, 0),
		Presence: make(map[string]string, 0),
	}

	for _, disk := range c.state.Unraid.Disks {
		for _, item := range chosen {
			name := strings.Replace(item, "/mnt/user/", "", -1)
			entry := filepath.Join(disk.Path, name)

			exists := true
			if _, err := os.Stat(entry); err != nil {
				exists = !os.IsNotExist(err)
			}

			mlog.Info("entry(%s)-exists(%t)", entry, exists)

			if exists {
				location.Disks[disk.Name] = disk

				presence := disk.Name
				if val, ok := location.Presence[name]; ok {
					presence += ", " + val
				}

				location.Presence[name] = presence
			}
		}
	}

	msg.Reply <- location
}

// SCATTER CALCULATE
func (c *Core) setupScatterCalculateOperation(msg *pubsub.Message) (*domain.Operation, error) {
	payload, ok := msg.Payload.(string)
	if !ok {
		return nil, errors.New("Unable to convert scatter calculate parameters")
	}

	var param domain.Operation
	err := json.Unmarshal([]byte(payload), &param)
	if err != nil {
		return nil, err
	}

	// get a fresh operation
	operation := resetOp(c.state.Unraid.Disks)

	operation.OpKind = common.OP_SCATTER_CALC
	operation.ChosenFolders = param.ChosenFolders

	for _, disk := range c.state.Unraid.Disks {
		operation.VDisks[disk.Path].Src = param.VDisks[disk.Path].Src
		operation.VDisks[disk.Path].Dst = param.VDisks[disk.Path].Dst
	}

	return operation, nil
}

func (c *Core) scatterCalculate(msg *pubsub.Message) {
	c.state.Status = common.OP_SCATTER_CALC

	operation, err := c.setupScatterCalculateOperation(msg)
	if err != nil {
		// send to front end the signal of operation finished
		outbound := &dto.Packet{Topic: common.WS_CALC_FINISHED, Payload: resetOp(c.state.Unraid.Disks)}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		outbound = &dto.Packet{Topic: common.WS_CALC_ISSUES, Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	calc := &pubsub.Message{Payload: &domain.State{
		Status:    c.state.Status,
		Unraid:    c.state.Unraid,
		Operation: operation,
	}}

	c.bus.Pub(calc, common.INT_SCATTER_CALCULATE)
}

func (c *Core) scatterCalculateFinished(msg *pubsub.Message) {
	operation := msg.Payload.(*domain.Operation)

	c.state.Status = common.OP_NEUTRAL
	c.state.Operation = operation

	// send to front end the signal of operation finished
	outbound := &dto.Packet{Topic: common.WS_CALC_FINISHED, Payload: operation}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
	// and there actually perm issues
	if c.state.Operation.BytesToTransfer > 0 && (c.state.Operation.OwnerIssue+c.state.Operation.GroupIssue+c.state.Operation.FolderIssue+c.state.Operation.FileIssue > 0) {
		outbound = &dto.Packet{Topic: common.WS_CALC_ISSUES, Payload: fmt.Sprintf("%d|%d|%d|%d", c.state.Operation.OwnerIssue, c.state.Operation.GroupIssue, c.state.Operation.FolderIssue, c.state.Operation.FileIssue)}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	mlog.Info(`Operation
		OpKind: %d
		Started: %s
		Finished: %s
		ChosenFolders: %v
		FolderNotTransferred: %v
		OwnerIssue: %d
		GroupIssue: %d
		FolderIssue: %d
		FileIssue: %d
		BytesToTransfer: %d
		DryRun: %t
		RsyncFlags: %v
		RsyncStrFlags: %s
		Commands: %v
		BytesTransferred: %d
	`, operation.OpKind, operation.Started, operation.Finished, operation.ChosenFolders,
		operation.FoldersNotTransferred, operation.OwnerIssue, operation.GroupIssue,
		operation.FolderIssue, operation.FileIssue, operation.BytesToTransfer, operation.DryRun,
		operation.RsyncFlags, operation.RsyncStrFlags, operation.Commands, operation.BytesTransferred,
	)

	for _, disk := range c.state.Unraid.Disks {
		vdisk := operation.VDisks[disk.Path]

		if vdisk.Bin != nil {
			mlog.Info(`VDisk
			Path: %s
			PlannedFree: %d
			Src: %t,
			Dst: %t
		`, vdisk.Path, vdisk.PlannedFree, vdisk.Src, vdisk.Dst)

			mlog.Info(`Bin
				Size: %d
			`, vdisk.Bin.Size)

			for _, item := range vdisk.Bin.Items {
				mlog.Info(`Item
					Name: %s
					Size: %d
					Path: %s
					Location: %s
				`, item.Name, item.Size, item.Path, item.Location)
			}
		} else {
			mlog.Info(`VDisk
				Path: %s
				PlannedFree: %d
				Src: %t,
				Dst: %t
			`, vdisk.Path, vdisk.PlannedFree, vdisk.Src, vdisk.Dst)
		}
	}
}

// GATHER CALCULATE
func (c *Core) setupGatherCalculateOperation(msg *pubsub.Message) (*domain.Operation, error) {
	data, ok := msg.Payload.(string)
	if !ok {
		return nil, errors.New("Unable to convert findTargets parameters")
	}

	var chosen []string
	err := json.Unmarshal([]byte(data), &chosen)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to bind findTargets parameters: %s", err))
	}

	operation := resetOp(c.state.Unraid.Disks)

	operation.OpKind = c.state.Status
	operation.ChosenFolders = chosen

	return operation, nil
}

func (c *Core) gatherCalculate(msg *pubsub.Message) {
	c.state.Status = common.OP_GATHER_CALC

	operation, err := c.setupGatherCalculateOperation(msg)
	if err != nil {
		mlog.Warning(err.Error())
		outbound := &dto.Packet{Topic: "opError", Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		return
	}

	// TODO: we should probably refresh unraid here (applies to scatterCalculate too)
	calc := &pubsub.Message{Payload: &domain.State{
		Status:    c.state.Status,
		Unraid:    c.state.Unraid,
		Operation: operation,
	}}

	c.bus.Pub(calc, common.INT_GATHER_CALCULATE)
}

func (c *Core) gatherCalculateFinished(msg *pubsub.Message) {
	operation := msg.Payload.(*domain.Operation)

	c.state.Status = common.OP_NEUTRAL
	c.state.Operation = operation

	// send to front end the signal of operation finished
	outbound := &dto.Packet{Topic: common.WS_CALC_FINISHED, Payload: operation}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
	// and there actually perm issues
	if c.state.Operation.BytesToTransfer > 0 && (c.state.Operation.OwnerIssue+c.state.Operation.GroupIssue+c.state.Operation.FolderIssue+c.state.Operation.FileIssue > 0) {
		outbound = &dto.Packet{Topic: common.WS_CALC_ISSUES, Payload: fmt.Sprintf("%d|%d|%d|%d", c.state.Operation.OwnerIssue, c.state.Operation.GroupIssue, c.state.Operation.FolderIssue, c.state.Operation.FileIssue)}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}
}

// SCATTER TRANSFER
func (c *Core) setupScatterTransferOperation(copyOperation *domain.Operation) *domain.Operation {
	operation := &domain.Operation{
		ID:              shortid.MustGenerate(),
		OpKind:          c.state.Status,
		BytesToTransfer: copyOperation.BytesToTransfer,
		DryRun:          c.settings.DryRun,
		RsyncFlags:      c.settings.RsyncFlags,
		VDisks:          copyOperation.VDisks,
	}

	// user may have changed rsync flags or dry-run setting, adjust for it
	if operation.DryRun {
		operation.RsyncFlags = append(operation.RsyncFlags, "--dry-run")
	}
	operation.RsyncStrFlags = strings.Join(operation.RsyncFlags, " ")

	operation.Commands = make([]*domain.Command, 0)

	for _, disk := range c.state.Unraid.Disks {
		vdisk := operation.VDisks[disk.Path]
		if vdisk.Bin == nil || vdisk.Src {
			continue
		}

		for _, item := range vdisk.Bin.Items {
			var entry string

			// this a double check. item.Path should be in the form of films/bluray or tvshows/Billions
			// if for some reason it starts with a "/", we strip it
			if item.Path[0] == filepath.Separator {
				entry = item.Path[1:]
			} else {
				entry = item.Path
			}

			operation.Commands = append(operation.Commands, &domain.Command{
				Src:   item.Location,
				Dst:   vdisk.Path + string(filepath.Separator),
				Entry: entry,
				Size:  item.Size,
			})
		}
	}

	return operation
}

func (c *Core) scatterMove(msg *pubsub.Message) {
	c.state.Status = common.OP_SCATTER_MOVE
	c.state.Operation = c.setupScatterTransferOperation(c.state.Operation)
	go c.runOperation("Move")
}

func (c *Core) scatterCopy(msg *pubsub.Message) {
	c.state.Status = common.OP_SCATTER_COPY
	c.state.Operation = c.setupScatterTransferOperation(c.state.Operation)
	go c.runOperation("Copy")
}

// GATHER TRANSFER
func (c *Core) setupGatherTransferOperation(msg *pubsub.Message) (*domain.Operation, error) {
	// mlog.Info("%+v", msg.Payload)

	data, ok := msg.Payload.(string)
	if !ok {
		return nil, errors.New("Unable to convert gather parameters")
	}

	var target domain.Disk
	err := json.Unmarshal([]byte(data), &target)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Unable to bind gather parameters: %s", err))
	}

	currentOp := c.state.Operation

	operation := &domain.Operation{
		ID:         shortid.MustGenerate(),
		OpKind:     c.state.Status,
		DryRun:     c.settings.DryRun,
		RsyncFlags: c.settings.RsyncFlags,
		VDisks:     currentOp.VDisks,
	}

	// user chose a target disk, adjust bytestotransfer to the size of its bin, since
	// that's the amount of data we need to transfer. Also remove bin from all other disks,
	// since only the target will have work to do
	for _, disk := range c.state.Unraid.Disks {
		if disk.Path == target.Path {
			operation.BytesToTransfer = operation.VDisks[target.Path].Bin.Size
		} else {
			operation.VDisks[disk.Path].Bin = nil
		}
	}

	// user may have changed rsync flags or dry-run setting, adjust for it
	if operation.DryRun {
		operation.RsyncFlags = append(operation.RsyncFlags, "--dry-run")
	}
	operation.RsyncStrFlags = strings.Join(operation.RsyncFlags, " ")

	operation.Commands = make([]*domain.Command, 0)

	for _, disk := range c.state.Unraid.Disks {
		vdisk := operation.VDisks[disk.Path]
		if vdisk.Bin == nil {
			continue
		}

		for _, item := range vdisk.Bin.Items {
			var entry string

			// this a double check. item.Path should be in the form of films/bluray or tvshows/Billions
			// if for some reason it starts with a "/", we strip it
			if item.Path[0] == filepath.Separator {
				entry = item.Path[1:]
			} else {
				entry = item.Path
			}

			operation.Commands = append(operation.Commands, &domain.Command{
				Src:   item.Location,
				Dst:   vdisk.Path + string(filepath.Separator),
				Entry: entry,
				Size:  item.Size,
			})
		}
	}

	return operation, nil
}

func (c *Core) gatherMove(msg *pubsub.Message) {
	c.state.Status = common.OP_GATHER_MOVE

	operation, err := c.setupGatherTransferOperation(msg)
	if err != nil {
		mlog.Warning(err.Error())
		return
	}

	c.state.Operation = operation

	go c.runOperation("Move")
}

// COMMON TRANSFER
func (c *Core) runOperation(opName string) {
	mlog.Info("Running %s operation ...", opName)

	operation := c.state.Operation
	operation.Started = time.Now()

	// notifyMove should be renamed to notifyTransfer
	if c.settings.NotifyMove == 2 {
		c.notifyCommandsToRun(opName, operation)
	}

	operation.Line = "Waiting to collect stats ..."
	outbound := &dto.Packet{Topic: "transferStarted", Payload: operation}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// Initialize local variables
	var calls int64
	var callsPerDelta int64
	var elapsed time.Duration

	commandsExecuted := make([]string, 0)

	for _, command := range operation.Commands {
		args := append(
			operation.RsyncFlags,
			command.Entry,
			command.Dst,
		)
		cmd := fmt.Sprintf(`rsync %s %s %s`, operation.RsyncStrFlags, strconv.Quote(command.Entry), strconv.Quote(command.Dst))
		mlog.Info("Command Started: (src: %s) %s ", command.Src, cmd)

		operation.Line = cmd
		outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		// tvshows/Billions/Season 01/banner.s01.jpg
		//              20 100%    0.00kB/s    0:00:00
		//              20 100%    0.00kB/s    0:00:00 (xfr#1, to-chk=4/8)
		// tvshows/Billions/Season 01/billions.s01e01.1080p
		//          32,768   0%    1.01MB/s    0:02:53
		//      48,005,120  26%   45.78MB/s    0:00:02
		//     100,696,064  56%   47.97MB/s    0:00:01
		//     157,810,688  88%   50.15MB/s    0:00:00
		//     179,306,496 100%   52.91MB/s    0:00:03 (xfr#2, to-chk=3/8)
		// tvshows/Billions/Season 01/billions.s01e02.1080p
		//          32,768   0%  137.34kB/s    0:21:52
		//      38,305,792  21%   36.35MB/s    0:00:03
		//     106,397,696  58%   50.58MB/s    0:00:01
		//     180,355,072 100%   64.83MB/s    0:00:02 (xfr#3, to-chk=2/8)
		// tvshows/Billions/Season 01/billions.s01e03.1080p
		//          32,768   0%   49.31kB/s    1:01:18
		//      41,811,968  23%   39.88MB/s    0:00:03
		//     157,810,688  86%   75.29MB/s    0:00:00
		//     181,403,648 100%   79.21MB/s    0:00:02 (xfr#4, to-chk=1/8)
		// tvshows/Billions/Season 01/billions.s01e04.1080p
		//          32,768   0%  171.12kB/s    0:17:46
		//      82,542,592  45%   78.72MB/s    0:00:01
		//     112,492,544  61%   53.45MB/s    0:00:01
		//     120,881,152  66%   38.11MB/s    0:00:01
		//     164,134,912  89%   38.30MB/s    0:00:00
		//     182,452,224 100%   39.94MB/s    0:00:04 (xfr#5, to-chk=0/8)
		//
		// rsync is very particular in how it reports progress: each line shows the total bytes transferred for a
		// particular file, then starts over with the next file
		// makes sense for them I guess, but it's a pita to track and get an overall total
		//
		// so this is what the following represent:
		//
		// - cmdTransferred holds the running total for the current command
		//
		// - accumTransferred holds the running total for all the files that have been transferred, not including the
		// current file, for the current command
		//
		// - perFileTransferred holds the running total for the file that is currently being transferred

		var cmdTransferred, accumTransferred, perFileTransferred int64

		// actual shell execution
		err := lib.ShellEx(func(text string) {
			line := strings.TrimSpace(text)

			if len(line) <= 0 {
				return
			}

			if callsPerDelta <= 50 {
				calls++
			}

			delta := int64(time.Since(operation.Started) / time.Second)
			if delta == 0 {
				delta = 1
			}
			callsPerDelta = calls / delta

			match := c.reProgress.FindStringSubmatch(line)

			// this is a regular output line from rsync
			if match == nil {
				// make sure it's available for the front end
				operation.Line = line

				// log it according to verbosity settings
				if c.settings.Verbosity == 1 {
					mlog.Info("%s", line)
				}

				if callsPerDelta <= 50 {
					outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
					c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
				}

				return
			}

			// this is a file transfer progress output line
			if match[1] == "" {
				// this happens when the file hasn't finished transferring
				moved := strings.Replace(match[2], ",", "", -1)
				perFileTransferred, _ = strconv.ParseInt(moved, 10, 64)
				cmdTransferred = accumTransferred + perFileTransferred
			} else {
				// the file has finished transferring
				moved := strings.Replace(match[1], ",", "", -1)
				perFileTransferred, _ = strconv.ParseInt(moved, 10, 64)
				cmdTransferred = accumTransferred + perFileTransferred
				accumTransferred += perFileTransferred
			}

			if callsPerDelta <= 50 {
				percent, left, speed := progress(operation.BytesToTransfer, operation.BytesTransferred+cmdTransferred, time.Since(operation.Started))

				operation.Completed = percent
				operation.Speed = speed
				operation.Remaining = fmt.Sprintf("%s", left)
				operation.DeltaTransfer = cmdTransferred
				command.Transferred = cmdTransferred

				outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			}

		}, mlog.Warning, command.Src, "rsync", args...)

		if err != nil {
			c.transferInterrupted(opName, operation, command, cmd, err, cmdTransferred, commandsExecuted)
			return
		}

		c.commandCompleted(operation, command, elapsed)

		commandsExecuted = append(commandsExecuted, cmd)
	}

	c.operationCompleted(opName, operation, commandsExecuted)
}

func (c *Core) transferInterrupted(opName string, operation *domain.Operation, command *domain.Command, cmd string, err error, cmdTransferred int64, commandsExecuted []string) {
	operation.Finished = time.Now()
	elapsed := time.Since(operation.Started)

	subject := fmt.Sprintf("unBALANCE - %s operation INTERRUPTED", strings.ToUpper(opName))
	headline := fmt.Sprintf("Command Interrupted: %s (%s)", cmd, err.Error()+" : "+getError(err.Error(), c.reRsync, c.rsyncErrors))

	mlog.Warning(headline)
	outbound := &dto.Packet{Topic: "opError", Payload: fmt.Sprintf("%s operation was interrupted. Check log (/boot/logs/unbalance.log) for details.", opName)}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	operation.BytesTransferred += cmdTransferred
	percent, left, speed := progress(operation.BytesToTransfer, operation.BytesTransferred, elapsed)

	operation.Completed = percent
	operation.Speed = speed
	operation.Remaining = fmt.Sprintf("%s", left)
	operation.DeltaTransfer = cmdTransferred
	command.Transferred = cmdTransferred

	c.endOperation(subject, headline, commandsExecuted, operation)
}

func (c *Core) commandCompleted(operation *domain.Operation, command *domain.Command, elapsed time.Duration) {
	text := "Command Finished"
	mlog.Info(text)

	operation.BytesTransferred += command.Size
	percent, left, speed := progress(operation.BytesToTransfer, operation.BytesTransferred, elapsed)

	operation.Completed = percent
	operation.Speed = speed
	operation.Remaining = fmt.Sprintf("%s", left)
	operation.DeltaTransfer = 0
	operation.Line = text
	command.Transferred = command.Size

	msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, operation.Remaining, speed)
	mlog.Info("Current progress: %s", msg)

	outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// this is just a heads up for the user, shows which folders would/wouldn't be pruned if run without dry-run
	if operation.DryRun && operation.OpKind == common.OP_GATHER_MOVE {
		parent := filepath.Dir(command.Entry)
		mlog.Info("parent(%s)-src(%s)-dst(%s)-entry(%s)", parent, command.Src, command.Dst, command.Entry)
		if parent != "." {
			mlog.Info(`Would delete empty folders starting from (%s) - (find "%s" -type d -empty -prune -exec rm -rf {} \;) `, filepath.Join(command.Src, parent), filepath.Join(command.Src, parent))
		} else {
			mlog.Info(`WONT DELETE: find "%s" -type d -empty -prune -exec rm -rf {} \;`, filepath.Join(command.Src, parent))
		}
	}

	// if it isn't a dry-run and the operation is Move or Gather, delete the source folder
	if !operation.DryRun && (operation.OpKind == common.OP_SCATTER_MOVE || operation.OpKind == common.OP_GATHER_MOVE) {
		exists, _ := lib.Exists(filepath.Join(command.Dst, command.Entry))
		if exists {
			rmrf := fmt.Sprintf("rm -rf \"%s\"", filepath.Join(command.Src, command.Entry))
			mlog.Info("Removing: %s", rmrf)
			err := lib.Shell(rmrf, mlog.Warning, "transferProgress:", "", func(line string) {
				mlog.Info(line)
			})

			if err != nil {
				msg := fmt.Sprintf("Unable to remove source folder (%s): %s", filepath.Join(command.Src, command.Entry), err)

				outbound := &dto.Packet{Topic: "transferProgress", Payload: msg}
				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

				mlog.Warning(msg)
			}

			if operation.OpKind == common.OP_GATHER_MOVE {
				parent := filepath.Dir(command.Entry)
				if parent != "." {
					rmdir := fmt.Sprintf(`find "%s" -type d -empty -prune -exec rm -rf {} \;`, filepath.Join(command.Src, parent))
					mlog.Info("Running %s", rmdir)

					err = lib.Shell(rmdir, mlog.Warning, "transferProgress:", "", func(line string) {
						mlog.Info(line)
					})

					if err != nil {
						msg := fmt.Sprintf("Unable to remove parent folder (%s): %s", filepath.Join(command.Src, parent), err)

						outbound := &dto.Packet{Topic: "transferProgress", Payload: msg}
						c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

						mlog.Warning(msg)
					}
				}
			}
		} else {
			mlog.Warning("Skipping deletion (file/folder not present in destination): %s", filepath.Join(command.Dst, command.Entry))
		}
	}
}

func (c *Core) operationCompleted(opName string, operation *domain.Operation, commandsExecuted []string) {
	operation.Finished = time.Now()
	elapsed := time.Since(operation.Started)

	subject := fmt.Sprintf("unBALANCE - %s operation completed", strings.ToUpper(opName))
	headline := fmt.Sprintf("%s operation has finished", opName)

	percent, left, speed := progress(operation.BytesToTransfer, operation.BytesTransferred, elapsed)
	operation.Completed = percent
	operation.Speed = speed
	operation.Remaining = fmt.Sprintf("%s", left)

	c.endOperation(subject, headline, commandsExecuted, operation)
}

func (c *Core) endOperation(subject, headline string, commands []string, operation *domain.Operation) {
	fstarted := operation.Started.Format(timeFormat)
	ffinished := operation.Finished.Format(timeFormat)
	elapsed := lib.Round(time.Since(operation.Started), time.Millisecond)

	outbound := &dto.Packet{Topic: "transferFinished", Payload: operation}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s\n\n%s\n\nTransferred %s at ~ %.2f MB/s",
		fstarted, ffinished, elapsed, headline, lib.ByteSize(operation.BytesTransferred), operation.Speed,
	)

	switch c.settings.NotifyMove {
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
		if sendErr := c.sendmail(c.settings.NotifyMove, subject, message, c.settings.DryRun); sendErr != nil {
			mlog.Error(sendErr)
		}
	}()

	mlog.Info("\n%s\n%s", subject, message)

	c.bus.Pub(&pubsub.Message{}, common.INT_OPERATION_FINISHED)
}

func (c *Core) operationFinished(msg *pubsub.Message) {
	count := len(c.state.History)
	if count == common.HISTORY_CAPACITY {
		c.state.History = append(c.state.History[1:], c.state.Operation)
	} else {
		c.state.History = append(c.state.History, c.state.Operation)
	}

	err := c.historyWrite(c.state.History)
	if err != nil {
		mlog.Warning("Unable to write history: %s", err)
	}

	c.state.Status = common.OP_NEUTRAL
	c.state.Operation = resetOp(c.state.Unraid.Disks)
}

// SETTINGS RELATED
func (c *Core) setNotifyCalc(msg *pubsub.Message) {
	fnotify := msg.Payload.(float64)
	notify := int(fnotify)

	mlog.Info("Setting notifyCalc to (%d)", notify)

	c.settings.NotifyCalc = notify
	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) setNotifyMove(msg *pubsub.Message) {
	fnotify := msg.Payload.(float64)
	notify := int(fnotify)

	mlog.Info("Setting notifyMove to (%d)", notify)

	c.settings.NotifyMove = notify
	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) setVerbosity(msg *pubsub.Message) {
	fverbosity := msg.Payload.(float64)
	verbosity := int(fverbosity)

	mlog.Info("Setting verbosity to (%d)", verbosity)

	c.settings.Verbosity = verbosity
	err := c.settings.Save()
	if err != nil {
		mlog.Warning("not right %s", err)
	}

	msg.Reply <- &c.settings.Config
}

func (c *Core) setCheckUpdate(msg *pubsub.Message) {
	fupdate := msg.Payload.(float64)
	update := int(fupdate)

	mlog.Info("Setting checkForUpdate to (%d)", update)

	c.settings.CheckForUpdate = update
	err := c.settings.Save()
	if err != nil {
		mlog.Warning("not right %s", err)
	}

	msg.Reply <- &c.settings.Config
}

func (c *Core) setReservedSpace(msg *pubsub.Message) {
	mlog.Warning("payload: %+v", msg.Payload)
	payload, ok := msg.Payload.(string)
	if !ok {
		mlog.Warning("Unable to convert Reserved Space parameters")
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert Reserved Space parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		msg.Reply <- &c.settings.Config

		return
	}

	var reserved dto.Reserved
	err := json.Unmarshal([]byte(payload), &reserved)
	if err != nil {
		mlog.Warning("Unable to bind reservedSpace parameters: %s", err)
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind reservedSpace parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		return
	}

	amount := int64(reserved.Amount)
	unit := reserved.Unit

	mlog.Info("Setting reservedAmount to (%d)", amount)
	mlog.Info("Setting reservedUnit to (%s)", unit)

	c.settings.ReservedAmount = amount
	c.settings.ReservedUnit = unit
	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) toggleDryRun(msg *pubsub.Message) {
	mlog.Info("Toggling dryRun from (%t)", c.settings.DryRun)

	c.settings.ToggleDryRun()
	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) setRsyncFlags(msg *pubsub.Message) {
	// mlog.Warning("payload: %+v", msg.Payload)
	payload, ok := msg.Payload.(string)
	if !ok {
		mlog.Warning("Unable to convert Rsync Flags parameters")
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert Rsync Flags parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		msg.Reply <- &c.settings.Config

		return
	}

	var rsync dto.Rsync
	err := json.Unmarshal([]byte(payload), &rsync)
	if err != nil {
		mlog.Warning("Unable to bind rsyncFlags parameters: %s", err)
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind rsyncFlags parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		return
		// mlog.Fatalf(err.Error())
	}

	mlog.Info("Setting rsyncFlags to (%s)", strings.Join(rsync.Flags, " "))

	c.settings.RsyncFlags = rsync.Flags
	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

func (c *Core) getUpdate(msg *pubsub.Message) {
	var newest string

	if c.settings.CheckForUpdate == 1 {
		latest, err := lib.GetLatestVersion("https://raw.githubusercontent.com/jbrodriguez/unbalance/master/VERSION")
		if err != nil {
			return
		}

		latest = strings.TrimSuffix(latest, "\n")
		if version.Compare(latest, c.settings.Version, ">") {
			newest = latest
		}
	}

	msg.Reply <- newest
}

// HELPERS
func (c *Core) sendmail(notify int, subject, message string, dryRun bool) (err error) {
	if notify == 0 {
		return nil
	}

	dry := ""
	if dryRun {
		dry = "-------\nDRY RUN\n-------\n"
	}

	msg := dry + message

	// strCmd := fmt.Sprintf("-s \"%s\" -m \"%s\"", mailCmd, subject, msg)
	cmd := exec.Command(mailCmd, "-e", "unBALANCE operation update", "-s", subject, "-m", msg)
	err = cmd.Run()

	return
}

func (c *Core) notifyCommandsToRun(opName string, operation *domain.Operation) {
	message := "\n\nThe following commands will be executed:\n\n"

	for _, command := range operation.Commands {
		cmd := fmt.Sprintf(`(src: %s) rsync %s %s %s`, command.Src, operation.RsyncStrFlags, strconv.Quote(command.Entry), strconv.Quote(command.Dst))
		message += cmd + "\n"
	}

	subject := fmt.Sprintf("unBALANCE - %s operation STARTED", strings.ToUpper(opName))

	go func() {
		if sendErr := c.sendmail(c.settings.NotifyMove, subject, message, c.settings.DryRun); sendErr != nil {
			mlog.Error(sendErr)
		}
	}()
}

func resetOp(disks []*domain.Disk) *domain.Operation {
	op := &domain.Operation{
		ID:     shortid.MustGenerate(),
		OpKind: common.OP_NEUTRAL,
		VDisks: make(map[string]*domain.VDisk, 0),
	}

	for _, disk := range disks {
		vdisk := &domain.VDisk{Path: disk.Path, PlannedFree: disk.Free, Src: false, Dst: false}
		op.VDisks[disk.Path] = vdisk
	}

	return op
}

func progress(bytesToTransfer, bytesTransferred int64, elapsed time.Duration) (percent float64, left time.Duration, speed float64) {
	bytesPerSec := float64(bytesTransferred) / elapsed.Seconds()
	speed = bytesPerSec / 1024 / 1024 // MB/s

	percent = (float64(bytesTransferred) / float64(bytesToTransfer)) * 100 // %

	left = time.Duration(float64(bytesToTransfer-bytesTransferred)/bytesPerSec) * time.Second

	return
}

func getError(line string, re *regexp.Regexp, errors map[int]string) string {
	result := re.FindStringSubmatch(line)
	status, _ := strconv.Atoi(result[1])
	msg, ok := errors[status]
	if !ok {
		msg = "unknown error"
	}

	return msg
}

func (c *Core) historyRead() ([]*domain.Operation, error) {
	history := make([]*domain.Operation, 0)
	fileName := filepath.Join(common.PLUGIN_LOCATION, common.HISTORY_FILENAME)

	file, err := os.Open(fileName)
	if err != nil {
		return history, err
	}
	defer file.Close()

	mlog.Info(`before decoding`)

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&history)
	if err != nil {
		return history, err
	}

	mlog.Info(`insider(%+v)`, history)

	return history, nil
}

func (c *Core) historyWrite(history []*domain.Operation) error {
	// b, err := json.Marshal(history)
	// if err != nil {
	// 	mlog.Warning("Unable to serialize op: %s", err)
	// }

	tmpName := filepath.Join(common.PLUGIN_LOCATION, common.HISTORY_FILENAME+"."+shortid.MustGenerate())

	file, err := os.Create(tmpName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.Encode(history)

	os.Rename(tmpName, filepath.Join(common.PLUGIN_LOCATION, common.HISTORY_FILENAME))

	return err
}

// func (c *Core) validate(msg *pubsub.Message) {
// 	c.operation.OpState = model.StateValidate
// 	c.operation.PrevState = model.StateValidate
// 	go c.checksum(msg)
// }

// func (c *Core) checksum(msg *pubsub.Message) {
// 	defer func() {
// 		c.operation.OpState = model.StateIdle
// 		c.operation.PrevState = model.StateIdle
// 		c.operation.Started = time.Time{}
// 		c.operation.BytesTransferred = 0
// 	}()

// 	opName := "Validate"

// 	mlog.Info("Running %s operation ...", opName)
// 	c.operation.Started = time.Now()

// 	outbound := &dto.Packet{Topic: "transferStarted", Payload: "Operation started"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	multiSource := false

// 	if !strings.HasPrefix(c.operation.RsyncStrFlags, "-a") {
// 		finished := time.Now()
// 		elapsed := time.Since(c.operation.Started)

// 		subject := fmt.Sprintf("unBALANCE - %s operation INTERRUPTED", strings.ToUpper(opName))
// 		headline := fmt.Sprintf("For proper %s operation, rsync flags MUST begin with -a", opName)

// 		mlog.Warning(headline)
// 		outbound := &dto.Packet{Topic: "opError", Payload: fmt.Sprintf("%s operation was interrupted. Check log (/boot/logs/unbalance.log) for details.", opName)}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		_, _, speed := progress(c.operation.BytesToTransfer, 0, elapsed)
// 		c.finishTransferOperation(subject, headline, make([]string, 0), c.operation.Started, finished, elapsed, 0, speed, multiSource)

// 		return
// 	}

// 	// Initialize local variables
// 	// we use the rsync flags that were created by the transfer operation,
// 	// but replace -a with -rc, to perform the validation
// 	checkRsyncFlags := make([]string, 0)
// 	for _, flag := range c.operation.RsyncFlags {
// 		checkRsyncFlags = append(checkRsyncFlags, strings.Replace(flag, "-a", "-rc", -1))
// 	}

// 	checkRsyncStrFlags := strings.Join(checkRsyncFlags, " ")

// 	// execute each rsync command created in the transfer phase
// 	c.runOperation(opName, checkRsyncFlags, checkRsyncStrFlags, multiSource)
// }

// func (c *Core) gather(msg *pubsub.Message) {
// 	mlog.Info("%+v", msg.Payload)
// 	data, ok := msg.Payload.(string)
// 	if !ok {
// 		mlog.Warning("Unable to convert gather parameters")
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert gather parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 	}

// 	var target model.Disk
// 	err := json.Unmarshal([]byte(data), &target)
// 	if err != nil {
// 		mlog.Warning("Unable to bind gather parameters: %s", err)
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind gather parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 		// mlog.Fatalf(err.Error())
// 	}

// 	// user chose a target disk, adjust bytestotransfer to the size of its bin, since
// 	// that's the amount of data we need to transfer. Also remove bin from all other disks,
// 	// since only the target will have work to do
// 	for _, disk := range c.storage.Disks {
// 		if disk.Path == target.Path {
// 			c.operation.BytesToTransfer = disk.Bin.Size
// 		} else {
// 			disk.Bin = nil
// 		}
// 	}

// 	c.operation.OpState = model.StateGather
// 	c.operation.PrevState = model.StateGather

// 	go c.transfer("Move", true, msg)
// }

// func (c *Core) getLog(msg *pubsub.Message) {
// 	log := c.storage.GetLog()

// 	outbound := &dto.Packet{Topic: "gotLog", Payload: log}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	return
// }
