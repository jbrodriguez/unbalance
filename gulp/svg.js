var gulp = require('gulp');

// SVG optimization task
gulp.task('svg', function () {
  return gulp.src('src/svg/*.svg')
//    .pipe(svgmin())
    .pipe(gulp.dest('client/img'));
});