package common

const (
	MailCmd = "/usr/local/emhttp/webGui/scripts/notify"

	CHAN_CAPACITY = 3

	OP_NEUTRAL      = 0
	OP_SCATTER_CALC = 1
	OP_GATHER_CALC  = 2
	OP_MOVE         = 3
	OP_COPY         = 4
	OP_VALIDATE     = 5

	StateIdle        = 0
	StateCalc        = 1
	StateMove        = 2
	StateCopy        = 3
	StateValidate    = 4
	StateFindTargets = 5
	StateGather      = 6

	API_GET_CONFIG       = "core/get/config"
	API_GET_STATUS       = "core/get/status"
	API_GET_STATE        = "core/get/state"
	API_RESET_OP         = "core/reset/op"
	INT_GET_ARRAY_STATUS = "int/array/get/status"
	API_GET_TREE         = "array/get/tree"
	API_LOCATE_FOLDER    = "core/locate/folder"

	API_CALCULATE_SCATTER          = "api/calculate/scatter"
	INT_CALCULATE_SCATTER          = "int/calc/scatter"
	INT_CALCULATE_SCATTER_FINISHED = "int/calc/scatter/finished"
	INT_CALCULATE_SCATTER_ERROR    = "int/calc/scatter/error"

	WS_CALC_STARTED  = "calcStarted"
	WS_CALC_PROGRESS = "calcProgress"
	WS_CALC_FINISHED = "calcFinished"
	WS_CALC_ISSUES   = "calcPermIssue"
)
