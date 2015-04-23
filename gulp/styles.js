var gulp = require('gulp');
var	sass = require('gulp-sass');
var	autoprefixer = require('gulp-autoprefixer');
var	minifycss = require('gulp-minify-css');

gulp.task('styles', function() {
	return gulp.src('src/styles/*.scss')
		.pipe(sass())
		.pipe(gulp.dest('client/css'))
});