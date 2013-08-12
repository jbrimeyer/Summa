CREATE TABLE "snippet_file" (
	"snippet_id" INTEGER,
	"filename" TEXT,
	"language" TEXT,
	PRIMARY KEY("snippet_id", "filename")
);