CREATE TABLE "snippet" (
	"snippet_id" TEXT PRIMARY KEY,
	"search_id" INTEGER NOT NULL DEFAULT 0,
	"username" TEXT NOT NULL DEFAULT '',
	"description" TEXT NOT NULL DEFAULT '',
	"created" INTEGER NOT NULL DEFAULT 0,
	"updated" INTEGER NOT NULL DEFAULT 0
);
CREATE UNIQUE INDEX "idx_snippet_search_id" ON "snippet" ("search_id");
CREATE INDEX "idx_snippet_username" ON "snippet" ("username");
CREATE INDEX "idx_snippet_created" ON "snippet" ("created");
CREATE INDEX "idx_snippet_updated" ON "snippet" ("updated");