var gulp = require('gulp');
var templatecache = require('gulp-angular-templatecache');
var gutil = require('gulp-util');
var config = require('./config.js');

gulp.task('templates', function() {
    gutil.log('Creating an AngularJS $templateCache ' + config.templates.src);

	return gulp.src(config.templates.src)
        .pipe(templatecache('templates.js', {
            module: 'app.core',
            standalone: false,
            root: 'app/'
        }))
        .pipe(gulp.dest(config.templates.dst));
});