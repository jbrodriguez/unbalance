module.exports = {
	getStorage,
	gotStorage,

	checkFrom,
	checkTo,

	calculate,
	calcStarted,
	calcProgress,
	calcFinished,

	move,
	moveStarted,
	moveProgress,
	moveFinished,

	opError,
}

function getStorage({state, actions, opts: {api}}) {
	actions.setOpInProgress("Getting storage info")

	api.getStorage()
		.then(json => {
			actions.gotStorage(json)
		})
	// here i can catch the error and show an appropriate message

	return state
}

function gotStorage({state}, unraid) {
	// console.log('unraid: ', unraid)

	// let toDisk = {}
	// let fromDisk = {}
	// let maxFreeSize = 0
	// let maxFreePath = ""

	// unraid.disks.map( disk => {
	// 	toDisk[disk.path] = true
	// 	fromDisk[disk.path] = false

	// 	if (disk.free > maxFreeSize) {
	// 		maxFreeSize = disk.free
	// 		maxFreePath = disk.path
	// 	}

	// 	return disk
	// })

	// if (maxFreePath != "") {
	// 	toDisk[maxFreePath] = false
	// 	fromDisk[maxFreePath] = true
	// }

	let toDisk = {}
	let fromDisk = {}

	unraid.disks.forEach( disk => {
		fromDisk[disk.path] = disk.src
		toDisk[disk.path] = disk.dst
	})

	return {
		...state,
		unraid,
		fromDisk,
		toDisk,
		opInProgress: null,
		moveDisabled: true,
		lines: [],
	}
}

function checkFrom({state}, path) {
	let fromDisk = Object.assign({}, state.fromDisk)
	let toDisk = Object.assign({}, state.toDisk)

	for (var key in fromDisk) {
		if (key !== path) {
			fromDisk[key] = false
		}
	}
	fromDisk[path] = true

	for (var key in toDisk) {
		toDisk[key] = !(key === path)
	}

	return {
		...state,
		fromDisk,
		toDisk,
	}
}

function checkTo({state}, path) {
	return {
		...state,
		toDisk: {
			...state.toDisk,
			[path]: !state.toDisk[path]
		}
	}
}

function calculate({state, actions, opts: {api, ws}}) {
	actions.setOpInProgress("Calculating")

	let srcDisk = ""
	for (var key in state.fromDisk) {
		if (state.fromDisk[key]) {
			srcDisk = key
			break
		}
	}			

	ws.send({topic: "calculate", payload: {srcDisk, dstDisks: state.toDisk}})

	return state
}

function calcStarted({state}, line) {
	return {
		...state,
		lines: [].concat('CALCULATE: ' + line),
		unraid: {
			...state.unraid,
			disks: state.unraid.disks.map( disk => {
				disk.newFree = disk.free
				return disk
			})
		}
	}


	// let newState = Object.assign({}, state)

	// // make sure we clean out the lines array
	// newState.lines = [].concat('CALCULATE: ' + line)

	// newState.unraid.disks.forEach( disk => {
	// 	disk.newFree = disk.free
	// })

	// return newState


	// return {
	// 	...state,
	// 	lines: [].concat('CALCULATE: ' + payload),
	// }
}

function calcProgress({state}, line) {
	return {
		...state,
		opInProgress: "calculate",
		moveDisabled: true,
		lines: state.lines.concat('CALCULATE: ' + line),
	}

	// let newState = Object.assign({}, state)

	// // make sure we disable the interface, in case another browser is open
	// // or even the initial browser is woken from sleep 
	// newState.opInProgress = actions.calculate
	// newState.moveDisabled = true

	// newState.lines.push('CALCULATE: ' + line)

	// return newState

	// return {
	// 	...state,
	// 	lines: state.lines.concat('CALCULATE: ' + payload),
	// }
}

function calcFinished({state, actions}, unraid) {
	let feedback = []
	if (unraid.bytesToMove === 0) {
		feedback.push("There isn't available space in any of the target disks, to move the folders you selected.")
		feedback.push("Check more disks in the TO column or go to the Settings page, to review the folders selected for moving or to change the amount of reserved space.")
	}

	window.setTimeout( _ => actions.removeFeedback(), 15*1000)

	return {
		...state,
		unraid,
		feedback,
		opInProgress: null,
		moveDisabled: unraid.bytesToMove === 0,
	}


	// let newState = Object.assign({}, state)

	// newState.unraid = unraid
	// newState.opInProgress = null
	// newState.moveDisabled = false

	// // if (newState.unraid.bytesToMove === 0) {
	// // 	newState.feedback.push("There's no space available to move any of the folders you selected.")
	// // 	newState.feedback.push("Check more disks in the TO column or go to the Settings page, to review the folders selected for moving.")
	// // }

	// return newState

	// return {
	// 	...state,
	// 	unraid,
	// 	opInProgress: null,
	// 	moveDisabled: false,
	// }
}

// // this message is received when the browser requests
// function calcIsRunning({state, actions, dispatch}, _, unraid) {
// 	let newState = Object.assign({}, state)

// 	newState.opInProgress = actions.calculate
// 	newState.moveDisabled = true
// 	// newState.lines.push('CALCULATE: ' + line)

// 	return newState
// }

function move({state, actions, opts: {api, ws}}) {
	actions.setOpInProgress("Moving")

	ws.send({topic: "move"})

	return state
}

function moveStarted({state}, line) {
	return {
		...state,
		lines: [].concat('MOVE: ' + line)
	}

	// let newState = Object.assign({}, state)

	// // make sure we clean out the lines array
	// newState.lines = [].concat('MOVE: ' + line)

	// return newState

	// return {
	// 	...state,
	// 	lines: [].concat('MOVE: ' + payload),
	// }
}

function moveProgress({state}, line) {
	return {
		...state,
		lines: state.lines.concat('MOVE: ' + line)
	}

	// return {
	// 	...state,
	// 	lines: state.lines.concat('MOVE: ' + payload),
	// }
}

function moveFinished({state}, unraid) {
	return {
		...state,
		unraid,
		opInProgress: null,
		moveDisabled: !state.config.dryRun
	}


	// let newState = Object.assign({}, state)

	// newState.unraid = unraid
	// newState.opInProgress = null
	// newState.moveDisabled = !state.config.dryRun

	// return newState
	
	// let moveDisabled = !state.config.dryRun
	// console.log('moveDisabled: ', moveDisabled)
	// return {
	// 	...state,
	// 	unraid,
	// 	opInProgress: null,
	// 	moveDisabled,
	// }
}

function opError({state, actions}, error) {
	actions.addFeedback(error)
	return state
	// return {
	// 	...state,
	// 	feedback: [].concat(error)
	// }
	// let newState = Object.assign({}, state)

	// newState.feedback.push(error)
	
	// return newState
}
