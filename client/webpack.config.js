const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')

module.exports = {
	devtool: 'eval',
	entry: [
		'webpack-hot-middleware/client?reload=true',
		'./src/main.js',
	],
	output: {
		path: path.join(__dirname, '/dist/'),
		filename: '[name].js',
		publicPath: '/',
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: 'index.tpl.html',
			inject: 'body',
			filename: 'index.html',
		}),
		new webpack.optimize.OccurenceOrderPlugin(),
		new webpack.HotModuleReplacementPlugin(),
		new webpack.NoErrorsPlugin(),
		new webpack.DefinePlugin({
			'process.env.NODE_ENV': JSON.stringify('development'),
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
			test: /\.scss$/,
			include: path.resolve(__dirname, 'src/styles'),
			loaders: [
				'style',
				'css?modules&localIdentName=[name]---[local]---[hash:base64:5]',
				'postcss',
				'sass',
			],
		}, {
			test: /\.css$/,
			loader: 'style!css?modules&localIdentName=[local]',
		}],
	},
}
