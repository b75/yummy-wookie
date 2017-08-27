"use strict";

var Conn = (function() {
	var state = "initial";
	var conn = null;
	var stateChangeHandler = null;

	var setState = function(newState) {
		if (!newState || typeof newState !== "string") {
			return;
		}
		state = newState;
		if (typeof stateChangeHandler === "function") {
			stateChangeHandler(newState);
		}
	};

	return {
		open: function() {
			if (state !== "initial" && state !== "closed") {
				return;
			}
			setState("checking");

			if (!window["WebSocket"]) {
				setState("failed");
				console.error("websockets not supported");
				return;
			}
			setState("connecting");

			conn = new WebSocket("ws://127.0.0.1:8081/connect");
			conn.onopen = function() {
				conn.send("Hello world!");
				setState("open");
			};
			conn.onclose = function(event) {
				console.log("close", event);
				setState("closed");
			};
			conn.onmessage = function(event) {
				if (!event.isTrusted) {
					console.error("untrusted event", event);
					return;
				}
				World.update(event.data);
			};
			conn.onerror = function(error) {
				console.error("connection error:", error);
				setState("failed");
			};
		},

		close: function() {
			if (state !== "open") {
				return;
			}

			conn.close();
			setState("closed");
			conn = null;
		},

		onStateChange: function(f) {
			if (typeof f !== "function") {
				console.error("Conn.onStateChange called with", typeof f);
				return;
			}

			stateChangeHandler = stateChangeHandler ? stateChangeHandler : f;
		}
	};
})();
