(function () {
	'use strict';

	/**
	 * Delete the snippet
	 *
	 * @private
	 */
	var _snippetDelete = function _snippetDelete() {
		if (confirm('Are you sure you want to delete this snippet?')) {
			summa.renderInlineView(summa.consts.INLINE_VIEW_LOADER);
			summa.postToApi('/api/snippet/delete', {data: {id: this.snippet.id}})
				.fail(function snippetDeleteFail(jqXhr) {
					summa.renderInlineView(jqXhr.status);
				})
				.done(function snippetDeleteDone() {
					summa.setHash('/profile');
				});
		}
	};

	/**
	 * Add a comment
	 */
	var _snippetAddComment = function _snippetAddComment() {
		var editor = ace.edit('comment-add-editor');
		var apiData = {
			snippet_id: this.snippet.id,
			message: editor.getValue().trim()
		};

		if (apiData.message === '') {
			alert('Please enter some text for your comment');
			return false;
		}

		var $btn = $('#btn-add-comment').attr('disabled', 'disabled');

		summa.postToApi('/api/comment/create', {data: apiData})
			.fail(function commentAddFail(jqXhr) {
				summa.renderInlineView(jqXhr.status);
			})
			.done(function commentAddDone(json) {
				var $clone = $('#comment-template').clone().removeAttr('id');
				var comment = json.data.comment;

				$clone.attr('data-id', comment.id);
				$clone.find('.comment-user').text(comment.displayName).attr('href', '#/profile/' + comment.username);
				$clone.find('.comment-ago').text(summa.ago(comment.created));
				$clone.find('.comment-body').html(comment.html);

				$clone.appendTo('#snip-view-comments');

				$btn.removeAttr('disabled');
				editor.setValue('');
			});
	};

	/**
	 * Delete a comment
	 */
	var _snippetDeleteComment = function _snippetDeleteComment() {
		if (confirm('Are you sure you want to delete this comment?')) {
			var $comment = $(this).parents('.comment-box');
			var id = $comment.attr('data-id');

			summa.postToApi('/api/comment/delete', {data: {id: id}})
				.fail(function commentDeleteFail(jqXhr) {
					summa.renderInlineView(jqXhr.status);
				})
				.done(function commentDeleteDone(json) {
					$comment.remove();
				});
		}
	};

	/**
	 * Create a new ACE editor for comments
	 *
	 * @param {Object} el The DOM element in which to inject the editor
	 * @param {Object} [options] Options for the editor
	 * @param {string} [options.value=""] The value to insert into the editor
	 * @return {Object} The editor object
	 * @private
	 */
	var _newCommentEditor = function _newEditor(el, options) {
		options = $.extend({
				value: ''
			},
			options
		);

		var editor = ace.edit(el);
		editor.setShowPrintMargin(false);
		editor.setShowFoldWidgets(false);
		editor.setHighlightActiveLine(false);
		editor.setTheme('ace/theme/chrome');
		editor.renderer.setShowGutter(false);

		var session = editor.getSession();
		session.setTabSize(3);
		session.setUseSoftTabs(false);
		session.setUseWorker(false);
		session.setMode('ace/mode/text');
		session.setValue(options.value);

		return editor;
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
		var apiData = {id: args.id, markRead: true};

		summa.postToApi('/api/snippet', {data: apiData})
			.fail(function snippetFetchFail(jqXhr) {
				summa.renderInlineView(jqXhr.status);
			})
			.done(function snippetFetchDone(json) {
				that.snippet = json.data.snippet;
				that._super.render.call(
					that,
					{
						user: summa.getUser(),
						snippet: that.snippet
					}
				);

				var $editors = $('.snippet-editor');
				for (var i = 0; i < that.snippet.files.length; i++) {
					var file = that.snippet.files[i];
					var editor = summa.newEditor(
						$editors.get(i),
						{
							mode: summa.languages[file.language].mode,
							readonly: true,
							value: file.contents
						}
					);
				}

				_newCommentEditor('comment-add-editor');

				$('#btn-snip-delete').click(function deleteClick() {
					_snippetDelete.call(that);
				});

				$('#btn-add-comment').click(function addCommentClick() {
					_snippetAddComment.call(that);
				});

				$('#snip-view-comments').on('click', '.icon-delete', _snippetDeleteComment);
			});
	};

	summa.registerView(new SnippetView());
})();