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
			Util.setScreenSize(width, height);

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

	window.addEventListener("keypress", function(event) {
		switch (event.code) {
		case "KeyF":
			World.toggleFollowMode();
		}
	});

	window.addEventListener("keydown", function(event) {
		switch (event.code) {
		case "KeyW":
			World.setKeyState("w", true);
			break;
		case "KeyA":
			World.setKeyState("a", true);
			break;
		case "KeyS":
			World.setKeyState("s", true);
			break;
		case "KeyD":
			World.setKeyState("d", true);
			break;
		}
	});

	window.addEventListener("keyup", function(event) {
		switch (event.code) {
		case "KeyW":
			World.setKeyState("w", false);
			break;
		case "KeyA":
			World.setKeyState("a", false);
			break;
		case "KeyS":
			World.setKeyState("s", false);
			break;
		case "KeyD":
			World.setKeyState("d", false);
			break;
		}
	});

	if (Game.ok()) {
		setInterval(World.sendActions, 20);
		setInterval(World.update, 20);
		setInterval(World.draw, 20);
	}
});
