var path = require('path')
var webpack = require('webpack')
var HtmlWebpackPlugin = require('html-webpack-plugin')
var ExtractTextPlugin = require('extract-text-webpack-plugin')
var StatsPlugin = require('stats-webpack-plugin')

module.exports = {
	entry: [
		path.join(__dirname, 'client/main.js')
	],
	output: {
		path: path.join(__dirname, '/dist/'),
		filename: '[name]-[hash:7].min.js'
	},
	plugins: [
		new webpack.optimize.OccurenceOrderPlugin(),
		new HtmlWebpackPlugin({
			template: 'client/index.tpl.html',
			inject: 'body',
			filename: 'index.html'
		}),
		new ExtractTextPlugin('[name]-[hash:7].min.css'),
		new webpack.optimize.UglifyJsPlugin({
			compressor: {
				warnings: false,
				screw_ie8: true
			}
		}),
		new StatsPlugin('webpack.stats.json', {
			source: false,
			modules: false
		}),
		new webpack.DefinePlugin({
			'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV)
		})
	],
	module: {
		loaders: [{
			test: /\.js?$/,
			exclude: /node_modules/,
			loader: 'babel'
		}, {
			test: /\.(jpe?g|png|gif|svg)$/i,
			include: path.resolve(__dirname, 'client/img'),
			loaders: [
				'file?hash=sha512&digest=hex&name=[name]-[hash:7].[ext]',
				'image-webpack?{progressive:true, optimizationLevel: 7, interlaced: false, pngquant:{quality: "65-90", speed: 4}}'
			]
		}, {
		// 	test: /\.scss$/,
		// 	loaders: [
		// 		'style',
		// 		'css?modules&localIdentName=[name]---[local]---[hash:base64:5]',
		// 		'postcss',
		// 		'sass'
		// 	]
		// }, {
		// 	test: /\.css$/,
		// 	loader: 'style!css?modules&localIdentName=[name]---[local]---[hash:base64:5]'
		// }]			
			test: /\.scss$/,
			loader: ExtractTextPlugin.extract('style', 'css?modules&localIdentName=[name]---[local]---[hash:base64:5]!postcss!sass'),
		}, {
			test: /\.css$/,
			loader: ExtractTextPlugin.extract('style', 'css?modules&localIdentName=[name]---[local]---[hash:base64:5]!postcss')
		}]
	},
	sassLoader: {
		includePaths: path.join(__dirname, '/client/styles/')
	},	
	postcss: [
		require('autoprefixer')
	]
}