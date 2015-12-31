module.exports = [
	{type: "opInProgress", fn: _opInProgress},
	{type: "removeAlert", fn: _removeAlert},
]

function _opInProgress({state, actions, dispatch}, _, action) {
	let newState = Object.assign({}, state)

	newState.opInProgress = action

	return newState
}

function _removeAlert({state, actions, dispatch}, _, action) {
	let newState = Object.assign({}, state)

	newState.alerts = []

	return newState
}
