"use strict";

var World = (function() {
	return {
		update: function() {
			if (!Game.ok()) {
				return;
			}
			//console.log("update");
		},

		draw: function() {
			if (!Game.ok()) {
				return;
			}
			var ctx = Game.getContext();
			var width = Game.getWidth();
			var height = Game.getHeight();

			ctx.clearRect(0, 0, width, height);

			ctx.fillStyle = "#0F0";
			ctx.fillRect(10, 10, 10, 10);
		}
	};
})();
