package common

const (
	PluginName      = "unbalance"
	APIEndpoint     = "/api"
	MailCmd         = "/usr/local/emhttp/webGui/scripts/notify" // MailCmd - location of notify command
	PluginLocation  = "/boot/config/plugins/unbalance"          // PluginLocation - plugin's base config path
	ChanCapacity    = 3
	HistoryCapacity = 25
	HistoryFilename = "unbalance.hist"
	HistoryVersion  = 2
	RsyncArgs       = "-avPR"
)

const ReservedSpace uint64 = 1024 * 1024 * 1024 // 1Gb

const (
	OpNeutral         = 0
	OpScatterPlan     = 1
	OpScatterMove     = 2
	OpScatterCopy     = 3
	OpScatterValidate = 4
	OpGatherPlan      = 5
	OpGatherMove      = 6
)

const (
	CommandScatterPlanStart  = "scatter:plan:start"
	EventScatterPlanStarted  = "scatter:plan:started"
	EventScatterPlanProgress = "scatter:plan:progress"
	EventScatterPlanEnded    = "scatter:plan:ended"
	CommandScatterMove       = "scatter:move"
	CommandScatterCopy       = "scatter:copy"
	CommandScatterValidate   = "scatter:validate"

	CommandGatherPlanStart  = "gather:plan:start"
	EventGatherPlanStarted  = "gather:plan:started"
	EventGatherPlanProgress = "gather:plan:progress"
	EventGatherPlanEnded    = "gather:plan:ended"
	CommandGatherMove       = "gather:move"

	EventTransferStarted  = "transfer:started"
	EventTransferProgress = "transfer:progress"
	EventTransferEnded    = "transfer:ended"

	EventOperationError = "operation:error"

	CommandRemoveSource = "remove:source"
	CommandReplay       = "replay"
	CommandStop         = "stop"
)

const (
	CmdCompleted = iota
	CmdPending
	CmdFlagged
	CmdStopped
	CmdSourceRemoval
	CmdInProgress
)
