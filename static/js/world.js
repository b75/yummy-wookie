"use strict";

var World = (function() {
	var particles = [];
	
	
	return {
		update: function(data) {
			if (!data || typeof data !== "string") {
				return;
			}
			particles = [];

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
				default:
					console.error("unknown object type", parts[0]);
					break;
				}
			}
		},


		parseParticle: function(parts) {
			if (typeof parts !== "object") {
				return false;
			}

			if (parts.length !== 5) {
				return false;
			}

			var p = {};
			p.x = parseFloat(parts[1]);
			p.y = parseFloat(parts[2]);
			p.vx = parseFloat(parts[3]);
			p.vy = parseFloat(parts[4]);

			if (isNaN(p.x) || isNaN(p.y) || isNaN(p.vx) || isNaN(p.vy)) {
				return false;
			}

			return p;
		},

		draw: function() {
			var ctx = Game.getContext();
			var width = Game.getWidth();
			var height = Game.getHeight();
			var i;

			ctx.clearRect(0, 0, width, height);

			ctx.fillStyle = "#0F0";
			for (i = 0; i < particles.length; i++) {
				ctx.fillRect(particles[i].x, particles[i].y, 2, 2);
			}
		}
	};
})();
