export default class WebSocketApi {
	constructor() {
		const hostw = "ws://" + document.location.host + "/skt"

		this.skt = new WebSocket(hostw)

		this.skt.onopen = function() {
		    console.log("Connection opened")
		}

		this.skt.onclose = function() {
		    console.log("Connection is closed...")
		}		
	}

	receive(fn) {
		this.skt.onmessage = fn
	}

	send({topic, payload}) {
		// console.log('topic: ', topic)
		const packet = {
			topic,
			payload: JSON.stringify(payload)
		}

		// console.log('packet: ', packet)

		this.skt.send(JSON.stringify(packet))			
	} 
}