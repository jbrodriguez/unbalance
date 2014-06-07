angular.module('unbalance', [
	'unbalance.services'
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

.controller('AppCtrl', ['core', '$scope', function AppCtrl(core, $scope) {
	$scope.getDisks = function() {
		console.log("modofoco");
		core.getDisks();
	}	
}])

;