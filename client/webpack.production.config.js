const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const ExtractTextPlugin = require('extract-text-webpack-plugin')

// fontawesome based on https://www.laurivan.com/load-fontawesome-fonts-with-webpack-2/

module.exports = {
	entry: ['./src/main.js'],
	output: {
		path: path.join(__dirname, '..', 'dist'),
		filename: 'app/[name]-[hash:7].min.js',
	},
	plugins: [
		new HtmlWebpackPlugin({
			template: 'index.tpl.html',
			inject: 'body',
			filename: 'index.html',
		}),
		new webpack.optimize.UglifyJsPlugin({
			compress: {
				warnings: false,
				drop_console: false,
			},
		}),
		new ExtractTextPlugin('app/[name]-[hash:7].min.css'),

		new webpack.DefinePlugin({
			'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV),
		}),
	],
	module: {
		rules: [
			{
				test: /\.jsx?$/,
				loader: 'babel-loader',
				exclude: /node_modules/,
			},
			{
				test: /\.(jpe?g|png|gif|svg)$/,
				include: path.resolve(__dirname, 'src/img'),
				use: [
					{
						loader: 'file-loader',
						options: {
							name: 'img/[name]-[hash:7].[ext]',
							publicPath: '/',
						},
					},
					{
						loader: 'image-webpack-loader',
						options: {
							mozjpeg: {
								progressive: true,
								quality: 65,
							},
							// optipng.enabled: false will disable optipng
							optipng: {
								enabled: false,
							},
							pngquant: {
								quality: '65-90',
								speed: 4,
							},
							gifsicle: {
								interlaced: false,
							},
						},
					},
				],
			},
			{
				test: /\.scss$/,
				include: path.resolve(__dirname, 'src/styles'),
				use: [
					'style-loader',
					{
						loader: 'css-loader',
						options: {
							// localIdentName: '[name]_[local]_[hash:base64:5]',
							minimize: true,
						},
					},
					{
						loader: 'postcss-loader',
						options: {
							ident: 'postcss',
							plugins: loader => [require('autoprefixer')()],
						},
					},
					'sass-loader',
				],
			},
			{
				test: /\.css$/,
				// include: path.resolve(__dirname, 'src/styles'),
				use: [
					'style-loader',
					{
						loader: 'css-loader',
						options: {
							localIdentName: '[name]_[local]_[hash:base64:5]',
							minimize: true,
						},
					},
				],
			},
			{
				test: /\.(ttf|otf|eot|svg|woff(2)?)(\?[a-z0-9]+)?$/,
				loader: 'file-loader?name=fonts/[name]-[hash:7].[ext]',
			},
		],
	},
}
