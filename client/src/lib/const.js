export const OP_NEUTRAL = 0
export const OP_SCATTER_PLAN = 1
export const OP_SCATTER_MOVE = 2
export const OP_SCATTER_COPY = 3
export const OP_SCATTER_VALIDATE = 4
export const OP_GATHER_PLAN = 5
export const OP_GATHER_MOVE = 6

export const opMap = {
	[OP_SCATTER_MOVE]: { name: 'SCATTER / MOVE', color: 'opWhite' },
	[OP_SCATTER_COPY]: { name: 'SCATTER / COPY', color: 'opYellow' },
	[OP_GATHER_MOVE]: { name: 'GATHER / MOVE', color: 'opOrange' },
	[OP_SCATTER_VALIDATE]: { name: 'SCATTER / VALIDATE', color: 'opBlue' },
}

export const API_SCATTER_PLAN = 'api/scatter/plan'
export const API_SCATTER_MOVE = 'api/scatter/move'
export const API_SCATTER_COPY = 'api/scatter/copy'

export const API_GATHER_PLAN = 'api/gather/plan'
export const API_GATHER_MOVE = 'api/gather/move'

export const API_GET_LOG = 'api/get/log'

export const API_VALIDATE = 'api/validate'
export const API_REPLAY = 'api/replay'
