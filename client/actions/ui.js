module.exports = {
	setOpInProgress,
	addFeedback,
	removeFeedback,
	// setConsole,
}

function setOpInProgress({state}, action) {
	return {
		...state,
		opInProgress: action
	}
}

function addFeedback({state, actions}, feedback) {
	window.setTimeout(() => actions.removeFeedback(), 15000)

	return {
		...state,
		feedback: [].concat(feedback)
	}
}

function removeFeedback({state}) {
	return {
		...state,
		feedback: []
	}
}

// function setConsole({state}, line) {
// 	return {
// 		...state,
// 		lines: [].concat(line)
// 	}
// }