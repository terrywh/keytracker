Vue.component("data-field-adder", function(resolve, reject) {
	Vue.http.get("/index/data-field-adder.html").then(function(response) {
		resolve({
			template: response.body,
			data: function() {
				return {"key":""};
			},
			mounted: function() {
				this.$input = this.$el.querySelector("input");
				this.$input.setCustomValidity("选项 KEY 必须存在，且仅允许使用 大小写字母、数字及下划线！");
			},
			methods: {
				trigger: function(val) {
					if(!this.key.match(/^[a-zA-Z0-9_]+$/)) {
						this.$input.reportValidity();
						return;
					}
					this.$emit("append", this.key, val);
					this.key = "";
				},
			}

		});
	});
});
