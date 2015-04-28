var gulp = require('gulp');
var gutil = require('gulp-util');
var path = require('path');
var	imagemin = require('gulp-imagemin');
var changed = require('gulp-changed');
var cache = require('gulp-cache');
var folder = require('./config.json');

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
    var dest = path.join(folder.dist,'img');
    var custom = new cache.Cache({ tmpDir: folder.staging, cacheDirName: 'img' })

    gutil.log('Compressing, caching, and copying images');

    return gulp
		.src(folder.images)
		// .pipe(changed(stage))
        .pipe(cache(imagemin({optimizationLevel: 3}), {fileCache: custom, name: ''}))
        .pipe(gulp.dest(dest));
});