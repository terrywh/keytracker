(function(exports) {
	// 路由
	var router = new VueRouter({
		"routes": [
			{
				path: "/namespace/:ns",
				component: Vue.component("router-data"),
				children: [
					{path: "edit/:key", component: Vue.component("router-data-edit")}
				]
			}
		]
	});

	// 应用
	var app = new Vue({
		el: "#app",
		data: {
			"data":  [],  // 右侧项列表数据
			"list": [], // 左侧导航条数据
		},
		router: router,
		methods: {
			
		},
	});

	Vue.mixin({
		beforeCreate: function() {
			this.$app = app;
		}
	});

	exports.$app = app;
})(window);
