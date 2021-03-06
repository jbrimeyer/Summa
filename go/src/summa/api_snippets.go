package summa

import (
	"database/sql"
	_ "go-sqlite3"
	"strings"
)

const (
	SNIPPETS_LIMIT_MAX     = 200
	SNIPPETS_LIMIT_DEFAULT = 100
)

var (
	snippetsOrderBy = map[string]string{
		"commentsAsc":     "num_comments",
		"commentsDesc":    "num_comments DESC",
		"filesAsc":        "num_files",
		"filesDesc":       "num_files DESC",
		"createdAsc":      "s.created",
		"createdDesc":     "s.created DESC",
		"updatedAsc":      "s.updated",
		"updatedDesc":     "s.updated DESC",
		"descriptionAsc":  "s.description",
		"descriptionDesc": "s.description DESC",
	}
)

func apiSnippets(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	start, _ := req.Data["start"].(float64)
	limit, _ := req.Data["limit"].(float64)
	orderBy, _ := req.Data["orderBy"].(string)
	username, _ := req.Data["username"].(string)

	if start < 1 {
		start = 1
	}

	switch {
	case limit < 1:
		limit = SNIPPETS_LIMIT_DEFAULT

	case limit > SNIPPETS_LIMIT_MAX:
		limit = SNIPPETS_LIMIT_MAX
	}

	orderBy, _ = snippetsOrderBy[strings.ToLower(orderBy)]
	if orderBy == "" {
		orderBy = snippetsOrderBy["updatedDesc"] + ", " + snippetsOrderBy["createdDesc"]
	}

	snips, err := snippetsFetch(db, start, limit, orderBy, username)
	if err != nil {
		return &internalServerError{"Could not fetch snippets", err}
	}

	resp["snippets"] = snips

	return nil
}

func apiSnippetsSearch(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	term, _ := req.Data["term"].(string)
	orderBy, _ := req.Data["orderBy"].(string)

	orderBy, _ = snippetsOrderBy[strings.ToLower(orderBy)]
	if orderBy == "" {
		orderBy = snippetsOrderBy["updatedDesc"] + ", " + snippetsOrderBy["createdDesc"]
	}

	snips, err := snippetsSearch(db, orderBy, term)
	if err != nil {
		return &internalServerError{"Could not fetch snippets", err}
	}

	resp["snippets"] = snips

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
