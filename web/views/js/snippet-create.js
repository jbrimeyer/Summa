(function () {
	'use strict';

	var $snippetFiles;

	/**
	 * Update the number of files present for this snippet
	 *
	 * @private
	 */
	var _updateFileCount = function _updateFileCount() {
		var numFiles = $snippetFiles.find('.snippet-container').length;
		$snippetFiles.attr('data-files', numFiles);
	};

	/**
	 * Handler triggered when user clicks the "Add Another File" button
	 * @private
	 */
	var _addFile = function _addFile(options) {
		var opts = $.extend({
				focus: true,
				scrollTo: true
			},
			options
		);

		var $file = $('#snippet-template')
			.clone()
			.appendTo($snippetFiles)
			.removeAttr('id');

		summa.newEditor($file.find('.snippet-editor').get(0));

		$file.find('.chosen').chosen();

		if (opts.scrollTo === true) {
			if (opts.focus === true) {
				opts.focus = $file.find('.snippet-name');
			}
			summa.scrollIntoView($file, {focus: opts.focus});
		}

		_updateFileCount();
	};

	/**
	 * Handler triggered when the user clicks on the remove icon
	 * for a given file
	 *
	 * @returns {boolean}
	 * @private
	 */
	var _removeFile = function _removeFile() {
		var $file = $(this).parents('.snippet-container');
		$file.remove();
		_updateFileCount();
		return false;
	};

	/**
	 * Handler triggered when a user selects a different language from
	 * the drop down selection field
	 *
	 * @private
	 */
	var _updateEditorMode = function _updateEditorMode(evt, opt) {
		var $editor = $(this).parents('.snippet').find('.snippet-editor');
		var editor = ace.edit($editor.get(0));
		var mode = summa.languages[summa.consts.DEFAULT_LANGUAGE].mode;

		if (summa.languages[opt.selected].mode !== '') {
			mode = summa.languages[opt.selected].mode;
		}

		editor.getSession().setMode('ace/mode/' + mode);
	};

	/**
	 * Gather up all of the files into an object that can be
	 * converted to JSON and submitted to the API
	 *
	 * @private
	 */
	var _gatherSnippet = function _gatherSnippet() {
		var snippet = {id: $('#snippet-id').val()};
		var $description = $('#snippet-description');

		snippet.description = $description.val().trim();
		if (snippet.description === '') {
			alert('Please enter a short description of your snippet');
			summa.scrollIntoView($description);
			return false;
		}

		snippet.files = [];
		$snippetFiles.find('.snippet').each(function () {
			var $snippet = $(this);
			var $name = $snippet.find('.snippet-name');
			var name = $name.val().trim();

			if (name === '') {
				alert('All files must have a name');
				// TODO: Check for valid filename (alphanumeric, dash, underscore, period)
				summa.scrollIntoView($name);
				snippet = false;
				return false;
			}

			if (typeof snippet.files[name] !== 'undefined') {
				alert('All file names must be unique');
				summa.scrollIntoView($name);
				snippet = false;
				return false;
			}

			var editor = ace.edit($snippet.find('.snippet-editor').get(0));

			snippet.files.push({
				filename: name,
				language: $snippet.find('.snippet-language').val(),
				contents: editor.getValue()
			});

			return true;
		});

		return snippet;
	};

	/**
	 * Handler triggered when the user clicks on the "Create Snippet" button
	 *
	 * @private
	 */
	var _createSnippet = function _createSnippet() {
		var snippet = _gatherSnippet();
		if (snippet === false) {
			return;
		}

		var $btn = $(this).attr('disabled', 'disabled');

		summa.postToApi('/api/snippet/create', {data: snippet})
			.done(function (json) {
				$btn.removeAttr('disabled');
				summa.setHash('/snippet/' + json.data.id);
			});
	};

	/**
	 * Handler triggered when the user clicks on the "Update Snippet" button
	 *
	 * @private
	 */
	var _updateSnippet = function _updateSnippet() {
		var snippet = _gatherSnippet();
		if (snippet === false) {
			return;
		}

		var $btn = $(this).attr('disabled', 'disabled');

		summa.postToApi('/api/snippet/update', {data: snippet})
			.done(function () {
				$btn.removeAttr('disabled');
				summa.setHash('/snippet/' + snippet.id);
			});
	};

	/**
	 * Initialize the user interface
	 *
	 * @private
	 */
	var _initUi = function _initUi() {
		$snippetFiles = $('#snippet-files');
		$('#btn-add-file').click(_addFile);
		$('#btn-create-snippet').click(_createSnippet);
		$('#btn-update-snippet').click(_updateSnippet);
		$snippetFiles.on('click', '.snippet-remove', _removeFile);
		$snippetFiles.on('change', '.snippet-language', _updateEditorMode);
	};

	/**
	 * Our view object
	 *
	 * @constructor
	 */
	var SnippetCreateView = function SnippetCreateView() {
		this._super.constructor.call(this);
		this.name = 'snippet-create';
	};
	summa.inherit(summa.View, SnippetCreateView);

	/**
	 * Render the view
	 */
	SnippetCreateView.prototype.render = function render(args) {
		var that = this;

		if (typeof args.id !== 'undefined') {
			var apiData = {id: args.id};

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
							snippet: that.snippet,
							languages: summa.languages
						}
					);

					var $editors = $('.snippet-editor');
					for (var i = 0; i < that.snippet.files.length; i++) {
						var file = that.snippet.files[i];
						var editor = summa.newEditor(
							$editors.get(i),
							{
								mode: summa.languages[file.language].mode,
								value: file.contents
							}
						);
					}

					$('#snippet-files').find('.chosen').chosen();

					_initUi();
					_updateFileCount();
				});
		}
		else {
			this._super.render.call(this, {languages: summa.languages});
			_initUi();
			_addFile({scrollTo: false});
		}
	};

	summa.registerView(new SnippetCreateView());
})();