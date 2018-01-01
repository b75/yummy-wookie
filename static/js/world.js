"use strict";

var World = (function() {
	var particles = [];
	var humans = [];
	var actions = [];
	var keyState = {
		wDown: false,
		sDown: false,
		aDown: false,
		dDown: false,
	};
	var focalPoint = {
		x: 0,
		y: 0,
		vx: 0,
		vy: 0
	};
	var followMode = true;
	
	
	return {
		serverUpdate: function(data) {
			if (!data || typeof data !== "string") {
				return;
			}
			particles = [];
			humans = [];

			var objects = data.split(":");
			for (var i = 0; i < objects.length - 1; i++) {
				var parts = objects[i].split(",");
				switch (parts[0]) {
				case "p":
					var particle = World.parseParticle(parts);
					if (!particle) {
						continue;
					}
					particles.push(particle);
					break;
				case "h":
					var human = World.parseHuman(parts);
					if (!human) {
						continue;
					}
					humans.push(human);
					if (followMode) {
						focalPoint.x = human.x;
						focalPoint.y = human.y;
					}
					break;
				default:
					console.error("unknown object type", parts[0]);
					break;
				}
			}
		},

		update: function() {
			if (followMode) {
				focalPoint.vx = 0;
				focalPoint.vy = 0;
			} else {
				if (keyState.wDown) {
					focalPoint.vy -= 5;
				}
				if (keyState.sDown) {
					focalPoint.vy += 5;
				}
				if (keyState.dDown) {
					focalPoint.vx += 5;
				}
				if (keyState.aDown) {
					focalPoint.vx -= 5;
				}
				focalPoint.x += focalPoint.vx;
				focalPoint.y += focalPoint.vy;
				focalPoint.vx *= 0.9;
				focalPoint.vy *= 0.9;
			}
		},

		parseHuman: function(parts) {
			if (typeof parts !== "object") {
				return false;
			}

			if (parts.length !== 4) {
				return false;
			}

			var h = {};
			h.x = parseFloat(parts[1]);
			h.y = parseFloat(parts[2]);
			h.dir = parseFloat(parts[3]);

			if (isNaN(h.x) || isNaN(h.y) || isNaN(h.dir)) {
				return false;
			}

			return h;
		},

		parseParticle: function(parts) {
			if (typeof parts !== "object") {
				return false;
			}

			if (parts.length !== 3) {
				return false;
			}

			var p = {};
			p.x = parseFloat(parts[1]);
			p.y = parseFloat(parts[2]);

			if (isNaN(p.x) || isNaN(p.y)) {
				return false;
			}

			return p;
		},

		draw: function() {
			var ctx = Game.getContext();
			var width = Game.getWidth();
			var height = Game.getHeight();
			var i;
			var tau = 2 * Math.PI;

			ctx.clearRect(0, 0, width, height);

			ctx.strokeStyle = "#FF0000";
			ctx.beginPath();
			for (i = 0; i < humans.length; i++) {
				var p = Util.worldToScreen(humans[i].x, humans[i].y, focalPoint.x, focalPoint.y);
				ctx.arc(p.x, p.y, 10, 0, tau);
				ctx.moveTo(p.x, p.y);
				var p2 = Util.directionPoint(p.x, p.y, humans[i].dir, 10);
				ctx.lineTo(p2.x, p2.y);
			}
			ctx.stroke();

			ctx.fillStyle = "#00FF00";
			for (i = 0; i < particles.length; i++) {
				var p = Util.worldToScreen(particles[i].x, particles[i].y, focalPoint.x, focalPoint.y);
				ctx.fillRect(p.x, p.y, 2, 2);
			}
		},

		action: function(act) {
			if (!act || typeof act !== "object") {
				return;
			}
			// TODO non-keystate related actions
		},

		setKeyState: function(key, state) {
			state = state ? true : false;
			switch (key) {
			case "w":
				keyState.wDown = state;
				break;
			case "a":
				keyState.aDown = state;
				break;
			case "s":
				keyState.sDown = state;
				break;
			case "d":
				keyState.dDown = state;
				break;
			}
		},

		sendActions: function() {
			var batch = actions;
			actions = [];

			if (followMode) {
				if (keyState.wDown) {
					batch.push("mf");
				}
				if (keyState.sDown) {
					batch.push("mb");
				}
				if (keyState.aDown) {
					batch.push("tl");
				}
				if (keyState.dDown) {
					batch.push("tr");
				}
			}

			if (batch.length === 0) {
				return;
			}
			Conn.send(batch.join(":"));
		},

		toggleFollowMode: function() {
			followMode = followMode ? false : true;
		}
	};
})();
