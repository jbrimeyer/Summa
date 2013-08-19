(function () {
	var view = new summa.View();
	view.render = function render(args) {
		$('#view').html(this.$html.clone());
		console.log('Render snippet EDIT', args);
	};

	summa.registerView('snippet-edit', view);
})();