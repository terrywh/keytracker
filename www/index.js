(function(exports) {
	var session = createSession();
	var app = new Vue({
		el: "#app",
		created: function() {
			this.navigate(this.path);
		},
		data: {
			columns: [],
			path: location.hash.substr(1),
		},
		methods: {
			navigate: function(path) {
				this.path = path;
				location.hash = "#" + path;
				path = path.split("/");
				var i;
				for(i=0;i<path.length-1;++i) {
					if(!this.columns[i]) {
						this.columns.push([]);
						session.list(i == 0 ? "/" : path.slice(0, i+1).join("/"));
					}
				}
				if(!this.columns[i]) {
					this.columns.push([]);
				}else{
					this.columns[i].splice(0, this.columns[i].length);
				}
				this.columns = this.columns.splice(0,path.length)
				console.log(this.consolelumns,path.length)
				session.list(i == 0 ? "/" : path.slice(0, i+1).join("/"));
			}
		}
	});
	var sortTimeout, sortData = [];
	var autoPosition = function(){
		document.getElementById('app').querySelectorAll('.column').forEach(function(el){
			if(el.offsetHeight<window.innerHeight){
				el.classList.add('auto-top');
			}else{
				el.classList.remove('auto-top');
			}
		})
	};
	session.ondata = function(data) {
		sortData.push(data);
		var sortDataLen = sortData.length;
		clearTimeout(sortTimeout);
		sortTimeout = setTimeout(function() {
			var item, index, sorted = {};
			while(item = sortData.pop()) {
				index = item.k.split("/").length - 2;
				app.columns[index].push(item);
				sorted[index] = true;
			}
			for(index in sorted) {
				app.columns[parseInt(index)].sort(function(a, b) {
					// 由于 Go map 遍历会乱序，这里单独对 key 进行排序
					return a.k < b.k ? -1 : a.k == b.k ? 0 : 1;
				});
			}
			setTimeout(autoPosition, 50);
		}, 200);
	};
	Vue.mixin({
		beforeCreate: function() {
			this.$app = app;
		}
	});

	exports.$app = app;
})(window);
