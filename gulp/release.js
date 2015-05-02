var gulp = require("gulp");
var gutil = require('gulp-util');
var path = require("path");
var series = require('run-sequence');

var config = require('./config.js');
var helper = require('./helper.js');

gulp.task('release', function(cb) {
	series('clean', ['build:server', 'build:client', 'build:static'], 'tar', cb);
});