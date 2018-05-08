package services

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
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

type converter func(*domain.History) *domain.History

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

	converters []converter

	stopped bool
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

	core.converters = []converter{
		convertToV2,
	}

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
		99: "Interrupted by the user",
	}

	return core
}

// Start -
func (c *Core) Start() (err error) {
	mlog.Info("starting service Core ...")

	msg := &pubsub.Message{Reply: make(chan interface{}, common.ChanCapacity)}
	c.bus.Pub(msg, common.IntGetArrayStatus)
	reply := <-msg.Reply
	message := reply.(dto.Message)
	if message.Error != nil {
		return message.Error
	}

	c.state.Status = common.OpNeutral
	c.state.Unraid = message.Data.(*domain.Unraid)

	history, err := c.historyRead()
	if err != nil {
		mlog.Warning("Unable to read history: %s", err)
	}

	if history.Version < common.HistoryVersion {
		history = runConverters(history, c.converters, common.HistoryVersion)

		history.Version = common.HistoryVersion

		err := c.historyWrite(history)
		if err != nil {
			mlog.Warning("Unable to write history: %s", err)
		}
	}

	c.state.History = history

	c.actor.Register(common.APIGetConfig, c.getConfig)
	c.actor.Register(common.APIGetState, c.getState)
	c.actor.Register(common.APIGetStorage, c.getStorage)
	c.actor.Register(common.APIGetOperation, c.getOperation)
	c.actor.Register(common.APIGetHistory, c.getHistory)
	c.actor.Register(common.APILocateFolder, c.locate)

	c.actor.Register(common.APIScatterPlan, c.scatterPlan)
	c.actor.Register(common.IntScatterPlanFinished, c.scatterPlanFinished)
	c.actor.Register(common.APIScatterMove, c.scatterMove)
	c.actor.Register(common.APIScatterCopy, c.scatterCopy)

	c.actor.Register(common.APIGatherPlan, c.gatherPlan)
	c.actor.Register(common.IntGatherPlanFinished, c.gatherPlanFinished)
	c.actor.Register(common.APIGatherMove, c.gatherMove)

	c.actor.Register(common.APIToggleDryRun, c.toggleDryRun)
	c.actor.Register(common.APINotifyPlan, c.setNotifyPlan)
	c.actor.Register(common.APINotifyTransfer, c.setNotifyTransfer)
	c.actor.Register(common.APISetReserved, c.setReservedSpace)
	c.actor.Register(common.APISetVerbosity, c.setVerbosity)
	c.actor.Register(common.APISetCheckUpdate, c.setCheckUpdate)
	c.actor.Register(common.APIGetUpdate, c.getUpdate)
	c.actor.Register(common.APISetRsyncArgs, c.setRsyncArgs)
	c.actor.Register(common.APIValidate, c.validate)
	c.actor.Register(common.APIReplay, c.replay)
	c.actor.Register(common.APIRemoveSource, c.removeSource)
	c.actor.Register(common.APIStopCommand, c.stopCommand)
	// c.actor.Register("getLog", c.getLog)

	go c.actor.React()

	return nil
}

// Stop -
func (c *Core) Stop() {
	mlog.Info("stopped service Core ...")
}

// COMMON API ENDPOINTS
func (c *Core) getConfig(msg *pubsub.Message) {
	mlog.Info("Sending config")
	msg.Reply <- &c.settings.Config
}

func (c *Core) getState(msg *pubsub.Message) {
	mlog.Info("Sending state")
	msg.Reply <- c.state
}

func (c *Core) getOperation(msg *pubsub.Message) {
	mlog.Info("Sending operation")
	msg.Reply <- c.state.Operation
}

func (c *Core) getStorage(msg *pubsub.Message) {
	mlog.Info("Sending storage")

	c.state.Unraid = c.refreshUnraid()

	msg.Reply <- c.state.Unraid
}

func (c *Core) getHistory(msg *pubsub.Message) {
	mlog.Info("Sending history")

	c.state.History.LastChecked = time.Now()
	msg.Reply <- c.state.History

	go func() {
		err := c.historyWrite(c.state.History)
		if err != nil {
			mlog.Warning("Unable to write history: %s", err)
		}
	}()
}

