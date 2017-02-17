Vue.component("index-path-data",function(resolve, reject) {
	Vue.http.get("/index-path-data.html").then(function(response) {
		resolve({
			template: response.body,
			props: ["data"],
			methods: {
				segment: function(path) {
					return path.split("/").pop();
				},
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
