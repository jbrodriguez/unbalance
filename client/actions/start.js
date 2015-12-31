module.exports = [
	{type: "start", fn: _start},
]

function _start({state, actions, dispatch}, {_, ws}) {
	ws.receive(event => {
		let data = JSON.parse(event.data)
		dispatch(data.topic, data.payload)
	})

	return state
}
