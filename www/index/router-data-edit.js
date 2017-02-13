Vue.component("router-data-edit", function(resolve, reject) {
	Vue.http.get("/index/router-data-edit.html").then(function(response) {
		resolve({
			template: response.body,
			data: function() {
				return {data: []};
			},
			beforeMount: function() {
				this.$retry();
			},
			created: function() {
				this.$find = function(data, key) {
					for(var i=0;i<data.length;++i) {
						if(data[i].key == key) {
							return data[i];
						}
					}
					return null;
				};
				this.$load = function(c) {
					var r = [],
						d = Object.assign({}, c); // 复制一份，不改变原始值

					delete d.remote_addr;
					delete d.key;
					var keys = Object.keys(d).sort(); // 稳定参数顺序
					keys.forEach(function(k) {
						if(typeof(d[k]) === "object") return;
						r.push({"key": k, "val": d[k]});
					});
					this.data = r;
					console.log(this.data);
				};
				this.$retry = function() {
					var c = this.$find(this.$app.data, this.$route.params.key), self = this;
					console.log(c)
					if(c) {
						this.$load(c);
					}else if(this.$app.data.length > 0) {
						this.$app.$refs.alert.warning("节点未找到。");
						setTimeout(function() {
							self.$router.push('/namespace/' + self.$route.params.ns);
						}, 200);
					}else{
						setTimeout(this.$retry.bind(this), 100);
					}
				}
			},
			methods: {
				confirm: function() {
					var c = {}, d = this.$find(this.$app.data, this.$route.params.key);

					for(var i=0;i<this.$children.length;++i) {
						if(this.$children[i].data !== this.$children[i].value) {
							c[this.$children[i].name] = this.$children[i].value;
							d[this.$children[i].name] = this.$children[i].value;
						}
					}
					this.$http.post("/update/" + this.$route.params.ns + "/" + this.$route.params.key, JSON.stringify(c)).then(function(response) {
						this.$app.$refs.alert.success("设置将会被推送到该节点。");
						this.$router.push('/namespace/' + this.$route.params.ns);
					}, function() {
						this.$app.$refs.alert.danger("设置将会被丢弃，请尝试刷新页面。");
					});
				},
				cancel: function() {
					this.$router.push('/namespace/' + this.$route.params.ns);
				}
			}
		});
	});
});
