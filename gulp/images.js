var gulp = require('gulp');
var	imagemin = require('gulp-imagemin');
var changed = require('gulp-changed');

gulp.task('images', function () {
  return gulp.src('src/img/*.*')
  	.pipe(changed('client/img'))
    .pipe(imagemin())
    .pipe(gulp.dest('client/img'));
});