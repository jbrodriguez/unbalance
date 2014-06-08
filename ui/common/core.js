angular.module('unbalance.services', [
	'unbalance.socket'
])

.factory('core', ['socket', function(socket) {
	var api = '/api/v1/'
	var bus = {};

	bus.getDisks = function() {
		var msg = {
			id: 0,
			method: api + 'get/disks',
			params: {},
			data: {},
		}
		return socket.request(msg)
	};

	return bus;
}])