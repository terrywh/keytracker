Vue.component("router-list", function(resolve, reject) {
	Vue.http.get("/index/router-list.html").then(function(response) {
		resolve({
			template: response.body,
			created: function() {
				this.reload();
			},
			computed: {
				data: function() {
					var data = {}, r = [];
					this.$app.list.forEach(function(item) {
						return data[item.ns] = data[item.ns] ? ++data[item.ns] : 1;
					});
					for(var ns in data) {
						r.push({name: ns, size: data[ns]});
					}
					return r;
				}
			},
			methods: {
				reload: function() {
					this.$app.list = [];
					this.$http.get("/data").then(function(response) {
						return response.json();
					}).then(function(data) {
						this.$app.list = data;
					});
				}
			}
		});
	});
});
