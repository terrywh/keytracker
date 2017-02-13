Vue.component("data-alert", function(resolve, reject) {
	Vue.http.get("/index/data-alert.html").then(function(response) {
		resolve({
			template: response.body,
			created: function() {
				this.$index = 0;
				this.$prepend = function(type, msg, timeout) {
					var self = this, index = ++this.$index;
					this.data.unshift({"index": index, "type": type, "text": msg});
					setTimeout(function() {
						for(var i=0;i<self.data.length;++i) {
							if(self.data[i].index === index) {
								self.data.splice(i, 1);
								break;
							}
						}
					}, timeout);
				}
			},
			data: function() {
				return {"data":[]};
			},
			methods: {
				success: function(text, timeout) {
					this.$prepend("success", text, timeout || 2000);
				},
				warning: function(text, timeout) {
					this.$prepend("warning", text, timeout || 3000);
				},
				danger: function(text, timeout) {
					this.$prepend("danger", text, timeout || 5000);
				},
			}
		});
	});
});
