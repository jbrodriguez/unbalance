const getConfig = ({ state, actions, opts: { api } }) => {
	actions.setBusy(true)

	api.getConfig().then(json => {
		actions.gotConfig(json)
		actions.setBusy(false)
	})
	// here i can catch the error and show an appropriate message

	return state
}

const gotConfig = ({ state }, config) => {
	return {
		...state,
		config,
	}
}

const setNotifyCalc = ({ state, actions, opts: { api } }, notify) => {
	if (state.config.notifyCalc !== notify) {
		api.setNotifyCalc(notify).then(json => actions.gotConfig(json))
	}

	return state
}

const setNotifyMove = ({ state, actions, opts: { api } }, notify) => {
	if (state.config.notifyMove !== notify) {
		api.setNotifyMove(notify).then(json => actions.gotConfig(json))
	}

	return state
}

const setReservedSpace = ({ state, actions, opts: { api } }, stringAmount, unit) => {
	const amount = Number(stringAmount)

	switch (unit) {
		case '%':
			if (amount < 0 || amount > 100) {
				actions.addFeedback('Percentage value must be between 0 and 100')
				return state
			}
			break

		case 'Gb':
			if (amount < 0.45) {
				actions.addFeedback('Gb value must be higher than 0.45')
				return state
			}
			break

		case 'Mb':
		default:
			if (amount < 450) {
				actions.addFeedback('Mb value must be higher than 450')
				return state
			}
			break
	}

	if (state.config.reservedAmount !== amount || state.config.reservedUnit !== unit) {
		actions.setBusy(true)

		api.setReservedSpace(amount, unit).then(json => {
			actions.gotConfig(json)
			actions.setBusy(false)
		})
	}

	return state
}

const toggleDryRun = ({ state, actions, opts: { api } }) => {
	actions.setBusy(true)

	api.toggleDryRun().then(json => {
		actions.dryRunToggled(json)
		actions.setBusy(false)
	})

	return state
}

const dryRunToggled = ({ state, actions }, config) => {
	return {
		...state,
		config,
	}
}

const setRsyncArgs = ({ state, actions, opts: { api } }, args) => {
	api.setRsyncArgs(args).then(json => actions.gotConfig(json))

	return state
}

const setVerbosity = ({ state, actions, opts: { api } }, verbosity) => {
	if (state.config.verbosity !== verbosity) {
		api.setVerbosity(verbosity).then(json => actions.gotConfig(json))
	}

	return state
}

const setUpdateCheck = ({ state, actions, opts: { api } }, checkForUpdate) => {
	if (state.config.checkForUpdate !== checkForUpdate) {
		api.setUpdateCheck(checkForUpdate).then(json => actions.gotConfig(json))
	}

	return state
}

export default {
	getConfig,
	gotConfig,

	setNotifyCalc,
	setNotifyMove,

	setReservedSpace,

	toggleDryRun,
	dryRunToggled,

	setRsyncArgs,

	setVerbosity,
	setUpdateCheck,
}
