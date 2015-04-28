var gulp = require('gulp');
var templatecache = require('gulp-angular-templatecache');
var gutil = require('gulp-util');
var path = require('path');
var folder = require('./config.json');

gulp.task('templates', function() {
    gutil.log('Creating an AngularJS $templateCache');

	return gulp.src(folder.templates)

        // // .pipe(plug.bytediff.start())
        // .pipe(plug.minifyHtml({
        //     empty: true
        // }))
        // .pipe(plug.bytediff.stop(bytediffFormatter))
        .pipe(templatecache('templates.js', {
            module: 'app.core',
            standalone: false,
            root: 'app/'
        }))
        .pipe(gulp.dest(path.join(folder.staging, 'scripts')));
});