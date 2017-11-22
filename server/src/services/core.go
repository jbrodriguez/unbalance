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
	c.state.History = make([]*domain.Operation, 0)

	c.actor.Register(common.API_GET_CONFIG, c.getConfig)
	c.actor.Register(common.API_GET_STATUS, c.getStatus)
	c.actor.Register(common.API_GET_STATE, c.getState)
	c.actor.Register(common.API_RESET_OP, c.resetOp)
	c.actor.Register(common.API_LOCATE_FOLDER, c.locate)

	c.actor.Register(common.API_SCATTER_CALCULATE, c.scatterCalculate)
	c.actor.Register(common.INT_SCATTER_CALCULATE_FINISHED, c.scatterCalculateFinished)

	c.actor.Register(common.API_SCATTER_MOVE, c.scatterMove)
	c.actor.Register(common.API_SCATTER_COPY, c.scatterCopy)

	c.actor.Register(common.INT_OPERATION_FINISHED, c.operationFinished)

	// c.actor.Register("/config/set/notifyCalc", c.setNotifyCalc)
	// c.actor.Register("/config/set/notifyMove", c.setNotifyMove)
	// c.actor.Register("/config/set/reservedSpace", c.setReservedSpace)
	// c.actor.Register("/config/set/verbosity", c.setVerbosity)
	// c.actor.Register("/config/set/checkUpdate", c.setCheckUpdate)
	// c.actor.Register("/get/storage", c.getStorage)
	// c.actor.Register("/get/update", c.getUpdate)
	c.actor.Register("/config/toggle/dryRun", c.toggleDryRun)
	// c.actor.Register("/get/tree", c.getTree)
	// c.actor.Register("/disks/locate", c.locate)
	// c.actor.Register("/config/set/rsyncFlags", c.setRsyncFlags)

	// c.actor.Register("calculate", c.calc)
	// c.actor.Register("move", c.move)
	// c.actor.Register("copy", c.copy)
	// c.actor.Register("validate", c.validate)
	// c.actor.Register("getLog", c.getLog)
	// c.actor.Register("findTargets", c.findTargets)
	// c.actor.Register("gather", c.gather)

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

func (c *Core) resetOp(msg *pubsub.Message) {
	mlog.Info("resetting op")

	c.state.Operation = resetOp(c.state.Unraid.Disks)

	msg.Reply <- c.state.Operation
}

func (c *Core) locate(msg *pubsub.Message) {
	chosen := msg.Payload.([]string)

	disks := make([]*domain.Disk, 0)

	for _, disk := range c.state.Unraid.Disks {
		for _, item := range chosen {
			location := filepath.Join(disk.Path, strings.Replace(item, "/mnt/user", "", -1))

			exists := true
			if _, err := os.Stat(location); err != nil {
				exists = !os.IsNotExist(err)
			}

			mlog.Info("location(%s)-exists(%t)", location, exists)

			if exists {
				disks = append(disks, disk)
			}
		}
	}

	msg.Reply <- disks
}

func getScatterParams(msg *pubsub.Message) (*domain.Operation, error) {
	payload, ok := msg.Payload.(string)
	if !ok {
		return nil, errors.New("Unable to convert scatter calculate parameters")
	}

	var param domain.Operation
	err := json.Unmarshal([]byte(payload), &param)
	if err != nil {
		return nil, err
	}

	return &param, nil
}

