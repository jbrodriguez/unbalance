const path = require('path')
const webpack = require('webpack')
const HtmlWebpackPlugin = require('html-webpack-plugin')

// fontawesome based on https://www.laurivan.com/load-fontawesome-fonts-with-webpack-2/
// based on https://medium.com/@chanonroy/webpack-2-and-font-awesome-icon-importing-59df3364f35c

module.exports = {
	devtool: 'cheap-module-eval-source-map',
	entry: ['webpack-hot-middleware/client?reload=true', './src/main.js'],
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
		new webpack.HotModuleReplacementPlugin(),
		new webpack.DefinePlugin({
			'process.env.NODE_ENV': JSON.stringify('development'),
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
							// hash: 'sha512',
							// digest: 'hex',
							name: 'img/[name]-[hash:7].[ext]',
							publicPath: '/',
							// publicPath: '../',
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
