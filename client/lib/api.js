import fetch from 'node-fetch'

export default class Api {
	constructor() {
		this.hostr = 'http://' + document.location.host + '/api/v1'
	}

	getConfig() {
		console.log('this.hostr', this.hostr)
		return fetch(this.hostr + '/config')
			.then(resp => resp.json())
	}

	addFolder(folder) {
		// console.log('sending: ', folder)
		return fetch(this.hostr + '/config/folder', {
			method: 'PUT',
			headers: {'content-type': 'application/json'},
			body: JSON.stringify({payload: folder})
		})
		.then(resp => resp.json())
	}

	getStorage() {
		// console.log('this.hostr', this.hostr)
		return fetch(this.hostr + '/storage')
			.then(resp => resp.json())
	}
}