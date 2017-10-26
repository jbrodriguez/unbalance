import * as constant from '../lib/const'

module.exports = {
	getState,
	gotState,

	setState,
}

function getState({ state, actions, opts: { api } }) {
	actions.setOpInProgress('Getting status')

	api.getState().then(json => {
		// console.log(`state(${JSON.stringify(json)})`)
		actions.gotState(json)
	})
	// here i can catch the error and show an appropriate message

	return state
}

function gotState({ state, actions }, backendState) {
	const lines = []

	let pathname = '/'
	let line = ''

	switch (backendState.operation.state) {
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

	const unraid = backendState.unraid
	const operation = backendState.operation

	return {
		...state,
		unraid,
		operation,
		lines,
	}
}

function setState({ state, actions }, backendState) {
	return {
		...state,
		unraid: backendState.unraid,
		operation: backendState.operation,
	}
}
