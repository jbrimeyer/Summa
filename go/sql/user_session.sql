CREATE TABLE "user_session" (
	"username" TEXT,
	"token" TEXT,
	"timestamp" INTEGER,
	PRIMARY KEY ("username", "token")
);