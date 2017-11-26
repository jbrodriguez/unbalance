package common

const (
	MailCmd         = "/usr/local/emhttp/webGui/scripts/notify"
	PLUGIN_LOCATION = "/boot/config/plugins/unbalance"

	CHAN_CAPACITY    = 3
	HISTORY_CAPACITY = 2
	HISTORY_FILENAME = "unbalance.hist"

	OP_NEUTRAL          = 0
	OP_SCATTER_PLAN     = 1
	OP_SCATTER_MOVE     = 2
	OP_SCATTER_COPY     = 3
	OP_SCATTER_VALIDATE = 4
	OP_GATHER_PLAN      = 5
	OP_GATHER_MOVE      = 6

	API_GET_CONFIG       = "core/get/config"
	API_GET_STATE        = "core/get/state"
	API_GET_STORAGE      = "array/get/storage"
	API_GET_OPERATION    = "core/get/operation"
	API_GET_HISTORY      = "core/get/history"
	INT_GET_ARRAY_STATUS = "int/array/get/status"
	API_GET_TREE         = "array/get/tree"
	API_LOCATE_FOLDER    = "core/locate/folder"
	API_GET_LOG          = "api/get/log"

	API_SCATTER_PLAN          = "api/scatter/plan"
	INT_SCATTER_PLAN          = "int/scatter/plan"
	INT_SCATTER_PLAN_FINISHED = "int/scatter/plan/finished"
	INT_SCATTER_PLAN_ERROR    = "int/scatter/plan/error"

	WS_SCATTERPLAN_STARTED  = "scatterPlanStarted"
	WS_SCATTERPLAN_PROGRESS = "scatterPlanProgress"
	WS_SCATTERPLAN_FINISHED = "scatterPlanFinished"
	WS_SCATTERPLAN_ISSUES   = "scatterPlanIssue"

	WS_GATHERPLAN_STARTED  = "gatherPlanStarted"
	WS_GATHERPLAN_PROGRESS = "gatherPlanProgress"
	WS_GATHERPLAN_FINISHED = "gatherPlanFinished"
	WS_GATHERPLAN_ISSUES   = "gatherPlanIssue"

	API_SCATTER_MOVE     = "api/scatter/move"
	API_SCATTER_COPY     = "api/scatter/copy"
	API_SCATTER_VALIDATE = "api/scatter/validate"

	INT_OPERATION_FINISHED = "core/operation/finished"

	API_GATHER_PLAN          = "api/gather/plan"
	INT_GATHER_PLAN          = "int/gather/plan"
	INT_GATHER_PLAN_FINISHED = "int/gather/plan/finished"

	API_GATHER_MOVE = "api/gather/move"

	API_TOGGLE_DRYRUN   = "config/toggle/dryrun"
	API_NOTIFY_CALC     = "config/notify/calc"
	API_NOTIFY_MOVE     = "config/notify/move"
	API_SET_RESERVED    = "config/set/reserved"
	API_SET_VERBOSITY   = "config/set/verbosity"
	API_SET_CHECKUPDATE = "config/set/checkupdate"
	API_GET_UPDATE      = "config/get/update"
)
