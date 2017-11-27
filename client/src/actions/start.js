const start = ({ state, actions, opts: { ws } }) => {
	ws.receive(event => {
		const data = JSON.parse(event.data)
		// console.log(`topic(${data.topic})`)
		actions[data.topic](data.payload)
	})

	actions.getConfig()
	actions.getState()
	// actions.getStatus()
	// actions.checkForUpdate()

	return state
}

export default {
	start,
}
