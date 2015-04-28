var gulp = require('gulp');
var gutil = require('gulp-util');
var path = require('path');
var concat = require('gulp-concat');
var annotate = require('gulp-ng-annotate');
var bytediff = require('gulp-bytediff');
var uglify = require('gulp-uglify');
var config = require('./config.js');
var debug = require('gulp-debug');

/**
 * Minify and bundle the app's JavaScript
 * @return {Stream}
 */
gulp.task('scripts', ['templates'], function() {
    gutil.log('Bundling, minifying, and copying the app\'s JavaScript');

    var source = [].concat(config.scripts.src, path.join(config.templates.dst, 'templates.js'));

    return gulp
        .src(source)
        // .pipe(plug.sourcemaps.init()) // get screwed up in the file rev process
        .pipe(concat('app.min.js'))
        .pipe(annotate({ add: true, single_quotes: true }))
        .pipe(bytediff.start())
        .pipe(uglify({ mangle: true }))
        .pipe(bytediff.stop(bytediffFormatter))
        // .pipe(plug.sourcemaps.write('./'))
        .pipe(gulp.dest(config.scripts.dst));
});

gulp.task('vendor-scripts', function() {
    gutil.log('Bundling, minifying, and copying the vendor JavaScript');

    return gulp
        .src(config.scripts.vendors)
        // .pipe(plug.sourcemaps.init()) // get screwed up in the file rev process
        .pipe(concat('vendor.min.js'))
        .pipe(bytediff.start())
        .pipe(uglify())
        .pipe(bytediff.stop(bytediffFormatter))
        // .pipe(plug.sourcemaps.write('./'))
        .pipe(gulp.dest(config.scripts.dst));
});

function bytediffFormatter(data) {
    var difference = (data.savings > 0) ? ' smaller.' : ' larger.';
    return data.fileName + ' went from ' +
        (data.startSize / 1000).toFixed(2) + ' kB to ' + (data.endSize / 1000).toFixed(2) + ' kB' +
        ' and is ' + formatPercent(1 - data.percent, 2) + '%' + difference;
};

function formatPercent(num, precision) {
    return (num * 100).toFixed(precision);
};