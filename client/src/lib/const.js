export const OP_NEUTRAL = 0
export const OP_SCATTER_PLAN = 1
export const OP_SCATTER_MOVE = 2
export const OP_SCATTER_COPY = 3
export const OP_SCATTER_VALIDATE = 4
export const OP_GATHER_PLAN = 5
export const OP_GATHER_MOVE = 6

export const opMap = {
	[OP_SCATTER_MOVE]: 'SCATTER / MOVE',
	[OP_SCATTER_COPY]: 'SCATTER / COPY',
	[OP_GATHER_MOVE]: 'GATHER / MOVE',
}

export const API_SCATTER_PLAN = 'api/scatter/plan'
export const API_SCATTER_MOVE = 'api/scatter/move'
export const API_SCATTER_COPY = 'api/scatter/copy'

export const API_GATHER_PLAN = 'api/gather/plan'
