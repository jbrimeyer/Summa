(function () {
	var SearchView = function SearchView() {
		this._super.constructor.call(this);
		this.name = 'search';
	};
	summa.inherit(summa.View, SearchView);

	/**
	 * Render the view
	 */
	SearchView.prototype.render = function render(args) {
		var that = this;
		var apiData = {term: args.term};

		summa.postToApi('/api/snippets/search', {data: apiData})
			.done(function searchLoadDone(json) {
				that._super.render.call(
					that,
					{
						snippets: json.data.snippets
					}
				);
			})
			.fail(function searchLoadFail(jqXhr) {
				summa.renderInlineView(jqXhr.status);
			});
	};

	summa.registerView(new SearchView());
})();