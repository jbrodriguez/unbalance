var gulp = require("gulp");
var exec = require('child_process').execSync;
var gutil = require('gulp-util');
var path = require('path');
var del = require('del');

function command(tag, cmd) {
	var result = exec(cmd, {encoding: 'utf-8'});
	gutil.log('tag: \n' + result);
}

gulp.task('build:server', function() {
	var src = path.join(process.cwd(), 'server');
	var dst = path.join(process.cwd(), 'dist');

	gutil.log('src: ' + src + ' dst: ' + dst);

	del.sync(dst);	
	command('build', 'GOOS=linux GOARCH=amd64 go build -v -o ' + path.join(dst, 'unbalance') + ' ' + path.join(src, 'boot.go'));
});

gulp.task('build:client', ['reference'], function() {
	gutil.log('Revved and injected');
});

// gulp.task('build:static', function() {
// 	gutil.log('Revved and injected');
// });
