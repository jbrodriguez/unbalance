var gulp = require('gulp');
var path = require('path');
var config = require('./config.js');


// SVG optimization task
gulp.task('svg', function () {
  return gulp.src(config.svg.src)
//    .pipe(svgmin())
    .pipe(gulp.dest(config.svg.dst));
});