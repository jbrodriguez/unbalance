var gulp = require('gulp');
var gutil = require('gulp-util');
var	sass = require('gulp-sass');
var	autoprefixer = require('gulp-autoprefixer');
var concat = require('gulp-concat');
var	minifyCss = require('gulp-minify-css');
var bytediff = require('gulp-bytediff');
var plumber = require('gulp-plumber');

var config = require('./config.js');
var helper = require('./helper.js')

gulp.task('styles', function() {
    gutil.log('Bundling, minifying, and copying the app\'s css');

    return gulp.src(config.styles.src)
        .pipe(plumber())
		.pipe(sass())
        .pipe(concat('app.min.css')) // Before bytediff or after
        .pipe(autoprefixer('last 2 version', '> 5%'))
        .pipe(bytediff.start())
//        .pipe(minifyCss())
        .pipe(bytediff.stop(helper.bytediffFormatter))
        //        .pipe(plug.concat('all.min.css')) // Before bytediff or after
        .pipe(plumber.stop())
        .pipe(gulp.dest(config.styles.dst));
});

gulp.task('vendor-styles', function() {
    gutil.log('Bundling, minifying, and copying the vendor css');

    return gulp.src(config.styles.vendors)
        .pipe(plumber())
        .pipe(concat('vendor.min.css')) // Before bytediff or after
        .pipe(plumber.stop())
        .pipe(gulp.dest(config.styles.dst));
});