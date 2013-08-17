CREATE TABLE "user_session" (
	"username" TEXT,
	"token" TEXT,
	"created" INTEGER,
	PRIMARY KEY ("username", "token")
);