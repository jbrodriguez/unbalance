import React from 'react'
import { render } from 'react-dom'
import { Router, Route, IndexRoute } from 'react-router'

import { createStore, combineActions } from 'reactorx'

import startActions from './actions/start'
import uiActions from './actions/ui'
import configActions from './actions/config'
import treeActions from './actions/tree'
import unraidActions from './actions/unraid'

import Api from './lib/api'
import WSApi from './lib/wsapi'

import App from './components/app'
import Home from './components/home'
import Settings from './components/settings'

// SAMPLE STATE
// state = {
// 	config: {
// 		folders: [
// 			"movies/films",
// 			"movies/tvshows"
// 		],
// 		dryRun: true
// 	}
// 	unraid: {
// 		condition: {
// 			numDisks: 24,
// 			numProtected: 0,
//			state: "STARTED",
// 		},
// 		disks: [
// 			{id: 1, name: "disk1", path: "/mnt/disk1"},
// 			{id: 2, name: "disk2", path: "/mnt/disk2"},
// 			{id: 3, name: "disk3", path: "/mnt/disk3"},
// 		],
// 		bytesToMove: 0,
// 	}
//  fromDisk: null,
//  toDisk: null,
//  opInProgress: null,
//  moveDisabled: true,
//  lines: [],
// 	tree: {
// 		items: 
// 			'/': [
// 				{type: 'folder', path: '/films'}
// 				{type: 'folder', path: '/tvshows'}
// 				{type: 'folder', path: '/storage'}
// 				{type: 'folder', path: '/data'}
// 			],
// 		selected: "",
// 		fetching: false,
// 	},
//	feedback: []
// }

const initialState = {
	config: null,
	unraid: null,
	fromDisk: null,
	toDisk: null,
	opInProgress: null,
	moveDisabled: true,
	lines: [],
	tree: {
		items: [],
		selected: '',
		fetching: false,
	},
	feedback: [],
}

const actions = combineActions(
	startActions, 
	uiActions, 
	configActions, 
	treeActions, 
	unraidActions
)

const api = new Api()
const ws = new WSApi()

const store = createStore(initialState, actions, {api, ws})

store.subscribe(
	store => {
		// console.log('main.store.state: ', store.state)

		function createElement(Component, props) {
			return <Component {...props} store={store} />
		}

		render(
			<Router createElement={createElement}>
				<Route path='/' component={App}>
					<IndexRoute component={Home} />
					<Route path='settings' component={Settings} />
				</Route>
			</Router>,
			document.getElementById('mnt')
		)
	}
)

store.actions.start()

// Promise.all([api.getConfig(), api.getTree('/')])
// 	.then( boot )

// function boot([config, entry]) {
// 	// console.log('config: ', config)
// 	// console.log('entry: ', entry)

// 	let treeItems = {}
// 	treeItems[entry.path] = entry.nodes

// 	let initialState = {
// 		config,
// 		unraid: null,
// 		fromDisk: null,
// 		toDisk: null,
// 		opInProgress: null,
// 		moveDisabled: true,
// 		lines: [],
// 		tree: {
// 			items: treeItems,
// 			selected: '',
// 			fetching: false,
// 		},
// 		feedback: [],
// 	}

// 	let actions = [].concat(startActions, uiActions, configActions, treeActions, unraidActions)

// 	const store = createStore(initialState, actions, {api, ws})

// 	store.subscribe(
// 		store => {
// 			// console.log('main.store: ', store)

// 			function createElement(Component, props) {
// 				return <Component {...props} store={store} />
// 			}

// 			render(
// 				<Router createElement={createElement}>
// 					<Route path='/' component={App}>
// 						<IndexRoute component={Home} />
// 						<Route path='settings' component={Settings} />
// 					</Route>
// 				</Router>,
// 				document.getElementById('mnt')
// 			)
// 		}
// 	)

// 	store.dispatch(store.actions.start)	
// }
