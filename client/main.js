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

api.getConfig().then( boot )

function boot(config) {
	console.log('config: ', config)

	let initialState = {
		config,
		unraid: null,
		fromDisk: null,
		toDisk: null,
		opInProgress: null,
		moveDisabled: true,
		lines: [],
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
