module.exports = [
	{type: "opInProgress", fn: _opInProgress},
	{type: "removeFeedback", fn: _removeFeedback},
]

function _opInProgress({state, actions, dispatch}, _, action) {
	let newState = Object.assign({}, state)

	newState.opInProgress = action

	return newState
}

function _removeFeedback({state, actions, dispatch}) {
	let newState = Object.assign({}, state)

	newState.feedback = []

	return newState
}
