(function(exports) {
	var app = new Vue({
		el: "#app",
		data: {
			path: "/" + location.hash.substr(1).split("/").filter(function(v) {
				return v != "";
			}).join("/"),
		},
		methods: {
			onEdit: function(key, val) {
				this.$refs.editor.set(key, val);
				console.log(key, val);
			}
		}
	});
	var session = createSession();
	Vue.mixin({
		beforeCreate: function() {
			this.$app = app;
			this.$session = session;
		}
	});

	exports.$app = app;
})(window);
