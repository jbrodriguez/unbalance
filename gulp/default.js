var gulp = require('gulp');
var series = require('run-sequence');

gulp.task('default', function(cb) {
	series('clean', 'publish:wopr', cb);
})