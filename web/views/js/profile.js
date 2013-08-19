(function () {
	var render = function render() {
		$('#view').html(this.$html.clone());
	};

	var view = new summa.View();
	view.render = render;

	summa.registerView('profile', view);
})();