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
			data: {},
		}
		return socket.request(msg)
	};

	return bus;
}])