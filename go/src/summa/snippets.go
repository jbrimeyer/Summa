package summa

import (
	"database/sql"
	"fmt"
	_ "go-sqlite3"
)

type snippets []snippet

// snippetsFetchGeneric will fetch snippets from the database
func snippetsFetchGeneric(db *sql.DB, query string, params []interface{}) (*snippets, error) {
	var snips snippets

	rows, err := db.Query(
		query,
		params...,
	)

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var snip snippet

		rows.Scan(
			&snip.ID,
			&snip.Username,
			&snip.DisplayName,
			&snip.Description,
			&snip.Created,
			&snip.Updated,
			&snip.NumFiles,
			&snip.NumComments,
		)

		snips = append(snips, snip)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &snips, nil
}

// snippetsSearch will fetch snippets using a search term, sorted by the given value
func snippetsSearch(db *sql.DB, orderBy, term string) (*snippets, error) {
	query := fmt.Sprintf(
		"SELECT s.snippet_id,s.username,u.display_name,s.description,s.created,s.updated,"+
			"COUNT(sf.snippet_id) files,COUNT(sc.snippet_id) comments FROM snippet s JOIN "+
			"user u ON u.username=s.username JOIN snippet_file sf ON s.snippet_id=sf.snippet_id "+
			"JOIN snippet_search ss ON ss.docid=s.search_id LEFT JOIN snippet_comment sc ON "+
			"s.snippet_id=sc.snippet_id WHERE ss.snippet MATCH(?) GROUP BY s.snippet_id ORDER BY %s",
		orderBy,
	)

	params := []interface{}{term}

	return snippetsFetchGeneric(db, query, params)
}

// snippetsFetch will fetch snippets in a given range, sorted by the given value and optionally
// filtered by username
func snippetsFetch(db *sql.DB, start, limit float64, orderBy, username string) (*snippets, error) {
	var params []interface{}

	whereClause := ""
	if username != "" {
		whereClause = "WHERE s.username=?"
		params = append(params, username)
	}
	query := fmt.Sprintf(
		"SELECT s.snippet_id,s.username,display_name,description,s.created,s.updated,"+
			"COUNT(sf.snippet_id) files,COUNT(sc.snippet_id) comments "+
			"FROM snippet s JOIN user u USING (username) JOIN snippet_file sf USING (snippet_id) "+
			"LEFT JOIN snippet_comment sc USING (snippet_id) %s GROUP BY s.snippet_id "+
			"ORDER BY %s LIMIT %d OFFSET %d",
		whereClause,
		orderBy,
		int(limit),
		int(start-1),
	)

	return snippetsFetchGeneric(db, query, params)
}

// snippetsUnread will return unread snippets for a specific user
func snippetsUnread(db *sql.DB, username string) (*snippets, error) {
	query := "SELECT s.snippet_id,s.username,u.display_name,s.description,s.created,s.updated," +
		"COUNT(sf.snippet_id) files,COUNT(sc.snippet_id) comments FROM snippet s JOIN " +
		"user u ON u.username=s.username JOIN snippet_file sf ON s.snippet_id=sf.snippet_id " +
		"LEFT JOIN snippet_comment sc ON s.snippet_id=sc.snippet_id LEFT JOIN snippet_view sv " +
		"ON s.snippet_id=sv.snippet_id AND sv.username=? WHERE sv.snippet_id IS NULL " +
		"GROUP BY s.snippet_id"

	params := []interface{}{username}

	return snippetsFetchGeneric(db, query, params)
}
