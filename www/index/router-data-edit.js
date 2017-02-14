Vue.component("router-data-edit", function(resolve, reject) {
	Vue.http.get("/index/router-data-edit.html").then(function(response) {
		resolve({
			template: response.body,
			data: function() {
				return {"data": []};
			},
			beforeMount: function() {
				this.$retry();
			},
			created: function() {
				this.$find = function(s) {
					for(var i=0;i<this.$app.data.length;++i) {
						if(this.$app.data[i].key == this.$route.params.key) {
							if(s) {
								this.$app.data.splice(i, 1, s);
								break;
							}else{
								return this.$app.data[i];
							}
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
				};
				this.$retry = function() {
					var c = this.$find(), self = this;
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
				};
				this.$removed = {};
			},
			methods: {
				confirm: function() {
					var c = {}, n = 0, d = this.$find(), r = {};
					for(var key in this.$removed) {// 删除的项
						++n;
						c[key] = null;
					}
					// transition-group 变为一个元素
					this.$children[0].$children.forEach(function(f) {
						if(f.data != f.value) {
							++n;
							c[f.name] = f.value;
						}
					});
					if(n > 0) {
						this.$http.post("/update/" + this.$route.params.ns + "/" + this.$route.params.key, JSON.stringify(c)).then(function(response) {
							this.$app.$refs.alert.success("共计变更 "+n+" 项设置将会被推送到该节点。");
							r = Object.assign({}, d);
							for(var key in c) {
								if(c[key] === null) {
									delete r[key]
								}else{
									r[key] = c[key];
								}
							}
							this.$find(r);
							this.$router.push('/namespace/' + this.$route.params.ns);
						}, function() {
							this.$app.$refs.alert.danger("未能设置，请重试！");
						});
					}else{
						this.$router.push('/namespace/' + this.$route.params.ns);
					}
					this.$removed = {};
				},
				append: function(key, val) {
					var self = this;
					this.data.push({"key": key, "val": val});
					setTimeout(function() {
						var dfs = self.$children[0].$children;
						// transition-group 变为一个元素
						dfs[dfs.length-1].focus();
					}, 150);
				},
				remove: function(index) {
					var item = this.data.splice(index, 1)[0];
					this.$removed[item.key] = null;
				},
				cancel: function() {
					this.$router.push('/namespace/' + this.$route.params.ns);
				}
			}
		});
	});
});