func (c *Core) locate(msg *pubsub.Message) {
	chosen := msg.Payload.([]string)

	location := &dto.Location{
		Disks:    make(map[string]*domain.Disk),
		Presence: make(map[string]string),
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

// SCATTER PLAN
func (c *Core) getScatterPlan(msg *pubsub.Message) (*domain.Plan, error) {
	payload, ok := msg.Payload.(string)
	if !ok {
		return nil, errors.New("Unable to convert scatter plan parameters")
	}

	var plan domain.Plan
	err := json.Unmarshal([]byte(payload), &plan)
	if err != nil {
		return nil, err
	}

	return &plan, nil
}

func (c *Core) scatterPlan(msg *pubsub.Message) {
	c.state.Status = common.OpScatterPlan
	c.state.Unraid = c.refreshUnraid()

	basePlan, err := c.getScatterPlan(msg)
	if err != nil {
		mlog.Warning("Unable to get scatter plan: %s", err)

		// send to front end the signal of operation finished
		outbound := &dto.Packet{Topic: common.WsScatterPlanFinished, Payload: basePlan}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		outbound = &dto.Packet{Topic: common.WsScatterPlanIssues, Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	plan := &domain.Plan{
		ChosenFolders: basePlan.ChosenFolders,
		VDisks:        basePlan.VDisks,
	}

	for _, disk := range c.state.Unraid.Disks {
		plan.VDisks[disk.Path].PlannedFree = disk.Free
		plan.VDisks[disk.Path].Bin = nil
	}

	param := &pubsub.Message{Payload: &domain.State{
		Status: c.state.Status,
		Unraid: c.state.Unraid,
		Plan:   plan,
	}}

	c.bus.Pub(param, common.IntScatterPlan)
}

func (c *Core) scatterPlanFinished(msg *pubsub.Message) {
	plan := msg.Payload.(*domain.Plan)

	c.state.Status = common.OpNeutral

	// send to front end the signal of operation finished
	outbound := &dto.Packet{Topic: common.WsScatterPlanFinished, Payload: plan}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
	// and there actually perm issues
	if plan.BytesToTransfer > 0 && (plan.OwnerIssue+plan.GroupIssue+plan.FolderIssue+plan.FileIssue > 0) {
		outbound = &dto.Packet{Topic: common.WsScatterPlanIssues, Payload: fmt.Sprintf("%d|%d|%d|%d", plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	// order := make([]string, 0)
	// for _, disk := range c.state.Unraid.Disks {
	// 	order = append(order, disk.Path)
	// }

	// plan.Print(order)
}

// GATHER PLAN
func (c *Core) getGatherPlan(msg *pubsub.Message) (*domain.Plan, error) {
	data, ok := msg.Payload.(string)
	if !ok {
		return nil, errors.New("Unable to convert gather plan parameters")
	}

	var plan domain.Plan
	err := json.Unmarshal([]byte(data), &plan)
	if err != nil {
		return nil, fmt.Errorf("Unable to bind gather plan parameters: %s", err)
	}

	return &plan, nil
}

func (c *Core) gatherPlan(msg *pubsub.Message) {
	c.state.Status = common.OpGatherPlan
	c.state.Unraid = c.refreshUnraid()

	basePlan, err := c.getGatherPlan(msg)
	if err != nil {
		mlog.Warning("Unable to get gather plan: %s", err)

		// send to front end the signal of operation finished
		outbound := &dto.Packet{Topic: common.WsGatherPlanFinished, Payload: basePlan}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		outbound = &dto.Packet{Topic: "opError", Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	plan := &domain.Plan{
		ChosenFolders: basePlan.ChosenFolders,
		VDisks:        basePlan.VDisks,
	}

	for _, disk := range c.state.Unraid.Disks {
		plan.VDisks[disk.Path].PlannedFree = disk.Free
		plan.VDisks[disk.Path].Bin = nil
	}

	param := &pubsub.Message{Payload: &domain.State{
		Status: c.state.Status,
		Unraid: c.state.Unraid,
		Plan:   plan,
	}}

	c.bus.Pub(param, common.IntGatherPlan)
}

func (c *Core) gatherPlanFinished(msg *pubsub.Message) {
	plan := msg.Payload.(*domain.Plan)

	c.state.Status = common.OpNeutral

	// send to front end the signal of operation finished
	outbound := &dto.Packet{Topic: common.WsGatherPlanFinished, Payload: plan}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// only send the perm issue msg if there's actually some work to do (BytesToTransfer > 0)
	// and there actually perm issues
	if plan.BytesToTransfer > 0 && (plan.OwnerIssue+plan.GroupIssue+plan.FolderIssue+plan.FileIssue > 0) {
		outbound = &dto.Packet{Topic: common.WsGatherPlanIssues, Payload: fmt.Sprintf("%d|%d|%d|%d", plan.OwnerIssue, plan.GroupIssue, plan.FolderIssue, plan.FileIssue)}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	// order := make([]string, 0)
	// for _, disk := range c.state.Unraid.Disks {
	// 	order = append(order, disk.Path)
	// }

	// plan.Print(order)
}

// SCATTER TRANSFER
func (c *Core) setupScatterTransferOperation(status int64, disks []*domain.Disk, plan *domain.Plan) *domain.Operation {
	operation := &domain.Operation{
		ID:              shortid.MustGenerate(),
		OpKind:          status,
		BytesToTransfer: plan.BytesToTransfer,
		DryRun:          c.settings.DryRun,
	}

	operation.RsyncArgs = append([]string{common.RsyncArgs}, c.settings.RsyncArgs...)

	// user may have changed dry-run setting, adjust for it
	if operation.DryRun {
		operation.RsyncArgs = append(operation.RsyncArgs, "--dry-run")
	}
	operation.RsyncStrArgs = strings.Join(operation.RsyncArgs, " ")

	operation.Commands = make([]*domain.Command, 0)

	for _, disk := range disks {
		vdisk := plan.VDisks[disk.Path]

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
				ID:     shortid.MustGenerate(),
				Src:    item.Location,
				Dst:    vdisk.Path + string(filepath.Separator),
				Entry:  entry,
				Size:   item.Size,
				Status: common.CmdPending,
			})
		}
	}

	return operation
}

func (c *Core) scatterMove(msg *pubsub.Message) {
	c.state.Status = common.OpScatterMove
	c.state.Unraid = c.refreshUnraid()

	plan, err := c.getScatterPlan(msg)
	if err != nil {
		mlog.Warning("Unable to get scatter plan: %s", err.Error())

		outbound := &dto.Packet{Topic: "opError", Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	c.state.Operation = c.setupScatterTransferOperation(c.state.Status, c.state.Unraid.Disks, plan)

	go c.runOperation("Move")
}

func (c *Core) scatterCopy(msg *pubsub.Message) {
	c.state.Status = common.OpScatterCopy
	c.state.Unraid = c.refreshUnraid()

	plan, err := c.getScatterPlan(msg)
	if err != nil {
		mlog.Warning("Unable to get scatter plan: %s", err)

		outbound := &dto.Packet{Topic: "opError", Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	c.state.Operation = c.setupScatterTransferOperation(c.state.Status, c.state.Unraid.Disks, plan)

	go c.runOperation("Copy")
}

// GATHER TRANSFER
func (c *Core) setupGatherTransferOperation(status int64, disks []*domain.Disk, plan *domain.Plan) *domain.Operation {
	operation := &domain.Operation{
		ID:     shortid.MustGenerate(),
		OpKind: status,
		DryRun: c.settings.DryRun,
	}

	operation.RsyncArgs = append([]string{common.RsyncArgs}, c.settings.RsyncArgs...)

	// user may have changed dry-run setting, adjust for it
	if operation.DryRun {
		operation.RsyncArgs = append(operation.RsyncArgs, "--dry-run")
	}
	operation.RsyncStrArgs = strings.Join(operation.RsyncArgs, " ")

	operation.Commands = make([]*domain.Command, 0)

	for _, disk := range disks {
		vdisk := plan.VDisks[disk.Path]

		// only one disk will be destination (target)
		if !vdisk.Dst {
			continue
		}

		// user chose a target disk, adjust bytestotransfer to the size of its bin, since
		// that's the amount of data we need to transfer. Also remove bin from all other disks,
		operation.BytesToTransfer = vdisk.Bin.Size

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
				ID:     shortid.MustGenerate(),
				Src:    item.Location,
				Dst:    vdisk.Path + string(filepath.Separator),
				Entry:  entry,
				Size:   item.Size,
				Status: common.CmdPending,
			})
		}
	}

	return operation
}

func (c *Core) gatherMove(msg *pubsub.Message) {
	c.state.Status = common.OpGatherMove
	c.state.Unraid = c.refreshUnraid()

	plan, err := c.getGatherPlan(msg)
	if err != nil {
		mlog.Warning("Unable to get gather plan: %s", err.Error())

		outbound := &dto.Packet{Topic: "opError", Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	c.state.Operation = c.setupGatherTransferOperation(c.state.Status, c.state.Unraid.Disks, plan)

	go c.runOperation("Move")
}

// VALIDATE TRANSFER
func (c *Core) getValidateOperation(msg *pubsub.Message) (*domain.Operation, error) {
	data, ok := msg.Payload.(string)
	if !ok {
		return nil, errors.New("Unable to convert validate parameters")
	}

	var operation domain.Operation
	err := json.Unmarshal([]byte(data), &operation)
	if err != nil {
		return nil, fmt.Errorf("Unable to bind validate parameters: %s", err)
	}

	return &operation, nil
}

func (c *Core) setupValidateOperation(originalOp *domain.Operation) *domain.Operation {
	operation := &domain.Operation{
		ID:              shortid.MustGenerate(),
		OpKind:          common.OpScatterValidate,
		BytesToTransfer: originalOp.BytesToTransfer,
		DryRun:          false,
	}

	operation.RsyncArgs = append([]string{strings.Replace(common.RsyncArgs, "-a", "-rc", -1)}, originalOp.RsyncArgs[1:]...)
	operation.RsyncStrArgs = strings.Join(operation.RsyncArgs, " ")

	operation.Commands = originalOp.Commands

	for _, command := range operation.Commands {
		command.ID = shortid.MustGenerate()
		command.Transferred = 0
		command.Status = common.CmdPending
	}

	return operation
}

func (c *Core) validate(msg *pubsub.Message) {
	c.state.Status = common.OpScatterValidate
	c.state.Unraid = c.refreshUnraid()

	originalOp, err := c.getValidateOperation(msg)
	if err != nil {
		mlog.Warning("Unable to get validate operation: %s", err)

		outbound := &dto.Packet{Topic: "opError", Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	c.state.Operation = c.setupValidateOperation(originalOp)

	go c.runOperation("Validate")
}

// REPLAY TRANSFER
func (c *Core) getReplayOperation(msg *pubsub.Message) (*domain.Operation, error) {
	data, ok := msg.Payload.(string)
	if !ok {
		return nil, errors.New("Unable to convert replay parameters")
	}

	var operation domain.Operation
	err := json.Unmarshal([]byte(data), &operation)
	if err != nil {
		return nil, fmt.Errorf("Unable to bind replay parameters: %s", err)
	}

	return &operation, nil
}

func (c *Core) setupReplayOperation(originalOp *domain.Operation) *domain.Operation {
	operation := &domain.Operation{
		ID:              shortid.MustGenerate(),
		OpKind:          originalOp.OpKind,
		BytesToTransfer: originalOp.BytesToTransfer,
		DryRun:          false,
	}

	operation.RsyncArgs = originalOp.RsyncArgs
	operation.RsyncStrArgs = strings.Join(operation.RsyncArgs, " ")

	operation.Commands = originalOp.Commands

	for _, command := range operation.Commands {
		command.ID = shortid.MustGenerate()
		command.Transferred = 0
		command.Status = common.CmdPending
	}

	return operation
}

func getOpName(status int64) string {
	if status == common.OpScatterMove || status == common.OpGatherMove {
		return "Move"
	}

	return "Copy"
}

func (c *Core) replay(msg *pubsub.Message) {
	originalOp, err := c.getReplayOperation(msg)
	if err != nil {
		mlog.Warning("Unable to get replay operation: %s", err)

		outbound := &dto.Packet{Topic: "opError", Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	c.state.Operation = c.setupReplayOperation(originalOp)
	c.state.Status = c.state.Operation.OpKind
	c.state.Unraid = c.refreshUnraid()

	go c.runOperation(getOpName(c.state.Status))
}

// COMMON TRANSFER
func (c *Core) runOperation(opName string) {
	mlog.Info("Running %s operation ...", opName)

	operation := c.state.Operation
	operation.Started = time.Now()

	if c.settings.NotifyTransfer == 2 {
		c.notifyCommandsToRun(opName, operation)
	}

	operation.Line = "Waiting to collect stats ..."
	outbound := &dto.Packet{Topic: "transferStarted", Payload: operation}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	commandsExecuted := make([]string, 0)

	for _, command := range operation.Commands {
		args := append(
			operation.RsyncArgs,
			command.Entry,
			command.Dst,
		)

		cmd := fmt.Sprintf(`rsync %s %s %s`, operation.RsyncStrArgs, strconv.Quote(command.Entry), strconv.Quote(command.Dst))
		mlog.Info("Command Started: (src: %s) %s ", command.Src, cmd)

		operation.Line = cmd
		outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		cmdTransferred, err := c.runCommand(operation, command, c.bus, c.reProgress, args, c.settings.Verbosity)
		if err != nil {
			c.commandInterrupted(opName, operation, command, cmd, err, cmdTransferred, commandsExecuted)
			return
		}

		c.commandCompleted(operation, command)

		commandsExecuted = append(commandsExecuted, cmd)
	}

	c.operationCompleted(opName, operation, commandsExecuted)
}

func (c *Core) runCommand(operation *domain.Operation, command *domain.Command, bus *pubsub.PubSub, re *regexp.Regexp, args []string, verbosity int) (uint64, error) {
	// make sure the command will run
	c.stopped = false

	// start rsync command
	cmd, err := lib.StartRsync(command.Src, mlog.Warning, args...)
	if err != nil {
		return 0, err
	}

	// give some time for /proc/pid to come alive
	time.Sleep(500 * time.Millisecond)

	// monitor rsync progress
	retcode, transferred, err := c.monitorRsync(operation, command, bus, cmd.Process.Pid)
	if err != nil {
		mlog.Warning("command:monitor(%s)", err)
	}

	// 99 is a magic number meaning the user clicked on the stop operation button
	if retcode == 99 {
		mlog.Info("command:received:StopCommand")
		err = lib.KillRsync(cmd)
		if err != nil {
			mlog.Warning("command:kill:error(%s)", err)
		}

		command.Status = common.CmdStopped
		return transferred, fmt.Errorf("exit status %d", retcode)
	}

	command.Status = common.CmdCompleted

	// end rsync process
	exitCode, err := lib.EndRsync(cmd)
	if err != nil {
		mlog.Warning("command:end:error(%s)", err)
		command.Status = common.CmdStopped
	}

	mlog.Info("command:retcode(%d):exitcode(%d)", retcode, exitCode)

	if exitCode == 23 {
		err = nil
		command.Status = common.CmdFlagged
	}

	return transferred, err
}

func (c *Core) monitorRsync(operation *domain.Operation, command *domain.Command, bus *pubsub.PubSub, procPid int) (int, uint64, error) {
	var transferred uint64
	var current string
	var retcode int
	var zombie bool
	var err error

	// started := time.Now()
	// throttled := false

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
			mlog.Warning("isZombie:err(%s)", err)
			break
		}

		// pid finished processing, exit the loop and Wait on the pid, to finish it properly
		if zombie {
			break
		}

		// throttle both /proc consumption and messaging to the front end
		time.Sleep(250 * time.Millisecond)

		// getReadBytes
		transferred, err = getReadBytes(procIo)
		if err != nil {
			mlog.Warning("getReadBytes:err(%s):xfer(%d)", err, transferred)
			break
		}

		// getCurrentTransfer
		current, err = getCurrentTransfer(procFd, filepath.Join(command.Src, command.Entry))
		if err != nil {
			// mlog.Warning("getCurrentTransfer:err(%s)", err)
			continue
		}

		// mlog.Info("read(%d)-current(%s)-size(%d)", transferred, current, command.Size)

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

		outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
		bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
	}

	return retcode, transferred, nil
}

func isZombie(proc string) (bool, int, error) {
	var zombie bool
	var retcode int

	b, e := ioutil.ReadFile(proc)
	if e != nil {
		return false, 0, e
	}

	fields := strings.Split(string(b), " ")
	state := fields[2]
	zombie = state == "Z"
	if zombie {
		retcode, _ = strconv.Atoi(fields[51])
	}

	return zombie, retcode, nil
}

func getReadBytes(proc string) (uint64, error) {
	var sRead string

	b, e := ioutil.ReadFile(proc)
	if e != nil {
		return 0, e
	}

	lines := strings.Split(string(b), "\n")

	for _, line := range lines {
		if strings.HasPrefix(line, "rchar:") {
			sRead = line[7:]
			break
		}
	}

	read, _ := strconv.ParseUint(sRead, 10, 64)

	return read, nil
}

func getCurrentTransfer(proc, prefix string) (string, error) {
	var current string

	name, e := os.Readlink(proc)
	if e != nil {
		return "", e
	}

	if strings.HasPrefix(name, prefix) {
		current = name
	}

	return current, nil
}

func (c *Core) commandInterrupted(opName string, operation *domain.Operation, command *domain.Command, cmd string, err error, cmdTransferred uint64, commandsExecuted []string) {
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
	operation.Remaining = left.String()
	operation.DeltaTransfer = cmdTransferred
	command.Transferred = cmdTransferred

	c.endOperation(subject, headline, commandsExecuted, operation)
}

func showPotentiallyPrunedItems(operation *domain.Operation, command *domain.Command) {
	if operation.DryRun && operation.OpKind == common.OpGatherMove {
		parent := filepath.Dir(command.Entry)
		// mlog.Info("parent(%s)-src(%s)-dst(%s)-entry(%s)", parent, command.Src, command.Dst, command.Entry)
		if parent != "." {
			mlog.Info(`Would delete empty folders starting from (%s) - (find "%s" -type d -empty -prune -exec rm -rf {} \;) `, filepath.Join(command.Src, parent), filepath.Join(command.Src, parent))
		} else {
			mlog.Info(`WONT DELETE: find "%s" -type d -empty -prune -exec rm -rf {} \;`, filepath.Join(command.Src, parent))
		}
	}
}

func handleItemDeletion(operation *domain.Operation, command *domain.Command, bus *pubsub.PubSub) {
	if !operation.DryRun && (operation.OpKind == common.OpScatterMove || operation.OpKind == common.OpGatherMove) {
		// the command was flagged due to an error 23, don't delete the source file/folder in these cases
		if command.Status == common.CmdFlagged {
			msg := fmt.Sprintf("skipping:deletion:(rsync command was flagged):(%s)", filepath.Join(command.Dst, command.Entry))
			operation.Line = msg

			outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
			bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			mlog.Warning(msg)

			return
		}

		exists, _ := lib.Exists(filepath.Join(command.Dst, command.Entry))
		if exists {
			rmrf := fmt.Sprintf("rm -rf \"%s\"", filepath.Join(command.Src, command.Entry))
			operation.Line = fmt.Sprintf("Removing source %s", filepath.Join(command.Src, command.Entry))

			outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
			bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
			mlog.Info("removing:(%s)", rmrf)

			err := lib.Shell(rmrf, mlog.Warning, "transferProgress:", "", func(line string) {
				mlog.Info(line)
			})

			if err != nil {
				msg := fmt.Sprintf("Unable to remove source folder (%s): %s", filepath.Join(command.Src, command.Entry), err)
				operation.Line = msg

				outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
				bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

				mlog.Warning(msg)
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

					outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
					bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

					mlog.Info("pruning:(%s)", rmdir)

					err = lib.Shell(rmdir, mlog.Warning, "transferProgress:", "", func(line string) {
						mlog.Info(line)
					})

					if err != nil {
						msg := fmt.Sprintf("Unable to remove parent folder (%s): %s", filepath.Join(command.Src, parent), err)
						operation.Line = msg

						outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
						bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

						mlog.Warning(msg)
					}
				} else {
					mlog.Warning("skipping:prune:(%s)", filepath.Join(command.Src, parent))
				}
			}
		} else {
			mlog.Warning("skipping:deletion:(file/folder not present in destination):(%s)", filepath.Join(command.Dst, command.Entry))
		}
	}
}

func (c *Core) commandCompleted(operation *domain.Operation, command *domain.Command) {
	text := "Command Finished"
	mlog.Info(text)

	operation.BytesTransferred += command.Size
	percent, left, speed := progress(operation.BytesToTransfer, operation.BytesTransferred, time.Since(operation.Started))

	operation.Completed = percent
	operation.Speed = speed
	operation.Remaining = left.String()
	operation.DeltaTransfer = 0
	operation.Line = text
	command.Transferred = command.Size

	msg := fmt.Sprintf("%.2f%% done ~ %s left (%.2f MB/s)", percent, operation.Remaining, speed)
	mlog.Info("Current progress: %s", msg)

	outbound := &dto.Packet{Topic: "transferProgress", Payload: operation}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// this is just a heads up for the user, shows which folders would/wouldn't be pruned if run without dry-run
	showPotentiallyPrunedItems(operation, command)

	// if it isn't a dry-run and the operation is Move or Gather, delete the source folder
	handleItemDeletion(operation, command, c.bus)
}

func (c *Core) operationCompleted(opName string, operation *domain.Operation, commandsExecuted []string) {
	operation.Finished = time.Now()
	elapsed := operation.Finished.Sub(operation.Started)

	subject := fmt.Sprintf("unBALANCE - %s operation completed", strings.ToUpper(opName))
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

	c.updateHistory(c.state.History, operation)

	c.state.Unraid = c.refreshUnraid()

	state := &domain.State{
		Operation: operation,
		History:   c.state.History,
		Unraid:    c.state.Unraid,
	}

	outbound := &dto.Packet{Topic: "transferFinished", Payload: state}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	message := fmt.Sprintf("\n\nStarted: %s\nEnded: %s\n\nElapsed: %s\n\n%s\n\nTransferred %s at ~ %.2f MB/s",
		fstarted, ffinished, elapsed, headline, lib.ByteSize(operation.BytesTransferred), operation.Speed,
	)

	switch c.settings.NotifyTransfer {
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
		if sendErr := c.sendmail(c.settings.NotifyTransfer, subject, message, c.settings.DryRun); sendErr != nil {
			mlog.Error(sendErr)
		}
	}()

	mlog.Info("\n%s\n%s", subject, message)

	// c.bus.Pub(&pubsub.Message{}, common.INT_OPERATION_FINISHED)
	c.state.Status = common.OpNeutral
	c.state.Operation = nil
}

// REMOVE SOURCE SCENARIO
func getRemoveSourceParams(msg *pubsub.Message) (*domain.Operation, string, error) {
	data, ok := msg.Payload.(string)
	if !ok {
		return nil, "", errors.New("Unable to convert removeSource parameters")
	}

	var rmsrc dto.RmSrc
	err := json.Unmarshal([]byte(data), &rmsrc)
	if err != nil {
		return nil, "", fmt.Errorf("Unable to bind removeSource parameters: %s", err)
	}

	return rmsrc.Operation, rmsrc.ID, nil
}

func (c *Core) removeSource(msg *pubsub.Message) {
	operation, cmdID, err := getRemoveSourceParams(msg)
	if err != nil {
		mlog.Warning("Unable to get rmsrc params: %s", err)

		outbound := &dto.Packet{Topic: "opError", Payload: err.Error()}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		return
	}

	c.state.Status = operation.OpKind
	c.state.Operation = operation
	c.state.Unraid = c.refreshUnraid()

	go c.performRemoveSource(cmdID)
}

func (c *Core) performRemoveSource(id string) {
	operation := c.state.Operation

	operation.Line = "Removing source ..."
	outbound := &dto.Packet{Topic: "transferStarted", Payload: operation}
	c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

	// do work
	// find command
	// var command *domain.Command
	for _, command := range operation.Commands {
		if command.ID != id {
			continue
		}

		mlog.Info("Removing source:(%s)", filepath.Join(command.Src, command.Entry))

		// status is cmdFlagged currently, let's change this to cmdSourceRemoval, so that handleItemDeletion works
		// correctly and the UI can display a proper feedback for this command
		command.Status = common.CmdSourceRemoval

		handleItemDeletion(operation, command, c.bus)

		command.Status = common.CmdCompleted

		text := "Removal completed"
		operation.Line = text
		mlog.Info(text)

		c.state.Unraid = c.refreshUnraid()
		c.state.History.Items[operation.ID] = operation

		state := &domain.State{
			Operation: operation,
			History:   c.state.History,
			Unraid:    c.state.Unraid,
		}

		outbound = &dto.Packet{Topic: "transferFinished", Payload: state}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		// TODO: how to handle this
		// c.updateHistory(c.state.History, operation)
		err := c.historyWrite(c.state.History)
		if err != nil {
			mlog.Warning("Unable to write history: %s", err)
		}

		c.state.Status = common.OpNeutral
		c.state.Operation = nil

		break
	}
}

// STOP COMMAND
func (c *Core) stopCommand(msg *pubsub.Message) {
	c.stopped = true
}

// SETTINGS RELATED
func (c *Core) setNotifyPlan(msg *pubsub.Message) {
	fnotify := msg.Payload.(float64)
	notify := int(fnotify)

	mlog.Info("Setting notifyPlan to (%d)", notify)

	c.settings.NotifyPlan = notify
	err := c.settings.Save()
	if err != nil {
		mlog.Warning("Unable to save settings: %s", err)
	}

	msg.Reply <- &c.settings.Config
}

func (c *Core) setNotifyTransfer(msg *pubsub.Message) {
	fnotify := msg.Payload.(float64)
	notify := int(fnotify)

	mlog.Info("Setting notifyTransfer to (%d)", notify)

	c.settings.NotifyTransfer = notify
	err := c.settings.Save()
	if err != nil {
		mlog.Warning("Unable to save settings: %s", err)
	}

	msg.Reply <- &c.settings.Config
}

func (c *Core) setVerbosity(msg *pubsub.Message) {
	fverbosity := msg.Payload.(float64)
	verbosity := int(fverbosity)

	mlog.Info("Setting verbosity to (%d)", verbosity)

	c.settings.Verbosity = verbosity
	err := c.settings.Save()
	if err != nil {
		mlog.Warning("Unable to save settings: %s", err)
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

	amount := uint64(reserved.Amount)
	unit := reserved.Unit

	mlog.Info("Setting reservedAmount to (%d)", amount)
	mlog.Info("Setting reservedUnit to (%s)", unit)

	c.settings.ReservedAmount = amount
	c.settings.ReservedUnit = unit
	err = c.settings.Save()
	if err != nil {
		mlog.Warning("Unable to save settings: %s", err)
	}

	msg.Reply <- &c.settings.Config
}

func (c *Core) toggleDryRun(msg *pubsub.Message) {
	mlog.Info("Toggling dryRun from (%t)", c.settings.DryRun)

	c.settings.ToggleDryRun()
	err := c.settings.Save()
	if err != nil {
		mlog.Warning("Unable to save settings: %s", err)
	}

	msg.Reply <- &c.settings.Config
}

func (c *Core) setRsyncArgs(msg *pubsub.Message) {
	// mlog.Warning("payload: %+v", msg.Payload)
	payload, ok := msg.Payload.(string)
	if !ok {
		mlog.Warning("Unable to convert Rsync arguments")
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to convert Rsync arguments"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")

		msg.Reply <- &c.settings.Config

		return
	}

	var rsync dto.Rsync
	err := json.Unmarshal([]byte(payload), &rsync)
	if err != nil {
		mlog.Warning("Unable to bind rsyncArgs parameters: %s", err)
		outbound := &dto.Packet{Topic: "opError", Payload: "Unable to bind rsyncArgs parameters"}
		c.bus.Pub(&pubsub.Message{Payload: outbound}, "socket:broadcast")
		return
		// mlog.Fatalf(err.Error())
	}

	mlog.Info("Setting rsyncArgs to (%s)", strings.Join(rsync.Args, " "))

	c.settings.RsyncArgs = rsync.Args
	err = c.settings.Save()
	if err != nil {
		mlog.Warning("Unable to save settings: %s", err)
	}

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
		cmd := fmt.Sprintf(`(src: %s) rsync %s %s %s`, command.Src, operation.RsyncStrArgs, strconv.Quote(command.Entry), strconv.Quote(command.Dst))
		message += cmd + "\n"
	}

	subject := fmt.Sprintf("unBALANCE - %s operation STARTED", strings.ToUpper(opName))

	go func() {
		if sendErr := c.sendmail(c.settings.NotifyTransfer, subject, message, c.settings.DryRun); sendErr != nil {
			mlog.Error(sendErr)
		}
	}()
}

func progress(bytesToTransfer, bytesTransferred uint64, elapsed time.Duration) (percent float64, left time.Duration, speed float64) {
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

func (c *Core) refreshUnraid() *domain.Unraid {
	unraid := c.state.Unraid

	param := &pubsub.Message{Reply: make(chan interface{}, common.ChanCapacity)}
	c.bus.Pub(param, common.IntGetArrayStatus)

	reply := <-param.Reply
	message := reply.(dto.Message)
	if message.Error != nil {
		mlog.Warning("Unable to get storage: %s", message.Error)
	} else {
		unraid = message.Data.(*domain.Unraid)
	}

	return unraid
}

// HISTORY HANDLERS/CONVERTERS
func (c *Core) historyRead() (*domain.History, error) {
	var history domain.History

	fileName := filepath.Join(common.PluginLocation, common.HistoryFilename)

	file, err := os.Open(fileName)
	if err != nil {
		empty := &domain.History{
			Version: common.HistoryVersion,
			Items:   make(map[string]*domain.Operation),
			Order:   make([]string, 0),
		}

		return empty, err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&history)
	if err != nil {
		empty := &domain.History{
			Version: common.HistoryVersion,
			Items:   make(map[string]*domain.Operation),
			Order:   make([]string, 0),
		}

		return empty, err
	}

	return &history, nil
}

func (c *Core) historyWrite(history *domain.History) error {
	tmpName := filepath.Join(common.PluginLocation, common.HistoryFilename+"."+shortid.MustGenerate())

	file, err := os.Create(tmpName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(history)
	if err != nil {
		return err
	}

	err = os.Rename(tmpName, filepath.Join(common.PluginLocation, common.HistoryFilename))

	return err
}

func (c *Core) updateHistory(history *domain.History, operation *domain.Operation) {
	count := len(history.Order)
	if count == common.HistoryCapacity {
		delete(history.Items, history.Order[count-1])
		// prepend item, remove oldest item
		history.Order = append([]string{c.state.Operation.ID}, history.Order[:count-1]...)
	} else {
		// prepend item
		history.Order = append([]string{c.state.Operation.ID}, history.Order...)
	}

	history.Items[operation.ID] = operation

	go func() {
		err := c.historyWrite(history)
		if err != nil {
			mlog.Warning("Unable to write history: %s", err)
		}
	}()
}

func runConverters(history *domain.History, converters []converter, version int) *domain.History {
	// converters is a zero-based array, we're currently at historyversion = 2, so we can do this math to get to
	// the first converter we need to run
	base := version - 2
	toRun := converters[base:]
	for _, converter := range toRun {
		history = converter(history)
	}

	return history
}

func convertToV2(history *domain.History) *domain.History {
	for _, item := range history.Items {
		for _, command := range item.Commands {
			if command.Transferred == 0 {
				command.Status = common.CmdPending
			} else if command.Transferred != command.Size {
				command.Status = common.CmdStopped
			} else {
				command.Status = common.CmdCompleted
			}
		}
	}

	return history
}
