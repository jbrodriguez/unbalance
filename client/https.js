const proxy = require('http-proxy-middleware')
const Bundler = require('parcel-bundler')
const express = require('express')
const process = require('process')
const https = require('https')

console.log(process.env.NODE_ENV)

const options = {
	target: 'https://lucy.apertoire.org:6237/',
	changeOrigin: true,
	secure: false,
	agent: https.globalAgent,
}

let bundler = new Bundler('./index.html', { minify: false })
let app = express()
let socket = proxy('/api/*', { ...options, ws: true })

app.use('/', proxy(options))
app.use(socket)
app.use(bundler.middleware())

let server = app.listen(Number(process.env.PORT || 1234))
server.on('upgrade', socket.upgrade)
