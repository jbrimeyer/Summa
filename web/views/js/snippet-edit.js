(function () {
	var SnippetEditView = function SnippetEditView() {
		this._super.constructor.call(this);
		this.name = 'snippet-edit';
	};
	summa.inherit(summa.View, SnippetEditView);

	summa.registerView(new SnippetEditView());
})();