package common

// MailCmd - location of notify command
const (
	MailCmd         = "/usr/local/emhttp/webGui/scripts/notify" // MailCmd - location of notify command
	PluginLocation  = "/boot/config/plugins/unbalance"          // PluginLocation - plugin's base config path
	ChanCapacity    = 3
	HistoryCapacity = 25
	HistoryFilename = "unbalance.hist"
	RsyncArgs       = "-avPRX"
)

// OpNeutral -
const (
	OpNeutral         = 0
	OpScatterPlan     = 1
	OpScatterMove     = 2
	OpScatterCopy     = 3
	OpScatterValidate = 4
	OpGatherPlan      = 5
	OpGatherMove      = 6
)

// APIGetConfig -
const (
	APIGetConfig      = "core/get/config"
	APIGetState       = "core/get/state"
	APIGetStorage     = "array/get/storage"
	APIGetOperation   = "core/get/operation"
	APIGetHistory     = "core/get/history"
	IntGetArrayStatus = "int/array/get/status"
	APIGetTree        = "array/get/tree"
	APILocateFolder   = "core/locate/folder"
	APIGetLog         = "api/get/log"

	APIScatterPlan         = "api/scatter/plan"
	IntScatterPlan         = "int/scatter/plan"
	IntScatterPlanFinished = "int/scatter/plan/finished"
	IntScatterPlanError    = "int/scatter/plan/error"

	WsScatterPlanStarted  = "scatterPlanStarted"
	WsScatterPlanProgress = "scatterPlanProgress"
	WsScatterPlanFinished = "scatterPlanFinished"
	WsScatterPlanIssues   = "scatterPlanIssue"

	WsGatherPlanStarted  = "gatherPlanStarted"
	WsGatherPlanProgress = "gatherPlanProgress"
	WsGatherPlanFinished = "gatherPlanFinished"
	WsGatherPlanIssues   = "gatherPlanIssue"

	APIScatterMove     = "api/scatter/move"
	APIScatterCopy     = "api/scatter/copy"
	APIScatterValidate = "api/scatter/validate"

	APIGatherPlan         = "api/gather/plan"
	IntGatherPlan         = "int/gather/plan"
	IntGatherPlanFinished = "int/gather/plan/finished"

	APIGatherMove = "api/gather/move"

	APIToggleDryRun   = "config/toggle/dryrun"
	APINotifyCalc     = "config/notify/calc"
	APINotifyMove     = "config/notify/move"
	APISetReserved    = "config/set/reserved"
	APISetVerbosity   = "config/set/verbosity"
	APISetCheckUpdate = "config/set/checkupdate"
	APIGetUpdate      = "config/get/update"
	APISetRsyncArgs   = "config/set/rsyncArgs"

	APIValidate = "api/validate"
	APIReplay   = "api/replay"
)
