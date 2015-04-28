var gulp = require('gulp');
var gutil = require('gulp-util');
var filter = require('gulp-filter');
var rev = require('gulp-rev');
var path = require('path');
var config = require('./config.js');
var debug = require('gulp-debug');
var plumber = require('gulp-plumber');

/**
 * Inject all the files into the new index.html
 * rev, but no map
 * @return {Stream}
 */
gulp.task('fingerprint', ['scripts', 'vendor-scripts', 'styles', 'vendor-styles', 'images', 'svg'], function() {
    gutil.log('Fingerprinting files ' + config.fingerprint.dst);

    var minFilter = filter(config.fingerprint.minFilter);
    var indexFilter = filter(config.fingerprint.index);

    gutil.log('src: ' + gutil.colors.green(config.fingerprint.src));
    gutil.log('dst: ' + gutil.colors.green(config.fingerprint.dst));    

    return gulp
        // Write the revisioned files
        .src(config.fingerprint.src) // add all built min files and index.html

        .pipe(debug({title: 'sources '}))

        .pipe(minFilter) // filter the stream to minified css and js

        .pipe(debug({title: 'minfilter '}))

        .pipe(rev()) // create files with rev's
        .pipe(gulp.dest(config.fingerprint.dst)) // write the rev files
        .pipe(minFilter.restore()) // remove filter, back to original stream

        // inject the files into index.html
        .pipe(indexFilter) // filter to index.html
        .pipe(gulp.dest(config.fingerprint.dst)) // write the rev files
        .pipe(indexFilter.restore()) // remove filter, back to original stream

//        .pipe(debug({title: 'files to be manifested'}))
        
        .pipe(rev.manifest()) // create the manifest (must happen last or we screw up the injection)
        .pipe(gulp.dest(config.fingerprint.dst)) // write the index.html file changes
});
