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
		case constant.stateCalc:
			line = 'CALCULATE: in progress ...'
			break
		case constant.stateMove:
			line = 'MOVE: in progress ...'
			break
		case constant.stateCopy:
			line = 'COPY: in progress ...'
			break
		case constant.stateValidate:
			line = 'VALIDATE: in progress ...'
			break
		case constant.stateFindTargets:
			line = 'FIND TARGET: in progress ...'
			pathname = '/gather/target'
			break
		case constant.stateGather:
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

function getState({ state, actions, opts: { api } }, mode) {
	actions.setBusy(true)

	api.getState(mode).then(json => {
		// console.log(`state(${JSON.stringify(json)})`)
		actions.gotState(json)
		actions.setBusy(false)
	})
	// here i can catch the error and show an appropriate message

	return state
}

function gotState({ state, actions }, core) {
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

	resetOperation,
	gotOperation,

	// setState,
}
