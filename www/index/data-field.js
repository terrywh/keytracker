Vue.component("data-field", function(resolve, reject) {
	Vue.http.get("/index/data-field.html").then(function(response) {
		resolve({
			template: response.body,
			props: ["index", "name", "data"],
			data: function() {
				return {"type":"text", "value": this.data};
			},
			created: function() {
				var t = typeof(this.value);
				switch(t) {
					case "boolean":
					case "number":
					case "string":
					this.type = t;
					break;
					default:
					this.type = "unkonwn";
				}
			},
			methods: {
				update: function(e) {
					switch(this.type) {
						case "boolean":
							this.value = e.currentTarget.checked;
						break;
						case "number":
							this.value = parseInt(e.currentTarget.value);
						break;
						case "string":
							this.value = e.currentTarget.value;
						break;
					}
				},
				focus: function() {
					this.$el.querySelector("input").focus();
				},
			}

		});
	});
});
