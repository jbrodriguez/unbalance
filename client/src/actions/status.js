import * as constant from '../lib/const'

module.exports = {
	getStatus,
	gotStatus,

	setStatus,
}

function getStatus({ state, actions, opts: { api } }) {
	actions.setOpInProgress('Getting status')

	api.getStatus().then(json => actions.gotStatus(json.status))
	// here i can catch the error and show an appropriate message

	return state
}

function gotStatus({ state, actions }, status) {
	const lines = []
	let pathname = '/'
	let opState = null

	switch (status) {
		case constant.stateCalc:
			opState = 'Calculate operation in progress ...'
			break
		case constant.stateMove:
			opState = 'Move operation in progress ...'
			break
		case constant.stateCopy:
			opState = 'Copy operation in progress ...'
			break
		case constant.stateValidate:
			opState = 'Validate operation in progress ...'
			break
		case constant.stateFindTargets:
			opState = 'Find target operation in progress ...'
			pathname = '/gather/target'
			break
		case constant.stateGather:
			opState = 'Move operation in progress ...'
			pathname = '/gather/move'
			break
		default:
			break
	}

	if (opState) {
		lines.push(opState)
	}

	state.history.replace({ pathname })

	return {
		...state,
		status,
		opInProgress: opState,
		lines,
	}
}

function setStatus({ state, actions }, status) {
	return {
		...state,
		status,
	}
}
