import React, { PureComponent } from 'react'
import { render } from 'react-dom'
import { PropTypes } from 'prop-types'

import { Router, Route } from 'react-router-dom'
import createBrowserHistory from 'history/createBrowserHistory'

import { createStore, combineActions } from 'reactorx'

import start from './actions/start'
import config from './actions/config'
import env from './actions/env'
import core from './actions/core'
import scatter from './actions/scatter'
import gather from './actions/gather'

import Api from './lib/api'
import WSApi from './lib/wsapi'

import App from './components/app'
import Scatter from './components/scatter'
import Gather from './components/gather'
import Transfer from './components/transfer'
import History from './components/history'
import Settings from './components/settings'
import Log from './components/log'

const history = createBrowserHistory()

const initialState = {
	config: null,
	core: null,
	env: {
		isBusy: false,
		log: [],
		feedback: [],
		timeout: null,
		latestVersion: '',
	},
	scatter: {
		cache: null,
		chosen: {},
		items: [],
		plan: null,
		lines: [],
	},
	gather: {
		cache: null,
		chosen: {},
		items: [],
		plan: null,
		lines: [],
		location: null,
	},
	history,
}

const actions = combineActions(start, config, env, core, scatter, gather)

const api = new Api()
const ws = new WSApi()

const appStore = createStore(initialState, actions, { api, ws })

class Layout extends PureComponent {
	static propTypes = {
		store: PropTypes.object.isRequired,
	}

	render() {
		const { store } = this.props

		// we wait for a valid config and a valid status before rendering the content
		if (!(store.state && store.state.config)) {
			return null
		}

		return (
			<Router history={store.state.history}>
				<App store={store}>
					<Route exact path="/" render={props => <Scatter store={store} {...props} />} />
					<Route path="/gather" render={props => <Gather store={store} {...props} />} />
					<Route exact path="/transfer" render={props => <Transfer store={store} {...props} />} />
					<Route exact path="/history" render={props => <History store={store} {...props} />} />
					<Route exact path="/settings" render={props => <Settings store={store} {...props} />} />
					<Route exact path="/log" render={props => <Log store={store} {...props} />} />
				</App>
			</Router>
		)
	}
}

appStore.subscribe(state => render(<Layout store={state} />, document.getElementById('app')))

appStore.actions.start()
