import B from 'baconjs'

import Dispatcher from './lib/dispatcher'
import * as C from './constant'

export default class Store {
	constructor(initialState = {}) {
		this.dispatch = Dispatcher.dispatch
		this.state = this._setup(initialState, Dispatcher.register)
	}

	get status() {
		return this.state
	}

	_setup(initialState, register) {
		const [
			start, 
			getConfig,
			gotConfig,
			addFolder,
			folderAdded,
			opInProgress,
			calculate,
			move,
			toggleDryRun,
			gotWsMessage,
		] = register(
			C.START,
			C.GET_CONFIG,
			C.GOT_CONFIG,
			C.ADD_FOLDER,
			C.FOLDER_ADDED,
			C.OP_IN_PROGRESS,
			C.CALCULATE,
			C.MOVE,
			C.TOGGLE_DRY_RUN,
			C.GOT_WS_MESSAGE,
		)

		start.onValue(value => {
			const ws = new WebSocket(WS_URL)
			Bacon.fromEventTarget(ws, "message").onValue(event => {
				this.dispatch(C.GOT_WS_MESSAGE, JSON.parse(event.data))
			})
		})

		return B.update(
			initialState,
			[getConfig], this._getConfig,
			[gotConfig], this._gotConfig,
			[addFolder], this._addFolder,
			[folderAdded], this._folderAdded,
			[gotWsMessage], this._gotWsMessage,
		)
	}

	_getConfig(state, _) {
		this.dispatch(C.OP_IN_PROGRESS, C.GET_CONFIG)

		B.fromPromise(api.getConfig).onValue(json => {
			this.dispatch(C.GOT_CONFIG, json)
		})

		return state
	}


	_gotConfig(state, config) {
		return {
			...state,
			config: config,
			opInProgress: null,
		}
	}

	_addFolder(state, _) {
		this.dispatch(C.OP_IN_PROGRESS, C.ADD_FOLDER)

		B.fromPromise(api.addFolder).onValue(json => {
			this.dispatch(C.FOLDER_ADDED, json)
		})

		return state
	}

	_folderAdded(state, config) {
		return {
			...state,
			config: config,
			opInProgress: null,
		}
	}

	_gotWsMessage(state, message) {
		return {
			...state,
			consoleLines: consoleLines.push(message)
		}
	}
}


// export default function createStore(initialState = {}, dispatcher) {
// 	let state = { ...initialState, dispatch: dispatcher.dispatch }

// 	return B.combineTemplate


// }