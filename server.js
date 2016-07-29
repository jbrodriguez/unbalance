import path 				from 'path'
import express 				from 'express'
import webpack 				from 'webpack'
import webpackMiddleware 	from 'webpack-dev-middleware'
import webpackHotMiddleware from 'webpack-hot-middleware'
import config 				from './webpack.config.js'
import httpProxy 			from 'http-proxy'

// var path = require('path');
// var express = require('express');
// var webpack = require('webpack');
// var webpackMiddleware = require('webpack-dev-middleware');
// var webpackHotMiddleware = require('webpack-hot-middleware');
// var config = require('./webpack.config.js');
// var httpProxy = require('http-proxy');


// We need to add a configuration to our proxy server,
// as we are now proxying outside localhost
var proxy = httpProxy.createProxyServer({
  changeOrigin: true,
  ws: true,
})

const isDeveloping = process.env.NODE_ENV !== 'production'
// const port = isDeveloping ? 3000 : process.env.PORT;

const app = express()
const server = require('http').createServer(app)

if (isDeveloping) {
	const compiler = webpack(config);
	const middleware = webpackMiddleware(compiler, {
		publicPath: config.output.publicPath,
		contentBase: 'src',
		stats: {
			colors: true,
			hash: false,
			timings: true,
			chunks: false,
			chunkModules: false,
			modules: false
		}
	});

	app.use(middleware);
	app.use(webpackHotMiddleware(compiler));

	app.all('/api/*', function(req, res) {
		proxy.web(req, res, {target: 'http://wopr.apertoire.org:6237'})
	})

	server.on('upgrade', function(req, socket, head) {
		proxy.ws(req, socket, head, {target: 'http://wopr.apertoire.org:6237'})
	});

	// app.all('/skt', function(req, res) {
	// 	proxy.ws(req, res, {
	// 		target: 'ws://wopr.apertoire.org:6237/skt'
	// 	})
	// })

	app.get('*', function response(req, res) {
		res.write(middleware.fileSystem.readFileSync(path.join(__dirname, 'dist/index.html')));
		res.end();
	});

} else {
	app.use(express.static(__dirname + '/dist'));
	app.get('*', function response(req, res) {
		res.sendFile(path.join(__dirname, 'dist/index.html'));
	});
}

export default server

// // var server = require('http').createServer(app);
// server.listen(port, '0.0.0.0', function onStart(err) {
// 	if (err) {
// 		console.log(err);
// 	}
// 	console.info('==> 🌎 Listening on port %s. Open up http://0.0.0.0:%s/ in your browser.', port, port);
// });
