var gulp = require("gulp");
var reload = require("browser-sync").reload;

gulp.task('build:content', ['html'], reload);

gulp.task('build:all', ['html', 'styles', 'images', 'svg'], reload);

