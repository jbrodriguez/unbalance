(function () {
    'use strict';

    angular
        .module('app.settings')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('settings', {
                    url: '/settings',
                    templateUrl: 'app/settings/settings.html',
                    controller: 'Settings as vm',
                })            
        });

})();