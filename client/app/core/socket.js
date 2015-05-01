(function () {
    'use strict';

    angular
        .module('app.core')
        .factory('socket', socket);

    /* @ngInject */
    function socket($location, $rootScope) {
		// Open a WebSocket connection
		var url = $location.host() + ":" + $location.port();

		var skt = new WebSocket("ws://" + url + "/ws");

		var actions = [];

		skt.onopen = function(evt) { 
			console.log("Connection open ...");
		};

		skt.onmessage = function(evt) {
			var msg = JSON.parse(event.data);
			console.log(msg);
			if (actions.hasOwnProperty(msg.topic))
				$rootScope.$apply(actions[msg.topic](msg.payload));
		};

		skt.onclose = function(evt) {
			console.log("Connection closed."); 
		};

    	var service = {
    		register: register,
    		send: send,
    	};

    	return service;

    	function register(action, fn) {
    		actions[action] = fn;
    	};

    	function send(topic, data) {
    		console.log("are we there yet: " + topic + " " + data)
    		var message = {
    			topic: topic,
    			data: angular.toJson(data)
    		};

    		skt.send(angular.toJson(message));
    	}
    }

})();