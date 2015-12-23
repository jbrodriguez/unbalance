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

module.exports = {
	humanBytes,
	percentage,
}