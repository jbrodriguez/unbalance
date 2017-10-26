package common

const (
	CHAN_CAPACITY = 3

	StateIdle        = 0
	StateCalc        = 1
	StateMove        = 2
	StateCopy        = 3
	StateValidate    = 4
	StateFindTargets = 5
	StateGather      = 6

	API_GET_CONFIG   = "core/getConfig"
	API_GET_STATE    = "core/getState"
	GET_ARRAY_STATUS = "array/getStatus"
	API_GET_TREE     = "array/getTree"
)
