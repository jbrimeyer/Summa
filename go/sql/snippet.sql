CREATE TABLE "snippet" (
	"snippet_id" INTEGER PRIMARY KEY,
	"id_base36" TEXT,
	"username" TEXT NOT NULL DEFAULT '',
	"description" TEXT NOT NULL DEFAULT '',
	"created" INTEGER NOT NULL DEFAULT 0,
	"updated" INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX "idx_snippet_id_base36" ON "snippet" ("id_base36");
CREATE INDEX "idx_snippet_username" ON "snippet" ("username");
CREATE INDEX "idx_snippet_created" ON "snippet" ("created");
CREATE INDEX "idx_snippet_updated" ON "snippet" ("updated");