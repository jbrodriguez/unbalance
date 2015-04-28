var gulp = require("gulp");
var gutil = require('gulp-util');

function publish(server) {
	gutil.log('Deployed to ' + gutil.colors.blue(server));
}

gulp.task('publish:wopr', ['build:server', 'build:client'], function() {
	publish('wopr');
});

gulp.task('publish:hal', ['build:server', 'build:client'], function() {
	publish('hal');
});

// gulp.task('publish:wopr', ['build:server', 'build:client', 'build:static'], function() {
// 	publish('wopr');
// });

// gulp.task('publish:hal', ['build:server', 'build:client', 'build:static'], function() {
// 	publish('hal');
// });