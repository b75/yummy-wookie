"use strict";

var Util = (function() {
	var screenWidth;
	var screenHeight;
	var screenHalfWidth;
	var screenHalfHeight;
	
	return {
		directionPoint: function(x, y, dir, len) {
			return {
				x: x + len * Math.cos(dir),
				y: y + len * Math.sin(dir)
			};
		},

		worldToScreen: function(x, y, fx, fy) {
			return {
				x: screenHalfWidth + x - fx,
				y: screenHalfHeight + y - fy
			};
		},

		setScreenSize: function(width, height) {
			var widthOk = width && typeof width === "number" && width > 0;
			var heightOk = height && typeof height === "number" && height > 0;
			if (!widthOk || !heightOk) {
				console.error("illegal screen size arguments", width, height);
				return;
			}
			screenWidth = width;
			screenHeight = height;
			screenHalfWidth = width / 2;
			screenHalfHeight = height / 2;
		}
	};
})();
