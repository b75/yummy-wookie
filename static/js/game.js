"use strict";

var Game = (function() {
	var canvas;
	var ctx;
	var statusElem;

	var ok = false;
	var width;
	var height;

	return {
		setCanvas: function(elem) {
			if (canvas) {
				console.error("canvas already set");
				return;
			}

			if (elem.length !== 1) {
				console.error("canvas element length", elem.length);
				return;
			}
			canvas = elem;
			ctx = elem.get(0).getContext("2d");
			if (!ctx) {
				console.error("error setting context");
				return;
			}

			width = elem.prop("width");
			if (typeof width !== "number" || width <= 0) {
				console.error("canvas width error");
				return;
			}
			height = elem.prop("height");
			if (typeof height !== "number" || height <= 0) {
				console.error("canvas height error");
				return;
			}

			console.log(width, "x", height);

			ok = true;
		},

		setStatusElem: function(elem) {
			if (!elem || typeof elem !== "object" || elem.length !== 1) {
				return;
			}
			statusElem = statusElem ? statusElem : elem;
		},

		setStatus: function(status) {
			if (typeof status !== "string" || !statusElem) {
				return;
			}
			statusElem.html(status);
		},

		getContext: function() {
			return ctx;
		},

		getWidth: function() {
			return width;
		},

		getHeight: function() {
			return height;
		},

		ok: function() {
			return ok;
		}
	};
})();

$(function() {
	Game.setCanvas($("#map"));
	Game.setStatusElem($("#status"));
	Conn.onStateChange(function(newState) {
		Game.setStatus("connection state " + newState);
	});

	$("#close-connection").on("click", function(event) {
		Conn.close();
	});

	$("#open-connection").on("click", function(event) {
		Conn.open();
	});

	if (Game.ok()) {
		setInterval(function() {
			World.draw();
		}, 200);
	}
});
