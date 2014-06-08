angular.module('unbalance.socket', [
])

.factory('socket', ['$q', '$rootScope', function($q, $rootScope) {
	var currentCallbackId = 0;
	var callbacks = {};

	var ws = new WebSocket("ws://blackbeard.apertoire.org:6237/api");

	ws.onopen = function() {
		console.log("socket has been opened");
	};

	ws.onmessage = function(message) {
		console.log("message is: ", message)
		listener(JSON.parse(message.data))
	}

	ws.onclose = function() {
		console.log("socket has been closed")
	}

	var request = function(req) {
		var defer = $q.defer();

		var callbackId = getId();
		callbacks[callbackId] = {
			time: new Date(),
			promise: defer
		};
		req.id = callbackId;

		console.log('sending request: ', JSON.stringify(req));
		ws.send(JSON.stringify(req));

		return defer.promise;
	};

	var signal = function(sig) {
		console.log('sending signal: ', sig);
		ws.send(JSON.stringify(sig));
	};

	function listener(data) {
		var msg = data;
		console.log("received data from websocket: ", msg);

		if (callbacks.hasOwnProperty(msg.Id)) {
			console.log("callback was: ", JSON.stringify(callbacks[msg.Id]));
			callbacks[msg.Id].promise.resolve(msg.result);
		} else {
			console.log("emitting event");
			$rootScope.emit(msg.method, msg);
		}
	}

	function getId() {
		currentCallbackId += 1;
		if (currentCallbackId > 10000) {
			currentCallbackId = 0;
		}
		return currentCallbackId;
	}

	return {
		request: request,
		signal: signal
	}
}]);