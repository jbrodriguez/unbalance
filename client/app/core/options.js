(function () {
    'use strict';

    angular
        .module('app.core')
        .factory('options', options);

    // api.$inject = ['$http', '$location', exception, logger];

    /* @ngInject */
    function options(exception, logger) {
        var config = {};

    	var service = {
            config: config
    	};

    	return service;
    }

})();