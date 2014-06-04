'use strict';

angular.module('unbalance.services', [
	'unbalance.socket'
])

.factory('core', ['socket', function(socket) {
	var api = '/v1/'
	var bus = {};

	bus.getDisks = function() {
		var msg = {
			type: api + 'get/disks'
		}
		return socket.request()
	};

	return bus;
}])