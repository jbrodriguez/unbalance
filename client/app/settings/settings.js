(function () {
    'use strict';

    angular
        .module('app.settings')
        .controller('Settings', Settings);

    /* @ngInject */
    function Settings($state, $scope, $timeout, api, options, logger) {

        /*jshint validthis: true */
        var vm = this;

        vm.options = options;
        vm.folder = '';
        vm.context = {};
        vm.running = false;

        vm.addFolder = addFolder;
        vm.removeFolder = removeFolder;
        vm.importer = importer;

        activate();

        function activate() {
            // console.log("behind petrified eyes", options);
        };

        function addFolder() {
            if (vm.folder === '') {
                return;
            };

            if (vm.options.config.folders.indexOf(vm.folder) != -1) {
                logger.warning('Folder already selected');
                return;
            }

            vm.options.config.folders.push(vm.folder);

            console.log('vm.options.config.folders: ' + vm.options.config.folders);
            console.log('options.config.folders: ' + options.config.folders);

            return api.saveConfig(vm.options.config).then(function(data) {
                logger.success('config saved succesfully');
            });
        };

        function removeFolder(index) {
            vm.options.config.folders.splice(index, 1);

            return api.saveConfig(vm.options.config).then(function(data) {
                logger.success('config saved succesfully');
            });
        };

        function importer() {
            return startImport().then(function() {
                logger.info('started import function');
                update();
            });
        };

        function startImport() {
            return api.startImport().then(function (data) {
                vm.context = null;
                vm.context = data;
                vm.running = true;
                return vm.context;
            });
        };

        function update() {
            getStatus();
            if (!vm.context.completed) {
                schedule(update, 1000);
            } else {
                vm.running = false;
                $state.go('cover');
            };
        };

        function getStatus() {
            return api.getStatus().then(function (data) {
                vm.context = null;
                vm.context = data;
                return vm.context;
            });
        };        

        function schedule(fn, delay) {
            var promise = $timeout(fn, delay);
            var deregister = $scope.$on('$destroy', function() {
                $timeout.cancel(promise);
            });
            promise.then(deregister);
        };        

    }
})();