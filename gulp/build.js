var gulp = require("gulp");
var exec = require('child_process').execSync;
var gutil = require('gulp-util');
var path = require('path');
var config = require('./config.js');

function command(tag, cmd) {
	var result = exec(cmd, {encoding: 'utf-8'});
	gutil.log(gutil.colors.yellow('tag: ' + tag) + gutil.colors.green(result));
}

gulp.task('build:server', ['tools'], function() {
	// var src = path.join(process.cwd(), 'server');
	// var dst = path.join(process.cwd(), 'dist');
	var src = config.build.server;
	var dst = config.build.dist;

	gutil.log('\n src: ' + src + '\n dst: ' + dst);
	command('build', 'GOOS=linux GOARCH=amd64 go build -v -o ' + path.join(dst, 'unbalance') + ' ' + path.join(src, 'boot.go'));
});

gulp.task('build:client', ['reference'], function() {
	gutil.log('Revved and injected');
});

// gulp.task('build:static', function() {
// 	gutil.log('Revved and injected');
// });
