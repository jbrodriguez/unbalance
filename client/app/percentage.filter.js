(function () {
    'use strict';

    angular
        .module('app')
        .filter('percentage', Percentage);

    function Percentage() {
		return function(input, decimals, suffix) {
			decimals = decimals || 2;
			suffix = suffix || '%';
			
			// if ($window.isNaN(input)) {
			// 	return '';
			// }

			return Math.round(input * Math.pow(10, decimals + 2))/Math.pow(10, decimals) + suffix
		}
    };    

})();