var gulp = require('gulp');
var gutil = require('gulp-util');
var filter = require('gulp-filter');
var rev = require('gulp-rev');
var revReplace = require('gulp-rev-replace');
var inject = require('gulp-inject');
var path = require('path');
var folder = require('./config.json');
var debug = require('gulp-debug');


/**
 * Inject all the files into the new index.html
 * rev, but no map
 * @return {Stream}
 */
gulp.task('rev-inject', ['scripts', 'styles', 'images', 'svg'], function() {
    gutil.log('Rev\'ing files and building index.html');

    var minified = path.join(folder.dist, '*.min.*');
    var index = path.join(folder.client, 'index.html');
    var minFilter = filter(['*.min.*', '!**/*.map']);
    var indexFilter = filter(['index.html']);

    var stream = gulp
        // Write the revisioned files
        .src([].concat(minified, index)) // add all built min files and index.html
        .pipe(debug({title: '1'}))        

        .pipe(minFilter) // filter the stream to minified css and js
        .pipe(debug({title: '2'}))        

        .pipe(rev()) // create files with rev's
        .pipe(debug({title: '3'}))        

        .pipe(gulp.dest(folder.dist)) // write the rev files
        .pipe(debug({title: '4'}))        

        .pipe(minFilter.restore()) // remove filter, back to original stream

        // inject the files into index.html
        .pipe(indexFilter) // filter to index.html
//        .pipe(inject('content/vendor.min.css', 'inject-vendor'))
        .pipe(doInject('app.min.css'))
//        .pipe(inject('vendor.min.js', 'inject-vendor'))
        .pipe(doInject('app.min.js'))
        .pipe(gulp.dest(folder.dist)) // write the rev files
        .pipe(indexFilter.restore()) // remove filter, back to original stream

        // replace the files referenced in index.html with the rev'd files
        .pipe(revReplace()) // Substitute in new filenames
        .pipe(gulp.dest(folder.dist)) // write the index.html file changes
        .pipe(rev.manifest()) // create the manifest (must happen last or we screw up the injection)

        .pipe(gulp.dest(folder.dist)); // write the manifest

    function doInject(folderpath, name) {
        var pathGlob = path.join(folder.dist, folderpath);
        var options = {
            ignorePath: folder.dist.substring(1),
            read: false
        };
        if (name) {
            options.name = name;
        }

        gutil.log('glob: ' + gutil.colors.red(pathGlob));
        gutil.log('options: ' + gutil.colors.red(options));

        return inject(gulp.src(pathGlob), options);
    }
});
