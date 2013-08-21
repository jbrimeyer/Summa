(function () {
	var UnreadView = function UnreadView() {
		this._super.constructor.call(this);
		this.name = 'unread';
	};
	summa.inherit(summa.View, UnreadView);

	summa.registerView(new UnreadView());
})();