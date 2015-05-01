(function() {
    'use strict';

    angular
        .module('blocks.storage')
        .factory('storage', storage);

    /* @ngInject */
    function storage($log, $window) {
        var prefix = 'aptn.';
        
        var service = {
            get: get,
            set: set,
            remove: remove,
            // clearAll: clearAll,
        };

        return service;
        /////////////////////

        function get(key) {
            var item = $window.localStorage.getItem(prefix + key);
            if (!item || item == 'null') {
                return null
            };

            if (item.charAt(0) === "{" || item.charAt(0) === "[") {
                return angular.fromJson(item);
            } ;

            return item;
        };

        function set(key, value) {
            if (angular.isObject(value) || angular.isArray(value)) {
                value = angular.toJson(value);
            };

            $window.localStorage.setItem(prefix + key, value);
            return true;
        };

        function remove(key) {
            $window.localStorage.removeItem(prefix + key);
            return true;
        };

        // function clearAll() {
        // }
    }
}());