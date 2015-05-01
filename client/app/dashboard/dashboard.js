(function () {
    'use strict';

    angular
        .module('app.dashboard')
        .controller('Dashboard', Dashboard);

    /* @ngInject */
    function Dashboard($scope, $state, $q, api, socket, logger, options) {

        /*jshint validthis: true */
        var vm = this;

        vm.options = options;

        vm.condition = {};
        vm.disks = [];
        vm.commands = [];

        vm.toDisk = [];
        vm.fromDisk = [];

        vm.maxFreeSize = 0;
        vm.maxFreePath = 0;             

        vm.getStatus = getStatus;
        vm.calculateBestFit = calculateBestFit;
        vm.move = move;
        vm.checkFrom = checkFrom;
        vm.checkTo = checkTo;
        vm.flipDryRun = flipDryRun;

        vm.moveStarted = false
        vm.moveInProgress = false;
        vm.lines = [];

        socket.register("storage:move:begin", storageMoveBegin);
        socket.register("storage:move:progress", storageMoveProgress);
        socket.register("storage:move:end", storageMoveEnd);

        activate();

        function activate() {
            vm.moveInProgress = false;
            vm.moveStarted = false;

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

        function calculateBestFit() {
            var srcDisk = "";

            for (var key in vm.fromDisk) {
                if (vm.fromDisk.hasOwnProperty(key)) {
                    if (vm.fromDisk[key]) {
                        srcDisk = key;
                        break;
                    }
                }
            }

            console.log("src: " + srcDisk);

            if (srcDisk === "") {
                alert("I won't take that !");
                return;
            }

            return api.calculateBestFit({"sourceDisk": srcDisk, "destDisk": ""}).then(function(data) {
                vm.condition = data.condition;

                // vm.maxFreeSize = 0;
                // vm.maxFreePath = 0;                

                // vm.disks = data.disks.map(function(disk) {
                //     vm.toDisk[disk.path] = true;
                //     vm.fromDisk[disk.path] = false;

                //     if (disk.free > vm.maxFreeSize) {
                //         vm.maxFreeSize = disk.free;
                //         vm.maxFreePath = disk.path;
                //     }

                //     return disk;
                // });

                // // if (vm.maxFreePath != "") {
                // //     vm.toDisk[vm.maxFreePath] = false;
                // //     vm.fromDisk[vm.maxFreePath] = true;
                // // }
                vm.disks = data.disks;

                return vm.disks;                
            });
        };

        function move() {
            socket.send("storage:move");

            // return api.move().then(function(data) {
            //     vm.commands = data;
            //     logger.info("Scroll down to see list of commands");
            //     return vm.commands;
            // });
        };

        function checkFrom(from) {
            console.log("something changed: " + from);
            for (var key in vm.fromDisk) {
                if (key !== from) {
                    vm.fromDisk[key] = false;
                }
            };

            for (var key in vm.toDisk) {
                vm.toDisk[key] = !(key === from);
            };            

//            vm.toDisk[from] = false;
        };

        function checkTo(to) {
            return;
        };

        function flipDryRun() {
            vm.options.config.dryRun != vm.options.config.dryRun;

            console.log('vm.options.config.dryRun: ' + vm.options.config.dryRun);

            return api.saveConfig(vm.options.config).then(function(data) {
                logger.success('config saved succesfully');
            });            
        };

        function storageMoveBegin(data) {
            vm.moveInProgress = true;
            vm.moveStarted = true;
        };

        function storageMoveProgress(data) {
            vm.lines.push(data);
        };

        function storageMoveEnd(data) {
            vm.moveInProgress = false;
        };
    }
})();