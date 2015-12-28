import React from 'react'
import { render } from 'react-dom'
import { Router, Route, IndexRoute } from 'react-router'

import Store from './store'
import Provider from './lib/provider'
import Dispatcher from './lib/dispatcher'
import Api from './lib/api'

import App from './module/app'
import Home from './module/home'
import Settings from './module/settings'

let api = new Api()

Promise.all([api.getConfig(), api.getTree('/')])
	.then( boot )

// api.getConfig().then( boot )

function boot([config, entry]) {
	console.log('config: ', config)
	console.log('entry: ', entry)

	let items = {}
	items[entry.path] = entry.nodes

	let initialState = {
		config,
		unraid: null,
		fromDisk: null,
		toDisk: null,
		opInProgress: null,
		moveDisabled: true,
		lines: [],
		tree: {
			items,
			selected: '',
			fetching: false,
		}
	}

	var store = new Store(initialState)

	// function requireConfig(nextState, replaceState, callback) {
	// 	// console.log('mofo')
	// 	store.status.onValue(state => {
	// 		// console.log('require.state: ', state)
	// 		console.log('require.config: ', state.config)
	// 		if (!state.config) {
	// 			console.log('do i replace')
	// 			replaceState({ nextPathname: nextState.location.pathname }, '/settings')
	// 		}

	// 		callback()
	// 	})
	// }

						// <IndexRoute component={Home} onEnter={requireConfig} />
						// <IndexRoute component={Home} />


	store.status.onValue(
		model => {
			// console.log('this.model: ', model)
			// since wrapping the router in a Provider doesn't work right now
			// I think I'm missing something
			function createElement(Component, props) {
				return <Component {...props} model={model} dispatch={Dispatcher.dispatch} />
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
}
