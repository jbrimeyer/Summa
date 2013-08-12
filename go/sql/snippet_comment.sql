CREATE TABLE "snippet_comment" (
	"comment_id" INTEGER PRIMARY KEY AUTOINCREMENT,
	"snippet_id" INTEGER,
	"username" TEXT,
	"message" TEXT,
	"created" INTEGER,
	"updated" INTEGER
);
CREATE INDEX "idx_snippet_comment_snippet_id" ON "snippet_comment" ("snippet_id");
CREATE INDEX "idx_snippet_comment_created" ON "snippet_comment" ("created");