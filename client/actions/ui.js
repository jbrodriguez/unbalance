module.exports = {
	setOpInProgress,
	addFeedback,
	removeFeedback,
}

function setOpInProgress({state}, action) {
	return {
		...state,
		opInProgress: action
	}
}

function addFeedback({state}, feedback) {
	return {
		...state,
		feedback: [].concat[feedback]
	}
}

function removeFeedback({state}) {
	return {
		...state,
		feedback: []
	}
}
