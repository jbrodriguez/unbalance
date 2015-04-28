var gulp = require('gulp');
var gutil = require('gulp-util');
var revReplace = require('gulp-rev-replace');
var path = require('path');
var config = require('./config.js');
var debug = require('gulp-debug');

gulp.task('reference', ['fingerprint'], function() {
    gutil.log('Reconstructing index.html');

	var manifest = gulp.src(config.reference.dst + "rev-manifest.json");

    gutil.log('src: ' + gutil.colors.green(config.reference.src));
    gutil.log('dst: ' + gutil.colors.green(config.reference.dst));
    gutil.log('src: ' + gutil.colors.green(config.reference.ext));

    return gulp.src(config.reference.src) // add all built min files and index.html

    	.pipe(debug({title: 'ref-sources'}))

        // replace the files referenced in index.html with the rev'd files
        .pipe(revReplace({manifest: manifest, replaceInExtensions: config.reference.ext})) // Substitute in new filenames

    	.pipe(debug({title: 'ref-replaced'}))

        .pipe(gulp.dest(config.reference.dst)); // write the manifest
});