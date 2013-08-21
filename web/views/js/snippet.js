(function () {
	var SnippetView = function SnippetView() {
		this._super.constructor.call(this);
		this.name = 'snippet';
	};
	summa.inherit(summa.View, SnippetView);

	SnippetView.prototype.render = function render(args) {
		var that = this;
		var apiData = {id: args.id};

		that._super.render.call(this, args);

		summa.postToApi('/api/snippet', {data: apiData})
			.fail(function snippetFetchFail(jqXhr) {
				summa.renderInlineView(jqXhr.status);
			})
			.done(function snippetFetchDone(json) {
				var snippet = json.data.snippet;
				$('#view').html(that.$html.clone());

				$('#snip-view-username').text(snippet.username);
				$('#snip-view-repo').text(snippet.files[0].filename);
				$('#snip-view-description').text(snippet.description);
				$('#snip-view-created').text('Created ' + summa.ago(snippet.created));
			});
	};

	summa.registerView(new SnippetView());
})();