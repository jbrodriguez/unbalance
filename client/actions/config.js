module.exports = {
	getConfig,
	gotConfig,

	setNotifyCalc,
	setNotifyMove,

	addFolder,
	folderAdded,
	deleteFolder,
	folderDeleted,

	toggleDryRun,
	dryRunToggled,
}

function getConfig({state, actions, opts: {api}}) {
	actions.setOpInProgress("Getting configuration")

	api.getConfig()
		.then(json => {
			actions.gotConfig(json)
		})
	// here i can catch the error and show an appropriate message

	return state
}

function gotConfig({state, actions}, config) {
	return {
		...state,
		config,
		opInProgress: null
	}
}

function setNotifyCalc({state, actions, opts: {api}}, notify) {
	if (state.config.notifyCalc !== notify) {
		api.setNotifyCalc(notify)
			.then(json => {
				actions.gotConfig(json)
			})
	}

	return state
}

function setNotifyMove({state, actions, opts: {api}}, notify) {
	if (state.config.notifyMove !== notify) {
		api.setNotifyMove(notify)
			.then(json => {
				actions.gotConfig(json)
			})
	}

	return state
}

function addFolder({state, actions, opts: {api}}, folder) {
	const exists = state.config.folders.some( chosen => {
		return (folder === chosen || chosen.indexOf(folder) > -1 || folder.indexOf(chosen) > -1)
	})

	if (exists) {
		// set a seven second timeout to remove the feedback panel
		window.setTimeout( _ => actions.removeFeedback(), 15*1000)

		return {
			...state,
			feedback: [].concat(["The folder you're trying to add is already selected or contains or is contained by a folder that you already added. Please choose another folder or remove one of the selected folders and try again."])
		}
	}

	// if (state.config.folders.indexOf(folder) !== -1) {
	// 	return state
	// }

	actions.setOpInProgress("Adding folder")

	api.addFolder(folder)
		.then(json => {
			actions.folderAdded(json)
		})

	return state
}

function folderAdded({state}, config) {
	return {
		...state,
		config,
		opInProgress: null
	}
}

function deleteFolder({state, actions, opts: {api}}, folder) {
	actions.setOpInProgress("Deleting folder")

	api.deleteFolder(folder)
		.then(json => {
			actions.folderDeleted(json)
		})

	return state
}

function folderDeleted({state}, config) {
	return {
		...state,
		config,
		opInProgress: null
	}
}

function toggleDryRun({state, actions, opts: {api}}) {
	actions.setOpInProgress("Toggling dry run")

	api.toggleDryRun()
		.then(json => {
			actions.dryRunToggled(json)
		})	

	return state			
}

function dryRunToggled({state}, config) {
	return {
		...state,
		config,
		opInProgress: null
	}	
}
