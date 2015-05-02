var exec = require('child_process').execSync;
var gutil = require('gulp-util');
var strings = require('string');

function bytediffFormatter(data) {
    var difference = (data.savings > 0) ? ' smaller.' : ' larger.';
    return data.fileName + ' went from ' +
        (data.startSize / 1000).toFixed(2) + ' kB to ' + (data.endSize / 1000).toFixed(2) + ' kB' +
        ' and is ' + formatPercent(1 - data.percent, 2) + '%' + difference;
};

function formatPercent(num, precision) {
    return (num * 100).toFixed(precision);
};

function command(tag, cmd) {
	gutil.log(gutil.colors.blue('executing ' + cmd))
	result = exec(cmd, {encoding: 'utf-8'});
	var output = strings(result).chompRight('\n').toString();
	gutil.log(gutil.colors.yellow('tag: [' + tag + '] ') + gutil.colors.green(output));
	return output;
}

module.exports = {
	command: command,
	bytediffFormatter: bytediffFormatter,
	formatPercent: formatPercent
}