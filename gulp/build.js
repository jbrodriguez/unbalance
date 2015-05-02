var gulp = require("gulp");
var path = require('path');
var gutil = require('gulp-util');

var config = require('./config.js');
var helper = require('./helper.js');


// function command(tag, cmd) {
// 	gutil.log(gutil.colors.blue('executing ' + cmd))
// 	result = exec(cmd, {encoding: 'utf-8'});
// 	var output = strings(result).chompRight('\n').toString();
// 	gutil.log(gutil.colors.yellow('tag: [' + tag + '] ') + gutil.colors.green(output));
// 	return output;
// }

gulp.task('build:server', ['tools'], function() {
	// var src = path.join(process.cwd(), 'server');
	// var dst = path.join(process.cwd(), 'dist');
	var src = config.build.server;
	var dst = config.build.dist;

	version = helper.command('version', 'cat VERSION');
	count = helper.command('count', 'git rev-list HEAD --count')
	hash = helper.command('hash', 'git rev-parse --short HEAD')

	gutil.log('\n src: ' + src + '\n dst: ' + dst);
	helper.command('build', 'GOOS=linux GOARCH=amd64 go build -ldflags \"-X main.Version ' + version + '-' + count + '.' + hash + '\" -v -o ' + path.join(dst, 'unbalance') + ' ' + path.join(src, 'boot.go'));
});

gulp.task('build:client', ['reference'], function() {
	gutil.log('Revved and injected');
});

gulp.task('build:static', function() {
	gulp.src(['./CHANGES', './LICENSE'])
	.pipe(gulp.dest(config.build.dist))
});
