angular.module('unbalance.services', [
	'unbalance.socket'
])

.factory('core', ['socket', function(socket) {
	var api = '/api/v1/'
	var bus = {};

	bus.getStatus = function() {
		var msg = {
			id: 0,
			method: api + 'get/status',
			params: {},
			data: {}
		}
		return socket.request(msg)
	};

	bus.getBestFit = function(fromDisk, toDisk) {
		var msg = {
			id: 0,
			method: api + 'get/bestFit',
			params: { "fromDisk": fromDisk, "toDisk": toDisk },
			data: {}
		}

		return socket.request(msg);
	}

	bus.move = function() {
		var msg = {
			id: 0,
			method: api + 'post/move',
			params: {},
			data: {}
		}

		return socket.signal(msg);
	}

	return bus;
}])