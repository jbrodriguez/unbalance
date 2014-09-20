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
                logger.info('this is what i got: ', data);
    			return data.data;
    		}
    	};

    }

})();