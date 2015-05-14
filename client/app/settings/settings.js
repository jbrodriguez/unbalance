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
        vm.submitted = false;

        vm.addFolder = addFolder;
        vm.removeFolder = removeFolder;
        vm.flipNotifications = flipNotifications;
        vm.saveNotifications = saveNotifications;
        vm.submit = submit;

        activate();

        function activate() {
            console.log("behind petrified eyes", options);
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

            return api.saveConfig(vm.options.config).then(function(data) {
                logger.success('Configuration saved succesfully.');
            });
        };

        function removeFolder(index) {
            vm.options.config.folders.splice(index, 1);

            return api.saveConfig(vm.options.config).then(function(data) {
                logger.success('Configuration saved succesfully.');
            });
        };

        function flipNotifications() {
            vm.options.config.notifications != vm.options.config.notifications
            console.log("notifications: " + vm.options.config.notifications);

            return api.saveConfig(vm.options.config).then(function(data) {
                logger.success('Configuration saved succesfully.');
            });
        }

        function saveNotifications() {

            return api.saveConfig(vm.options.config).then(function(data) {
                logger.success('Configuration saved succesfully.');
            });
        } 

        function submit(isValid) {
            vm.submitted = !isValid;

            if (!isValid) {
                console.log("invalid")
            } else {
                return api.saveConfig(vm.options.config).then(function(data) {
                    vm.submitted = false;
                    logger.success('Configuration saved succesfully.');
                });
            };
        }

    }
})();