(function () {
    'use strict';

    angular
        .module('app.core')
        .factory('api', api);

    // api.$inject = ['$http', '$location', exception, logger];

    /* @ngInject */
    function api($http, exception, logger) {
    	var ep = "/api/v1";

    	var service = {
            getStatus: getStatus,
            calculateBestFit: calculateBestFit,
            move: move,
    	};

    	return service;

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