var gulp = require("gulp");
var gutil = require('gulp-util');

var config = require('./config.js');
var helper = require('./helper.js');

gulp.task('tar', function() {
	gutil.log(gutil.colors.red(config.release.src));

	version = helper.command('version', 'cat VERSION');
	tar = "tar czvf ./unbalance-" + version + "-linux-amd64.tar.gz --exclude rev-manifest.json " + config.release.src

	helper.command('release', tar);
});