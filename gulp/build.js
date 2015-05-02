var gulp = require("gulp");
var exec = require('child_process').execSync;
var gutil = require('gulp-util');
var path = require('path');
var config = require('./config.js');
var strings = require('string');

function command(tag, cmd) {
	gutil.log(gutil.colors.blue('executing ' + cmd))
	result = exec(cmd, {encoding: 'utf-8'});
	var output = strings(result).chompRight('\n').toString();
	gutil.log(gutil.colors.yellow('tag: [' + tag + '] ') + gutil.colors.green(output));
	return output;
}

gulp.task('build:server', ['tools'], function() {
	// var src = path.join(process.cwd(), 'server');
	// var dst = path.join(process.cwd(), 'dist');
	var src = config.build.server;
	var dst = config.build.dist;

	version = command('version', 'cat VERSION');
	count = command('count', 'git rev-list HEAD --count')
	hash = command('hash', 'git rev-parse --short HEAD')

	gutil.log('\n src: ' + src + '\n dst: ' + dst);
	command('build', 'GOOS=linux GOARCH=amd64 go build -ldflags \"-X main.Version ' + version + '-' + count + '.' + hash + '\" -v -o ' + path.join(dst, 'unbalance') + ' ' + path.join(src, 'boot.go'));
});

gulp.task('build:client', ['reference'], function() {
	gutil.log('Revved and injected');
});

// gulp.task('build:static', function() {
// 	gutil.log('Revved and injected');
// });
