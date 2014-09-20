(function () {
    'use strict';

    angular
        .module('app')
        .filter('replace', Replace);

    function Replace() {
        return function(text, fromString, toString) {
            return text.split(fromString).join(toString);
        }
    };    

})();