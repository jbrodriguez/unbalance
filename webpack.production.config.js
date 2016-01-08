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
		path: path.join(__dirname, '/dist'),
		filename: 'app/[name]-[hash:7].min.js'
	},
	plugins: [
		new webpack.optimize.OccurenceOrderPlugin(),
		new HtmlWebpackPlugin({
			template: 'client/index.tpl.html',
			inject: 'body',
			filename: 'index.html'
		}),
		new ExtractTextPlugin('app/[name]-[hash:7].min.css'),

		new webpack.DefinePlugin({
			'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV)
		})
	],
	module: {
		loaders: [{
			loader: 'babel-loader',
			exclude: /node_modules/,
			test: /\.js$/,
			query: {
				plugins: ['transform-runtime'],
				presets: ['react', 'es2015', 'stage-2']
			}
		}, {
			test: /\.json?$/,
			loader: 'json'
		}, { 
			test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/, 
			loader: "url-loader?limit=10000&minetype=application/font-woff&name=img/[name]-[hash:7].[ext]"
		}, { 
			test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/, 
			loader: "file?hash=sha512&digest=hex&name=img/[name]-[hash:7].[ext]"
		}, {
			test: /\.(jpe?g|png|gif|svg)$/i,
			include: path.resolve(__dirname, 'client/img'),
			loaders: [
				'file?hash=sha512&digest=hex&name=img/[name]-[hash:7].[ext]',
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