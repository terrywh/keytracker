Vue.component("index-form",function(resolve, reject) {
	Vue.http.get("/index-form.html").then(function(response) {
		resolve({
			template: response.body,
			data: function() {
				return {
					k: "",
					v: null,
					x: 2,
					type:  "删除",
					value: "",
				};
			},
			// computed: {
			// 	type: function() {
			// 		switch(typeof this.v) {
			// 		case "string":
			// 			return "文本";
			// 		case "number":
			// 			return "数值";
			// 		case "boolean":
			// 			return "布尔";
			// 		}
			// 		return "删除";
			// 	}
			// },
			methods: {
				set: function(key, val) {
					this.k = key;
					this.v = val;
					this.typeSet();
					this.value = this.v.toString();

					var $el = this.$el.querySelector("[name=val]");
					$el.focus();
					$el.select();
				},
				onSubmit: function(e) {
					this.$session.set(this.k, this.v, this.x);
					this.value = "";
					this.v = null;
					this.typeSet();
					var $el = this.$el.querySelector("[name=val]");
					$el.focus();
					$el.select();
				},
				typeSet: function() {
					switch(typeof this.v) {
					case "string":
						this.type = "文本";
						break;
					case "number":
						this.type = "数值";
						break;
					case "boolean":
						this.type = "布尔";
						break;
					default:
						this.type = "删除";
					}
				},
				debouncePersistent: function() {
					var self = this;
					clearTimeout(this.$dpTimeout);
					this.$dpTimeout = setTimeout(function() {
						self.x = self.x | 0x02;
					}, 50)
				},
				debounceSuffix: function() {
					var self = this;
					clearTimeout(this.$dsTimeout);
					this.$dsTimeout = setTimeout(function() {
						self.x = self.x | 0x04;
					}, 50);
				},
			},
			watch: {
				"k": function() {
					this.k = "/" + this.k.trim().split("/").filter(function(v) { return v!=""; }).join("/");
				},
				"value": function(n, o) {
					if(n.trim() == n) {
						if(n == "") this.v = null;
						else if(n.match(/^[+\-]?(\d+|0x[0-9a-fA-F]+)$/)) this.v = parseInt(n);
						else if(n.match(/^[+\-]?\d+\.\d*$/)) this.v = parseFloat(n);
						else if(n.match(/^ok|yes|on|true$/i)) this.v = true;
						else if(n.match(/^no|off|false$/i)) this.v = false;
						else this.v = n;
						this.typeSet()
					}else{
						this.value = n.trim();
					}
				},
			}
		});
	}, reject);
});
