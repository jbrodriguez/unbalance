import fetch from 'isomorphic-fetch'

export default class Api {
	constructor() {
		this.hostr = `${document.location.protocol}//${document.location.host}/api/v1`
	}

	getConfig() {
		return fetch(`${this.hostr}/config`).then(resp => resp.json())
	}

	getState() {
		return fetch(`${this.hostr}/state`).then(resp => resp.json())
	}

	getStorage() {
		return fetch(`${this.hostr}/storage`).then(resp => resp.json())
	}

	getOperation() {
		return fetch(`${this.hostr}/operation`).then(resp => resp.json())
	}

	getHistory() {
		return fetch(`${this.hostr}/history`).then(resp => resp.json())
	}

	resetOperation() {
		return fetch(`${this.hostr}/resetOp`).then(resp => resp.json())
	}

	setNotifyPlan(notify) {
		return fetch(`${this.hostr}/config/notifyPlan`, {
			method: 'PUT',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ payload: notify }),
		}).then(resp => resp.json())
	}

	setNotifyTransfer(notify) {
		return fetch(`${this.hostr}/config/notifyTransfer`, {
			method: 'PUT',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ payload: notify }),
		}).then(resp => resp.json())
	}

	setReservedSpace(amount, unit) {
		return fetch(`${this.hostr}/config/reservedSpace`, {
			method: 'PUT',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ payload: JSON.stringify({ amount, unit }) }),
		}).then(resp => resp.json())
	}

	getTree(path) {
		return fetch(`${this.hostr}/tree`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ payload: path }),
		}).then(resp => resp.json())
	}

	toggleDryRun(dryRun) {
		return fetch(`${this.hostr}/config/dryRun`, {
			method: 'PUT',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ payload: dryRun }),
		}).then(resp => resp.json())
	}

	setRsyncArgs(rsyncArgs) {
		return fetch(`${this.hostr}/config/rsyncArgs`, {
			method: 'PUT',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ payload: JSON.stringify({ rsyncArgs }) }),
		}).then(resp => resp.json())
	}

	setVerbosity(verbosity) {
		return fetch(`${this.hostr}/config/verbosity`, {
			method: 'PUT',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ payload: verbosity }),
		}).then(resp => resp.json())
	}

	setUpdateCheck(checkUpdate) {
		return fetch(`${this.hostr}/config/checkUpdate`, {
			method: 'PUT',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ payload: checkUpdate }),
		}).then(resp => resp.json())
	}

	locate(path) {
		return fetch(`${this.hostr}/locate`, {
			method: 'POST',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ payload: path }),
		}).then(resp => resp.json())
	}

	checkForUpdate() {
		return fetch(`${this.hostr}/update`).then(resp => resp.json())
	}

	setRefreshRate(refreshRate) {
		return fetch(`${this.hostr}/config/refreshRate`, {
			method: 'PUT',
			headers: { 'content-type': 'application/json' },
			body: JSON.stringify({ payload: refreshRate }),
		}).then(resp => resp.json())
	}
}
