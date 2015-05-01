(function () {
    'use strict';

    angular
        .module('app.dashboard')
        .config(function($stateProvider, $urlRouterProvider) {
            $stateProvider
                .state('dashboard', {
                    url: '/dashboard',
                    templateUrl: 'app/template/main.html',
                    controller: 'Dashboard as vm',
                })            
        });

})();