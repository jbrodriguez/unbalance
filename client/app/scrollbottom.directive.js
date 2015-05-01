(function () {
    'use strict';

    angular
        .module('app')
        .directive('unbScrollBottom', scrollBottom);

    function scrollBottom() {
    	return {
    		scope: {
    			unbScrollBottom: "="
    		},
    		link: function(scope, element) {
    			scope.$watchCollection('unbScrollBottom', function(newValue) {
    				if (newValue) {
    					$(element).scrollTop($(element)[0].scrollHeight);
    				}
    			});
    		}
    	}
    };    

})();