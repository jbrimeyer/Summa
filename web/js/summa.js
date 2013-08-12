(function () {
	var defaultLang = 'Text';
	var allLangs = {
		'ABAP': {mode: 'abap', exts: ''},
		'ActionScript': {mode: 'actionscript', exts: ''},
		'Ada': {mode: 'ada', exts: ''},
		'ApacheConf': {mode: '', exts: ''},
		'Apex': {mode: '', exts: ''},
		'AppleScript': {mode: '', exts: ''},
		'Arc': {mode: '', exts: ''},
		'Arduino': {mode: '', exts: ''},
		'ASP': {mode: '', exts: ''},
		'Assembly': {mode: 'assembly_x86', exts: ''},
		'Augeas': {mode: '', exts: ''},
		'AutoHotkey': {mode: 'autohotkey', exts: ''},
		'Awk': {mode: '', exts: ''},
		'Batchfile': {mode: 'batchfile', exts: ''},
		'Befunge': {mode: '', exts: ''},
		'BlitzMax': {mode: '', exts: ''},
		'Boo': {mode: '', exts: ''},
		'Brainfuck': {mode: '', exts: ''},
		'Bro': {mode: '', exts: ''},
		'C': {mode: 'c_cpp', exts: 'c,h'},
		'C-ObjDump': {mode: '', exts: ''},
		'C#': {mode: 'csharp', exts: 'cs'},
		'C++': {mode: 'c_cpp', exts: 'C,cc,cpp,cxx,hpp,hxx'},
		'C2hs Haskell': {mode: '', exts: ''},
		'Ceylon': {mode: '', exts: ''},
		'ChucK': {mode: '', exts: ''},
		'CLIPS': {mode: '', exts: ''},
		'Clojure': {mode: 'clojure', exts: ''},
		'CMake': {mode: '', exts: ''},
		'CoffeeScript': {mode: 'coffee', exts: ''},
		'ColdFusion': {mode: 'coldfusion', exts: ''},
		'Common Lisp': {mode: '', exts: ''},
		'Coq': {mode: '', exts: ''},
		'Cpp-ObjDump': {mode: '', exts: ''},
		'CSS': {mode: 'css', exts: 'css'},
		'Cucumber': {mode: '', exts: ''},
		'Cython': {mode: '', exts: ''},
		'D': {mode: 'd', exts: ''},
		'D-ObjDump': {mode: '', exts: ''},
		'Darcs Patch': {mode: '', exts: ''},
		'Dart': {mode: 'dart', exts: 'dart'},
		'DCPU-16 ASM': {mode: '', exts: ''},
		'Delphi': {mode: '', exts: ''},
		'Diff': {mode: 'diff', exts: ''},
		'DOT': {mode: 'dot', exts: ''},
		'Dylan': {mode: '', exts: ''},
		'eC': {mode: '', exts: ''},
		'Ecere Projects': {mode: '', exts: ''},
		'Ecl': {mode: '', exts: ''},
		'edn': {mode: '', exts: ''},
		'Eiffel': {mode: '', exts: ''},
		'Elixir': {mode: '', exts: ''},
		'Elm': {mode: '', exts: ''},
		'Emacs Lisp': {mode: '', exts: ''},
		'Erlang': {mode: 'erlang', exts: ''},
		'F#': {mode: '', exts: ''},
		'Factor': {mode: '', exts: ''},
		'Fancy': {mode: '', exts: ''},
		'Fantom': {mode: '', exts: ''},
		'fish': {mode: '', exts: ''},
		'Forth': {mode: 'forth', exts: ''},
		'FORTRAN': {mode: '', exts: ''},
		'GAS': {mode: '', exts: ''},
		'Genshi': {mode: '', exts: ''},
		'Gentoo Ebuild': {mode: '', exts: ''},
		'Gentoo Eclass': {mode: '', exts: ''},
		'Gettext Catalog': {mode: '', exts: ''},
		'Go': {mode: 'golang', exts: 'go'},
		'Gosu': {mode: '', exts: ''},
		'Groff': {mode: '', exts: ''},
		'Groovy': {mode: 'groovy', exts: ''},
		'Groovy Server Pages': {mode: '', exts: ''},
		'Haml': {mode: 'haml', exts: ''},
		'Handlebars': {mode: '', exts: ''},
		'Haskell': {mode: 'haskell', exts: ''},
		'Haxe': {mode: 'haxe', exts: ''},
		'HTML': {mode: 'html', exts: 'htm,html'},
		'HTML+Django': {mode: '', exts: ''},
		'HTML+ERB': {mode: '', exts: ''},
		'HTML+PHP': {mode: '', exts: ''},
		'HTTP': {mode: '', exts: ''},
		'INI': {mode: 'ini', exts: ''},
		'Io': {mode: '', exts: ''},
		'Ioke': {mode: '', exts: ''},
		'IRC log': {mode: '', exts: ''},
		'Java': {mode: 'java', exts: 'java'},
		'Java Server Pages': {mode: '', exts: ''},
		'JavaScript': {mode: 'javascript', exts: 'js'},
		'JSON': {mode: 'json', exts: 'json'},
		'Julia': {mode: 'julia', exts: ''},
		'Kotlin': {mode: '', exts: ''},
		'Lasso': {mode: '', exts: ''},
		'Less': {mode: 'less', exts: ''},
		'LilyPond': {mode: '', exts: ''},
		'Literate CoffeeScript': {mode: '', exts: ''},
		'Literate Haskell': {mode: '', exts: ''},
		'LiveScript': {mode: 'livescript', exts: ''},
		'LLVM': {mode: '', exts: ''},
		'Logos': {mode: '', exts: ''},
		'Logtalk': {mode: '', exts: ''},
		'Lua': {mode: 'lua', exts: ''},
		'Makefile': {mode: 'makefile', exts: ''},
		'Mako': {mode: '', exts: ''},
		'Markdown': {mode: 'markdown', exts: 'md'},
		'Matlab': {mode: 'matlab', exts: ''},
		'Max': {mode: '', exts: ''},
		'MiniD': {mode: '', exts: ''},
		'Mirah': {mode: '', exts: ''},
		'Monkey': {mode: '', exts: ''},
		'Moocode': {mode: '', exts: ''},
		'MoonScript': {mode: '', exts: ''},
		'mupad': {mode: '', exts: ''},
		'Myghty': {mode: '', exts: ''},
		'Nemerle': {mode: '', exts: ''},
		'Nginx': {mode: '', exts: ''},
		'Nimrod': {mode: '', exts: ''},
		'NSIS': {mode: '', exts: ''},
		'Nu': {mode: '', exts: ''},
		'NumPy': {mode: '', exts: ''},
		'ObjDump': {mode: '', exts: ''},
		'Objective-C': {mode: 'objectivec', exts: ''},
		'Objective-J': {mode: '', exts: ''},
		'OCaml': {mode: 'ocaml', exts: ''},
		'Omgrofl': {mode: '', exts: ''},
		'ooc': {mode: '', exts: ''},
		'Opa': {mode: '', exts: ''},
		'OpenCL': {mode: '', exts: ''},
		'OpenEdge ABL': {mode: '', exts: ''},
		'Parrot': {mode: '', exts: ''},
		'Parrot Assembly': {mode: '', exts: ''},
		'Parrot Internal Representation': {mode: '', exts: ''},
		'Perl': {mode: 'perl', exts: 'pl'},
		'PHP': {mode: 'php', exts: 'php'},
		'Pike': {mode: '', exts: ''},
		'PogoScript': {mode: '', exts: ''},
		'PowerShell': {mode: 'powershell', exts: ''},
		'Prolog': {mode: 'prolog', exts: ''},
		'Puppet': {mode: '', exts: ''},
		'Pure Data': {mode: '', exts: ''},
		'Python': {mode: 'python', exts: 'py'},
		'Python traceback': {mode: '', exts: ''},
		'R': {mode: 'r', exts: ''},
		'Racket': {mode: '', exts: ''},
		'Ragel in Ruby Host': {mode: '', exts: ''},
		'Raw token data': {mode: '', exts: ''},
		'Rebol': {mode: '', exts: ''},
		'Redcode': {mode: '', exts: ''},
		'reStructuredText': {mode: '', exts: ''},
		'RHTML': {mode: '', exts: ''},
		'Rouge': {mode: '', exts: ''},
		'Ruby': {mode: 'ruby', exts: 'rb'},
		'Rust': {mode: 'rust', exts: ''},
		'Sage': {mode: '', exts: ''},
		'Sass': {mode: 'sass', exts: ''},
		'Scala': {mode: 'scala', exts: ''},
		'Scheme': {mode: 'scheme', exts: ''},
		'Scilab': {mode: '', exts: ''},
		'SCSS': {mode: 'scss', exts: ''},
		'Self': {mode: '', exts: ''},
		'Shell': {mode: 'sh', exts: ''},
		'Smalltalk': {mode: '', exts: ''},
		'Smarty': {mode: '', exts: ''},
		'SQL': {mode: 'sql', exts: 'sql'},
		'Standard ML': {mode: '', exts: ''},
		'SuperCollider': {mode: '', exts: ''},
		'Tcl': {mode: 'tcl', exts: ''},
		'Tcsh': {mode: '', exts: ''},
		'Tea': {mode: '', exts: ''},
		'TeX': {mode: 'tex', exts: ''},
		'Text': {mode: 'text', exts: ''},
		'Textile': {mode: 'textile', exts: ''},
		'TOML': {mode: 'toml', exts: ''},
		'Turing': {mode: '', exts: ''},
		'Twig': {mode: 'twig', exts: ''},
		'TXL': {mode: '', exts: ''},
		'TypeScript': {mode: 'typscript', exts: ''},
		'Vala': {mode: '', exts: ''},
		'Verilog': {mode: '', exts: ''},
		'VHDL': {mode: '', exts: ''},
		'VimL': {mode: '', exts: ''},
		'Visual Basic': {mode: '', exts: ''},
		'XML': {mode: 'xml', exts: 'xml'},
		'XProc': {mode: '', exts: ''},
		'XQuery': {mode: '', exts: ''},
		'XS': {mode: '', exts: ''},
		'XSLT': {mode: '', exts: ''},
		'Xtend': {mode: '', exts: ''},
		'YAML': {mode: 'yaml', exts: 'yaml'}
	};
	var langByExt = {};
	var $snippetFiles;

	/**
	 * Start a POST request to the Summa API
	 *
	 * @param {string} url
	 * @param {Object} data
	 * @returns {jqXHR}
	 * @private
	 */
	var _postToApi = function _postToApi(url, data) {
		return $.ajax({
			type: 'POST',
			url: url,
			data: JSON.stringify(data),
			contentType: 'application/json; charset=utf-8',
			dataType: 'json'
		});
	};

	/**
	 * Scroll an element into view if it is not currently in view
	 *
	 * @param {jQuery} $el The jQuery object containing the element to scroll into view
	 * @param {Object} [options] Options
	 * @param {boolean} [options.focus=true] Focus the element after scroll
	 * @private
	 */
	var _scrollIntoView = function _scrollIntoView($el, options) {
		var opts = $.extend({
				focus: true,
				topMargin: 20
			},
			options
		);

		function focus() {
			if (opts.focus === true) {
				$el.focus();
			}
			else if (opts.focus !== false) {
				$(opts.focus).focus();
			}
		}

		var $window = $(window);
		var scrollTop = $window.scrollTop();
		var docViewBottom = scrollTop + $window.height();
		var offset = $el.offset();
		offset.bottom = offset.top + $el.height();
		var inView = ((offset.bottom <= docViewBottom) && (offset.top >= scrollTop));

		if (!inView) {
			var newScrollTop = Math.max(offset.top - opts.topMargin, 0);

			$('html, body').animate(
				{scrollTop: newScrollTop},
				focus
			);
		}
		else {
			focus();
		}
	};

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

		var editor = ace.edit($file.find('.snippet-editor').get(0));
		editor.setShowPrintMargin(false);
		editor.setShowFoldWidgets(false);
		editor.setTheme('ace/theme/chrome');

		var session = editor.getSession();
		session.setTabSize(3);
		session.setUseSoftTabs(false);
		session.setUseWorker(false);
		session.setMode('ace/mode/text');

		$file.find('.chosen').chosen();
;
		if (opts.scrollTo === true) {
			if (opts.focus === true) {
				opts.focus = $file.find('.snippet-name');
			}
			_scrollIntoView($file, {focus: opts.focus});
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
	 * Initialize the available languages that can be selected
	 * from when adding/editing a file
	 *
	 * @private
	 */
	var _initLanguages = function _initLanguages() {
		var $select = $('#snippet-template').find('.snippet-language');
		for (var lang in allLangs) {
			$select.append(
				'<option value="' + lang + '">' + lang + '</option>'
			);

			var exts = allLangs[lang].exts.split(',');
			for (var i = 0; i < exts.length; i++) {
				if (exts[i] !== '') {
					langByExt[exts[i]] = lang;
				}
			}
		}

		$select.children('[value="' + defaultLang + '"]').attr('selected', 'selected');
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
		var mode = allLangs[defaultLang].mode;

		if (allLangs[opt.selected].mode !== '') {
			mode = allLangs[opt.selected].mode;
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
		var snippet = {};
		var $description = $('#snippet-description');

		snippet.description = $description.val().trim();
		if (snippet.description === '') {
			alert('Please enter a short description of your snippet');
			_scrollIntoView($description);
			return false;
		}

		snippet.files = {};
		$snippetFiles.find('.snippet').each(function () {
			var $snippet = $(this);
			var $name = $snippet.find('.snippet-name');
			var name = $name.val().trim();

			if (name === '') {
				alert('All files must have a name');
				// TODO: Check for valid filename (alphanumeric, dash, underscore, period)
				_scrollIntoView($name);
				snippet = false;
				return false;
			}

			if (typeof snippet.files[name] !== 'undefined') {
				alert('All file names must be unique');
				_scrollIntoView($name);
				snippet = false;
				return false;
			}

			var editor = ace.edit($snippet.find('.snippet-editor').get(0));

			snippet.files[name] = {
				lang: $snippet.find('.snippet-language').val(),
				content: editor.getValue()
			};

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

		_postToApi('/api/snippet/create', snippet)
			.fail(function () {
				console.log('FAIL', arguments);
			})
			.done(function (data) {
				console.log('DONE', arguments);
			});
	};


	// Initialization on page ready
	$(function () {
		$snippetFiles = $('#snippet-files');
		$('[data-toggle="tooltip"]').tooltip();
		$('#btn-add-file').click(_addFile);
		$('#btn-create-snippet').click(_createSnippet);
		_initLanguages();
		_addFile({scrollTo: false});

		$snippetFiles.on('click', '.snippet-remove', _removeFile);
		$snippetFiles.on('change', '.snippet-language', _updateEditorMode);
	});
})();