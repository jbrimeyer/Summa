CREATE VIRTUAL TABLE "snippet_search" USING fts4(
	tokenize=porter,
	"snippet" TEXT
);