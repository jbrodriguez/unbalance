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

        activate();

        function activate() {
            return getConfig().then(function() {
                logger.info('initialized state');
            })
        };

        function getConfig() {
            return api.getConfig().then(function(data) {
                vm.options.config = data;

                if (vm.options.config.mediaPath === []) {
                    $state.go('settings');
                } else {
                    $state.go('dashboard');
                };
            });
        };        
    };

})();