package common

const (
	MailCmd         = "/usr/local/emhttp/webGui/scripts/notify"
	PLUGIN_LOCATION = "/boot/config/plugins/unbalance"

	CHAN_CAPACITY    = 3
	HISTORY_CAPACITY = 2
	HISTORY_FILENAME = "unbalance.hist"

	OP_NEUTRAL          = 0
	OP_SCATTER_CALC     = 1
	OP_SCATTER_MOVE     = 2
	OP_SCATTER_COPY     = 3
	OP_SCATTER_VALIDATE = 4
	OP_GATHER_CALC      = 5
	OP_GATHER_MOVE      = 6

	API_GET_CONFIG       = "core/get/config"
	API_GET_STATUS       = "core/get/status"
	API_GET_STATE        = "core/get/state"
	API_GET_HISTORY      = "core/get/history"
	API_RESET_OP         = "core/reset/op"
	INT_GET_ARRAY_STATUS = "int/array/get/status"
	API_GET_TREE         = "array/get/tree"
	API_LOCATE_FOLDER    = "core/locate/folder"
	API_GET_LOG          = "api/get/log"

	API_SCATTER_CALCULATE          = "api/scatter/calculate"
	INT_SCATTER_CALCULATE          = "int/scatter/calc"
	INT_SCATTER_CALCULATE_FINISHED = "int/scatter/calc/finished"
	INT_SCATTER_CALCULATE_ERROR    = "int/scatter/calc/error"

	WS_CALC_STARTED  = "calcStarted"
	WS_CALC_PROGRESS = "calcProgress"
	WS_CALC_FINISHED = "calcFinished"
	WS_CALC_ISSUES   = "calcPermIssue"

	API_SCATTER_MOVE     = "api/scatter/move"
	API_SCATTER_COPY     = "api/scatter/copy"
	API_SCATTER_VALIDATE = "api/scatter/validate"

	INT_OPERATION_FINISHED = "core/operation/finished"

	API_GATHER_CALCULATE          = "api/gather/calculate"
	INT_GATHER_CALCULATE          = "int/gather/calc"
	INT_GATHER_CALCULATE_FINISHED = "int/gather/calc/finished"

	API_GATHER_MOVE = "api/gather/move"

	API_TOGGLE_DRYRUN   = "config/toggle/dryrun"
	API_NOTIFY_CALC     = "config/notify/calc"
	API_NOTIFY_MOVE     = "config/notify/move"
	API_SET_RESERVED    = "config/set/reserved"
	API_SET_VERBOSITY   = "config/set/verbosity"
	API_SET_CHECKUPDATE = "config/set/checkupdate"
	API_GET_UPDATE      = "config/get/update"
)
