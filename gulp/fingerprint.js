var gulp = require('gulp');
var gutil = require('gulp-util');
var filter = require('gulp-filter');
var rev = require('gulp-rev');
var path = require('path');
var config = require('./config.js');


gulp.task('fingerprint', ['scripts', 'vendor-scripts', 'styles', 'vendor-styles', 'images', 'svg'], function() {
    gutil.log('Fingerprinting files ' + config.fingerprint.dst);

    var revFilter = filter(config.fingerprint.revFilter);
    var indexFilter = filter(config.fingerprint.index);

    gutil.log('src: ' + gutil.colors.green(config.fingerprint.src));
    gutil.log('dst: ' + gutil.colors.green(config.fingerprint.dst));    

    return gulp
        // Write the revisioned files
        .src(config.fingerprint.src) // add all built min files and index.html

        .pipe(revFilter) // filter the stream to minified css and js
        .pipe(rev()) // create files with rev's
        .pipe(gulp.dest(config.fingerprint.dst)) // write the rev files
        .pipe(revFilter.restore()) // remove filter, back to original stream

        // copy index.html to destination
        .pipe(indexFilter) // filter to index.html
        .pipe(gulp.dest(config.fingerprint.dst)) // write the file
        .pipe(indexFilter.restore()) // remove filter, back to original stream

        .pipe(rev.manifest()) // create the manifest (must happen last or we screw up the injection)
        .pipe(gulp.dest(config.fingerprint.dst)) // write the index.html file changes
});
