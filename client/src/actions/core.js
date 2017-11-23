import * as constant from '../lib/const'

const getStatus = ({ state, actions, opts: { api } }) => {
	actions.setBusy(true)

	api.getStatus().then(json => {
		actions.gotStatus(json)
		actions.setBusy(false)
	})

	return state
}

const gotStatus = ({ state }, status) => {
	const lines = []

	let pathname = '/'
	let line = ''

	switch (status) {
		case constant.OP_SCATTER_CALC:
			line = 'CALCULATE: in progress ...'
			break
		case constant.OP_SCATTER_MOVE:
			line = 'MOVE: in progress ...'
			break
		case constant.OP_SCATTER_COPY:
			line = 'COPY: in progress ...'
			break
		case constant.OP_SCATTER_VALIDATE:
			line = 'VALIDATE: in progress ...'
			break
		case constant.OP_GATHER_CALC:
			line = 'FIND TARGET: in progress ...'
			pathname = '/gather/target'
			break
		case constant.OP_GATHER_MOVE:
			line = 'MOVE: in progress ...'
			pathname = '/gather/move'
			break
		default:
			break
	}

	if (line !== '') {
		lines.push(line)
	}

	state.history.replace({ pathname })

	return {
		...state,
		core: {
			...state.core,
			status,
		},
		env: {
			...state.env,
			lines,
		},
	}
}

const getState = ({ state, actions, opts: { api } }, mode) => {
	actions.setBusy(true)

	api.getState(mode).then(json => {
		// console.log(`state(${JSON.stringify(json)})`)
		actions.gotState(json)
		actions.setBusy(false)
	})
	// here i can catch the error and show an appropriate message

	return state
}

const gotState = ({ state, actions }, core) => {
	// const lines = []

	// let pathname = '/'
	// let line = ''

	// switch (core.status) {
	// 	case constant.stateCalc:
	// 		line = 'CALCULATE: in progress ...'
	// 		break
	// 	case constant.stateMove:
	// 		line = 'MOVE: in progress ...'
	// 		break
	// 	case constant.stateCopy:
	// 		line = 'COPY: in progress ...'
	// 		break
	// 	case constant.stateValidate:
	// 		line = 'VALIDATE: in progress ...'
	// 		break
	// 	case constant.stateFindTargets:
	// 		line = 'FIND TARGET: in progress ...'
	// 		pathname = '/gather/target'
	// 		break
	// 	case constant.stateGather:
	// 		line = 'MOVE: in progress ...'
	// 		pathname = '/gather/move'
	// 		break
	// 	default:
	// 		break
	// }

	// if (line !== '') {
	// 	lines.push(line)
	// }

	// state.history.replace({ pathname })

	return {
		...state,
		core,
		// env: {
		// 	...state.env,
		// 	lines,
		// },
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

const gotHistory = ({ state }, history) => {
	return {
		...state,
		core: {
			...state.core,
			history,
		},
	}
}

const resetOperation = ({ state, actions, opts: { api } }) => {
	actions.setBusy(true)

	api.resetOperation().then(json => {
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

const scatterCalculate = ({ state, actions, opts: { ws } }, fromDisk, toDisk) => {
	actions.setBusy(true)

	const operation = { ...state.core.operation }

	let srcDisk = ''

	state.core.unraid.disks.forEach(disk => {
		operation.vdisks[disk.path].src = fromDisk[disk.path]
		operation.vdisks[disk.path].dst = toDisk[disk.path]

		if (operation.vdisks[disk.path].src) {
			srcDisk = disk.path
		}
	})

	console.log(`chosen(${JSON.stringify(state.scatter.chosen)})`)

	const folders = Object.keys(state.scatter.chosen).map(folder => folder.slice(srcDisk.length + 1))

	operation.chosenFolders = folders

	ws.send({ topic: 'api/scatter/calculate', payload: operation })

	return state
}

const calcStarted = ({ state }, line) => {
	return {
		...state,
		env: {
			...state.env,
			lines: [].concat(`CALCULATE: ${line}`),
		},
	}
}

const calcProgress = ({ state }, line) => {
	const lines = state.env.lines.length > 1000 ? [] : state.env.lines

	return {
		...state,
		env: {
			...state.env,
			lines: lines.concat(`CALCULATE: ${line}`),
		},
	}
}

const calcFinished = ({ state, actions }, operation) => {
	if (operation.bytesToTransfer === 0) {
		const feedback = []

		feedback.push('The calculate operation found that no folders/files can be moved/copied.')
		feedback.push('')
		feedback.push('This might be due to one of the following reasons:')
		feedback.push(
			'- The source share(s)/folder(s) you selected are either empty or do not exist in the source disk',
		)
		feedback.push(
			"- There isn't available space in any of the target disks, to move/copy the share(s)/folder(s) you selected",
		)
		feedback.push('')
		feedback.push(
			'Check more disks in the TO column or go to the Settings page, to review the share(s)/folder(s) selected for moving/copying or to change the amount of reserved space.',
		)

		actions.addFeedback(feedback)
	}

	actions.setBusy(false)

	return {
		...state,
		core: {
			...state.core,
			operation,
		},
	}
}

const calcPermIssue = ({ state, actions }, permStats) => {
	const permIssues = permStats.split('|')

	const feedback = []

	feedback.push('There are some permission issues with the folders/files you want to move')
	feedback.push(`${permIssues[0]} file(s)/folder(s) with an owner other than 'nobody'`)
	feedback.push(`${permIssues[1]} file(s)/folder(s) with a group other than 'users'`)
	feedback.push(`${permIssues[2]} folder(s) with a permission other than 'drwxrwxrwx'`)
	feedback.push(`${permIssues[3]} files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'`)
	feedback.push('You can find more details about which files have issues in the log file (/boot/logs/unbalance.log)')
	feedback.push('')
	feedback.push(
		'At this point, you can move the folders/files if you want, but be advised that it can cause errors in the operation',
	)
	feedback.push('')
	feedback.push(
		'You are STRONGLY suggested to install the Fix Common Problems plugin, then run the Docker Safe New Permissions command',
	)

	actions.addFeedback(feedback)

	return state
}

const scatterMove = ({ state, actions, opts: { ws } }) => {
	actions.setBusy(true)

	ws.send({ topic: 'api/scatter/move' })

	state.history.replace({ pathname: '/transfer' })

	return state
}

const scatterCopy = ({ state, actions, opts: { ws } }) => {
	actions.setBusy(true)

	ws.send({ topic: 'api/scatter/copy' })

	state.history.replace({ pathname: '/transfer' })

	return state
}

function gatherMove({ state, actions, opts: { ws } }, drive) {
	actions.setBusy(true)

	ws.send({ topic: 'api/gather/move', payload: drive })

	state.history.replace({ pathname: '/transfer' })

	return state
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

// function setState({ state, actions }, backendState) {
// 	return {
// 		...state,
// 		unraid: backendState.unraid,
// 		operation: backendState.operation,
// 	}
// }

export default {
	getStatus,
	gotStatus,

	getState,
	gotState,

	getHistory,
	gotHistory,

	resetOperation,
	gotOperation,

	scatterCalculate,
	calcStarted,
	calcProgress,
	calcFinished,
	calcPermIssue,
	// setState,

	scatterMove,
	scatterCopy,

	gatherMove,

	transferStarted,
	transferProgress,
	transferFinished,
}
