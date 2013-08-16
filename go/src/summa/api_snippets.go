package summa

import (
	"database/sql"
	_ "go-sqlite3"
)

func apiSnippetsSearch(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	// TODO: apiSnippetsSearch()
	return nil
}

func apiSnippetsUnread(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	snippets, err := snippetsUnread(db, req.Username)
	if err != nil {
		return &internalServerError{"Could not fetch snippets", err}
	}

	resp["snippets"] = snippets

	return nil
}
