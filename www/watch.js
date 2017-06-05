(function(exports) {
	// 应用
	var app = new Vue({
		el: "#app",
		data: {
			"dataList": [],
			"watcherAppend": "",
			"watcherList": [],
			"dataAppend": {
				"path": "",
				"value": "",
				"persistent": false,
				"suffix": false,
				"type": "文本",
			},
		},
		beforeCreate: function() {

			this.$watchers   = [];
			var self = this;
			location.hash.substr(1).split("{|}").forEach(function(path) {
				if(self.$watchers.indexOf(path) === -1) {
					self.$watchers.push(path);
				}
			});
		},
		created: function() {
			var self = this;
			this.$session = createSession();
			this.$session.ondata = function(data) {
				if(data.v === null) return;
				
				self.dataSet({
					"path": data.k,
					"value": data.v,
					"type": translateType(data.v),
				});
			};
			this.$session.onready = function() { // 重连也会触发，重新建立监测
				self.$watchers.forEach(function(path) {
					self.$session.watch(path);
				});
			};
		},
		methods: {
			dataSet: function(item) {
				for(i=0;i<this.dataList.length;++i) {
					if(this.dataList[i].path == item.path) {
						Object.assign(this.dataList[i], item);
						break;
					}
				}
				if(i==this.dataList.length) {
					this.dataList.unshift(item);
				}
			},
			dataDel: function(item) {
				for(i=0;i<this.dataList.length;++i) {
					if(this.dataList[i].path == this.dataAppend.path) {
						this.dataList.splice(i, 1);
						break;
					}
				}
			},
			dataAdd: function() {
				switch(this.dataAppend.type) {
					case "文本":
						this.dataAppend.value = normalizeString(this.dataAppend.value);
					break;
					case "数值":
						this.dataAppend.value = normalizeNumber(this.dataAppend.value);
					break;
					case "布尔":
						this.dataAppend.value = normalizeBoolean(this.dataAppend.value);
					break;
				}
				var item = {
					path : this.dataAppend.path,
					value: this.dataAppend.value,
					type : this.dataAppend.type,
				}, i;

				if(item.value != null) { // 添加、更改
					this.dataSet(item);
				}else{ // 删除
					this.dataDel(item);
				}
				var x = 1;
				if(this.dataAppend.persistent) x+=2;
				if(this.dataAppend.suffix) x+=4;
				this.$session.set(this.dataAppend.path, this.dataAppend.value, x);

				this.dataAppend.type = "文本";
				this.dataAppend.path = "";
				this.dataAppend.value = "";
				this.dataAppend.persistent = false;
			},
			watcherAdd: function() {
				// 添加监控
				if(this.watcherAppend != "" && this.$watchers.indexOf(this.watcherAppend) == -1) {
					this.$session.watch(this.watcherAppend);
					this.$watchers.push(this.watcherAppend);
					location.hash = "#" + this.$watchers.join("{|}");
				}

				this.watcherAppend = "";
			}
		},
		watch: {
			"dataAppend.path": function() {
				this.dataAppend.path = normalizePath(this.dataAppend.path);
			},
			"dataAppend.value": function() {
				this.dataAppend.type = detectType(this.dataAppend.value);
			},
			"watcherAppend": function() {
				if(!this.watcherAppend) return;
				this.watcherAppend = normalizePath(this.watcherAppend);
			}
		},
	});
	function normalizePath(path) {
		path = path.split("/").join("/");
		if(path[0] != "/") path = "/" + path;
		return path;
	}
	function normalizeString(value) {
		return value.toString().trim() || null;
	}
	function normalizeNumber(value) {
		if(value.toString().indexOf(".") > -1) {
			return parseFloat(value);
		}else{
			return parseInt(value);
		}
	}
	function normalizeBoolean(value) {
		switch(value.toString().toLowerCase()) {
			case "ok":
			case "on":
			case "1":
			case "yes":
			case "true":
			case "done":
				return true;
			default:
				return false;
		}
	}
	function detectType(value) {
		value = value.toString().trim();
		if(value.match(/^[+\-]?\d+(\.\d*)?$/)) {
			return "数值";
		}
		switch(value.toLowerCase()) {
		case "ok":
		case "on":
		case "off":
		case "no":
		case "yes":
		case "true":
		case "false":
		case "done":
		case "fail":
			return "布尔";
		}
		return "文本";
	}
	function translateType(value) {
		switch(typeof value) {
		case "string":
			return "文本";
		case "number":
			return "数值";
		case "boolean":
			return "布尔";
		}
		// 理论上不会触发
		return "未知";
	}
	function initDropdownMenu() {
		window.addEventListener("click", function() {
			var menus = document.querySelectorAll(".dropdown-menu");
			menus.forEach(function(menu) {
				if(menu.parentElement.classList.contains("show")) {
					menu.parentElement.classList.remove("show");
				}
			});
		},false);
		var toggles = document.querySelectorAll(".dropdown-toggle");
		toggles.forEach(function(toggle) {
			toggle.addEventListener("click", function(e) {
				e.stopPropagation();
				toggle.parentElement.classList.add("show");
			});
		});
	}
	initDropdownMenu();
})(window);
