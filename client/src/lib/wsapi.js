export default class WebSocketApi {
	constructor() {
		const proto = document.location.protocol === 'https:' ? 'wss' : 'ws'
		const hostw = `${proto}://${document.location.host}/skt`

		this.skt = new WebSocket(hostw)

		this.skt.onopen = () => {
			console.log('Connection opened')
		}

		this.skt.onclose = () => {
			console.log('Connection is closed...')
		}
	}

	receive(fn) {
		this.skt.onmessage = fn
	}

	send({ topic, payload }) {
		const packet = {
			topic,
			payload: JSON.stringify(payload),
		}

		this.skt.send(JSON.stringify(packet))
	}
}
