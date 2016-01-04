import path from 'path'

module.exports = [
	{type: "treeItemClicked", fn: _treeItemClicked},
	{type: "getTree", fn: _getTree},
	{type: "gotTree", fn: _gotTree},
]

function _treeItemClicked({state, actions, dispatch}, {api, _}, item) {
	let newState = Object.assign({}, state)

	let open = newState.tree.items[item.path]

	// if (item.type !== 'folder') dispatch(C.TREE_FILE_SELECTED, item)

	if (open) {
		// dispatch(C.CLOSE_TREE_ITEM, item)
		delete newState.tree.items[item.path]
		Object.keys(newState.tree.items).forEach( p => {
			if (path.join(p, '/').indexOf(path.join(item.path, '/')) === 0) delete newState.tree.items[p]
		})

	} else {
		dispatch(actions.getTree, item.path)
		// api.getTree(item.path)
		// 	.then(json => {
		// 		dispatch(actions.gotTree, json)
		// 	})
	}

	newState.tree.selected = item.path
	newState.tree.fetching = !open

	return newState
	// return {
	// 	...state,
	// 	tree: {items, selected: item.path, fetching},
	// }
}

function _getTree({state, actions, dispatch}, {api, _}, path) {
	api.getTree(path)
		.then(json => {
			dispatch(actions.gotTree, json)
		})

	return state
}

function _gotTree({state, actions, dispatch}, _, newTree) {
	// console.log('newTree: ', newTree)

	let newState = Object.assign({}, state)

	newState.tree.items[newTree.path] = newTree.nodes
	newState.tree.fetching = false

	return newState
	// return {
	// 	...state,
	// 	tree
	// }
}
