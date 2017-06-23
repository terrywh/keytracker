(function(exports) {
	// require("WebSocket")
	if (typeof Object.assign != 'function') {
		Object.assign = function(target, varArgs) { // .length of function is 2
			'use strict';
			if (target == null) { // TypeError if undefined or null
				throw new TypeError('Cannot convert undefined or null to object');
			}

			var to = Object(target);

			for (var index = 1; index < arguments.length; index++) {
				var nextSource = arguments[index];

				if (nextSource != null) { // Skip over if undefined or null
					for (var nextKey in nextSource) {
						// Avoid bugs when hasOwnProperty is shadowed
						if (Object.prototype.hasOwnProperty.call(nextSource, nextKey)) {
							to[nextKey] = nextSource[nextKey];
						}
					}
				}
			}
			return to;
		};
	}
	function createWebSocket(options, session, resolve, reject) {
		var websocket = new WebSocket("ws://" + options.addr + ":" + options.port + "/session"),
			closing = false;

		websocket.onopen = function(e) {
			console.log("%c[session] %cconnection established.", "color: #888;", "color: #aaa;");
			var obj;
			while(obj = session._watcher.shift()) {
				websocket.send(JSON.stringify(obj)+"\n");
			}
			while(obj = session._cache.shift()) {
				websocket.send(JSON.stringify(obj)+"\n");
			}

			resolve(session);
		};
		websocket.onerror = function(e) {
			reject(e);
		};
		websocket.onclose = function() {
			if(!closing) {
				console.log("%c[session] %cconnection retry ...", "color: #666;", "color: #aaa;");
				session.onretry && session.onretry.call(session);
				setTimeout(createWebSocket.bind(null, options, session), 10000);
			}else{
				console.log("%c[session] %cconnection closed.", "color: #666;", "color: #aaa;");
				session.onclose && session.onclose.call(session);
			}
		};
		websocket.onmessage = function(e) {
			var data;
			try{
				data = JSON.parse(e.data)
			}catch(e){
				return;
			}
			if(!data.k) return;

			console.log("%c[session] %cdata recevied: %c"+e.data, "color: #666;", "color: #aaa;", "color: #888;");
			session.ondata && session.ondata.call(session, data, data.y == 1 ? true: false);
		};

		session.close = function() {
			closing = true;
			console.log("%c[session] %cconnection closing ...", "color: #666;", "color: #aaa;");
			websocket.close();
		};
		session._write = function(obj) {
			if(websocket.readyState == 1) {
				websocket.send(JSON.stringify(obj) + "\n");
			} else {
				session._cache.push(obj);
			}
		};
	}
	var optionDefaults = {
		"addr": location.hostname,
		"port": parseInt(location.port) || 7472,
	};
	function createSession(options) {
		options = Object.assign({}, optionDefaults, options);
		var session = {};
		session._cache   = [];
		session._watcher = [];
		session.watchSelf = function(key) {
			session._watcher.unshift({"k":key,"v":true,"x":257});
			session._write(session._watcher[0]);
		}
		session.watch   = function(key, r) {
			session._watcher.unshift({"k":key,"v":true,"x":r?258:256});
			session._write(session._watcher[0]);
			return session;
		};
		session.unwatch = function(key, val, tmp) {
			for(var i=0;i<session._watcher.length;++i) {
				if(session._watcher[i].k == key) {
					session._watcher.splice(i, 1);
					break;
				}
			}
			session._write({"k":key,"v":null,"x":256});
			return session;
		};
		session.get   = function(key) {
			session._write({"k":key,"v":null,"x":512});
			return session;
		};
		session.list  = function(key, r) {
			session._write({"k":key,"v":null,"x":r?514:513});
			return session;
		};
		session.set = function(key, val, x) {
			if(x === undefined) x = 1; // 默认临时非后缀数据
			session._write({"k":key,"v":val,"x":x});
			return session;
		};
		session.remove = function(key) {
			session._write({"k":key,"v":null,"x":0});
			return session;
		};
		session.ready = new Promise(function(resolve, reject) {
			createWebSocket(options, session, resolve, reject)
		});
		return session;
	}


	exports.createSession = createSession
})(typeof window === "undefined" ? module.exports : window);
