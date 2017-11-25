import * as constant from '../lib/const'

// const getStatus = ({ state, actions, opts: { api } }) => {
// 	actions.setBusy(true)

// 	api.getStatus().then(json => {
// 		actions.gotStatus(json)
// 		actions.setBusy(false)
// 	})

// 	return state
// }

// const gotStatus = ({ state }, status) => {
// 	const lines = []

// 	let pathname = '/'
// 	let line = ''

// 	switch (status) {
// 		case constant.OP_SCATTER_PLAN:
// 			line = 'PLANNING: in progress ...'
// 			break
// 		case constant.OP_SCATTER_MOVE:
// 			line = 'MOVE: in progress ...'
// 			break
// 		case constant.OP_SCATTER_COPY:
// 			line = 'COPY: in progress ...'
// 			break
// 		case constant.OP_SCATTER_VALIDATE:
// 			line = 'VALIDATE: in progress ...'
// 			break
// 		case constant.OP_GATHER_PLAN:
// 			line = 'FIND TARGET: in progress ...'
// 			pathname = '/gather/target'
// 			break
// 		case constant.OP_GATHER_MOVE:
// 			line = 'MOVE: in progress ...'
// 			pathname = '/gather/move'
// 			break
// 		default:
// 			break
// 	}

// 	if (line !== '') {
// 		lines.push(line)
// 	}

// 	state.history.replace({ pathname })

// 	return {
// 		...state,
// 		core: {
// 			...state.core,
// 			status,
// 		},
// 		env: {
// 			...state.env,
// 			lines,
// 		},
// 	}
// }

const getState = ({ state, actions, opts: { api } }, mode) => {
	actions.setBusy(true)

	api.getState(mode).then(json => {
		actions.gotState(json)
		actions.setBusy(false)
	})
	// here i can catch the error and show an appropriate message

	return state
}

const gotState = ({ state, actions }, data) => {
	const lines = []

	let pathname = '/'
	let line = ''

	switch (status) {
		case constant.OP_SCATTER_PLAN:
			line = 'PLANNING: in progress ...'
			break
		case constant.OP_SCATTER_MOVE:
			line = 'MOVE: in progress ...'
			pathname = '/transfer'
			break
		case constant.OP_SCATTER_COPY:
			line = 'COPY: in progress ...'
			pathname = '/transfer'
			break
		case constant.OP_SCATTER_VALIDATE:
			line = 'VALIDATE: in progress ...'
			pathname = '/transfer'
			break
		case constant.OP_GATHER_PLAN:
			line = 'FIND TARGET: in progress ...'
			pathname = '/gather/target'
			break
		case constant.OP_GATHER_MOVE:
			line = 'MOVE: in progress ...'
			pathname = '/transfer'
			break
		default:
			break
	}

	if (line !== '') {
		lines.push(line)
	}

	state.history.replace({ pathname })

	const { history, historyOrder } = buildHistory(data.history)

	return {
		...state,
		core: {
			status: data.status,
			unraid: data.unraid,
			operation: data.operation,
			history,
			historyOrder,
		},
		scatter: {
			...state.scatter,
			plan: initPlan(data.unraid.disks),
		},
		gather: {
			...state.gather,
			plan: initPlan(data.unraid.disks),
		},
		env: {
			...state.env,
			lines,
		},
	}
}

const getStorage = ({ state, actions, opts: { api } }, history) => {
	actions.setBusy(true)

	api.getStorage().then(json => {
		actions.setBusy(false)
		actions.gotStorage(json)
	})

	return state
}

const gotStorage = ({ state }, unraid) => {
	return {
		...state,
		core: {
			...state.core,
			unraid,
		},
	}
}

const getHistory = ({ state, actions, opts: { api } }, history) => {
	actions.setBusy(true)

	api.getHistory().then(json => {
		actions.setBusy(false)
		actions.gotHistory(json)
	})

	return state
}

const gotHistory = ({ state }, data) => {
	const { history, historyOrder } = buildHistory(data)

	return {
		...state,
		core: {
			...state.core,
			history,
			historyOrder,
		},
	}
}

// const resetOperation = ({ state, actions, opts: { api } }) => {
// 	actions.setBusy(true)

// 	api.resetOperation().then(json => {
// 		actions.gotOperation(json)
// 		actions.setBusy(false)
// 	})

// 	return state
// }

// const gotOperation = ({ state }, operation) => {
// 	return {
// 		...state,
// 		core: {
// 			...state.core,
// 			operation,
// 		},
// 	}
// }

// function gatherMove({ state, actions, opts: { ws } }, drive) {
// 	actions.setBusy(true)

// 	ws.send({ topic: 'api/gather/move', payload: drive })

// 	state.history.replace({ pathname: '/transfer' })

// 	return state
// }

const transferStarted = ({ state }, operation) => {
	return {
		...state,
		core: {
			...state.core,
			operation,
		},
	}
}

const transferProgress = ({ state }, operation) => {
	return {
		...state,
		core: {
			...state.core,
			operation,
		},
	}
}

const transferFinished = ({ state, actions }, operation) => {
	actions.setBusy(false)

	return {
		...state,
		core: {
			...state.core,
			operation,
		},
	}
}

const flipOperation = ({ state }, id) => {
	const operation = { ...state.core.history[id] }
	operation.open = !operation.open

	return {
		...state,
		core: {
			...state.core,
			history: {
				...state.core.history,
				[id]: operation,
			},
		},
	}
}

const initPlan = disks => {
	const vdisks = disks.reduce((map, disk) => {
		map[disk.path] = {
			path: disk.path,
			plannedFree: disk.free,
			src: false,
			dst: false,
		}
		return map
	}, {})

	const plan = {
		chosenFolders: [],
		foldersNotTransferred: [],
		ownerIssue: 0,
		groupIssue: 0,
		folderIssue: 0,
		fileIssue: 0,
		bytesToTransfer: 0,
		vdisks,
	}

	return plan
}

const buildHistory = data => {
	const history = data.reduce((map, operation) => {
		map[operation.id] = operation
		return map
	}, {})

	const historyOrder = Object.keys(history).reverse()

	return { history, historyOrder }
}

export default {
	// getStatus,
	// gotStatus,

	getState,
	gotState,

	getStorage,
	gotStorage,

	getHistory,
	gotHistory,

	// resetOperation,
	// gotOperation,

	// gatherMove,

	transferStarted,
	transferProgress,
	transferFinished,

	flipOperation,
}
