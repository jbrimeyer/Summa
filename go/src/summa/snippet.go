package summa

import (
	"database/sql"
	_ "go-sqlite3"
)

type snippetComment struct {
	ID          int64  `json:"id"`
	SnippetID   int64  `json:"-"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Message     string `json:"message"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
}

type snippetFile struct {
	SnippetID int64  `json:"-"`
	Filename  string `json:"filename"`
	Language  string `json:"language"`
	Contents  string `json:"contents,omitempty"`
}

type snippetComments []snippetComment
type snippetFiles []snippetFile

type snippet struct {
	ID          int64            `json:"-"`
	ID36        string           `json:"id"`
	Username    string           `json:"username"`
	DisplayName string           `json:"display_name"`
	Description string           `json:"description"`
	Created     int64            `json:"created"`
	Updated     int64            `json:"updated"`
	Files       *snippetFiles    `json:"files,omitempty"`
	Comments    *snippetComments `json:"comments,omitempty"`
	Revisions   []string         `json:"revisions,omitempty"`
}

type snippets []snippet

func snippetsUnread(db *sql.DB, username string) (*snippets, error) {
	var snips snippets

	rows, err := db.Query(
		"SELECT id_base36,username,display_name,description,created,updated "+
			"FROM snippet JOIN user USING (username) WHERE snippet_id NOT IN "+
			"(SELECT snippet_id FROM snippet_view WHERE username=?)",
		username,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var snip snippet

		rows.Scan(
			&snip.ID36,
			&snip.Username,
			&snip.DisplayName,
			&snip.Description,
			&snip.Created,
			&snip.Updated,
		)

		snips = append(snips, snip)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &snips, nil
}

func snippetFetch(db *sql.DB, id int64) (*snippet, error) {
	var snip snippet

	row := db.QueryRow(
		"SELECT id_base36,username,display_name,description,created,updated "+
			"FROM snippet JOIN user USING (username) WHERE snippet_id=?",
		id,
	)

	err := row.Scan(
		&snip.ID36,
		&snip.Username,
		&snip.DisplayName,
		&snip.Description,
		&snip.Created,
		&snip.Updated,
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	snip.Files, err = snippetFetchFiles(db, id)
	if err != nil {
		return nil, err
	}

	snip.Comments, err = snippetFetchComments(db, id)
	if err != nil {
		return nil, err
	}

	return &snip, nil
}

func snippetFetchComments(db *sql.DB, id int64) (*snippetComments, error) {
	var comments snippetComments

	rows, err := db.Query(
		"SELECT comment_id,username,display_name,message,created,updated FROM "+
			"snippet_comment JOIN user USING (username) WHERE snippet_id=? ORDER BY created",
		id,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment snippetComment

		rows.Scan(
			&comment.ID,
			&comment.Username,
			&comment.DisplayName,
			&comment.Message,
			&comment.Created,
			&comment.Updated,
		)

		comments = append(comments, comment)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &comments, nil
}

func snippetFetchFiles(db *sql.DB, id int64) (*snippetFiles, error) {
	var files snippetFiles

	rows, err := db.Query(
		"SELECT filename,language FROM "+
			"snippet_file WHERE snippet_id=? ORDER BY filename",
		id,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var file snippetFile

		rows.Scan(
			&file.Filename,
			&file.Language,
		)

		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &files, nil
}
