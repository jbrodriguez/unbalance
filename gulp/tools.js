var gulp = require('gulp');
var debug = require('gulp-debug');

var config = require('./config.js');


gulp.task('tools', function () {
  return gulp.src(config.tools.src)
  	.pipe(debug({title: "1"}))
    .pipe(gulp.dest(config.tools.dst));
});