(function () {
	'use strict';

	/**
	 * Delete the snippet
	 *
	 * @param snippet
	 * @private
	 */
	var _snippetDelete = function _snippetDelete(snippet) {
		if (confirm('Are you sure you want to delete this snippet?')) {
			summa.postToApi('/api/snippet/delete', {data: {id: snippet.id}})
				.fail(function snippetFetchFail(jqXhr) {
					summa.renderInlineView(jqXhr.status);
				})
				.done(function snippetFetchDone() {
					summa.setHash('/profile');
				});
		}
	};

	/**
	 * Our view object
	 *
	 * @constructor
	 */
	var SnippetView = function SnippetView() {
		this._super.constructor.call(this);
		this.name = 'snippet';
	};
	summa.inherit(summa.View, SnippetView);

	/**
	 * Render the view
	 *
	 * @param args
	 */
	SnippetView.prototype.render = function render(args) {
		var that = this;
		var apiData = {id: args.id};

		summa.postToApi('/api/snippet', {data: apiData})
			.fail(function snippetFetchFail(jqXhr) {
				summa.renderInlineView(jqXhr.status);
			})
			.done(function snippetFetchDone(json) {
				var snippet = json.data.snippet;
				that._super.render.call(
					that,
					{
						user: summa.getUser(),
						snippet: snippet
					}
				);

				var $editors = $('.snippet-editor');
				for (var i = 0; i < snippet.files.length; i++) {
					var file = snippet.files[i];
					var editor = summa.newEditor(
						$editors.get(i),
						{
							mode: summa.languages[file.language].mode,
							readonly: true,
							value: file.contents
						}
					);
				}

				$('#btn-snip-delete').click(function deleteClick() {
					_snippetDelete(snippet);
				});
			});
	};

	summa.registerView(new SnippetView());
})();