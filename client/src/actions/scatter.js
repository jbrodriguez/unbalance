import { markChosen, getNode } from '../lib/utils'
import * as constant from '../lib/const'

const checkFrom = ({ state, actions }, path) => {
	console.log(`keys(${JSON.stringify(Object.keys(state.scatter.plan.vdisks))})`)

	const vdisks = Object.keys(state.scatter.plan.vdisks).reduce((map, id) => {
		const vdisk = state.scatter.plan.vdisks[id]
		map[id] = {
			...vdisk,
			src: id === path,
			dst: id !== path,
		}
		return map
	}, {})

	// console.log(`vdisks(${JSON.stringify(vdisks)})`)

	actions.scatterGetTree(path)

	return {
		...state,
		scatter: {
			...state.scatter,
			cache: null,
			chosen: {},
			items: [{ label: 'Loading ...' }],
			plan: {
				...state.scatter.plan,
				vdisks,
			},
		},
	}
}

const checkTo = ({ state, actions }, path) => {
	return {
		...state,
		scatter: {
			...state.scatter,
			plan: {
				...state.scatter.plan,
				vdisks: {
					...state.scatter.plan.vdisks,
					[path]: {
						...state.scatter.plan.vdisks[path],
						dst: !state.scatter.plan.vdisks[path],
					},
				},
			},
		},
	}
}

const scatterGetTree = ({ state, actions, opts: { api } }, path) => {
	actions.setBusy(true)

	api.getTree(path).then(json => {
		actions.scatterGotTree(json)
		actions.setBusy(false)
	})

	return state
}

const scatterGotTree = ({ state }, newTree) => {
	let items = [].concat(state.scatter.items)

	if (state.scatter.cache) {
		const node = state.scatter.cache
		node.children = newTree.nodes
	} else {
		items = [].concat(newTree.nodes)
	}

	return {
		...state,
		scatter: {
			...state.scatter,
			items,
		},
	}
}

const scatterTreeCollapsed = ({ state, actions }, lineage) => {
	const tree = [].concat(state.scatter.items)
	const node = getNode(tree, lineage)

	node.collapsed = !node.collapsed

	const notRetrieved = node.children && node.children.length === 1 && node.children[0].label === 'Loading ...'
	notRetrieved && actions.scatterGetTree(node.path)

	return {
		...state,
		scatter: {
			...state.scatter,
			cache: node,
			items: tree,
		},
	}
}

const scatterTreeChecked = ({ state, actions }, lineage) => {
	const items = [].concat(state.scatter.items)
	const chosen = Object.assign({}, state.scatter.chosen)

	markChosen(items, lineage, chosen)

	return {
		...state,
		scatter: {
			...state.scatter,
			chosen,
			items,
		},
	}
}

const scatterPlanStarted = ({ state }, line) => {
	return {
		...state,
		env: {
			...state.env,
			lines: [].concat(`PLANNING: ${line}`),
		},
	}
}

const scatterPlanProgress = ({ state }, line) => {
	const lines = state.env.lines.length > 1000 ? [] : state.env.lines

	return {
		...state,
		env: {
			...state.env,
			lines: lines.concat(`PLANNING: ${line}`),
		},
	}
}

const scatterPlanFinished = ({ state, actions }, plan) => {
	if (plan.bytesToTransfer === 0) {
		const feedback = []

		feedback.push('The planning stage found that no folders/files can be moved/copied.')
		feedback.push('')
		feedback.push('This might be due to one of the following reasons:')
		feedback.push(
			'- The source share(s)/folder(s) you selected are either empty or do not exist in the source disk',
		)
		feedback.push(
			"- There isn't available space in any of the target disks, to move/copy the share(s)/folder(s) you selected",
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
		scatter: {
			...state.scatter,
			plan,
		},
	}
}

const scatterPlanIssue = ({ state, actions }, permStats) => {
	const permIssues = permStats.split('|')

	const feedback = []

	feedback.push('There are some permission issues with the folders/files you want to move')
	feedback.push(`${permIssues[0]} file(s)/folder(s) with an owner other than 'nobody'`)
	feedback.push(`${permIssues[1]} file(s)/folder(s) with a group other than 'users'`)
	feedback.push(`${permIssues[2]} folder(s) with a permission other than 'drwxrwxrwx'`)
	feedback.push(`${permIssues[3]} files(s) with a permission other than '-rw-rw-rw-' or '-r--r--r--'`)
	feedback.push('You can find more details about which files have issues in the log file (/boot/logs/unbalance.log)')
	feedback.push('')
	feedback.push(
		'At this point, you can move the folders/files if you want, but be advised that it can cause errors in the operation',
	)
	feedback.push('')
	feedback.push(
		'You are STRONGLY suggested to install the Fix Common Problems plugin, then run the Docker Safe New Permissions command',
	)

	actions.addFeedback(feedback)

	return state
}

const scatterPlan = ({ state, actions, opts: { ws } }) => {
	actions.setBusy(true)

	const srcDisk = Object.keys(state.scatter.plan.vdisks).find(vdisk => state.scatter.plan.vdisks[vdisk].src)
	const chosenFolders = Object.keys(state.scatter.chosen).map(folder =>
		folder.slice(state.scatter.plan.vdisks[srcDisk].path.length + 1),
	)

	const plan = {
		...state.scatter.plan,
		chosenFolders,
	}

	ws.send({ topic: constant.API_SCATTER_PLAN, payload: plan })

	return {
		...state,
		scatter: {
			...state.scatter,
			plan,
		},
	}
}

const scatterMove = ({ state, actions, opts: { ws } }) => {
	actions.setBusy(true)

	actions.resetState()

	ws.send({ topic: constant.API_SCATTER_MOVE, payload: state.scatter.plan })

	state.history.replace({ pathname: '/transfer' })

	return state
}

const scatterCopy = ({ state, actions, opts: { ws } }) => {
	actions.setBusy(true)

	actions.resetState()

	ws.send({ topic: constant.API_SCATTER_COPY, payload: state.scatter.plan })

	state.history.replace({ pathname: '/transfer' })

	return state
}

export default {
	scatterGetTree,
	scatterGotTree,

	checkFrom,
	checkTo,

	scatterTreeCollapsed,
	scatterTreeChecked,

	scatterPlanStarted,
	scatterPlanProgress,
	scatterPlanFinished,
	scatterPlanIssue,

	scatterPlan,
	scatterMove,
	scatterCopy,
}
