(function () {
	var SearchView = function SearchView() {
		this._super.constructor.call(this);
		this.name = 'search';
	};
	summa.inherit(summa.View, SearchView);

	summa.registerView(new SearchView());
})();