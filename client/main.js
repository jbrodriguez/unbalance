import React from 'react'
import { render } from 'react-dom'
import { Router, Route, IndexRoute } from 'react-router'

import Store from './store'
import Provider from './lib/provider'
import Dispatcher from './lib/dispatcher'

import App from './module/app'
import Home from './module/home'
import Settings from './module/settings'

let initialState = {
	unraid: null,
	config: null,
	opInProgress: null,
	consoleLines: [],
}

var store = new Store(initialState)

function requireConfig(nextState, replaceState, callback) {
	console.log('mofo')
	store.state.onValue(state => {
		console.log('require.state.config: ', state.config)
		// console.log('require.config: ', config)
		if (!state.config) {
			replaceState({ nextPathname: nextState.location.pathname }, '/settings')
		}

		callback()
	})
}


store.status.onValue(
	model => {
		// since wrapping the router in a Provider doesn't work right now
		// I think I'm missing something
		function createElement(Component, props) {
			return <Component {...props} model={model} dispatch={Dispatcher.dispatch} />
		}

		render(
			<Router createElement={createElement}>
				<Route path='/' component={App}>
					<IndexRoute component={Home} onEnter={requireConfig} />
					<Route path='settings' component={Settings} />
				</Route>
			</Router>,
			document.getElementById('mnt')
		)
	}
)