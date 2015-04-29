var gulp = require('gulp');
var gutil = require('gulp-util');
var	imagemin = require('gulp-imagemin');
var cache = require('gulp-cache');

var config= require('./config.js');


gulp.task('images', function() {
    gutil.log('Compressing, caching, and copying images ');

    gutil.log('cache: ' + gutil.colors.green(config.images.cache));
    gutil.log('src: ' + gutil.colors.green(config.images.src));
    gutil.log('dst: ' + gutil.colors.green(config.images.dst));

    var custom = new cache.Cache({ tmpDir: config.images.cache, cacheDirName: '' })

    return gulp
		.src(config.images.src)
        .pipe(cache(imagemin({optimizationLevel: 3}), {fileCache: custom, name: ''}))
        .pipe(gulp.dest(config.images.dst));
});