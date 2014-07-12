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

	$scope.toDisk = {};
	$scope.fromDisk = {};

	$scope.size = 0;
	$scope.free = 0;

	maxFreeSize = 0;
	maxFreePath = "";

	$scope.getStatus = function() {
		core.getStatus()
			.then(function(data) {
				console.log(data)

				$scope.size = data.box.size;
				$scope.free = data.box.free;
				$scope.newFree = data.box.newFree;

				$scope.disks = data.disks.map(function(disk) {
					console.log(disk);

					$scope.toDisk[disk.path] = true;
					$scope.fromDisk[disk.path] = false;

					console.log(disk.free, " > ", maxFreeSize);

					if (disk.free > maxFreeSize) {
						maxFreeSize = disk.free;
						maxFreePath = disk.path;
					}

					return new model.Disk(disk)
				});

				console.log("mother: ", maxFreePath);

				if (maxFreePath != "") {
					console.log('marrano');
					$scope.toDisk[maxFreePath] = false;
					$scope.fromDisk[maxFreePath] = true;
				}

				console.log($scope.disks)
			});
	};

	$scope.getBestFit = function() {
		fromDisk = "";

		for (var key in $scope.fromDisk) {
			if ($scope.fromDisk.hasOwnProperty(key)) {
				if ($scope.fromDisk[key]) {
					fromDisk = key;
					break;
				}
			}
		}

		if (fromDisk == "") {
			alert("I won't take that !");
			return;
		}

		core.getBestFit(fromDisk, "")
			.then(function(data) {
				console.log(data)

				$scope.size = data.box.size;
				$scope.free = data.box.free;
				$scope.newFree = data.box.newFree;

				$scope.disks = data.disks.map(function(disk) {
					console.log(disk);
					return new model.Disk(disk)
				});

				console.log($scope.disks)				
			});
	}

	$scope.move = function() {
		core.move();
	}
	
	// event handlers
	var onSocketOpened = function() {
		console.log("modofoco");
		$scope.getStatus();
	};

	$scope.$onRootScope("/api/v1/put/socketOpened", onSocketOpened);
}])

;