const proxy = require('http-proxy-middleware')
const Bundler = require('parcel-bundler')
const express = require('express')
const process = require('process')

console.log(process.env.NODE_ENV)

const options = {
	target: 'http://wopr.apertoire.org:6237/',
	changeOrigin: true,
}

let bundler = new Bundler('./index.html', { minify: false })
let app = express()
let socket = proxy('/skt', options)

app.use('/', proxy(options))
app.use(socket)
app.use(bundler.middleware())

let server = app.listen(Number(process.env.PORT || 1234))
server.on('upgrade', socket.upgrade)
