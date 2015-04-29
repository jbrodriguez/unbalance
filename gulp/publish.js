var gulp = require("gulp");
var gutil = require('gulp-util');
var rsync = require('rsync');
var path = require('path');
var config = require('./config.js');

function publish(server) {
//	src = path.join(config.publish.src, 'public/');
	dst = server + ":" + config.publish.dst;

	cmd = new rsync()
//		.dry()
	    .flags('rlptDvzhPX')
	    .shell('ssh')
	    .delete()
	    .exclude('.htaccess')
	    .exclude('cgi-bin')
	    .exclude('rev-manifest.json')
	    .compress()
//	    .set('log-file', '/Volumes/Users/kayak/tmp/jbrodriguezio.log')
	    .source(config.publish.src)
	    .destination(dst)
	    .output(
	    	function(data) {
	    		console.log(data.toString());
	    	},
	    	function(data) {
	    		console.log(gutil.colors.magenta(data.toString()));
	    		console.beep();
	    	}
	    );

	cmd.execute(function(error, code, cmd) {
	    gutil.log('rysnc: All done executing', cmd);
		gutil.log('Deployed to ' + gutil.colors.blue(server));
	});		
}

gulp.task('publish:wopr', ['build:server', 'build:client'], function() {
	publish('unraid6');
});

gulp.task('publish:hal', ['build:server', 'build:client'], function() {
	publish('hal');
});
