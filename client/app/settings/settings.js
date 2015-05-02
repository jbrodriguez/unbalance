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

        vm.addFolder = addFolder;
        vm.removeFolder = removeFolder;

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
    }
})();