Vue.component("router-data", function(resolve, reject) {
	Vue.http.get("/index/router-data.html").then(function(response) {
		resolve({
			template: response.body,
			created: function() {
				this.$reload = function() {
					this.$http.get("/data/" + this.$route.params.ns).then(function(response) {
						return response.json();
					}).then(function(data) {
						this.$app.data = data;
					});
					this.reload > 0 && (this.$reloadTimeout = setTimeout(this.$reload.bind(this), this.reload * 1000));
				};
				this.$app.data.length || this.$reload();
			},
			destroyed: function() {
				clearTimeout(this.$reloadTimeout);
			},
			data: function() {
				return {"reload": 0};
			},
			watch: {
				reload: function() {
					clearTimeout(this.$reloadTimeout);
					this.$reloadTimeout = setTimeout(this.$reload.bind(this), 500);
				},
			},
			methods: {
				format: function(item) {
					var r = Object.assign({}, item), h = '{ ', self = this;
					delete r.remote_addr;
					delete r.key;
					var keys = Object.keys(r).sort(); // 稳定参数顺序
					keys.forEach(function(k) {
						if(h !== '{ ') h += ', '
						h += '<span class="hljs-attr">"' + k + '"</span>: ' + self.highlight(r[k]);
					});
					h += ' }';
					return h;
				},
				highlight: function(v) {
					switch(typeof v) {
						case "number":
						return '<span class="hljs-number">'+v+'</span>';
						break;
						case "string":
						return '<span class="hljs-string">"'+v+'"</span>';
						break;
						case "boolean":
						return '<span class="hljs-literal">'+v+'</span>';
						break;
						default:
						return '<span class="hljs-comment">' + JSON.stringify(v)+ '</span>'
					}
				}
			}
		});
	});
});
