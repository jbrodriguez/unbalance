import { markChosen, getNode } from '../lib/utils'
import * as constant from '../lib/const'

const getEntries = ({ state, actions }) => {
	actions.getGatherTree('/mnt/user')

	return {
		...state,
		gather: {
			...state.gather,
			cache: null,
			items: [{ label: 'Loading ...' }],
			chosen: {},
			lines: [],
			location: null,
		},
	}
}

const clearGatherConsole = ({ state }) => {
	return {
		...state,
		gather: {
			...state.gather,
			lines: [],
		},
	}
}

const getGatherTree = ({ state, actions, opts: { api } }, path) => {
	actions.setBusy(true)

	api.getTree(path).then(json => {
		actions.setBusy(false)
		actions.gotGatherTree(json)
	})

	return state
}

const gotGatherTree = ({ state }, newTree) => {
	let items = [].concat(state.gather.items)

	if (state.gather.cache) {
		const node = state.gather.cache
		node.children = newTree.nodes

		// console.log(`node-${JSON.stringify(state.cache)}`)
		// console.log(`gotTree-${JSON.stringify(newTree.nodes)}`)
	} else {
		items = [].concat(newTree.nodes)
	}

	// items = Utils.getNewTreeState(lineage, items, "collapsed")

	return {
		...state,
		gather: {
			...state.gather,
			items,
		},
	}
}

const checkTarget = ({ state }, path) => {
	const vdisks = Object.keys(state.gather.plan.vdisks).reduce((map, id) => {
		const vdisk = state.gather.plan.vdisks[id]
		map[id] = {
			...vdisk,
			src: false,
			dst: id === path,
		}
		return map
	}, {})

	return {
		...state,
		gather: {
			...state.gather,
			plan: {
				...state.gather.plan,
				vdisks,
			},
		},
	}
}

const gatherTreeCollapsed = ({ state, actions }, lineage) => {
	const tree = [].concat(state.gather.items)
	const node = getNode(tree, lineage)
	// console.log(`node-${JSON.stringify(node)}`)

	// console.log(`utils-${JSON.stringify(Utils)}`)
	// return state

	node.collapsed = !node.collapsed

	const notRetrieved = node.children && node.children.length === 1 && node.children[0].label === 'Loading ...'
	notRetrieved && actions.getGatherTree(node.path)

	return {
		...state,
		gather: {
			...state.gather,
			cache: node,
			items: tree,
		},
	}
}

const gatherTreeChecked = ({ state, actions }, lineage) => {
	const items = [].concat(state.gather.items)
	const chosen = Object.assign({}, state.gather.chosen)

	markChosen(items, lineage, chosen)

	actions.gatherTreeLocate(Object.keys(chosen))

	return {
		...state,
		gather: {
			...state.gather,
			chosen,
			items,
			location: null,
		},
	}
}

const gatherTreeLocate = ({ state, actions, opts: { api } }, chosen) => {
	api.locate(chosen).then(json => actions.gatherTreeLocated(json))

	return state
}

const gatherTreeLocated = ({ state }, location) => {
	return {
		...state,
		gather: {
			...state.gather,
			location,
		},
	}
}

const gatherPlanStarted = ({ state }, line) => {
	return {
		...state,
		gather: {
			...state.gather,
			lines: [`PLANNING: ${line}`],
		},
	}
}

const gatherPlanProgress = ({ state }, line) => {
	const lines = state.gather.lines.length > 1000 ? [] : state.gather.lines

	return {
		...state,
		gather: {
			...state.gather,
			lines: lines.concat(`PLANNING: ${line}`),
		},
	}
}

const gatherPlanFinished = ({ state, actions }, plan) => {
	if (plan.bytesToTransfer === 0) {
		const feedback = []

		feedback.push('The planning stage found that no folders/files can be transferred.')
		feedback.push('')
		feedback.push('This might be due to one of the following reasons:')
		feedback.push(
			'- The source share(s)/folder(s) you selected are either empty or do not exist in the source disk',
		)
		feedback.push(
			"- There isn't available space in any of the target disks, to transfer the share(s)/folder(s) you selected",
		)
		feedback.push('')
		feedback.push(
			'Check more disks in the TO column or go to the Settings page, to review the share(s)/folder(s) selected for moving/copying or to change the amount of reserved space.',
		)

		actions.addFeedback(feedback)
	}

	actions.setBusy(false)

	return {
		...state,
		gather: {
			...state.gather,
			plan,
		},
	}
}

const gatherPlanIssue = ({ state, actions }, permStats) => {
	const permIssues = permStats.split('|')

	const feedback = []

	feedback.push('There are some permission issues with the folders/files you want to transfer')
	feedback.push(`${permIssues[0]} file(s)/folder(s) with an owner other than 'nobody'`)
	feedback.push(`${permIssues[1]} file(s)/folder(s) with a group other than 'users'`)
	feedback.push(`${permIssues[2]} folder(s) with a permission other than 'drwxrwxrwx'`)
	feedback.push(`${permIssues[3]} files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'`)
	feedback.push('You can find more details about which files have issues in the log file (/boot/logs/unbalance.log)')
	feedback.push('')
	feedback.push(
		'At this point, you can transfer the folders/files if you want, but be advised that it can cause errors in the operation',
	)
	feedback.push('')
	feedback.push(
		'You are STRONGLY suggested to install the Fix Common Problems plugin, then run the Docker Safe New Permissions command',
	)

	actions.addFeedback(feedback)

	return state
}

const gatherPlan = ({ state, actions, opts: { ws } }) => {
	actions.setBusy(true)

	const chosenFolders = Object.keys(state.gather.chosen).map(folder => folder.slice(10)) // remove /mnt/user/

	const plan = {
		...state.gather.plan,
		chosenFolders,
	}

	ws.send({ topic: constant.API_GATHER_PLAN, payload: plan })

	return {
		...state,
		gather: {
			...state.gather,
			plan,
		},
	}
}

const gatherMove = ({ state, actions, opts: { ws } }) => {
	actions.setBusy(true)

	actions.resetState()

	ws.send({ topic: constant.API_GATHER_MOVE, payload: state.gather.plan })

	state.history.replace({ pathname: '/transfer' })

	return state
}

export default {
	getEntries,
	clearGatherConsole,

	getGatherTree,
	gotGatherTree,

	checkTarget,

	gatherTreeCollapsed,
	gatherTreeChecked,

	gatherTreeLocate,
	gatherTreeLocated,

	gatherPlanStarted,
	gatherPlanProgress,
	gatherPlanFinished,
	gatherPlanIssue,

	gatherPlan,
	gatherMove,
}
