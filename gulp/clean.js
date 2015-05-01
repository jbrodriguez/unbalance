var gulp = require("gulp");
var gutil = require('gulp-util');
var del = require('del');
var config = require('./config.js');

/**
 * Remove all files from the build folder
 * One way to run clean before all tasks is to run
 * from the cmd line: gulp clean && gulp build
 * @return {Stream}
 */
gulp.task('clean', function() {
    gutil.log('Cleaning: ' + gutil.colors.blue(config.clean.dist) + ' ' + gutil.colors.blue(config.clean.staging));

    del.sync(config.clean.staging);
    del.sync(config.clean.dist);
});
