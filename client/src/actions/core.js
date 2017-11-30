import * as constant from '../lib/const'

const getState = ({ state, actions, opts: { api } }, mode) => {
	actions.setBusy(true)

	api.getState(mode).then(json => {
		actions.gotState(json)
		actions.setBusy(false)
	})
	// here i can catch the error and show an appropriate message

	return state
}

const gotState = ({ state, actions }, core) => {
	let pathname = '/'

	let scatterLine = []
	let gatherLine = []

	switch (status) {
		case constant.OP_SCATTER_PLAN:
			scatterLine = ['PLANNING: in progress ...']
			break
		case constant.OP_SCATTER_MOVE:
			pathname = '/transfer'
			break
		case constant.OP_SCATTER_COPY:
			pathname = '/transfer'
			break
		case constant.OP_SCATTER_VALIDATE:
			pathname = '/transfer'
			break
		case constant.OP_GATHER_PLAN:
			gatherLine = ['PLANNING: in progress ...']
			pathname = '/gather/target'
			break
		case constant.OP_GATHER_MOVE:
			pathname = '/transfer'
			break
		default:
			break
	}

	state.history.replace({ pathname })

	return {
		...state,
		core,
		scatter: {
			...state.scatter,
			plan: initPlan(core.unraid.disks),
			lines: scatterLine,
		},
		gather: {
			...state.gather,
			plan: initPlan(core.unraid.disks),
			lines: gatherLine,
		},
	}
}

const resetState = ({ state }) => {
	return {
		...state,
		scatter: {
			cache: null,
			chosen: {},
			items: [],
			plan: initPlan(state.core.unraid.disks),
			lines: [],
		},
		gather: {
			cache: null,
			chosen: {},
			items: [],
			plan: initPlan(state.core.unraid.disks),
			lines: [],
			location: null,
			target: null,
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

const getOperation = ({ state, actions, opts: { api } }) => {
	actions.setBusy(true)

	api.getOperation().then(json => {
		actions.gotOperation(json)
		actions.setBusy(false)
	})

	return state
}

const gotOperation = ({ state }, operation) => {
	return {
		...state,
		core: {
			...state.core,
			operation,
		},
	}
}

const getHistory = ({ state, actions, opts: { api } }) => {
	actions.setBusy(true)

	api.getHistory().then(json => {
		actions.setBusy(false)
		actions.gotHistory(json)
	})

	return state
}

const gotHistory = ({ state }, history) => {
	return {
		...state,
		core: {
			...state.core,
			history,
		},
	}
}

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

const transferFinished = ({ state, actions }, bState) => {
	actions.setBusy(false)

	return {
		...state,
		core: {
			...state.core,
			operation: bState.operation,
			history: bState.history,
		},
	}
}

const flipOperation = ({ state }, id) => {
	const operation = { ...state.core.history.items[id] }
	operation.open = !operation.open

	return {
		...state,
		core: {
			...state.core,
			history: {
				...state.core.history,
				items: {
					...state.core.history.items,
					[id]: operation,
				},
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
		ownerIssue: 0,
		groupIssue: 0,
		folderIssue: 0,
		fileIssue: 0,
		bytesToTransfer: 0,
		vdisks,
	}

	return plan
}

const validate = ({ state, actions, opts: { ws } }, id) => {
	actions.setBusy(true)
	actions.resetState()

	ws.send({ topic: constant.API_VALIDATE, payload: state.core.history.items[id] })

	state.history.replace({ pathname: '/transfer' })

	return state
}

const replay = ({ state, actions, opts: { ws } }, id) => {
	actions.setBusy(true)
	actions.resetState()

	ws.send({ topic: constant.API_REPLAY, payload: state.core.history.items[id] })

	state.history.replace({ pathname: '/transfer' })

	return state
}

export default {
	getState,
	gotState,

	resetState,

	getStorage,
	gotStorage,

	getOperation,
	gotOperation,

	getHistory,
	gotHistory,

	transferStarted,
	transferProgress,
	transferFinished,

	flipOperation,

	validate,
	replay,
}
