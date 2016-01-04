module.exports = {
	humanBytes,
	percentage,
	scramble,
}

function humanBytes(bytes) {
	if (bytes == 0) return '0 Byte';

	var k = 1000;
	var sizes = ['Bytes', 'KB', 'MB', 'GB', 'TB', 'PB', 'EB', 'ZB', 'YB'];
	var i = Math.floor(Math.log(bytes) / Math.log(k));
	
	return (bytes / Math.pow(k, i)).toPrecision(3) + ' ' + sizes[i];
}

function percentage(input, decimals, suffix) {
	decimals = decimals || 2;
	suffix = suffix || '%';

	return Math.round(input * Math.pow(10, decimals + 2))/Math.pow(10, decimals) + suffix
}

function scramble(serial) {
	if (serial.startsWith('WDC_WD30EZRX'))
		return 'WDC_WD30EZRX-' + makeid()
	else if (serial.startsWith('ST3000DM001'))
		return 'ST3000DM001-' + makeid()
	else if (serial.startsWith('TOSHIBA_DT01ACA300'))
		return 'TOSHIBA_DT01ACA300-' + makeid()
	else if (serial.startsWith('ST4000DM000'))
		return 'ST4000DM000-' + makeid()
	else if (serial.startsWith('WDC_WD40EZRX'))
		return 'WDC_WD40EZRX-' + makeid()
	else if (serial.startsWith('WDC_WD30EFRX'))
		return 'WDC_WD30EFRX-' + makeid()
	else if (serial.startsWith('ST4000VN000'))
		return 'ST4000VN000-' + makeid()
	else if (serial.startsWith('ST4000DM000'))
		return 'ST4000DM000-' + makeid()
	else
		return serial
}


function makeid() {
    var text = ""
    var possible = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

    for( var i=0; i < 5; i++ )
        text += possible.charAt(Math.floor(Math.random() * possible.length))

    return text
}
