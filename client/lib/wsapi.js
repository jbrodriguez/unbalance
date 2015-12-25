import Bacon from 'baconjs'

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

	get stream() {
		return Bacon.fromEventTarget(this.skt, "message").map(function(event) {
		    return JSON.parse(event.data)
		})
	}

	send({topic, payload}) {
		console.log('topic: ', topic)
		const packet = {
			topic,
			payload: JSON.stringify(payload)
		}

		console.log('packet: ', packet)

		this.skt.send(JSON.stringify(packet))			
	} 
}