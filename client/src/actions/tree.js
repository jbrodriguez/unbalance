import path from 'path'

module.exports = {
	treeItemClicked,
	getTree,
	gotTree,
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
	return {
		...state,
		tree: {
			...state.tree,
			fetching: false,
			items: {
				...state.tree.items,
				[newTree.path]: newTree.nodes
			},
		},
	}
}
