import B from 'baconjs'

import Dispatcher from './lib/dispatcher'
import * as C from './constant'

import Api from './lib/api'

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
// 		},
// 		disks: [
// 			{id: 1, name: "disk1", path: "/mnt/disk1"},
// 			{id: 2, name: "disk2", path: "/mnt/disk2"},
// 			{id: 3, name: "disk3", path: "/mnt/disk3"},
// 		],
// 		bytesToMove: 0,
// 		inProgress: false, // need to review this variable
// 	}
// 	opInProgress: null,
// }

export default class Store {
	constructor(initialState = {}) {
		const api = new Api()

		this.state = this._setup(initialState, api, Dispatcher.dispatch, Dispatcher.register)
	}

	get status() {
		return this.state
	}

	_setup(initialState, api, dispatch, register) {
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
			B.fromEventTarget(ws, "message").onValue(event => {
				dispatch(C.GOT_WS_MESSAGE, JSON.parse(event.data))
			})
		})

		return B.update(
			initialState,
			getConfig, _getConfig,
			[gotConfig], _gotConfig,
			[addFolder], _addFolder,
			[opInProgress], _opInProgress,
			[folderAdded], _folderAdded,
			[gotWsMessage], _gotWsMessage,
		)

		function _getConfig(state, _) {
			dispatch(C.OP_IN_PROGRESS, C.GET_CONFIG)

			B.fromPromise(api.getConfig()).onValue(json => {
				dispatch(C.GOT_CONFIG, json)
			})

			return state
		}

		function _gotConfig(state, config) {
			console.log('whatisconfig: ', config)
			return {
				...state,
				config: config,
				opInProgress: null,
			}
		}

		function _addFolder(state, _) {
			dispatch(C.OP_IN_PROGRESS, C.ADD_FOLDER)

			B.fromPromise(api.addFolder).onValue(json => {
				dispatch(C.FOLDER_ADDED, json)
			})

			return state
		}

		function _folderAdded(state, config) {
			return {
				...state,
				config: config,
				opInProgress: null,
			}
		}

		function _opInProgress(state, action) {
			return {
				...state,
				opInProgress: action
			}
		}

		function _gotWsMessage(state, message) {
			return {
				...state,
				consoleLines: consoleLines.push(message)
			}
		}		
	}
}


// export default function createStore(initialState = {}, dispatcher) {
// 	let state = { ...initialState, dispatch: dispatcher.dispatch }

// 	return B.combineTemplate


// }