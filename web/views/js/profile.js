(function () {
	var ProfileView = function ProfileView() {
		this._super.constructor.call(this);
		this.name = 'profile';
	};
	summa.inherit(summa.View, ProfileView);

	/**
	 * Render the view
	 */
	ProfileView.prototype.render = function render(args) {
		var that = this;
		var apiData = {username: summa.getUser().username, orderBy: 'created'};

		if (typeof args !== 'undefined' && args.user) {
			apiData.username = args.user;
		}

		$.when(
			summa.postToApi('/api/profile', {data: apiData}),
			summa.postToApi('/api/snippets', {data: apiData})
		)
		.done(function profileLoadDone(user, snip) {
			that._super.render.call(
				that,
				{
					user: summa.getUser(),
					profile: user[0].data.user,
					snippets: snip[0].data.snippets
				}
			);
		})
		.fail(function routeLoadFail(jqXhr) {
			summa.renderInlineView(jqXhr.status);
		});
	};

	summa.registerView(new ProfileView());
})();