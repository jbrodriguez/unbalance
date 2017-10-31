const phpRound = (number, precision = 0) => {
	const factor = Math.pow(10, precision)
	const tmp = number * factor
	const roundedTmp = Math.round(tmp)
	return roundedTmp / factor
}

const toFixedFix = (n, prec) => {
	const k = Math.pow(10, prec)
	return `${(Math.round(n * k) / k).toFixed(prec)}`
}

const numberFormat = (number, decimals, decPoint, thousandsSep) => {
	const value = `${number}`.replace(/[^0-9+\-Ee.]/g, '')

	const n = !isFinite(+value) ? 0 : +value
	const prec = !isFinite(+decimals) ? 0 : Math.abs(decimals)
	const sep = typeof thousandsSep === 'undefined' ? ',' : thousandsSep
	const dec = typeof decPoint === 'undefined' ? '.' : decPoint
	let s = []

	// Fix for IE parseFloat(0.55).toFixed(0) = 0;
	s = (prec ? toFixedFix(n, prec) : `${Math.round(n)}`).split('.')

	if (s[0].length > 3) {
		s[0] = s[0].replace(/\B(?=(?:\d{3})+(?!\d))/g, sep)
	}

	if ((s[1] || '').length < prec) {
		s[1] = s[1] || ''
		s[1] += new Array(prec - s[1].length + 1).join('0')
	}

	return s.join(dec)
}

const k = 1000
const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']

const humanBytes = bytes => {
	if (bytes === 0) return '0 Byte'

	let base = bytes ? Math.floor(Math.log(bytes) / Math.log(k)) : 0
	bytes = bytes / Math.pow(k, base)

	let precision = bytes >= 100 ? 0 : bytes >= 10 ? 1 : phpRound(bytes * 100) % 100 === 0 ? 0 : 2

	if (phpRound(bytes, precision) === k) {
		bytes = 1
		precision = 2
		base += 1
	}

	// return `${(bytes / Math.pow(k, i)).toPrecision(3)} ${sizes[i]}` // eslint-disable-line

	return `${numberFormat(bytes, precision, '.', bytes >= 10000 ? ',' : '')} ${sizes[base]}`
}

// function humanBytes(bytes) {
// 	if (bytes === 0) return '0 Byte'

// 	const k = 1000
// 	const sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB']
// 	const i = Math.floor(Math.log(bytes) / Math.log(k))

// 	return `${(bytes / Math.pow(k, i)).toPrecision(3)} ${sizes[i]}` // eslint-disable-line
// }

const percentage = (input, decimals = 2, suffix = '%') => {
	return `${Math.round(input * Math.pow(10, decimals + 2)) / Math.pow(10, decimals)}${suffix}` // eslint-disable-line
}

const scramble = serial => {
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

const makeid = () => {
	let text = ''
	const possible = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789'

	for (let i = 0; i < 5; i += 1) {
		text += possible.charAt(Math.floor(Math.random() * possible.length))
	}

	return text
}

const isValid = obj => {
	if (typeof obj === 'undefined') return false
	if (!obj) return false
	return true
}

// tree helper functions
const getNode = (tree, lineage) => {
	if (lineage.length === 0) {
		return null
	} else if (lineage.length === 1) {
		return tree[lineage[0]]
	}

	const node = lineage.shift()
	return getNode(tree[node].children, lineage)
}

const markChosen = (tree, lineage, chosen) => {
	if (lineage.length === 0) {
		// no-op
	} else if (lineage.length === 1) {
		const node = tree[lineage[0]]

		if (node.checked) {
			delete chosen[node.path]
		} else {
			uncheckChildren(node.children, chosen)
			chosen[node.path] = true
		}

		node.checked = !node.checked
	} else {
		const index = lineage.shift() // this mutates lineage
		const node = tree[index]

		if (node.checked) {
			delete chosen[node.path]
			node.checked = false
		}

		markChosen(node.children, lineage, chosen)
	}
}

const uncheckChildren = (tree, chosen) => {
	if (!tree) return

	tree.forEach(node => {
		delete chosen[node.path]
		node.checked = false

		uncheckChildren(node.children, chosen)
	})
}

export { humanBytes, percentage, scramble, isValid, getNode, markChosen, uncheckChildren }
