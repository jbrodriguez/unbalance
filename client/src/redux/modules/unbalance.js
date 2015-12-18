import { createAction, handleActions } from 'redux-actions'

// ------------------------------------
// Constants
// ------------------------------------
export const RECEIVE_CONFIG = 'RECEIVE_CONFIG'

export const OP_IN_PROGRESS = 'OP_IN_PROGRESS'
export const CALCULATE = 'CALCULATE'
export const MOVE = 'MOVE'

export const RECEIVE_CALCULATION = 'RECEIVE_CALCULATION'
export const RECEIVE_MOVE = 'RECEIVE_MOVE'



// ------------------------------------
// Actions
// ------------------------------------
export const actions = {
	getConfig,
	calculate,
	move,
	toggleDryRun,
	addFolder,
}

export const getConfig = () = {
	return (dispatch, getState) => {
		return fetch(API + `/config`)
			.then(response => response.json())
			.then(json => dispatch(receiveConfig(json)))		
	}
}

export const calculate = () => {
	return (dispatch, getState) => {
		dispatch(opInProgress(CALCULATE))

		return fetch(API + `/calculate`)
			.then(response => response.json())
			.then(json => dispatch(receiveCalculation(json)))		
	}
}

export const move = () => {
	return (dispatch, getState) => {
		dispatch(opInProgress(MOVE))

		return fetch(API + `/move`)
			.then(response => response.json())
			.then(json => dispatch(receiveMove(json)))		
	}
}

export const move = () => {
	return (dispatch, getState) => {
		dispatch(opInProgress(TOGGLE_DRY_RUN))

		return fetch(API + `/config/dryrun`)
			.then(response => response.json())
			.then(json => dispatch(receiveDryRun(json)))		
	}
}



export const receiveConfig = createAction(RECEIVE_CONFIG)

export const opInProgress = createAction(OP_IN_PROGRESS)

export const receiveCalculation = createAction(RECEIVE_CALCULATION)
export const receiveMove = createAction(RECEIVE_MOVE)
export const receiveDryRun = createAction(RECEIVE_DRY_RUN)



// ------------------------------------
// Reducer
// ------------------------------------

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

export default handleActions({
	RECEIVE_CONFIG: doReceiveConfig,
	OP_IN_PROGRESS: doOpInProgress,
	RECEIVE_CALCULATION: doReceiveCalculation,
	RECEIVE_MOVE: doReceiveMove,
}, {})

const doReceiveConfig = (state, action) => ({
	...state,
	config: action.payload,
})

const doOpInProgress = (state, action) => ({
	...state,
	opInProgress: action.payload
})

const doReceiveCalculation = (state, action) => ({
	...state,
	unraid: action.payload,
	opInProgress: null,
})

const doReceiveMove = (state, action) => ({
	...state,
	unraid: action.payload,
	opInProgress: null,
})