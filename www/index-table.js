Vue.component("index-table",function(resolve, reject) {
	Vue.http.get("/index-table.html").then(function(response) {
		resolve({
			created: function() {
				var $app = this.$app,
					$session = this.$session,
					self = this;

				$session.ondata = function(item) {
					// 下属逻辑针对初始数据有序的情况编写逻辑
					let i;
					for(i=0;i<self.data.length;++i) {
						if(self.data[i].k == item.k) {
							self.data.splice(i, 1);
							break;
						}else if(self.data[i].k > item.k) {
							break;
						}
					}
					item.blink = true;
					setTimeout(function() {
						item.blink = null;
						setTimeout(function() {
							delete item.blink;
						}, 0);
					}, 1000);
					if(item.v !== null) self.data.splice(i, 0, item);
				};
				$session.ready.then(function() {
					$app.path.split('{|}').forEach(function(v) {
						$session.watch(v, true);
					});
				});
			},
			template: response.body,
			data: function() {
				return {
					data: [],
				};
			},
			methods: {
				type: function(value) {
					switch(typeof value) {
					case "string":
						return "文本";
					case "number":
						return "数值";
					case "boolean":
						return "布尔";
					}
					return "对象";
				}
			},
		});
	}, reject);
});
