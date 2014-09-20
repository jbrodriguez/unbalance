(function () {
    'use strict';

    angular
        .module('app.dashboard')
        .controller('Dashboard', Dashboard);

    /* @ngInject */
    function Dashboard($state, $q, api, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.condition = {};
        vm.disks = [];

        vm.toDisk = [];
        vm.fromDisk = [];

        vm.maxFreeSize = 0;
        vm.maxFreePath = 0;             

        vm.getStatus = getStatus;
        vm.getBestFit = getBestFit;

        activate();

        function activate() {
            return getStatus().then(function() {
                logger.info('activated dashboard view');
            });
        };

        function getStatus() {
            return api.getStatus().then(function (data) {
                logger.info('what is: ', data)

                vm.condition = data.condition;

                vm.maxFreeSize = 0;
                vm.maxFreePath = 0;                

                vm.disks = data.disks.map(function(disk) {
                    vm.toDisk[disk.path] = true;
                    vm.fromDisk[disk.path] = false;

                    if (disk.free > vm.maxFreeSize) {
                        vm.maxFreeSize = disk.free;
                        vm.maxFreePath = disk.path;
                    }

                    return disk;
                });

                if (vm.maxFreePath != "") {
                    vm.toDisk[vm.maxFreePath] = false;
                    vm.fromDisk[vm.maxFreePath] = true;
                }                

                return vm.disks;
            });
        };

        function getBestFit() {
            return api.getBestFit().then(function(data) {
                
            });
        }
    }
})();