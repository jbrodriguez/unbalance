import * as constant from '../lib/const'

const setBusy = ({ state }, isBusy) => {
	return {
		...state,
		env: {
			...state.env,
			isBusy,
		},
	}
}

const clearConsole = ({ state }) => {
	return {
		...state,
		env: {
			...state.env,
			lines: [],
		},
	}
}

const addFeedback = ({ state, actions }, feedback) => {
	if (state.env.timeout) {
		window.clearTimeout(state.env.timeout)
	}
	const timeout = window.setTimeout(() => actions.removeFeedback(), 15000)

	return {
		...state,
		env: {
			...state.env,
			timeout,
			feedback: [].concat(feedback),
		},
	}
}

const removeFeedback = ({ state }) => {
	return {
		...state,
		env: {
			...state.env,
			feedback: [],
			timeout: null,
		},
	}
}

const getLog = ({ state, actions, opts: { ws } }) => {
	actions.setBusy(true)

	ws.send({ topic: constant.API_GET_LOG })

	return state
}

const gotLog = ({ state, actions }, log) => {
	actions.setBusy(false)

	return {
		...state,
		env: {
			...state.env,
			log,
		},
	}
}

// function checkForUpdate({ state, actions, opts: { api } }) {
// 	// console.log(`checking`)
// 	api.checkForUpdate().then(json => actions.updateAvailable(json))
// 	return state
// }

// function updateAvailable({ state }, version) {
// 	// console.log(`version-${JSON.stringify(version)}`)
// 	return {
// 		...state,
// 		latestVersion: version,
// 	}
// }

// function removeUpdateAvailable({ state }) {
// 	return {
// 		...state,
// 		latestVersion: '',
// 	}
// }

export default {
	setBusy,

	clearConsole,

	addFeedback,
	removeFeedback,

	getLog,
	gotLog,
	// checkForUpdate,
	// updateAvailable,
	// removeUpdateAvailable,
}
