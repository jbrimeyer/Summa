(function () {
	var ProfileView = function ProfileView() {
		this._super.constructor.call(this);
		this.name = 'profile';
	};
	summa.inherit(summa.View, ProfileView);

	summa.registerView(new ProfileView());
})();