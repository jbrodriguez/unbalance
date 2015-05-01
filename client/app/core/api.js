(function () {
    'use strict';

    angular
        .module('app.core')
        .factory('api', api);

    // api.$inject = ['$http', '$location', exception, logger];

    /* @ngInject */
    function api($http, $location, exception, logger) {
    	var ep = "/api/v1";

    	var service = {
            getConfig: getConfig,
            saveConfig: saveConfig,
            getStatus: getStatus,
            calculateBestFit: calculateBestFit,
            move: move,
    	};

    	return service;

        function getConfig() {
            return $http.get(ep + '/config')
                .then(getConfigEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for getConfig')(message);
                    $location.url('/');
                });

            function getConfigEnd(data, status, header, config) {
                return data.data
            }
        };

        function saveConfig(arg) {
            return $http.put(ep + '/config', arg)
                .then(saveConfigEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for saveConfig')(message);
                    $location.url('/');
                });

            function saveConfigEnd(data, status, header, config) {
                return data.data
            };
        };        

    	function getStatus() {
    		return $http.get(ep + '/storage')
    			.then(getStatusEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for getStatus')(message);
                    $location.url('/');
                });

    		function getStatusEnd(data, status, headers, config) {
    			return data.data;
    		}
    	};

        function calculateBestFit(params) {
            return $http.post(ep + '/storage/bestfit', params)
                .then(calculateBestFitEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for calculateBestFit')(message);
                    $location.url('/');
                });

            function calculateBestFitEnd(data, status, headers, config) {
                return data.data;
            }
        };

        function move() {
            return $http.post(ep + '/storage/move')
                .then(moveEnd)
                .catch(function(message) {
                    exception.catcher('XHR Failed for move')(message);
                    $location.url('/');
                });

            function moveEnd(data, status, headers, config) {
                console.log("this is what i got: ", data.data);
                return data.data;
            }
        }

    }

})();