var gulp = require("gulp");
var gutil = require('gulp-util');
var del = require('del');
var folder = require('./config.json');

/**
 * Remove all files from the build folder
 * One way to run clean before all tasks is to run
 * from the cmd line: gulp clean && gulp build
 * @return {Stream}
 */
gulp.task('clean', function() {
    gutil.log('Cleaning: ' + gutil.colors.blue(folder.dist));

    del.sync(folder.dist);
});