func (c *Core) scatterCalculate(msg *pubsub.Message) {
	c.state.Status = common.OP_SCATTER_CALC

	params, err := getScatterParams(msg)
	if err != nil {
		// send to front end the signal of operation finished
		outbound := &dto.Packet{Topic: common.WS_CALC_FINISHED, Payload: resetOp(c.state.Unraid.Disks)}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		outbound = &dto.Packet{Topic: common.WS_CALC_ISSUES, Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	// get a fresh operation
	operation := resetOp(c.state.Unraid.Disks)

	operation.OpKind = common.OP_SCATTER_CALC
	operation.ChosenFolders = params.ChosenFolders

	for _, disk := range c.state.Unraid.Disks {
		operation.VDisks[disk.Path].Src = params.VDisks[disk.Path].Src
		operation.VDisks[disk.Path].Dst = params.VDisks[disk.Path].Dst
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

	op2 := c.setupOperation(common.OP_SCATTER_COPY, operation)
	mlog.Info(`Operation
		RSyncFlags: %v
		RSyncStrFlags: %s
		`, op2.RsyncFlags, op2.RsyncStrFlags)
	for _, command := range op2.Commands {
		mlog.Info(`Command
		Src: %s
		Dst: %s
		Entry: %s
		Size: %d
		Transferred: %d
		`, command.Src, command.Dst, command.Entry, command.Size, command.Transferred)
	}
}

func (c *Core) scatterMove(msg *pubsub.Message) {
	c.state.Status = common.OP_SCATTER_MOVE
	c.state.Operation = c.setupOperation(c.state.Status, c.state.Operation)
	go c.runOperation("Move")
}

func (c *Core) scatterCopy(msg *pubsub.Message) {
	c.state.Status = common.OP_SCATTER_COPY
	c.state.Operation = c.setupOperation(c.state.Status, c.state.Operation)
	go c.runOperation("Copy")
}

func (c *Core) setupOperation(status int64, copyOperation *domain.Operation) *domain.Operation {
	operation := &domain.Operation{
		OpKind:          status,
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
	c.state.History = append(c.state.History, c.state.Operation)

	b, err := json.Marshal(c.state.Operation)
	if err != nil {
		mlog.Warning("Unable to serialize op: %s", err)
	}

	mlog.Info(`Serialized
	%s
	`, string(b))

	c.state.Status = common.OP_NEUTRAL
	c.state.Operation = resetOp(c.state.Unraid.Disks)
}

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

func resetOp(disks []*domain.Disk) *domain.Operation {
	op := &domain.Operation{
		OpKind: common.OP_NEUTRAL,
		VDisks: make(map[string]*domain.VDisk, 0),
	}

	for _, disk := range disks {
		vdisk := &domain.VDisk{Path: disk.Path, PlannedFree: disk.Free, Src: false, Dst: false}
		op.VDisks[disk.Path] = vdisk
	}

	return op
}

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

// func (c *Core) getStorage(msg *pubsub.Message) {
// 	var stats string

// 	if c.operation.OpState == model.StateIdle {
// 		c.storage.Refresh()
// 	} else if c.operation.OpState == model.StateMove || c.operation.OpState == model.StateCopy || c.operation.OpState == model.StateGather {
// 		percent, left, speed := progress(c.operation.BytesToTransfer, c.operation.BytesTransferred, time.Since(c.operation.Started))
// 		stats = fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)
// 	}

// 	c.storage.Stats = stats
// 	c.storage.OpState = c.operation.OpState
// 	c.storage.PrevState = c.operation.PrevState
// 	c.storage.BytesToTransfer = c.operation.BytesToTransfer

// 	msg.Reply <- c.storage
// }

// func (c *Core) getUpdate(msg *pubsub.Message) {
// 	var newest string

// 	if c.settings.CheckForUpdate == 1 {
// 		latest, err := lib.GetLatestVersion("https://raw.githubusercontent.com/jbrodriguez/unbalance/master/VERSION")
// 		if err != nil {
// 			return
// 		}

// 		latest = strings.TrimSuffix(latest, "\n")
// 		if version.Compare(latest, c.settings.Version, ">") {
// 			newest = latest
// 		}
// 	}

// 	msg.Reply <- newest
// }

func (c *Core) toggleDryRun(msg *pubsub.Message) {
	mlog.Info("Toggling dryRun from (%t)", c.settings.DryRun)

	c.settings.ToggleDryRun()
	c.settings.Save()

	msg.Reply <- &c.settings.Config
}

// func (c *Core) getTree(msg *pubsub.Message) {
// 	path := msg.Payload.(string)

// 	msg.Reply <- c.storage.GetTree(path)
// }

// func (c *Core) locate(msg *pubsub.Message) {
// 	chosen := msg.Payload.([]string)
// 	msg.Reply <- c.storage.Locate(chosen)
// }

// func (c *Core) setRsyncFlags(msg *pubsub.Message) {
// 	// mlog.Warning("payload: %+v", msg.Payload)
// 	payload, ok := msg.Payload.(string)
// 	if !ok {
// 		mlog.Warning("Unable to convert Rsync Flags parameters")
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert Rsync Flags parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		msg.Reply <- &c.settings.Config

// 		return
// 	}

// 	var rsync dto.Rsync
// 	err := json.Unmarshal([]byte(payload), &rsync)
// 	if err != nil {
// 		mlog.Warning("Unable to bind rsyncFlags parameters: %s", err)
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind rsyncFlags parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 		// mlog.Fatalf(err.Error())
// 	}

// 	mlog.Info("Setting rsyncFlags to (%s)", strings.Join(rsync.Flags, " "))

// 	c.settings.RsyncFlags = rsync.Flags
// 	c.settings.Save()

// 	msg.Reply <- &c.settings.Config
// }

// func (c *Core) calc(msg *pubsub.Message) {
// 	c.operation = model.Operation{OpState: model.StateCalc, PrevState: model.StateIdle}

// 	go c._calc(msg)
// }

// func (c *Core) _calc(msg *pubsub.Message) {
// 	defer func() { c.operation.OpState = model.StateIdle }()

// 	payload, ok := msg.Payload.(string)
// 	if !ok {
// 		mlog.Warning("Unable to convert calculate parameters")
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert calculate parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 	}

// 	var dtoCalc dto.Calculate
// 	err := json.Unmarshal([]byte(payload), &dtoCalc)
// 	if err != nil {
// 		mlog.Warning("Unable to bind calculate parameters: %s", err)
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind calculate parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 		// mlog.Fatalf(err.Error())
// 	}

// 	mlog.Info("Running calculate operation ...")
// 	c.operation.Started = time.Now()

// 	outbound := &dto.Packet{Topic: "calcStarted", Payload: "Operation started"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	disks := make([]*model.Disk, 0)

// 	// create array of destination disks
// 	var srcDisk *model.Disk
// 	for _, disk := range c.storage.Disks {
// 		// reset disk
// 		disk.NewFree = disk.Free
// 		disk.Bin = nil
// 		disk.Src = false
// 		disk.Dst = dtoCalc.DestDisks[disk.Path]

// 		if disk.Path == dtoCalc.SourceDisk {
// 			disk.Src = true
// 			srcDisk = disk
// 		} else {
// 			// add it to the target disk list, only if the user selected it
// 			if val, ok := dtoCalc.DestDisks[disk.Path]; ok && val {
// 				// double check, if it's a cache disk, make sure it's the main cache disk
// 				if disk.Type == "Cache" && len(disk.Name) > 5 {
// 					continue
// 				}

// 				disks = append(disks, disk)
// 			}
// 		}
// 	}

// 	mlog.Info("_calc:Begin:srcDisk(%s); dstDisks(%d)", srcDisk.Path, len(disks))

// 	for _, disk := range disks {
// 		mlog.Info("_calc:elegibleDestDisk(%s)", disk.Path)
// 	}

// 	sort.Sort(model.ByFree(disks))

// 	srcDiskWithoutMnt := srcDisk.Path[5:]

// 	owner := ""
// 	lib.Shell("id -un", mlog.Warning, "owner", "", func(line string) {
// 		owner = line
// 	})

// 	group := ""
// 	lib.Shell("id -gn", mlog.Warning, "group", "", func(line string) {
// 		group = line
// 	})

// 	c.operation.OwnerIssue = 0
// 	c.operation.GroupIssue = 0
// 	c.operation.FolderIssue = 0
// 	c.operation.FileIssue = 0

// 	// Check permission and gather folders to be transferred from
// 	// source disk
// 	folders := make([]*model.Item, 0)
// 	for _, path := range dtoCalc.Folders {
// 		msg := fmt.Sprintf("Scanning %s on %s", path, srcDiskWithoutMnt)
// 		outbound = &dto.Packet{Topic: "calcProgress", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		c.checkOwnerAndPermissions(&c.operation, dtoCalc.SourceDisk, path, owner, group)

// 		msg = "Checked permissions ..."
// 		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		mlog.Info("_calc:%s", msg)

// 		list := c.getFolders(dtoCalc.SourceDisk, path)
// 		if list != nil {
// 			folders = append(folders, list...)
// 		}
// 	}

// 	mlog.Info("_calc:totalIssues=owner(%d),group(%d),folder(%d),file(%d)", c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)
// 	mlog.Info("_calc:foldersToBeTransferredTotal(%d)", len(folders))

// 	for _, v := range folders {
// 		mlog.Info("_calc:toBeTransferred:Path(%s); Size(%s)", v.Path, lib.ByteSize(v.Size))
// 	}

// 	willBeTransferred := make([]*model.Item, 0)
// 	if len(folders) > 0 {
// 		// Initialize fields
// 		// c.storage.BytesToTransfer = 0
// 		// c.storage.SourceDiskName = srcDisk.Path
// 		c.operation.BytesToTransfer = 0
// 		c.operation.SourceDiskName = srcDisk.Path

// 		for _, disk := range disks {
// 			diskWithoutMnt := disk.Path[5:]
// 			msg := fmt.Sprintf("Trying to allocate folders to %s ...", diskWithoutMnt)
// 			outbound = &dto.Packet{Topic: "calcProgress", Payload: msg}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			mlog.Info("_calc:%s", msg)
// 			// time.Sleep(2 * time.Second)

// 			if disk.Path != srcDisk.Path {
// 				// disk.NewFree = disk.Free

// 				var reserved int64
// 				switch c.settings.ReservedUnit {
// 				case "%":
// 					fcalc := disk.Size * c.settings.ReservedAmount / 100
// 					reserved = int64(fcalc)
// 					break
// 				case "Mb":
// 					reserved = c.settings.ReservedAmount * 1000 * 1000
// 					break
// 				case "Gb":
// 					reserved = c.settings.ReservedAmount * 1000 * 1000 * 1000
// 					break
// 				default:
// 					reserved = lib.ReservedSpace
// 				}

// 				ceil := lib.Max(lib.ReservedSpace, reserved)
// 				mlog.Info("_calc:FoldersLeft(%d):ReservedSpace(%d)", len(folders), ceil)

// 				packer := algorithm.NewKnapsack(disk, folders, ceil)
// 				bin := packer.BestFit()
// 				if bin != nil {
// 					srcDisk.NewFree += bin.Size
// 					disk.NewFree -= bin.Size
// 					c.operation.BytesToTransfer += bin.Size
// 					// c.storage.BytesToTransfer += bin.Size

// 					willBeTransferred = append(willBeTransferred, bin.Items...)
// 					folders = c.removeFolders(folders, bin.Items)

// 					mlog.Info("_calc:BinAllocated=[Disk(%s); Items(%d)];Freespace=[original(%s); final(%s)]", disk.Path, len(bin.Items), lib.ByteSize(srcDisk.Free), lib.ByteSize(srcDisk.NewFree))
// 				} else {
// 					mlog.Info("_calc:NoBinAllocated=Disk(%s)", disk.Path)
// 				}
// 			}
// 		}
// 	}

// 	c.operation.Finished = time.Now()
// 	elapsed := lib.Round(time.Since(c.operation.Started), time.Millisecond)

// 	fstarted := c.operation.Started.Format(timeFormat)
// 	ffinished := c.operation.Finished.Format(timeFormat)

// 	// Send to frontend console started/ended/elapsed times
// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	if len(willBeTransferred) == 0 {
// 		mlog.Info("_calc:No folders can be transferred.")
// 	} else {
// 		mlog.Info("_calc:%d folders will be transferred.", len(willBeTransferred))
// 		for _, folder := range willBeTransferred {
// 			mlog.Info("_calc:willBeTransferred(%s)", folder.Path)
// 		}
// 	}

// 	// send to frontend the folders that will not be transferred, if any
// 	// notTransferred holds a string representation of all the folders, separated by a '\n'
// 	c.operation.FoldersNotTransferred = make([]string, 0)
// 	notTransferred := ""
// 	if len(folders) > 0 {
// 		outbound = &dto.Packet{Topic: "calcProgress", Payload: "The following folders will not be transferred, because there's not enough space in the target disks:\n"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		mlog.Info("_calc:%d folders will NOT be transferred.", len(folders))
// 		for _, folder := range folders {
// 			c.operation.FoldersNotTransferred = append(c.operation.FoldersNotTransferred, folder.Path)

// 			notTransferred += folder.Path + "\n"

// 			outbound = &dto.Packet{Topic: "calcProgress", Payload: folder.Path}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			mlog.Info("_calc:notTransferred(%s)", folder.Path)
// 		}
// 	}

// 	// send mail according to user preferences
// 	subject := "unBALANCE - CALCULATE operation completed"
// 	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", fstarted, ffinished, elapsed)
// 	if notTransferred != "" {
// 		switch c.settings.NotifyCalc {
// 		case 1:
// 			message += "\n\nSome folders will not be transferred because there's not enough space for them in any of the destination disks."
// 		case 2:
// 			message += "\n\nThe following folders will not be transferred because there's not enough space for them in any of the destination disks:\n\n" + notTransferred
// 		}
// 	}

// 	if c.operation.OwnerIssue > 0 || c.operation.GroupIssue > 0 || c.operation.FolderIssue > 0 || c.operation.FileIssue > 0 {
// 		message += fmt.Sprintf(`
// 			\n\nThere are some permission issues:
// 			\n\n%d file(s)/folder(s) with an owner other than 'nobody'
// 			\n%d file(s)/folder(s) with a group other than 'users'
// 			\n%d folder(s) with a permission other than 'drwxrwxrwx'
// 			\n%d files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'
// 			\n\nCheck the log file (/boot/logs/unbalance.log) for additional information
// 			\n\nIt's strongly suggested to install the Fix Common Plugins and run the Docker Safe New Permissions command
// 		`, c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)
// 	}

// 	if sendErr := c.sendmail(c.settings.NotifyCalc, subject, message, false); sendErr != nil {
// 		mlog.Error(sendErr)
// 	}

// 	// some local logging
// 	mlog.Info("_calc:FoldersLeft(%d)", len(folders))
// 	mlog.Info("_calc:src(%s):Listing (%d) disks ...", srcDisk.Path, len(c.storage.Disks))
// 	for _, disk := range c.storage.Disks {
// 		// mlog.Info("the mystery of the year(%s)", disk.Path)
// 		disk.Print()
// 	}

// 	mlog.Info("=========================================================")
// 	mlog.Info("Results for %s", srcDisk.Path)
// 	mlog.Info("Original Free Space: %s", lib.ByteSize(srcDisk.Free))
// 	mlog.Info("Final Free Space: %s", lib.ByteSize(srcDisk.NewFree))
// 	mlog.Info("Gained Space: %s", lib.ByteSize(srcDisk.NewFree-srcDisk.Free))
// 	mlog.Info("Bytes To Move: %s", lib.ByteSize(c.operation.BytesToTransfer))
// 	mlog.Info("---------------------------------------------------------")

// 	c.storage.Print()
// 	// msg.Reply <- c.storage

// 	mlog.Info("_calc:End:srcDisk(%s)", srcDisk.Path)

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: "Operation Finished"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	c.storage.BytesToTransfer = c.operation.BytesToTransfer
// 	c.storage.OpState = c.operation.OpState
// 	c.storage.PrevState = c.operation.PrevState

// 	// send to front end the signal of operation finished
// 	outbound = &dto.Packet{Topic: "calcFinished", Payload: c.storage}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
// 	// and there actually perm issues
// 	if c.operation.BytesToTransfer > 0 && (c.operation.OwnerIssue+c.operation.GroupIssue+c.operation.FolderIssue+c.operation.FileIssue > 0) {
// 		outbound = &dto.Packet{Topic: "calcPermIssue", Payload: fmt.Sprintf("%d|%d|%d|%d", c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	}
// }

// func (c *Core) getFolders(src string, folder string) (items []*model.Item) {
// 	srcFolder := filepath.Join(src, folder)

// 	mlog.Info("getFolders:Scanning disk(%s):folder(%s)", src, folder)

// 	var fi os.FileInfo
// 	var err error
// 	if fi, err = os.Stat(srcFolder); os.IsNotExist(err) {
// 		mlog.Warning("getFolders:Folder does not exist: %s", srcFolder)
// 		return nil
// 	}

// 	if !fi.IsDir() {
// 		mlog.Info("getFolder-found(%s)-size(%d)", srcFolder, fi.Size())

// 		item := &model.Item{Name: folder, Size: fi.Size(), Path: folder, Location: src}
// 		items = append(items, item)

// 		msg := fmt.Sprintf("Found %s (%s)", item.Name, lib.ByteSize(item.Size))
// 		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		return
// 	}

// 	dirs, err := ioutil.ReadDir(srcFolder)
// 	if err != nil {
// 		mlog.Warning("getFolders:Unable to readdir: %s", err)
// 	}

// 	mlog.Info("getFolders:Readdir(%d)", len(dirs))

// 	if len(dirs) == 0 {
// 		mlog.Info("getFolders:No subdirectories under %s", srcFolder)
// 		return nil
// 	}

// 	scanFolder := srcFolder + "/."
// 	cmdText := fmt.Sprintf("find \"%s\" ! -name . -prune -exec du -bs {} +", scanFolder)

// 	mlog.Info("getFolders:Executing %s", cmdText)

// 	lib.Shell(cmdText, mlog.Warning, "getFolders:find/du:", "", func(line string) {
// 		mlog.Info("getFolders:find(%s): %s", scanFolder, line)

// 		result := c.reItems.FindStringSubmatch(line)
// 		// mlog.Info("[%s] %s", result[1], result[2])

// 		size, _ := strconv.ParseInt(result[1], 10, 64)

// 		item := &model.Item{Name: result[2], Size: size, Path: filepath.Join(folder, filepath.Base(result[2])), Location: src}
// 		items = append(items, item)

// 		msg := fmt.Sprintf("Found %s (%s)", filepath.Base(item.Name), lib.ByteSize(size))
// 		outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	})

// 	return
// }

// func (c *Core) checkOwnerAndPermissions(operation *model.Operation, src, folder, ownerName, groupName string) {
// 	srcFolder := filepath.Join(src, folder)

// 	// outbound := &dto.Packet{Topic: "calcProgress", Payload: "Checking permissions ..."}
// 	// c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	ownerIssue, groupIssue, folderIssue, fileIssue := 0, 0, 0, 0

// 	mlog.Info("perms:Scanning disk(%s):folder(%s)", src, folder)

// 	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
// 		mlog.Warning("perms:Folder does not exist: %s", srcFolder)
// 		return
// 	}

// 	scanFolder := srcFolder + "/."
// 	cmdText := fmt.Sprintf(`find "%s" -exec stat --format "%%A|%%U:%%G|%%F|%%n" {} \;`, scanFolder)

// 	mlog.Info("perms:Executing %s", cmdText)

// 	lib.Shell(cmdText, mlog.Warning, "perms:find/stat:", "", func(line string) {
// 		result := c.reStat.FindStringSubmatch(line)
// 		if result == nil {
// 			mlog.Warning("perms:Unable to parse (%s)", line)
// 			return
// 		}

// 		u := result[1]
// 		g := result[2]
// 		o := result[3]
// 		user := result[4]
// 		group := result[5]
// 		kind := result[6]
// 		name := result[7]

// 		perms := u + g + o

// 		if user != "nobody" {
// 			if c.settings.Verbosity == 1 {
// 				mlog.Info("perms:User != nobody: [%s]: %s", user, name)
// 			}

// 			operation.OwnerIssue++
// 			ownerIssue++
// 		}

// 		if group != "users" {
// 			if c.settings.Verbosity == 1 {
// 				mlog.Info("perms:Group != users: [%s]: %s", group, name)
// 			}

// 			operation.GroupIssue++
// 			groupIssue++
// 		}

// 		if kind == "directory" {
// 			if perms != "rwxrwxrwx" {
// 				if c.settings.Verbosity == 1 {
// 					mlog.Info("perms:Folder perms != rwxrwxrwx: [%s]: %s", perms, name)
// 				}

// 				operation.FolderIssue++
// 				folderIssue++
// 			}
// 		} else {
// 			match := strings.Compare(perms, "r--r--r--") == 0 || strings.Compare(perms, "rw-rw-rw-") == 0
// 			if !match {
// 				if c.settings.Verbosity == 1 {
// 					mlog.Info("perms:File perms != rw-rw-rw- or r--r--r--: [%s]: %s", perms, name)
// 				}

// 				operation.FileIssue++
// 				fileIssue++
// 			}
// 		}
// 	})

// 	mlog.Info("perms:issues=owner(%d),group(%d),folder(%d),file(%d)", ownerIssue, groupIssue, folderIssue, fileIssue)

// 	return
// }

// func (c *Core) removeFolders(folders []*model.Item, list []*model.Item) []*model.Item {
// 	w := 0 // write index

// loop:
// 	for _, fld := range folders {
// 		for _, itm := range list {
// 			if itm.Name == fld.Name {
// 				continue loop
// 			}
// 		}
// 		folders[w] = fld
// 		w++
// 	}

// 	return folders[:w]
// }

// func (c *Core) move(msg *pubsub.Message) {
// 	c.operation.OpState = model.StateMove
// 	c.operation.PrevState = model.StateMove
// 	go c.transfer("Move", false, msg)
// }

// func (c *Core) copy(msg *pubsub.Message) {
// 	c.operation.OpState = model.StateCopy
// 	c.operation.PrevState = model.StateCopy
// 	go c.transfer("Copy", false, msg)
// }

// func (c *Core) validate(msg *pubsub.Message) {
// 	c.operation.OpState = model.StateValidate
// 	c.operation.PrevState = model.StateValidate
// 	go c.checksum(msg)
// }

// func (c *Core) transfer(opName string, multiSource bool, msg *pubsub.Message) {
// 	defer func() {
// 		c.operation.OpState = model.StateIdle
// 		c.operation.Started = time.Time{}
// 		c.operation.BytesTransferred = 0
// 		c.operation.Target = ""
// 	}()

// 	mlog.Info("Running %s operation ...", opName)
// 	c.operation.Started = time.Now()

// 	outbound := &dto.Packet{Topic: "transferStarted", Payload: "Operation started"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "progressStats", Payload: "Waiting to collect stats ..."}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// user may have changed rsync flags or dry-run setting, adjust for it
// 	c.operation.RsyncFlags = c.settings.RsyncFlags
// 	c.operation.DryRun = c.settings.DryRun
// 	if c.operation.DryRun {
// 		c.operation.RsyncFlags = append(c.operation.RsyncFlags, "--dry-run")
// 		mlog.Info("dry-run ON")
// 	} else {
// 		mlog.Info("dry-run OFF")
// 	}
// 	c.operation.RsyncStrFlags = strings.Join(c.operation.RsyncFlags, " ")

// 	workdir := c.operation.SourceDiskName

// 	c.operation.Commands = make([]model.Command, 0)
// 	for _, disk := range c.storage.Disks {
// 		if disk.Bin == nil || disk.Src {
// 			continue
// 		}

// 		for _, item := range disk.Bin.Items {
// 			var src, dst string
// 			if strings.Contains(c.operation.RsyncStrFlags, "R") {
// 				if item.Path[0] == filepath.Separator {
// 					src = item.Path[1:]
// 				} else {
// 					src = item.Path
// 				}

// 				dst = disk.Path + string(filepath.Separator)
// 			} else {
// 				src = filepath.Join(c.operation.SourceDiskName, item.Path)
// 				dst = filepath.Join(disk.Path, filepath.Dir(item.Path)) + string(filepath.Separator)
// 			}

// 			if multiSource {
// 				workdir = item.Location
// 			}

// 			c.operation.Commands = append(c.operation.Commands, model.Command{
// 				Src:     src,
// 				Dst:     dst,
// 				Path:    item.Path,
// 				Size:    item.Size,
// 				WorkDir: workdir,
// 			})
// 		}
// 	}

// 	if c.settings.NotifyMove == 2 {
// 		c.notifyCommandsToRun(opName)
// 	}

// 	// execute each rsync command created in the step above
// 	c.runOperation(opName, c.operation.RsyncFlags, c.operation.RsyncStrFlags, multiSource)
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

// func (c *Core) runOperation(opName string, rsyncFlags []string, rsyncStrFlags string, multiSource bool) {
// 	// Initialize local variables
// 	var calls int64
// 	var callsPerDelta int64

// 	var finished time.Time
// 	var elapsed time.Duration

// 	commandsExecuted := make([]string, 0)

// 	c.operation.BytesTransferred = 0

// 	for _, command := range c.operation.Commands {
// 		args := append(
// 			rsyncFlags,
// 			command.Src,
// 			command.Dst,
// 		)
// 		cmd := fmt.Sprintf(`rsync %s %s %s`, rsyncStrFlags, strconv.Quote(command.Src), strconv.Quote(command.Dst))
// 		mlog.Info("Command Started: (src: %s) %s ", command.WorkDir, cmd)

// 		outbound := &dto.Packet{Topic: "transferProgress", Payload: cmd}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		bytesTransferred := c.operation.BytesTransferred

// 		var deltaMoved int64

// 		// actual shell execution
// 		err := lib.ShellEx(func(text string) {
// 			line := strings.TrimSpace(text)

// 			if len(line) <= 0 {
// 				return
// 			}

// 			if callsPerDelta <= 50 {
// 				calls++
// 			}

// 			delta := int64(time.Since(c.operation.Started) / time.Second)
// 			if delta == 0 {
// 				delta = 1
// 			}
// 			callsPerDelta = calls / delta

// 			match := c.reProgress.FindStringSubmatch(line)
// 			if match == nil {
// 				// this is a regular output line from rsync, print it
// 				// according to verbosity settings
// 				if c.settings.Verbosity == 1 {
// 					mlog.Info("%s", line)
// 				}

// 				if callsPerDelta <= 50 {
// 					outbound := &dto.Packet{Topic: "transferProgress", Payload: line}
// 					c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 				}

// 				return
// 			}

// 			// this is a file transfer progress output line
// 			if match[1] == "" {
// 				// this happens when the file hasn't finished transferring
// 				moved := strings.Replace(match[2], ",", "", -1)
// 				deltaMoved, _ = strconv.ParseInt(moved, 10, 64)
// 			} else {
// 				// the file has finished transferring
// 				moved := strings.Replace(match[1], ",", "", -1)
// 				deltaMoved, _ = strconv.ParseInt(moved, 10, 64)
// 				bytesTransferred += deltaMoved
// 			}

// 			percent, left, speed := progress(c.operation.BytesToTransfer, bytesTransferred+deltaMoved, time.Since(c.operation.Started))
// 			msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)

// 			if callsPerDelta <= 50 {
// 				outbound := &dto.Packet{Topic: "progressStats", Payload: msg}
// 				c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			}

// 		}, mlog.Warning, command.WorkDir, "rsync", args...)

// 		finished = time.Now()
// 		elapsed = time.Since(c.operation.Started)

// 		if err != nil {
// 			subject := fmt.Sprintf("unBALANCE - %s operation INTERRUPTED", strings.ToUpper(opName))
// 			headline := fmt.Sprintf("Command Interrupted: %s (%s)", cmd, err.Error()+" : "+getError(err.Error(), c.reRsync, c.rsyncErrors))

// 			mlog.Warning(headline)
// 			outbound := &dto.Packet{Topic: "opError", Payload: fmt.Sprintf("%s operation was interrupted. Check log (/boot/logs/unbalance.log) for details.", opName)}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 			_, _, speed := progress(c.operation.BytesToTransfer, bytesTransferred+deltaMoved, elapsed)
// 			c.finishTransferOperation(subject, headline, commandsExecuted, c.operation.Started, finished, elapsed, bytesTransferred+deltaMoved, speed, multiSource)

// 			return
// 		}

// 		mlog.Info("Command Finished")

// 		c.operation.BytesTransferred = c.operation.BytesTransferred + command.Size
// 		percent, left, speed := progress(c.operation.BytesToTransfer, c.operation.BytesTransferred, elapsed)

// 		msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, left, speed)
// 		mlog.Info("Current progress: %s", msg)

// 		outbound = &dto.Packet{Topic: "progressStats", Payload: msg}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 		commandsExecuted = append(commandsExecuted, cmd)

// 		if c.operation.DryRun && c.operation.OpState == model.StateGather {
// 			parent := filepath.Dir(command.Src)
// 			// mlog.Info("parent(%s)-workdir(%s)-src(%s)-dst(%s)-path(%s)", parent, command.WorkDir, command.Src, command.Dst, command.Path)
// 			if parent != "." {
// 				mlog.Info(`Would delete empty folders starting from (%s) - (find "%s" -type d -empty -prune -exec rm -rf {} \;) `, filepath.Join(command.WorkDir, parent), filepath.Join(command.WorkDir, parent))
// 			} else {
// 				mlog.Info(`WONT DELETE: find "%s" -type d -empty -prune -exec rm -rf {} \;`, filepath.Join(command.WorkDir, parent))
// 			}
// 		}

// 		// if it isn't a dry-run and the operation is Move or Gather, delete the source folder
// 		if !c.operation.DryRun && (c.operation.OpState == model.StateMove || c.operation.OpState == model.StateGather) {
// 			exists, _ := lib.Exists(filepath.Join(command.Dst, command.Src))
// 			if exists {
// 				rmrf := fmt.Sprintf("rm -rf \"%s\"", filepath.Join(command.WorkDir, command.Path))
// 				mlog.Info("Removing: %s", rmrf)
// 				err = lib.Shell(rmrf, mlog.Warning, "transferProgress:", "", func(line string) {
// 					mlog.Info(line)
// 				})

// 				if err != nil {
// 					msg := fmt.Sprintf("Unable to remove source folder (%s): %s", filepath.Join(command.WorkDir, command.Path), err)

// 					outbound := &dto.Packet{Topic: "transferProgress", Payload: msg}
// 					c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 					mlog.Warning(msg)
// 				}

// 				if c.operation.OpState == model.StateGather {
// 					parent := filepath.Dir(command.Src)
// 					if parent != "." {
// 						rmdir := fmt.Sprintf(`find "%s" -type d -empty -prune -exec rm -rf {} \;`, filepath.Join(command.WorkDir, parent))
// 						mlog.Info("Running %s", rmdir)

// 						err = lib.Shell(rmdir, mlog.Warning, "transferProgress:", "", func(line string) {
// 							mlog.Info(line)
// 						})

// 						if err != nil {
// 							msg := fmt.Sprintf("Unable to remove parent folder (%s): %s", filepath.Join(command.WorkDir, parent), err)

// 							outbound := &dto.Packet{Topic: "transferProgress", Payload: msg}
// 							c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 							mlog.Warning(msg)
// 						}
// 					}
// 				}
// 			} else {
// 				mlog.Warning("Skipping deletion (file/folder not present in destination): %s", filepath.Join(command.Dst, command.Src))
// 			}
// 		}
// 	}

// 	subject := fmt.Sprintf("unBALANCE - %s operation completed", strings.ToUpper(opName))
// 	headline := fmt.Sprintf("%s operation has finished", opName)

// 	_, _, speed := progress(c.operation.BytesToTransfer, c.operation.BytesTransferred, elapsed)
// 	c.finishTransferOperation(subject, headline, commandsExecuted, c.operation.Started, finished, elapsed, c.operation.BytesTransferred, speed, multiSource)
// }

// func (c *Core) finishTransferOperation(subject, headline string, commands []string, started, finished time.Time, elapsed time.Duration, transferred int64, speed float64, multiSource bool) {
// 	fstarted := started.Format(timeFormat)
// 	ffinished := finished.Format(timeFormat)
// 	elapsed = lib.Round(time.Since(started), time.Millisecond)

// 	outbound := &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: fmt.Sprintf("Transferred %s at ~ %.2f MB/s", lib.ByteSize(transferred), speed)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: headline}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: "These are the commands that were executed:"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	printedCommands := ""
// 	for _, command := range commands {
// 		printedCommands += command + "\n"
// 		outbound = &dto.Packet{Topic: "transferProgress", Payload: command}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	}

// 	outbound = &dto.Packet{Topic: "transferProgress", Payload: "Operation Finished"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// send to front end the signal of operation finished
// 	if c.settings.DryRun {
// 		outbound = &dto.Packet{Topic: "transferProgress", Payload: "--- IT WAS A DRY RUN ---"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	}

// 	finishMsg := "transferFinished"
// 	if multiSource {
// 		finishMsg = "gatherFinished"
// 	}

// 	outbound = &dto.Packet{Topic: finishMsg, Payload: ""}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s\n\n%s\n\nTransferred %s at ~ %.2f MB/s", fstarted, ffinished, elapsed, headline, lib.ByteSize(transferred), speed)
// 	switch c.settings.NotifyMove {
// 	case 1:
// 		message += fmt.Sprintf("\n\n%d commands were executed.", len(commands))
// 	case 2:
// 		message += "\n\nThese are the commands that were executed:\n\n" + printedCommands
// 	}

// 	go func() {
// 		if sendErr := c.sendmail(c.settings.NotifyMove, subject, message, c.settings.DryRun); sendErr != nil {
// 			mlog.Error(sendErr)
// 		}
// 	}()

// 	mlog.Info("\n%s\n%s", subject, message)
// }

// func (c *Core) findTargets(msg *pubsub.Message) {
// 	c.operation = model.Operation{OpState: model.StateFindTargets, PrevState: model.StateIdle}
// 	go c._findTargets(msg)
// }

// func (c *Core) _findTargets(msg *pubsub.Message) {
// 	defer func() { c.operation.OpState = model.StateIdle }()

// 	data, ok := msg.Payload.(string)
// 	if !ok {
// 		mlog.Warning("Unable to convert findTargets parameters")
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert findTargets parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 	}

// 	var chosen []string
// 	err := json.Unmarshal([]byte(data), &chosen)
// 	if err != nil {
// 		mlog.Warning("Unable to bind findTargets parameters: %s", err)
// 		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind findTargets parameters"}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 		return
// 		// mlog.Fatalf(err.Error())
// 	}

// 	mlog.Info("Running findTargets operation ...")
// 	c.operation.Started = time.Now()

// 	outbound := &dto.Packet{Topic: "calcStarted", Payload: "Operation started"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// disks := make([]*model.Disk, 0)
// 	c.storage.Refresh()

// 	owner := ""
// 	lib.Shell("id -un", mlog.Warning, "owner", "", func(line string) {
// 		owner = line
// 	})

// 	group := ""
// 	lib.Shell("id -gn", mlog.Warning, "group", "", func(line string) {
// 		group = line
// 	})

// 	c.operation.OwnerIssue = 0
// 	c.operation.GroupIssue = 0
// 	c.operation.FolderIssue = 0
// 	c.operation.FileIssue = 0

// 	entries := make([]*model.Item, 0)

// 	// Check permission and look for the chosen folders on every disk
// 	for _, disk := range c.storage.Disks {
// 		for _, path := range chosen {
// 			msg := fmt.Sprintf("Scanning %s on %s", path, disk.Path)
// 			outbound = &dto.Packet{Topic: "calcProgress", Payload: msg}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			mlog.Info("_find:%s", msg)

// 			c.checkOwnerAndPermissions(&c.operation, disk.Path, path, owner, group)

// 			msg = "Checked permissions ..."
// 			outbound := &dto.Packet{Topic: "calcProgress", Payload: msg}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			mlog.Info("_find:%s", msg)

// 			list := c.getFolders(disk.Path, path)
// 			if list != nil {
// 				entries = append(entries, list...)
// 			}
// 		}
// 	}

// 	mlog.Info("_find:elegibleFolders(%d)", len(entries))

// 	var totalSize int64
// 	for _, entry := range entries {
// 		totalSize += entry.Size
// 		mlog.Info("_find:elegibleFolder:Location(%s); Size(%s)", filepath.Join(entry.Location, entry.Path), lib.ByteSize(entry.Size))
// 	}

// 	mlog.Info("_find:potentialSizeToBeTransferred(%s)", lib.ByteSize(totalSize))

// 	if len(entries) > 0 {
// 		// Initialize fields
// 		// c.storage.BytesToTransfer = 0
// 		// c.storage.SourceDiskName = srcDisk.Path
// 		c.operation.BytesToTransfer = 0
// 		// c.operation.SourceDiskName = mntUser

// 		for _, disk := range c.storage.Disks {
// 			diskWithoutMnt := disk.Path[5:]
// 			msg := fmt.Sprintf("Trying to allocate folders to %s ...", diskWithoutMnt)
// 			outbound = &dto.Packet{Topic: "calcProgress", Payload: msg}
// 			c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 			mlog.Info("_find:%s", msg)
// 			// time.Sleep(2 * time.Second)

// 			var reserved int64
// 			switch c.settings.ReservedUnit {
// 			case "%":
// 				fcalc := disk.Size * c.settings.ReservedAmount / 100
// 				reserved = int64(fcalc)
// 				break
// 			case "Mb":
// 				reserved = c.settings.ReservedAmount * 1000 * 1000
// 				break
// 			case "Gb":
// 				reserved = c.settings.ReservedAmount * 1000 * 1000 * 1000
// 				break
// 			default:
// 				reserved = lib.ReservedSpace
// 			}

// 			ceil := lib.Max(lib.ReservedSpace, reserved)
// 			mlog.Info("_find:FoldersLeft(%d):ReservedSpace(%d)", len(entries), ceil)

// 			packer := algorithm.NewGreedy(disk, entries, totalSize, ceil)
// 			bin := packer.FitAll()
// 			if bin != nil {
// 				disk.NewFree -= bin.Size
// 				disk.Src = false
// 				disk.Dst = false
// 				c.operation.BytesToTransfer += bin.Size
// 				mlog.Info("_find:BinAllocated=[Disk(%s); Items(%d)]", disk.Path, len(bin.Items))
// 			} else {
// 				mlog.Info("_find:NoBinAllocated=Disk(%s)", disk.Path)
// 			}
// 		}
// 	}

// 	c.operation.Finished = time.Now()
// 	elapsed := lib.Round(time.Since(c.operation.Started), time.Millisecond)

// 	fstarted := c.operation.Started.Format(timeFormat)
// 	ffinished := c.operation.Finished.Format(timeFormat)

// 	// Send to frontend console started/ended/elapsed times
// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Started: %s", fstarted)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Ended: %s", ffinished)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: fmt.Sprintf("Elapsed: %s", elapsed)}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// send to frontend the folders that will not be transferred, if any
// 	// notTransferred holds a string representation of all the folders, separated by a '\n'
// 	c.operation.FoldersNotTransferred = make([]string, 0)

// 	// send mail according to user preferences
// 	subject := "unBALANCE - CALCULATE operation completed"
// 	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s", fstarted, ffinished, elapsed)

// 	if c.operation.OwnerIssue > 0 || c.operation.GroupIssue > 0 || c.operation.FolderIssue > 0 || c.operation.FileIssue > 0 {
// 		message += fmt.Sprintf(`
// 			\n\nThere are some permission issues:
// 			\n\n%d file(s)/folder(s) with an owner other than 'nobody'
// 			\n%d file(s)/folder(s) with a group other than 'users'
// 			\n%d folder(s) with a permission other than 'drwxrwxrwx'
// 			\n%d files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'
// 			\n\nCheck the log file (/boot/logs/unbalance.log) for additional information
// 			\n\nIt's strongly suggested to install the Fix Common Plugins and run the Docker Safe New Permissions command
// 		`, c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)
// 	}

// 	if sendErr := c.sendmail(c.settings.NotifyCalc, subject, message, false); sendErr != nil {
// 		mlog.Error(sendErr)
// 	}

// 	// some local logging
// 	mlog.Info("_find:Listing (%d) disks ...", len(c.storage.Disks))
// 	for _, disk := range c.storage.Disks {
// 		// mlog.Info("the mystery of the year(%s)", disk.Path)
// 		disk.Print()
// 	}

// 	c.storage.Print()

// 	outbound = &dto.Packet{Topic: "calcProgress", Payload: "Operation Finished"}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	c.storage.BytesToTransfer = c.operation.BytesToTransfer
// 	c.storage.OpState = c.operation.OpState
// 	c.storage.PrevState = c.operation.PrevState

// 	// send to front end the signal of operation finished
// 	outbound = &dto.Packet{Topic: "findFinished", Payload: c.storage}
// 	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

// 	// only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
// 	// and there actually perm issues
// 	if c.operation.BytesToTransfer > 0 && (c.operation.OwnerIssue+c.operation.GroupIssue+c.operation.FolderIssue+c.operation.FileIssue > 0) {
// 		outbound = &dto.Packet{Topic: "calcPermIssue", Payload: fmt.Sprintf("%d|%d|%d|%d", c.operation.OwnerIssue, c.operation.GroupIssue, c.operation.FolderIssue, c.operation.FileIssue)}
// 		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
// 	}
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

// func (c *Core) sendmail(notify int, subject, message string, dryRun bool) (err error) {
// 	if notify == 0 {
// 		return nil
// 	}

// 	dry := ""
// 	if dryRun {
// 		dry = "-------\nDRY RUN\n-------\n"
// 	}

// 	msg := dry + message

// 	// strCmd := fmt.Sprintf("-s \"%s\" -m \"%s\"", mailCmd, subject, msg)
// 	cmd := exec.Command(mailCmd, "-e", "unBALANCE operation update", "-s", subject, "-m", msg)
// 	err = cmd.Run()

// 	return
// }

// func progress(bytesToTransfer, bytesTransferred int64, elapsed time.Duration) (percent float64, left time.Duration, speed float64) {
// 	bytesPerSec := float64(bytesTransferred) / elapsed.Seconds()
// 	speed = bytesPerSec / 1024 / 1024 // MB/s

// 	percent = (float64(bytesTransferred) / float64(bytesToTransfer)) * 100 // %

// 	left = time.Duration(float64(bytesToTransfer-bytesTransferred)/bytesPerSec) * time.Second

// 	return
// }

// func getError(line string, re *regexp.Regexp, errors map[int]string) string {
// 	result := re.FindStringSubmatch(line)
// 	status, _ := strconv.Atoi(result[1])
// 	msg, ok := errors[status]
// 	if !ok {
// 		msg = "unknown error"
// 	}

// 	return msg
// }

// func (c *Core) notifyCommandsToRun(opName string) {
// 	message := "\n\nThe following commands will be executed:\n\n"

// 	for _, command := range c.operation.Commands {
// 		cmd := fmt.Sprintf(`(src: %s) rsync %s %s %s`, command.WorkDir, c.operation.RsyncStrFlags, strconv.Quote(command.Src), strconv.Quote(command.Dst))
// 		message += cmd + "\n"
// 	}

// 	subject := fmt.Sprintf("unBALANCE - %s operation STARTED", strings.ToUpper(opName))

// 	go func() {
// 		if sendErr := c.sendmail(c.settings.NotifyMove, subject, message, c.settings.DryRun); sendErr != nil {
// 			mlog.Error(sendErr)
// 		}
// 	}()
// }
