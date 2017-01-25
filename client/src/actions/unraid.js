module.exports = {
	getStorage,
	gotStorage,

	checkFrom,
	checkTo,

	calculate,
	calcStarted,
	calcProgress,
	calcFinished,
	calcPermIssue,

	move,
	moveStarted,
	moveProgress,
	moveFinished,

	opError,
	progressStats,

	getLog,
	gotLog,
}

function getStorage({ state, actions, opts: { api } }) {
	actions.setOpInProgress('Getting storage info')

	api.getStorage()
		.then(json => actions.gotStorage(json))
	// here i can catch the error and show an appropriate message

	return state
}

function gotStorage({ state, actions }, unraid) {
	const toDisk = {}
	const fromDisk = {}
	let sourceDisk = null

	unraid.disks.forEach((disk) => {
		fromDisk[disk.path] = disk.src
		toDisk[disk.path] = disk.dst
		if (disk.src) {
			sourceDisk = disk
		}
	})

	const lines = []
	let opState = null
	switch (unraid.opState) {
	case 1:
		opState = 'Calculate operation in progress ...'
		break
	case 2:
		opState = 'Move operation in progress ...'
		break
	default:
		break
	}

	const tree = Object.assign({}, state.tree)
	if (opState) {
		lines.push(opState)
	} else {
		// console.log(`sourceDisk-${JSON.stringify(sourceDisk)}`)
		actions.getTree(sourceDisk.path)

		tree.cache = null
		tree.items = [{ label: 'Loading ...' }]
	}

	return {
		...state,
		unraid,
		fromDisk,
		toDisk,
		opInProgress: opState,
		stats: unraid.stats,
		moveDisabled: true,
		lines,
		tree,
	}
}

function checkFrom({ state, actions }, path) {
	const fromDisk = Object.assign({}, state.fromDisk)
	const toDisk = Object.assign({}, state.toDisk)

	Object.keys(fromDisk).forEach(key => (fromDisk[key] = key === path))
	Object.keys(toDisk).forEach(key => (toDisk[key] = !(key === path)))

	actions.changeDisk(path)

	return {
		...state,
		fromDisk,
		toDisk,
	}
}

function checkTo({ state }, path) {
	return {
		...state,
		toDisk: {
			...state.toDisk,
			[path]: !state.toDisk[path],
		},
	}
}

function calculate({ state, actions, opts: { ws } }) {
	actions.setOpInProgress('Calculating')

	let srcDisk = ''
	Object.keys(state.fromDisk).some((key) => {
		const isSource = state.fromDisk[key]
		isSource && (srcDisk = key)
		return isSource
	})

	const folders = Object.keys(state.tree.chosen).map(folder => folder.slice(srcDisk.length + 1))

	ws.send({ topic: 'calculate', payload: { srcDisk, folders, dstDisks: state.toDisk } })

	return state
}

function calcStarted({ state }, line) {
	return {
		...state,
		lines: [].concat(`CALCULATE: ${line}`),
		unraid: {
			...state.unraid,
			disks: state.unraid.disks.map((disk) => {
				disk.newFree = disk.free
				return disk
			}),
		},
	}
}

function calcProgress({ state }, line) {
	const lines = state.lines.length > 1000 ? [] : state.lines

	return {
		...state,
		opInProgress: 'calculate',
		moveDisabled: true,
		lines: lines.concat(`CALCULATE: ${line}`),
	}
}

function calcFinished({ state, actions} , unraid) {
	const feedback = []
	if (unraid.bytesToMove === 0) {
		feedback.push('The calculate operation found that no folders/files can be moved.')
		feedback.push('')
		feedback.push('This might be due to one of the following reasons:')
		feedback.push('- The source share(s)/folder(s) you selected are either empty or do not exist in the source disk')
		feedback.push('- There isn\'t available space in any of the target disks, to move the share(s)/folder(s) you selected')
		feedback.push('')
		feedback.push('Check more disks in the TO column or go to the Settings page, to review the share(s)/folder(s) selected for moving or to change the amount of reserved space.')
	}

	if (state.timeout) {
		window.clearTimeout(state.timeout)
	}
	const timeout = window.setTimeout(() => actions.removeFeedback(), 15 * 1000)

	return {
		...state,
		unraid,
		feedback,
		timeout,
		opInProgress: null,
		moveDisabled: unraid.bytesToMove === 0,
	}
}

function calcPermIssue({ state, actions }, permStats) {
	const permIssues = permStats.split('|')

	const feedback = []

	feedback.push('There are some permission issues with the folders/files you want to move')
	feedback.push(`${permIssues[0]} file(s)/folder(s) with an owner other than 'nobody'`)
	feedback.push(`${permIssues[1]} file(s)/folder(s) with a group other than 'users'`)
	feedback.push(`${permIssues[2]} folder(s) with a permission other than 'drwxrwxrwx'`)
	feedback.push(`${permIssues[3]} files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'`)
	feedback.push('You can find more details about which files have issues in the log file (/boot/logs/unbalance.log)')
	feedback.push('')
	feedback.push('At this point, you can move the folders/files if you want, but be advised that it can cause errors in the operation')
	feedback.push('')
	feedback.push('You are STRONGLY suggested to install the Fix Common Problems plugin, then run the Docker Safe New Permissions command')

	if (state.timeout) {
		window.clearTimeout(state.timeout)
	}
	const timeout = window.setTimeout(() => actions.removeFeedback(), 5 * 60 * 1000) // 60s timeout

	return {
		...state,
		feedback,
		timeout,
		opInProgress: null,
		moveDisabled: false,
	}
}

function move({ state, actions, opts: { ws } }) {
	actions.setOpInProgress('Moving')

	ws.send({ topic: 'move' })

	return state
}

function moveStarted({ state }, line) {
	return {
		...state,
		lines: [].concat(`MOVE: ${line}`),
	}
}

function moveProgress({ state }, line) {
	const lines = state.lines.length > 1000 ? [] : state.lines

	return {
		...state,
		lines: lines.concat(`MOVE: ${line}`),
	}
}

function moveFinished({ state }, unraid) {
	return {
		...state,
		unraid,
		opInProgress: null,
		stats: '',
		moveDisabled: !state.config.dryRun,
	}
}

function opError({ state, actions }, error) {
	actions.addFeedback(error)
	return state
}

function progressStats({ state, actions }, stats) {
	return {
		...state,
		stats,
	}
}

function getLog({ state, actions, opts: { ws } }) {
	actions.setOpInProgress('Getting logs ...')

	ws.send({ topic: 'getLog' })

	return state
}

function gotLog({ state }, log) {
	return {
		...state,
		opInProgress: null,
		log,
	}
}
