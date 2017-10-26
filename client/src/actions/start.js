module.exports = {
	start,
}

function start({ state, actions, opts: { ws } }) {
	ws.receive(event => {
		const data = JSON.parse(event.data)
		actions[data.topic](data.payload)
	})

	actions.getConfig()
	actions.getState()
	// actions.checkForUpdate()

	return state
}
