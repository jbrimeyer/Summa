(function () {
	var UnreadView = function UnreadView() {
		this._super.constructor.call(this);
		this.name = 'unread';
	};
	summa.inherit(summa.View, UnreadView);

	/**
	 * Render the view
	 */
	UnreadView.prototype.render = function render() {
		var that = this;

		summa.postToApi('/api/snippets/unread')
			.done(function unreadLoadDone(json) {
				that._super.render.call(
					that,
					{
						snippets: json.data.snippets
					}
				);
			})
			.fail(function unreadLoadFail(jqXhr) {
				summa.renderInlineView(jqXhr.status);
			});
	};

	summa.registerView(new UnreadView());
})();