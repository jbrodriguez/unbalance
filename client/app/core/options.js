(function () {
    'use strict';

    angular
        .module('app.core')
        .factory('options', options);

    // api.$inject = ['$http', '$location', exception, logger];

    /* @ngInject */
    function options(exception, logger) {
    	var searchTerm = '';

        var filterByOptions = ['title', 'genre'];
        var filterBy = '';

        var sortByOptions = ['title', 'runtime', 'added', 'last_watched'];
        var sortBy = '';

        var sortOrderOptions = ['asc', 'desc'];
        var sortOrder = '';

    	var service = {
            searchTerm: searchTerm,
            filterByOptions: filterByOptions,
            filterBy: filterBy,
            sortByOptions: sortByOptions,
            sortBy: sortBy,
            sortOrderOptions: sortOrderOptions,
            sortOrder: sortOrder
    	};

    	return service;
    }

})();