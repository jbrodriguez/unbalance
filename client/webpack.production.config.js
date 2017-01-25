const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const ExtractTextPlugin = require('extract-text-webpack-plugin')

module.exports = {
	entry: [
		'./src/main.js',
	],
	output: {
		path: path.join(__dirname, '..', 'dist'),
		filename: 'app/[name]-[hash:7].min.js',
	},
	plugins: [
		new webpack.optimize.OccurenceOrderPlugin(),
		new HtmlWebpackPlugin({
			template: 'index.tpl.html',
			inject: 'body',
			filename: 'index.html',
		}),
		new ExtractTextPlugin('app/[name]-[hash:7].min.css'),

		new webpack.DefinePlugin({
			'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV),
		}),
	],
	module: {
		loaders: [{
			test: /\.jsx?$/,
			loader: 'babel',
			include: path.join(__dirname, 'src'),
		}, {
			test: /\.json?$/,
			loader: 'json',
		}, {
			test: /\.woff(2)?(\?v=[0-9]\.[0-9]\.[0-9])?$/,
			loader: 'url-loader?limit=10000&minetype=application/font-woff&name=img/[name]-[hash:7].[ext]',
		}, {
			test: /\.(ttf|eot|svg)(\?v=[0-9]\.[0-9]\.[0-9])?$/,
			loader: 'file?hash=sha512&digest=hex&name=img/[name]-[hash:7].[ext]',
		}, {
			test: /\.(jpe?g|png|gif|svg)$/i,
			include: path.resolve(__dirname, 'src/img'),
			loaders: [
				'file?hash=sha512&digest=hex&name=img/[name]-[hash:7].[ext]',
				'image-webpack?{progressive:true, optimizationLevel: 7, interlaced: false, pngquant:{quality: "65-90", speed: 4}}',
			],
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
			include: path.join(__dirname, 'src/styles'),
			loader: ExtractTextPlugin.extract('style', 'css?modules&localIdentName=[name]---[local]---[hash:base64:5]!postcss!sass'),
		}, {
			test: /\.css$/,
			loader: ExtractTextPlugin.extract('style', 'css?modules&localIdentName=[local]'),
		}],
	},
	postcss: [
		require('autoprefixer'), // eslint-disable-line
	],
}
