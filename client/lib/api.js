import fetch from 'node-fetch'

export default class Api {
	constructor() {
		this.hostr = 'http://' + document.location.host + '/api/v1'
	}

	getConfig() {
		return fetch(this.hostr + '/config')
			.then(resp => resp.json())
	}

	addFolder(folder) {
		return fetch(this.hostr + '/config/folder', {
			method: 'PUT',
			body: JSON.stringify({topic: "", payload: folder})
		})
		.then(resp => resp.json())
	}
}