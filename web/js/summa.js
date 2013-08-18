(function () {
	'use strict';

	var _consts = {
		ROUTE_DEFAULT: '/',
		COOKIE_NAME: 'summa',
		PATH_VIEWS: '/views/'
	};
	var _routes;
	var _user;

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
	 * Show a modal
	 *
	 * @param {string} modal The selector of the modal to show
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

		$('#auth-button').click(_authenticate);
		$('#auth-modal')
			.on('hide.bs.modal', function (e) {
				if (!_isSignedIn()) {
					e.preventDefault();
				}
			});

		$('#email-button').click(_saveEmail);
		$('#email-modal')
			.on('hide.bs.modal', function (e) {
				if (!_hasEmail()) {
					e.preventDefault();
				}
			});
	};

	/**
	 * Fetch a route based on the path
	 *
	 * @param {string} path
	 * @returns {Object|null}
	 * @private
	 */
	var _fetchRoute = function _fetchRoute(path) {
		if (typeof _routes[path] === 'undefined') {
			// TODO: Deep lookup of route using regexp
			return null;
		}

		return _routes[path];
	};

	/**
	 * Setup handler for hash change events
	 */
	var _hashChange = function _hashChange() {
		_updateAuthStatus();

		if (!_isSignedIn()) {
			_showModal('#auth-modal');
		}
		else if (!_hasEmail()) {
			_showModal('#email-modal');
		}
		else {
			var hash = _getHash();
			var route = _fetchRoute(hash);

			if (route === null) {
				// TODO: 404
				console.log('No route for ' + hash);
				return;
			}

			if (typeof route.handler === 'function') {
				route.handler();
			}
			else {
				// TODO: Do something with route.view
			}
		}
	};

	_routes = {
		'/': {
			view: 'index.html'
		},
		'/signout': {
			handler: _signOut
		},
		'/profile': {
			view: 'profile.html'
		},
		'/search': {
			view: 'search.html'
		},
		'/snippet/{snippetId}': {
			view: 'snippet.html'
		},
		'/unread': {
			view: 'search.html'
		}
	};

	/**
	 * Initialization on document ready
	 */
	$(function () {
		_initUi();
		_clearUser();
		_restoreUserInfo();

		if (_getHash() === '') {
			_setHash(_consts.ROUTE_DEFAULT);
		}
		else {
			_hashChange();
		}
	});
})();