var gulp = require('gulp');
var gutil = require('gulp-util');
var revReplace = require('gulp-rev-replace');
var config = require('./config.js');

gulp.task('reference', ['fingerprint'], function() {
    gutil.log('Changing fingerprinted assets references in html and js (templates)');

	var manifest = gulp.src(config.reference.dst + "rev-manifest.json");

    gutil.log('src: ' + gutil.colors.green(config.reference.src));
    gutil.log('dst: ' + gutil.colors.green(config.reference.dst));
    gutil.log('src: ' + gutil.colors.green(config.reference.ext));

    return gulp.src(config.reference.src) // add all built min files and index.html
        // replace the files referenced in index.html with the rev'd files
        .pipe(revReplace({manifest: manifest, replaceInExtensions: config.reference.ext})) // Substitute in new filenames
        .pipe(gulp.dest(config.reference.dst)); // write the manifest
});