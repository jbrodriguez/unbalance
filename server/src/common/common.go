package common

const (
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

	API_GET_CONFIG   = "core/getConfig"
	API_GET_STATUS   = "core/getStatus"
	API_GET_STATE    = "core/getState"
	API_RESET_OP     = "core/resetOp"
	GET_ARRAY_STATUS = "array/getStatus"
	API_GET_TREE     = "array/getTree"
)
