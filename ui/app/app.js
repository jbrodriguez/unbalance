angular.module('unbalance', [
	'unbalance.services',
	'unbalance.models',
	'unbalance.filters'
])

.config(['$provide', function($provide) {
  $provide.decorator('$rootScope', ['$delegate', function($delegate) {
    $delegate.constructor.prototype.$onRootScope = function(name, listener) {
      var unsubscribe = $delegate.$on(name, listener);
      this.$on('$destroy', unsubscribe);
    };

    return $delegate;
  }]);
}])

.run( function run () {
	console.log("Im alive");
})

.controller('AppCtrl', ['core', 'model', '$scope', function AppCtrl(core, model, $scope) {
	$scope.disks = [];

	$scope.getStatus = function() {
		core.getStatus()
			.then(function(data) {
				console.log(data)
				$scope.disks = data.disks.map(function(disk) {
					console.log(disk);
					return new model.Disk(disk)
				});
				console.log($scope.disks)
			});
	};

	var onSocketOpened = function() {
		console.log("modofoco");
		$scope.getStatus();
	};

	$scope.$onRootScope("/api/v1/put/socketOpened", onSocketOpened);
}])

;