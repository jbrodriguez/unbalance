'use strict';

angular.module('unbalance'. [
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
})

.controller('AppCtrl', function AppCtrl() {
})

);