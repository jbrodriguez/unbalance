package common

const (
	APIEndpoint     = "/api"
	MailCmd         = "/usr/local/emhttp/webGui/scripts/notify" // MailCmd - location of notify command
	PluginLocation  = "/boot/config/plugins/unbalance"          // PluginLocation - plugin's base config path
	ChanCapacity    = 3
	HistoryCapacity = 25
	HistoryFilename = "unbalance.hist"
	HistoryVersion  = 2
	RsyncArgs       = "-avPR"
)

const ReservedSpace uint64 = 512 * 1024 * 1024 // 512Mb

const (
	OpNeutral         = 0
	OpScatterPlanning = 1
	OpScatterPlan     = 2
	OpScatterMove     = 3
	OpScatterCopy     = 4
	OpScatterValidate = 5
	OpGatherPlanning  = 6
	OpGatherPlan      = 7
	OpGatherMove      = 8
)

const (
	CommandScatterPlanStart  = "scatter:plan:start"
	EventScatterPlanStarted  = "scatter:plan:started"
	EventScatterPlanProgress = "scatter:plan:progress"
	EventScatterPlanEnded    = "scatter:plan:ended"
	CommandScatterMove       = "scatter:move"
	CommandScatterCopy       = "scatter:copy"

	CommandGatherPlanStart  = "gather:plan:start"
	EventGatherPlanStarted  = "gather:plan:started"
	EventGatherPlanProgress = "gather:plan:progress"
	EventGatherPlanEnded    = "gather:plan:ended"
)

const (
	CmdCompleted = iota
	CmdPending
	CmdFlagged
	CmdStopped
	CmdSourceRemoval
	CmdInProgress
)
