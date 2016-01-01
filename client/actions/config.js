module.exports = [
	{type: "getConfig", fn: _getConfig},
	{type: "gotConfig", fn: _gotConfig},

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

function _addFolder({state, actions, dispatch}, {api, _}, folder) {
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
