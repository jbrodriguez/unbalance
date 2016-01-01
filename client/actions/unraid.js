module.exports = [
	{type: "getStorage", fn: _getStorage},
	{type: "gotStorage", fn: _gotStorage},

	{type: "checkFrom", fn: _checkFrom},
	{type: "checkTo", fn: _checkTo},

	{type: "calculate", fn: _calculate},
	{type: "calcStarted", fn: _calcStarted},
	{type: "calcProgress", fn: _calcProgress},
	{type: "calcFinished", fn: _calcFinished},

	{type: "move", fn: _move},
	{type: "moveStarted", fn: _moveStarted},
	{type: "moveProgress", fn: _moveProgress},
	{type: "moveFinished", fn: _moveFinished},
]

function _getStorage({state, actions, dispatch}, {api, _}) {
	dispatch(actions.opInProgress, actions.getStorage)

	api.getStorage()
		.then(json => {
			dispatch(actions.gotStorage, json)
		})
	// here i can catch the error and show an appropriate message

	return state
}

function _gotStorage({state, actions, dispatch}, _, unraid) {
	// console.log('unraid: ', unraid)

	let toDisk = {}
	let fromDisk = {}
	let maxFreeSize = 0
	let maxFreePath = ""

	unraid.disks.map( disk => {
		toDisk[disk.path] = true
		fromDisk[disk.path] = false

		if (disk.free > maxFreeSize) {
			maxFreeSize = disk.free
			maxFreePath = disk.path
		}

		return disk
	})

	if (maxFreePath != "") {
		toDisk[maxFreePath] = false
		fromDisk[maxFreePath] = true
	}

	let newState = Object.assign({}, state)

	newState.opInProgress = null
	newState.moveDisabled = true
	newState.lines = []
	newState.unraid = unraid
	newState.fromDisk = fromDisk
	newState.toDisk = toDisk

	return newState
	// return {
	// 	...state,
	// 	opInProgress: null,
	// 	moveDisabled: true,
	// 	lines: [],
	// 	unraid,
	// 	fromDisk,
	// 	toDisk,
	// }
}

function _checkFrom({state, actions, dispatch}, _, path) {
	let newState = Object.assign({}, state)

	for (var key in newState.fromDisk) {
		if (key !== path) {
			newState.fromDisk[key] = false
		}
	}
	newState.fromDisk[path] = true

	for (var key in newState.toDisk) {
		newState.toDisk[key] = !(key === path)
	}

	return newState
}

function _checkTo({state, actions, dispatch}, _, path) {
	let newState = Object.assign({}, state)

	newState.toDisk[path] = !newState.toDisk[path]

	return newState
}

function _calculate({state, actions, dispatch}, {_, ws}) {
	dispatch(actions.opInProgress, actions.calculate)

	let srcDisk = ""

	for (var key in state.fromDisk) {
		if (state.fromDisk[key]) {
			srcDisk = key
			break
		}
	}			

	ws.send({topic: actions.calculate, payload: {srcDisk, dstDisks: state.toDisk}})

	return state
}

function _calcStarted({state, actions, dispatch}, _, line) {
	let newState = Object.assign({}, state)

	// make sure we clean out the lines array
	newState.lines = [].concat('CALCULATE: ' + line)

	newState.unraid.disks.forEach( disk => {
		disk.newFree = disk.free
	})

	return newState


	// return {
	// 	...state,
	// 	lines: [].concat('CALCULATE: ' + payload),
	// }
}

function _calcProgress({state, actions, dispatch}, _, line) {
	let newState = Object.assign({}, state)

	newState.lines.push('CALCULATE: ' + line)

	return newState

	// return {
	// 	...state,
	// 	lines: state.lines.concat('CALCULATE: ' + payload),
	// }
}

function _calcFinished({state, actions, dispatch}, _, unraid) {
	let newState = Object.assign({}, state)

	newState.unraid = unraid
	newState.opInProgress = null
	newState.moveDisabled = false

	return newState

	// return {
	// 	...state,
	// 	unraid,
	// 	opInProgress: null,
	// 	moveDisabled: false,
	// }
}

function _move({state, actions, dispatch}, {_, ws}) {
	dispatch(actions.opInProgress, actions.move)

	ws.send({topic: actions.move})

	return state
}

function _moveStarted({state, actions, dispatch}, _, line) {
	let newState = Object.assign({}, state)

	// make sure we clean out the lines array
	newState.lines = [].concat('MOVE: ' + line)

	return newState

	// return {
	// 	...state,
	// 	lines: [].concat('MOVE: ' + payload),
	// }
}

function _moveProgress({state, actions, dispatch}, _, line) {
	let newState = Object.assign({}, state)

	newState.lines.push('MOVE: ' + line)

	return newState
	// return {
	// 	...state,
	// 	lines: state.lines.concat('MOVE: ' + payload),
	// }
}

function _moveFinished({state, actions, dispatch}, _, unraid) {
	let newState = Object.assign({}, state)

	newState.unraid = unraid
	newState.opInProgress = null
	newState.moveDisabled = !state.config.dryRun

	return newState
	
	// let moveDisabled = !state.config.dryRun
	// console.log('moveDisabled: ', moveDisabled)
	// return {
	// 	...state,
	// 	unraid,
	// 	opInProgress: null,
	// 	moveDisabled,
	// }
}
