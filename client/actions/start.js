module.exports = [
	{type: "start", fn: _start},
]

function _start({state, actions, dispatch}, {api, ws}) {
	ws.receive(event => {
		let data = JSON.parse(event.data)
		dispatch(data.topic, data.payload)
	})

	dispatch(actions.getConfig)

	dispatch(actions.getTree, "/")

	dispatch(actions.getStorage)

	return state
}
