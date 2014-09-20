(function () {
    'use strict';

    angular
        .module('app')
        .controller('Home', Home)

    /* @ngInject */
    function Home($state, $scope, $rootScope, options) {

        /*jshint validthis: true */
        var vm = this;

        vm.options = options;
    };

})();