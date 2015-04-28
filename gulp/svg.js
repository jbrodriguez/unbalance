var gulp = require('gulp');
var path = require('path');
var folder = require('./config.json');


// SVG optimization task
gulp.task('svg', function () {
  return gulp.src(folder.svg)
//    .pipe(svgmin())
    .pipe(gulp.dest(path.join(folder.dist, 'img')));
});