require('babel-core/register')({})

const server = require('./server').default

const port = process.env.PORT || 3000

server.listen(port, '0.0.0.0', err => {
	if (err) {
		console.log(err)
	}
	console.info('==> 🌎 Listening on port %s. Open up http://0.0.0.0:%s/ in your browser.', port, port)
})
