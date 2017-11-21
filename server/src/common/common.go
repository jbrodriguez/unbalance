package common

const (
	MailCmd = "/usr/local/emhttp/webGui/scripts/notify"

	CHAN_CAPACITY = 3

	OP_NEUTRAL          = 0
	OP_SCATTER_CALC     = 1
	OP_SCATTER_MOVE     = 2
	OP_SCATTER_COPY     = 3
	OP_SCATTER_VALIDATE = 4
	OP_GATHER_CALC      = 5
	OP_GATHER_MOVE      = 6

	// StateIdle        = 0
	// StateCalc        = 1
	// StateMove        = 2
	// StateCopy        = 3
	// StateValidate    = 4
	// StateFindTargets = 5
	// StateGather      = 6

	API_GET_CONFIG       = "core/get/config"
	API_GET_STATUS       = "core/get/status"
	API_GET_STATE        = "core/get/state"
	API_RESET_OP         = "core/reset/op"
	INT_GET_ARRAY_STATUS = "int/array/get/status"
	API_GET_TREE         = "array/get/tree"
	API_LOCATE_FOLDER    = "core/locate/folder"

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
)
