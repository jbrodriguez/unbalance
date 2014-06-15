angular.module('unbalance.models', [
])

.factory('model', function() {
	var Disk = function(json) {
		var self = this;

		angular.extend(self, json)
		// self.id = json.id;
		// self.name = json.name;
		// self.path = json.path;
		// self.device = json.device;
		// self.free = json.free;
		// self.size = json.size;
		// self.serial = json.serial;
		// self.status
	}

	return {
		Disk: Disk
	}
})