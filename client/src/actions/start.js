module.exports = {
	start,
}

function start({state, actions, opts: {ws}}) {
	ws.receive(event => {
		let data = JSON.parse(event.data)
		actions[data.topic](data.payload)
	})

	actions.getConfig()
	actions.getTree("/")
	actions.getStorage()

	return state
}
