module.exports = {
	getGatherTree,
	gotGatherTree,

	gatherTreeCollapsed,
	gatherTreeChecked,

	gatherTreeLocate,
	gatherTreeLocated,
}

// utilities
const getNode = (tree, lineage) => {
	if (lineage.length === 0) {
		return null
	} else if (lineage.length === 1) {
		return tree[lineage[0]]
	}

	const node = lineage.shift()
	return getNode(tree[node].children, lineage)
}

const markChosen = (tree, lineage, chosen) => {
	if (lineage.length === 0) {
		// no-op
	} else if (lineage.length === 1) {
		const node = tree[lineage[0]]

		if (node.checked) {
			delete chosen[node.path]
		} else {
			uncheckChildren(node.children, chosen)
			chosen[node.path] = true
		}

		node.checked = !node.checked
	} else {
		const index = lineage.shift() // this mutates lineage
		const node = tree[index]

		if (node.checked) {
			delete chosen[node.path]
			node.checked = false
		}

		markChosen(node.children, lineage, chosen)
	}
}

const uncheckChildren = (tree, chosen) => {
	if (!tree) return

	tree.forEach(node => {
		delete chosen[node.path]
		node.checked = false

		uncheckChildren(node.children, chosen)
	})
}

// actions
function getGatherTree({ state, actions, opts: { api } }, path) {
	api.getTree(path).then(json => actions.gotGatherTree(json))

	return state
}

function gotGatherTree({ state }, newTree) {
	let items = [].concat(state.gatherTree.items)

	if (state.gatherTree.cache) {
		const node = state.gatherTree.cache
		node.children = newTree.nodes

		// console.log(`node-${JSON.stringify(state.tree.cache)}`)
		// console.log(`gotTree-${JSON.stringify(newTree.nodes)}`)
	} else {
		items = [].concat(newTree.nodes)
	}

	// items = Utils.getNewTreeState(lineage, items, "collapsed")

	return {
		...state,
		gatherTree: {
			...state.gatherTree,
			items,
		},
	}
}

function gatherTreeCollapsed({ state, actions }, lineage) {
	const tree = [].concat(state.gatherTree.items)
	const node = getNode(tree, lineage)
	// console.log(`node-${JSON.stringify(node)}`)

	// console.log(`utils-${JSON.stringify(Utils)}`)
	// return state

	node.collapsed = !node.collapsed

	const notRetrieved = node.children && node.children.length === 1 && node.children[0].label === 'Loading ...'
	notRetrieved && actions.getGatherTree(node.path)

	return {
		...state,
		gatherTree: {
			...state.gatherTree,
			cache: node,
			items: tree,
		},
	}
}

function gatherTreeChecked({ state, actions }, lineage) {
	const items = [].concat(state.gatherTree.items)
	const chosen = Object.assign({}, state.gatherTree.chosen)
	// const lineage2 = [].concat(lineage)

	markChosen(items, lineage, chosen)

	// trigger search for disks where this share/folder is present
	actions.gatherTreeLocate(Object.keys(chosen))
	// actions.gatherTreeLocate(lineage2)

	return {
		...state,
		gatherTree: {
			...state.gatherTree,
			chosen,
			items,
			present: [],
		},
	}
}

function gatherTreeLocate({ state, actions, opts: { api } }, chosen) {
	// const tree = [].concat(state.gatherTree.items)
	// const node = getNode(tree, lineage)

	// node.path is in the form of /mnt/user/tvshows/Breaking Bad
	// so we remove /mnt/user
	// const path = node.path.slice(10)

	api.locate(chosen).then(json => actions.gatherTreeLocated(json))

	return state
}

function gatherTreeLocated({ state }, disks) {
	console.log(`disks-(${JSON.stringify(disks)})`)
	return {
		...state,
		gatherTree: {
			...state.gatherTree,
			present: disks,
		},
	}
}
