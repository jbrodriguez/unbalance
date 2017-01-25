module.exports = {
	humanBytes,
	percentage,
	scramble,
}

function humanBytes(bytes) {
	if (bytes === 0) return '0 Byte'

	const k = 1000
	const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
	const i = Math.floor(Math.log(bytes) / Math.log(k))

	return `${(bytes / Math.pow(k, i)).toPrecision(3)} ${sizes[i]}` // eslint-disable-line
}

function percentage(input, decimals = 2, suffix = '%') {
	return `${Math.round(input * Math.pow(10, decimals + 2)) / Math.pow(10, decimals)}${suffix}` // eslint-disable-line
}

function scramble(serial) {
	if (serial.startsWith('WDC_WD30EZRX')) {
		return `WDC_WD30EZRX-${makeid()}`
	} else if (serial.startsWith('ST3000DM001')) {
		return `ST3000DM001-${makeid()}`
	} else if (serial.startsWith('TOSHIBA_DT01ACA300')) {
		return `TOSHIBA_DT01ACA300-${makeid()}`
	} else if (serial.startsWith('ST4000DM000')) {
		return `ST4000DM000-${makeid()}`
	} else if (serial.startsWith('WDC_WD40EZRX')) {
		return `WDC_WD40EZRX-${makeid()}`
	} else if (serial.startsWith('WDC_WD30EFRX')) {
		return `WDC_WD30EFRX-${makeid()}`
	} else if (serial.startsWith('ST4000VN000')) {
		return `ST4000VN000-${makeid()}`
	} else if (serial.startsWith('ST4000DM000')) {
		return `ST4000DM000-${makeid()}`
	}

	return serial
}


function makeid() {
	let text = ''
	const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'

	for (let i = 0; i < 5; i += 1) {
		text += possible.charAt(Math.floor(Math.random() * possible.length))
	}

	return text
}
