var gulp = require('gulp');
var gutil = require('gulp-util');
var path = require('path');
var concat = require('gulp-concat');
var annotate = require('gulp-ng-annotate');
var bytediff = require('gulp-bytediff');
var uglify = require('gulp-uglify');
var plumber = require('gulp-plumber');

var config = require('./config.js');
var helper = require('./helper.js')

/**
 * Minify and bundle the app's JavaScript
 * @return {Stream}
 */
gulp.task('scripts', ['templates'], function() {
    gutil.log('Bundling, minifying, and copying the app\'s js');

    var source = [].concat(config.scripts.src, path.join(config.templates.dst, 'templates.js'));

    return gulp
        .src(source)
        .pipe(plumber())
        // .pipe(plug.sourcemaps.init()) // get screwed up in the file rev process
        .pipe(concat('app.min.js'))
        .pipe(annotate({ add: true, single_quotes: true }))
        .pipe(bytediff.start())
//        .pipe(uglify({ mangle: true }))
        .pipe(bytediff.stop(helper.bytediffFormatter))
        // .pipe(plug.sourcemaps.write('./'))

        .pipe(plumber.stop())
        .pipe(gulp.dest(config.scripts.dst));
});

gulp.task('vendor-scripts', function() {
    gutil.log('Bundling, minifying, and copying the vendor\'s js');

    return gulp
        .src(config.scripts.vendors)
        // .pipe(plug.sourcemaps.init()) // get screwed up in the file rev process
        .pipe(concat('vendor.min.js'))
        .pipe(bytediff.start())
        .pipe(uglify())
        .pipe(bytediff.stop(helper.bytediffFormatter))
        // .pipe(plug.sourcemaps.write('./'))
        .pipe(gulp.dest(config.scripts.dst));
});

