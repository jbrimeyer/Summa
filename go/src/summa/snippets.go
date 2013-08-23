package summa

import (
	"database/sql"
	"fmt"
	_ "go-sqlite3"
)

type snippets []snippet

func snippetsFetch(db *sql.DB, start, limit float64, orderBy, username string) (*snippets, error) {
	var snips snippets

	whereClause := ""
	if username != "" {
		whereClause = "WHERE snippet.username=?"
	}
	query := fmt.Sprintf(
		"SELECT snippet.snippet_id,snippet.username,display_name,description,snippet.created,snippet.updated,"+
			"COUNT(snippet_file.snippet_id) files,COUNT(snippet_comment.snippet_id) comments "+
			"FROM snippet JOIN user USING (username) JOIN snippet_file USING (snippet_id) "+
			"LEFT JOIN snippet_comment USING (snippet_id) %s GROUP BY snippet.snippet_id "+
			"ORDER BY %s LIMIT %d OFFSET %d",
		whereClause,
		orderBy,
		int(limit),
		int(start-1),
	)

	var rows *sql.Rows
	var err error
	if whereClause == "" {
		rows, err = db.Query(
			query,
		)
	} else {
		rows, err = db.Query(
			query,
			username,
		)
	}

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

// snippetsUnread will return unread snippets for a specific user
func snippetsUnread(db *sql.DB, username string) (*snippets, error) {
	var snips snippets

	rows, err := db.Query(
		"SELECT snippet.snippet_id,snippet.username,display_name,description,snippet.created,snippet.updated,"+
			"COUNT(snippet_file.snippet_id) files,COUNT(snippet_comment.snippet_id) comments "+
			"FROM snippet JOIN user USING (username) JOIN snippet_file USING (snippet_id) "+
			"LEFT JOIN snippet_comment USING (snippet_id) WHERE snippet.snippet_id NOT IN "+
			"(SELECT snippet_id FROM snippet_view sv WHERE sv.username=?)",
		username,
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
