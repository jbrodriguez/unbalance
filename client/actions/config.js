module.exports = [
	{type: "getConfig", fn: _getConfig},
	{type: "gotConfig", fn: _gotConfig},

	{type: "setNotifyCalc", fn: _setNotifyCalc},
	{type: "setNotifyMove", fn: _setNotifyMove},

	{type: "addFolder", fn: _addFolder},
	{type: "folderAdded", fn: _folderAdded},
	{type: "deleteFolder", fn: _deleteFolder},
	{type: "folderDeleted", fn: _folderDeleted},

	{type: "toggleDryRun", fn: _toggleDryRun},
	{type: "dryRunToggled", fn: _dryRunToggled},
]

function _getConfig({state, actions, dispatch}, {api, _}) {
	dispatch(actions.opInProgress, actions.getConfig)

	api.getConfig()
		.then(json => {
			dispatch(actions.gotConfig, json)
		})
	// here i can catch the error and show an appropriate message

	return state
}

function _gotConfig({state, actions, dispatch}, _, config) {
	let newState = Object.assign({}, state)

	newState.config = config
	newState.opInProgress = null

	return newState
}

function _setNotifyCalc({state, actions, dispatch}, {api, _}, notify) {
	if (state.config.notifyCalc !== notify) {
		api.setNotifyCalc(notify)
			.then(json => {
				dispatch(actions.gotConfig, json)
			})
	}

	return state
}

function _setNotifyMove({state, actions, dispatch}, {api, _}, notify) {
	if (state.config.notifyMove !== notify) {
		api.setNotifyMove(notify)
			.then(json => {
				dispatch(actions.gotConfig, json)
			})
	}

	return state
}


function _addFolder({state, actions, dispatch}, {api, _}, folder) {
	let proceed = true
	state.config.folders.some( chosen => {
		if (folder === chosen || chosen.indexOf(folder) > -1 || folder.indexOf(chosen) > -1) {
			proceed = false
			return true
		}
	})

	if (!proceed) {
		let newState = Object.assign({}, state)
		newState.feedback = [].concat(["The folder you're trying to add is already selected, contains or is contained by an already selected folder. Please choose another folder or remove one of the selected folders and try again."])

		// set a seven second timeout for the feedback panel
		window.setTimeout( _ => dispatch(actions.removeFeedback), 7*1000)

		return newState
	}

	if (state.config.folders.indexOf(folder) !== -1) {
		return state
	}

	dispatch(actions.opInProgress, actions.addFolder)

	api.addFolder(folder)
		.then(json => {
			dispatch(actions.folderAdded, json)
		})

	return state
}

function _folderAdded({state, actions, dispatch}, _, config) {
	let newState = Object.assign({}, state)

	newState.config = config
	newState.opInProgress = null

	return newState
}

function _deleteFolder({state, actions, dispatch}, {api, _}, folder) {
	dispatch(actions.opInProgress, actions.deleteFolder)

	api.deleteFolder(folder)
		.then(json => {
			dispatch(actions.folderDeleted, json)
		})

	return state
}

function _folderDeleted({state, actions, dispatch}, _, config) {
	let newState = Object.assign({}, state)

	newState.config = config
	newState.opInProgress = null

	return newState
}

function _toggleDryRun({state, actions, dispatch}, {api, _}) {
	dispatch(actions.opInProgress, actions.toggleDryRun)

	api.toggleDryRun()
		.then(json => {
			dispatch(actions.dryRunToggled, json)
		})	

	return state			
}

function _dryRunToggled({state, actions, dispatch}, _, config) {
	let newState = Object.assign({}, state)

	newState.config = config
	newState.opInProgress = null

	return newState
}
