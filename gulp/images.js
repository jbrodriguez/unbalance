var gulp = require('gulp');
var gutil = require('gulp-util');
var path = require('path');
var	imagemin = require('gulp-imagemin');
var changed = require('gulp-changed');
var cache = require('gulp-cache');
var config= require('./config.js');
var debug = require('gulp-debug');

// gulp.task('images', function () {
//   return gulp.src('src/img/*.*')
//   	.pipe(changed('client/img'))
//     .pipe(imagemin())
//     .pipe(gulp.dest('client/img'));
// });

/**
 * Compress images
 * @return {Stream}
 */
gulp.task('images', function() {
    gutil.log('Compressing, caching, and copying images ');

    gutil.log('cache: ' + gutil.colors.green(config.images.cache));
    gutil.log('src: ' + gutil.colors.green(config.images.src));
    gutil.log('dst: ' + gutil.colors.green(config.images.dst));

    var custom = new cache.Cache({ tmpDir: config.images.cache, cacheDirName: '' })

    return gulp
		.src(config.images.src)
		.pipe(debug({title: '1'}))
		// .pipe(changed(stage))
        .pipe(cache(imagemin({optimizationLevel: 3}), {fileCache: custom, name: ''}))
        .pipe(gulp.dest(config.images.dst));
});