(function () {
    'use strict';

    angular
        .module('app')
        .controller('Home', Home)

    /* @ngInject */
    function Home($state, $scope, $rootScope, options, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.options = options;

        // vm.calculateBestFit = calculateBestFit;
        // vm.move = move;        

        activate();

        function activate() {
            return getConfig().then(function() {
                console.log('Activated Home controller');
            })
        };

        function getConfig() {
            return api.getConfig().then(function(data) {
                options.config = data;

                if (options.config.folders.length === 0) {
                    $state.go('settings');
                } else {
                    $state.go('dashboard');
                };
            });
        };
    };

})();