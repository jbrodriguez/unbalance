var path = require('path')
var webpack = require('webpack')
var HtmlWebpackPlugin = require('html-webpack-plugin')

module.exports = {
	devtool: 'eval-source-map',
	entry: [
		'webpack-hot-middleware/client?reload=true',
		path.join(__dirname, 'client/main.js')
	],
	output: {
		path: path.join(__dirname, '/dist/'),
		filename: 'js/[name].js',
		publicPath: '/'
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: 'client/index.tpl.html',
			inject: 'body',
			filename: 'index.html'
		}),
		new webpack.optimize.OccurenceOrderPlugin(),
		new webpack.HotModuleReplacementPlugin(),
		new webpack.NoErrorsPlugin(),
		new webpack.DefinePlugin({
			'process.env.NODE_ENV': JSON.stringify('development')
		})
	],
	module: {
		loaders: [{
			loader: 'babel-loader',
			exclude: /node_modules/,
			test: /\.js$/,
			query: {
				plugins: ['transform-runtime'],
				presets: ['es2015', 'react']
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
			test: /\.scss$/,
			loaders: [
				'style',
				'css?modules&localIdentName=[name]---[local]---[hash:base64:5]',
				'postcss',
				'sass'
			]
		}, {
			test: /\.css$/,
			loader: 'style!css?modules&localIdentName=[name]---[local]---[hash:base64:5]'
		}]
	},
	sassLoader: {
		includePaths: path.join(__dirname, '/client/styles/')
	}
}
