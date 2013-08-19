var summa = (function () {
	'use strict';

	var _consts = {
		ROUTE_DEFAULT: '/',
		COOKIE_NAME: 'summa',
		PATH_VIEWS: '/views/',
		PATH_JS_VIEWS: '/views/js/',
		DEFAULT_LANGUAGE: 'Text',
		INLINE_VIEW_400: '400',
		INLINE_VIEW_403: '403',
		INLINE_VIEW_404: '404',
		INLINE_VIEW_500: '500',
		INLINE_VIEW_LOADER: 'loader'
	};
	var _exports = {};
	var _routes = [];
	var _routesLookup = {};
	var _views = {};
	var _user;

	/**
	 * View object
	 *
	 * @constructor
	 */
	var View = function View() {
		this.$html = null;
		this.initHtml = null;
		this.render = null;
	};

	/**
	 * Set the view's HTML
	 *
	 * @param {string} [html]
	 */
	View.prototype.setHtml = function setHtml(html) {
		this.$html = $(html);
		if (typeof this.initHtml === 'function') {
			this.initHtml();
		}
	};

	/**
	 * Register a view
	 *
	 * @param {string} name
	 * @param {View} view
	 * @private
	 */
	var _registerView = function _registerView(name, view) {
		_views[name] = view;
	};

	/**
	 * Start a POST request to the Summa API
	 *
	 * @param {string} url
	 * @param {Object} data
	 * @returns {jqXHR}
	 * @private
	 */
	var _postToApi = function _postToApi(url, data) {
		if (typeof data === 'undefined') {
			data = {};
		}

		if (!data.username) {
			data.username = _user.username;
			data.token = _user.token;
		}

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
	 * Update the view to indicate that data is being loaded in the background
	 *
	 * @private
	 */
	var _pageLoading = function _pageLoading() {
		_renderInlineView(_consts.INLINE_VIEW_LOADER);
	};

	/**
	 * Show a modal
	 *
	 * @param {string} modal The selector of the modal to show
	 * @param {string} focus The selector of the item to receive focus
	 * @private
	 */
	var _showModal = function _showModal(modal) {
		$(modal).modal();
	};

	/**
	 * Hide a modal
	 *
	 * @param {string} modal The selector of the modal to hide
	 * @private
	 */
	var _hideModal = function _hideModal(modal) {
		$(modal).modal('hide');
	};

	/**
	 * Start the user authentication process, generally when
	 * the sign in button is clicked on the authentication modal
	 *
	 * @private
	 */
	var _saveEmail = function _saveEmail() {
		_modalLoading(true, '#email-modal');

		var apiData = {data: {}};

		apiData.data.email = $('#email-address').val().trim();
		if (apiData.data.email === '') {
			_setModalError('E-mail Address is required', '#email-modal', '#email-address');
			return;
		}

		apiData.data.displayName = _user.displayName;

		_postToApi('/api/profile/update', apiData)
			.done(function authDone(json) {
				_user.hasEmail = true;
				_saveUserInfo();
				_modalLoading(false, '#email-modal');
				_hideModal('#email-modal');
				_hashChange();
			})
			.fail(function authFail(xhr) {
				// TODO: Handle all types of errors
			});
	};

	/**
	 * Check if the user is signed in
	 *
	 * @returns {boolean}
	 * @private
	 */
	var _isSignedIn = function _isSignedIn() {
		return _user.username !== null;
	};

	/**
	 * Check if the user has an e-mail address set
	 *
	 * @returns {boolean}
	 * @private
	 */
	var _hasEmail = function _hasEmail() {
		return _user.hasEmail;
	};

	/**
	 * Update the page authentication status
	 * @private
	 */
	var _updateAuthStatus = function _updateAuthStatus() {
		var dataAuth = 0;
		var displayName = '';

		if (_isSignedIn()) {
			dataAuth = 1;
			displayName = _user.displayName;
		}

		$('#header-display-name').text(displayName);
		$('body').attr('data-auth', dataAuth);
	};

	/**
	 * Clear the _user variable and cookie
	 *
	 * @private
	 */
	var _clearUser = function _clearUser() {
		_user = {
			username: null,
			displayName: null,
			token: null,
			hasEmail: false
		};
	};

	/**
	 * Set the error message in a modal
	 *
	 * @param {string} message The message to display
	 * @param {string} modal The selector of the modal
	 * @param {string} [focus] The selector of the element to focus
	 * @private
	 */
	var _setModalError = function _setModalError(message, modal, focus) {
		$(modal).find('.error').text(message);

		if (typeof focus !== 'undefined') {
			$(focus).focus();
		}

		_modalLoading(false, modal);
	};

	/**
	 * Update a modal to indicate if a
	 * background task is being executed
	 *
	 * @param {boolean} isLoading
	 * @param {string} modal The selector of the model to update
	 * @private
	 */
	var _modalLoading = function _modalLoading(isLoading, modal) {
		isLoading = (isLoading !== false);

		if (isLoading) {
			_setModalError('', modal);
			$(modal).attr('data-loading', 1);
			$(modal).find('.btn-primary').attr('disabled', 'disabled');
		}
		else {
			$(modal).attr('data-loading', 0);
			$(modal).find('.btn-primary').removeAttr('disabled');
		}
	};

	/**
	 * Start the user authentication process, generally when
	 * the sign in button is clicked on the authentication modal
	 *
	 * @private
	 */
	var _authenticate = function _authenticate() {
		_modalLoading(true, '#auth-modal');

		var apiData = {};

		apiData.username = $('#auth-username').val().trim();
		if (apiData.username === '') {
			_setModalError('Username is required', '#auth-modal', '#auth-username');
			return;
		}

		apiData.password = $('#auth-password').val().trim();
		if (apiData.password === '') {
			_setModalError('Password is required', '#auth-modal', '#auth-password');
			return;
		}

		_postToApi('/api/auth/signin', apiData)
			.done(function authDone(json) {
				_user.username = json.data.user.username;
				_user.displayName = json.data.user.displayName;
				_user.hasEmail = !json.data.needEmail;
				_user.token = json.token;

				_saveUserInfo();
				_updateAuthStatus();
				_modalLoading(false, '#auth-modal');
				_hideModal('#auth-modal');

				if (json.data.needEmail) {
					_showModal('#email-modal');
				}
				else {
					_hashChange();
				}
			})
			.fail(function authFail(xhr) {
				// TODO: Handle all types of errors
				if (xhr.status === 401) {
					_setModalError('Authentication failed!', '#auth-modal');
				}
			});
	};

	/**
	 * Sign user out of Summa
	 *
	 * @private
	 */
	var _signOut = function _signOut() {
		_postToApi('/api/auth/signout');
		_deleteUserInfo();
		_updateAuthStatus();
		_setHash(_consts.ROUTE_DEFAULT);
	};

	/**
	 * Save the user login information stored in _user
	 * to a browser cookie
	 *
	 * @private
	 */
	var _saveUserInfo = function _saveUserInfo() {
		var userInfo = btoa(JSON.stringify(_user));

		document.cookie =
			_consts.COOKIE_NAME + '=' + encodeURIComponent(userInfo);
	};

	/**
	 * Restore the user login information stored in a
	 * browser cookie to the _user variable
	 *
	 * @private
	 */
	var _restoreUserInfo = function _restoreUserInfo() {
		var cookies = document.cookie.split('; ');
		for (var i = 0; i < cookies.length; i++) {
			var cookie = cookies[i].split('=');

			if (cookie[0] === _consts.COOKIE_NAME) {
				var value = decodeURIComponent(cookie[1]);
				_user = JSON.parse(atob(value));
				break;
			}
		}
	};

	/**
	 * Delete the user login information stored in
	 * _user and the browser cookie
	 * @private
	 */
	var _deleteUserInfo = function _deleteUserInfo() {
		_clearUser();
		document.cookie =
			_consts.COOKIE_NAME + '=deleted; expires=' + new Date(0).toUTCString();
	};

	/**
	 * Get the current page hash string
	 *
	 * @returns {string}
	 * @private
	 */
	var _getHash = function _getHash() {
		return window.location.hash.slice(1);
	};

	/**
	 * Set the page hash string
	 *
	 * @param value
	 * @private
	 */
	var _setHash = function _setHash(value) {
		window.location.hash = value;
	};

	/**
	 * Initialize the user interface
	 *
	 * @private
	 */
	var _initUi = function _initUi() {
		$(window).on('hashchange', _hashChange);
		$('[data-toggle="tooltip"]').tooltip();

		$('.modal-body').find('input').keydown(function modalKeyDown(e) {
			if (e.keyCode === 13) {
				$(this).parents('.modal-dialog').find('.btn-primary').click();
			}
		});

		$('#auth-button').click(_authenticate);
		$('#auth-modal')
			.on('hide.bs.modal', function authModalHide(e) {
				if (!_isSignedIn()) {
					e.preventDefault();
				}
			})
			.on('shown.bs.modal', function authModalShown(e) {
				$('#auth-username').focus();
			});

		$('#email-button').click(_saveEmail);
		$('#email-modal')
			.on('hide.bs.modal', function emailModalHide(e) {
				if (!_hasEmail()) {
					e.preventDefault();
				}
			})
			.on('shown.bs.modal', function emailModalShown(e) {
				$('#email-address').focus();
			});
	};

	/**
	 * Map file extensions to languages
	 *
	 * @private
	 */
	var _initLanguages = function _initLanguages() {
		for (var lang in _languages) {
			var exts = _languages[lang].exts.split(',');
			for (var i = 0; i < exts.length; i++) {
				if (exts[i] !== '') {
				_languagesByExt[exts[i]] = lang;
				}
			}
		}
	};

	/**
	 * Render one of the inline page views
	 *
 	 * @param view
	 * @private
	 */
	var _renderInlineView = function _renderInlineView(view) {
		$('#view').html($('[data-view="' + view + '"]').clone());
	};

	/**
	 * Add a route handler
	 *
	 * @param {string} path
	 * @param {string|function} handler
	 * @private
	 */
	var _addRoute = function _addRoute(path, handler) {
		if (path.indexOf('{') !== -1) {
			var parts = [];
			path = path.replace(/\{(.*?)\}/g, function (full, part) {
				parts.push(part);
				return '([^/]+)';
			});

			_routes.push({
				handler: handler,
				parts: parts,
				regex: new RegExp(path)
			})
		}
		else {
			_routesLookup[path] = handler;
		}
	};

	/**
	 * Fetch a route based on the path
	 *
	 * @param {string} path
	 * @returns {Object|null}
	 * @private
	 */
	var _fetchRoute = function _fetchRoute(path) {
		if (typeof _routesLookup[path] === 'undefined') {
			var matches;

			for (var i = 0; i < _routes.length; i++) {
				if ((matches = _routes[i].regex.exec(path)) !== null) {
					var args = {};

					for (var j = 1; j < matches.length; j++) {
						args[_routes[i].parts[j-1]] = matches[j]
					}

					return {
						handler: _routes[i].handler,
						args: args
					};
				}
			}
			return null;
		}

		return _routesLookup[path];
	};

	/**
	 * Handler for hash change events
	 */
	var _hashChange = function _hashChange() {
		_updateAuthStatus();

		if (!_isSignedIn()) {
			_showModal('#auth-modal', '');
		}
		else if (!_hasEmail()) {
			_showModal('#email-modal');
		}
		else {
			_pageLoading();
			var hash = _getHash();
			var route = _fetchRoute(hash);

			if (route === null) {
				_renderInlineView(_consts.INLINE_VIEW_404);
				return;
			}

			if (typeof route === 'function') {
				route();
			}
			else {
				_processRoute(route);
			}
		}
	};

	/**
	 * Process the route by loading assets and calling the render method
	 *
	 * @param route
	 * @private
	 */
	var _processRoute = function _processRoute(route) {
		var args = {};

		if (typeof route === 'object') {
			args = route.args;
			route = route.handler;
		}

		if (typeof _views[route] === 'undefined') {
			var jsFile = _consts.PATH_JS_VIEWS + route + '.js';
			var htmlFile = _consts.PATH_VIEWS + route + '.html';

			$.when(
				$.ajax({
					url: jsFile,
					dataType: 'script'
				}),
				$.ajax({
					url: htmlFile,
					dataType: 'html'
				})
			)
			.done(function routeLoadDone(js, html) {
				_views[route].setHtml(html[0]);
				_views[route].render(args);
			})
			.fail(function routeLoadFail() {
				// TODO: Handle route load failure
			});
		}
		else {
			_views[route].render(args);
		}
	};

	/**
	 * Generate a friendly "x time ago" string
	 *
	 * @param ms
	 * @private
	 */
	var _ago = function _ago(ms) {
		var friendly;
		var agoMs = Date.now() - ms;
		var seconds = agoMs / 1000;
		var minutes = seconds / 60;
		var hours = minutes / 60;
		var days = hours / 24;
		var years = days / 365;
		var s = {
			seconds: "less than a minute",
			minute: "about a minute",
			minutes: "%d minutes",
			hour: "about an hour",
			hours: "about %d hours",
			day: "a day",
			days: "%d days",
			month: "about a month",
			months: "%d months",
			year: "about a year",
			years: "%d years"
		};


		switch (true) {
			case seconds < 45:
				friendly = 'less than a minute';
				break;
			case seconds < 90:
				friendly = 'about a minute';
				break;
			case minutes < 45:
				friendly = Math.round(minutes) +  ' minutes';
				break;
			case minutes < 90:
				friendly = 'about an hour';
				break;
			case hours < 24:
				friendly = 'about ' + Math.round(hours) + ' hours';
				break;
			case hours < 42:
				friendly = 'a day';
				break;
			case days < 30:
				friendly = Math.round(days) + ' days';
				break;
			case days < 45:
				friendly = 'about a month';
				break;
			case days < 365:
				friendly = Math.round(days / 30) + ' months';
				break;
			case years < 1.5:
				friendly = 'about a year';
				break;
			default:
				friendly = Math.round(years) + ' years';
				break;
		}

		return friendly + ' ago';
	};

	// Define all of our known routes
	_addRoute('/', 'index');
	_addRoute('/signout', _signOut);
	_addRoute('/profile', 'profile');
	_addRoute('/search', 'search');
	_addRoute('/snippet/{id}/edit', 'snippet-edit');
	_addRoute('/snippet/{id}', 'snippet');
	_addRoute('/unread', 'unread');

	/**
	 * Initialization on document ready
	 */
	$(function () {
		_initUi();
		_initLanguages();
		_clearUser();
		_restoreUserInfo();

		if (_getHash() === '') {
			_setHash(_consts.ROUTE_DEFAULT);
		}
		else {
			_hashChange();
		}
	});

	var _languages = {
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
	var _languagesByExt = {};

	_exports.consts = _consts;
	_exports.View = View;
	_exports.postToApi = _postToApi;
	_exports.registerView = _registerView;
	_exports.languages = _languages;
	_exports.languagesByExt = _languagesByExt;
	_exports.scrollIntoView = _scrollIntoView;
	_exports.setHash = _setHash;
	_exports.renderInlineView = _renderInlineView;
	_exports.pageLoading = _pageLoading;
	_exports.ago = _ago;

	return _exports;
})();