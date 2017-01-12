import path from 'path'
import { Utils } from 'react-tree-menu'

module.exports = {
	treeItemClicked,
	getTree,
	gotTree,

	treeCollapsed,
	treeChecked,
}

function treeItemClicked({state, actions, opts: {api}}, item) {
	let items = Object.assign({}, state.tree.items)

	const open = items[item.path]
	if (open) {
		// dispatch(C.CLOSE_TREE_ITEM, item)
		delete items[item.path]
		Object.keys(items).forEach( p => {
			if (path.join(p, '/').indexOf(path.join(item.path, '/')) === 0) delete items[p]
		})

	} else {
		actions.getTree(item.path)
		// api.getTree(item.path)
		// 	.then(json => {
		// 		dispatch(actions.gotTree, json)
		// 	})
	}

	return {
		...state,
		tree: {
			items,
			selected: item.path,
			fetching: !open,
		},
	}

	// tree.selected = item.path
	// tree.fetching = !open
	// return newState

	// return {
	// 	...state,
	// 	tree: {items, selected: item.path, fetching},
	// }
}

function getTree({state, actions,  opts: {api}}, path) {
	api.getTree(path)
		.then(json => {
			actions.gotTree(json)
		})

	return state
}

function gotTree({state}, newTree) {
	let items = [].concat(state.tree.items)

	if (state.tree.cache) {
		const node = state.tree.cache
		node.children = newTree.nodes

		// console.log(`node-${JSON.stringify(state.tree.cache)}`)
		// console.log(`gotTree-${JSON.stringify(newTree.nodes)}`)
	} else {
		items = [].concat(newTree.nodes)
	}

	// items = Utils.getNewTreeState(lineage, items, "collapsed")

	return {
		...state,
		tree: {
			source: newTree.path,
			items,
		},
	}
}

function treeCollapsed({state, actions}, lineage) {
	let tree = [].concat(state.tree.items)
	const node = getNode(tree, lineage)
	// console.log(`node-${JSON.stringify(node)}`)

	// console.log(`utils-${JSON.stringify(Utils)}`)
	// return state

	node.collapsed = !node.collapsed

	actions.getTree(node.path)

	return {
		...state,
		tree: {
			...state.tree,
			cache: node,
			items: tree
		}
	}
}

function treeChecked({state, actions}, lineage) {
	let items = [].concat(state.tree.items)
	const node = getNode(items, lineage)

	node.checked = !node.checked
	// node.children.forEach( child => {
	// 	child.checked = node.checked
	// })

	// items = Utils.getNewTreeState(lineage, items, "checked")

	return {
		...state,
		tree: {
			...state.tree,
			items
		}
	}

}

const getNode = (tree, lineage) => {
	// console.log(`lineage-${JSON.stringify(lineage)}`)
	if (lineage.length === 0) {
		return null
	} else if (lineage.length === 1) {
		// console.log(`uno-${JSON.stringify(tree[lineage[0]])}`)
		return tree[lineage[0]]
	} else {
		const node = lineage.shift()
		return getNode(tree[node].children, lineage)
	}
}
