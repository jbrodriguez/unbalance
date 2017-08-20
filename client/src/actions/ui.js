module.exports = {
	setOpInProgress,
	addFeedback,
	removeFeedback,
	clearConsole,
	checkForUpdate,
	updateAvailable,
	removeUpdateAvailable,
}

function setOpInProgress({ state }, action) {
	return {
		...state,
		opInProgress: action,
	}
}

function addFeedback({ state, actions }, feedback) {
	if (state.timeout) {
		window.clearTimeout(state.timeout)
	}
	const timeout = window.setTimeout(() => actions.removeFeedback(), 15000)

	return {
		...state,
		timeout,
		feedback: [].concat(feedback),
	}
}

function removeFeedback({ state }) {
	return {
		...state,
		feedback: [],
		timeout: null,
	}
}

function clearConsole({ state }) {
	return {
		...state,
		lines: [],
	}
}

function checkForUpdate({ state, actions, opts: { api } }) {
	console.log(`checking`)
	api.checkForUpdate().then(json => actions.updateAvailable(json))
	return state
}

function updateAvailable({ state }, version) {
	console.log(`version-${JSON.stringify(version)}`)
	return {
		...state,
		latestVersion: version,
	}
}

function removeUpdateAvailable({ state }) {
	return {
		...state,
		latestVersion: '',
	}
}
