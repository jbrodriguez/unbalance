import { markChosen, getNode } from '../lib/utils'

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

		// console.log(`node-${JSON.stringify(state.tree.cache)}`)
		// console.log(`gotTree-${JSON.stringify(newTree.nodes)}`)
	} else {
		items = [].concat(newTree.nodes)
	}

	// items = Utils.getNewTreeState(lineage, items, "collapsed")

	return {
		...state,
		scatter: {
			...state.scatter,
			items,
		},
	}
}

const checkFrom = ({ state, actions }, path) => {
	actions.scatterGetTree(path)

	return {
		...state,
		scatter: {
			...state.scatter,
			cache: null,
			chosen: {},
			items: [{ label: 'Loading ...' }],
		},
	}
}

const scatterTreeCollapsed = ({ state, actions }, lineage) => {
	const tree = [].concat(state.scatter.items)
	const node = getNode(tree, lineage)
	// console.log(`node-${JSON.stringify(node)}`)

	// console.log(`utils-${JSON.stringify(Utils)}`)
	// return state

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

	console.log(`chosenActions(${JSON.stringify(chosen)})`)

	return {
		...state,
		scatter: {
			...state.scatter,
			chosen,
			items,
		},
	}
}

export default {
	scatterGetTree,
	scatterGotTree,

	checkFrom,

	scatterTreeCollapsed,
	scatterTreeChecked,
}
