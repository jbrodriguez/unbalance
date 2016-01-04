import React, { Component } from 'react'
import path from 'path'
import 'font-awesome-webpack'

import styles from '../styles/core.scss'
import classNames from 'classnames/bind'

let cx = classNames.bind(styles)

export default function TreePanel({tree, actions, dispatch}) {
	// console.log('treepanel.tree: ', tree)
	return (
		<div className={cx('treeView')}>
			{ renderTree('/', tree, actions, dispatch) }
		</div>
	)
}

function renderTree(entry, {items, fetching}, actions, dispatch) {
	// console.log('entry: ', entry)
	// let { items, selected, fetching } = state.tree
	let entries = items[entry]

	if (!entries) return (<ul></ul>)

	let list = entries.map( item => {
		// console.log('item: ', item)
		let name = path.basename(item.path)
		let open = items[item.path]
		let isFolder = item.type === 'folder'

		// let itemClass = cx({
		// 	'master': true,
		// 	'entry': true,
		// 	'file': !isFolder,
		// 	'folder': isFolder,
		// 	'closed': isFolder && !open,
		// 	'open': isFolder && open,
		// 	'selected': item.path === selected,
		// })

		let chevron = cx({
			'chevron': true,
			'fa': true,
			'fa-chevron-right': isFolder && !open,
			'fa-chevron-down': isFolder && open,
		})

		// console.log('about to render fetch: ', fetching)


		let spinner = null
		let spacer = (
			<span>&nbsp; &nbsp;</span>
		)

		if (fetching) {
			spinner = (
				<i className="fa fa-circle-o-notch fa-spin" />
			)
		} else {
			spinner = (
				<span>&nbsp; &nbsp;</span>
			)
		}

		return (
			<li key={item.path} className={cx()}>
				<div className={cx('flex', 'listItem')}>
					<div className={cx('floating')}>
						<button className={cx('btn', 'btn-alert')} onClick={onAdd.bind(null, actions, dispatch, item)}>add</button>
						{ spacer }
						{ spinner }
					</div>
					<div onClick={onClick.bind(null, actions, dispatch, item)}>
						&nbsp; &nbsp;
						<i className={chevron} />
						&nbsp; &nbsp;
						<i className="fa fa-folder" />
						&nbsp; &nbsp;
						<span className={cx('name')} >{name}</span>
					</div>
				</div>
				{ renderTree(item.path, {items, fetching}, actions, dispatch) }
			</li>
		)

	})

	return (
		<ul> 
			{ list }
		</ul>
	)
}

// function onClick({actions, dispatch}, item, e) {
function onClick(actions, dispatch, item, e) {
	// console.log('args: ', ...args)

	// console.log('alpha: ', alpha)
	// console.log('actions: ', actions)
	// console.log('dispatch: ', dispatch)
	// console.log('onClick.item: ', item)
	// console.log('e: ', e)
	dispatch(actions.treeItemClicked, item)
	// let { tree, onFolder, onFile, onClose }
	// let open = tree[item.path]

	// this.props.dispatch(C.SELECT_ITEM, item.path)

	// if (item.type !== 'folder') this.props.dispatch(C.FILE_SELECTED, item)

	// if (open) {
	// 	this.props.dispatch(C.CLOSE_ITEM, item)
	// } else {
	// 	this.props.dispatch(C.FOLDER_SELECTED, item)
	// }
}

function onAdd(actions, dispatch, item, e) {
	console.log('onAdd.item: ', item)
	e.preventDefault()

	dispatch(actions.addFolder, item.path)
	// this.props.dispatch(C.ADD_FOLDER, item.path)
}

// export default class TreeView extends Component {
// 	render() {
// 		return (
// 			<div className={cx('treeView')}>
// 				{ this._renderTree('/') }
// 			</div>
// 		)
// 	}

// 				// <li key={item.path} className={itemClass}>
// 				// 	<div className={cx('listItem')}>
// 				// 		<button className={cx('btn', 'btn-alert', 'detail')} onClick={this._onAdd.bind(this, item)}>add</button>
// 				// 		<span className={cx('name')} onClick={this._onClick.bind(this, item)}>{name}</span>
// 				// 	</div>
// 				// 	{ this._renderTree(item.path) }
// 				// </li>

// 					// <Icon name={chevron} />


// 	_renderTree(entry) {
// 		// console.log('entry: ', entry)
// 		let { items, selected, fetching } = this.props
// 		let entries = items[entry]

// 		if (!entries) return (<ul></ul>)

// 		let list = entries.map( item => {
// 			// console.log('item: ', item)
// 			let name = path.basename(item.path)
// 			let open = items[item.path]
// 			let isFolder = item.type === 'folder'

// 			// let itemClass = cx({
// 			// 	'master': true,
// 			// 	'entry': true,
// 			// 	'file': !isFolder,
// 			// 	'folder': isFolder,
// 			// 	'closed': isFolder && !open,
// 			// 	'open': isFolder && open,
// 			// 	'selected': item.path === selected,
// 			// })

// 			let chevron = cx({
// 				'chevron': true,
// 				'fa': true,
// 				'fa-chevron-right': isFolder && !open,
// 				'fa-chevron-down': isFolder && open,
// 			})

// 			// console.log('about to render fetch: ', fetching)


// 			let spinner = null
// 			let spacer = (
// 				<span>&nbsp; &nbsp;</span>
// 			)

// 			if (fetching) {
// 				spinner = (
// 					<i className="fa fa-circle-o-notch fa-spin" />
// 				)
// 			} else {
// 				spinner = (
// 					<span>&nbsp; &nbsp;</span>
// 				)
// 			}

// 			return (
// 				<li key={item.path} className={cx()}>
// 					<div className={cx('flex', 'listItem')}>
// 						<div className={cx('floating')}>
// 							<button className={cx('btn', 'btn-alert', )} onClick={this._onAdd.bind(this, item)}>add</button>
// 							{ spacer }
// 							{ spinner }
// 						</div>
// 						<div  onClick={this._onClick.bind(this, item)}>
// 							&nbsp; &nbsp;
// 							<i className={chevron} />
// 							&nbsp; &nbsp;
// 							<i className="fa fa-folder" />
// 							&nbsp; &nbsp;
// 							<span className={cx('name')} >{name}</span>
// 						</div>
// 					</div>
// 					{ this._renderTree(item.path) }
// 				</li>
// 			)

// 		})

// 		return (
// 			<ul> 
// 				{ list }
// 			</ul>
// 		)
// 	}

// 	_onClick(item, e) {
// 		console.log('_onClick.item: ', item)
// 		this.props.dispatch(C.TREE_ITEM_CLICKED, item)
// 		// let { tree, onFolder, onFile, onClose }
// 		// let open = tree[item.path]

// 		// this.props.dispatch(C.SELECT_ITEM, item.path)

// 		// if (item.type !== 'folder') this.props.dispatch(C.FILE_SELECTED, item)

// 		// if (open) {
// 		// 	this.props.dispatch(C.CLOSE_ITEM, item)
// 		// } else {
// 		// 	this.props.dispatch(C.FOLDER_SELECTED, item)
// 		// }
// 	}

// 	_onAdd(item, e) {
// 		e.preventDefault()
// 		console.log('_onAdd.item: ', item)

// 		this.props.dispatch(C.ADD_FOLDER, item.path)
// 	}
// }