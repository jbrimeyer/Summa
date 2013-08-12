package summa

import (
	"database/sql"
	_ "go-sqlite3"
)

type SnippetComment struct {
	ID          int64  `json:"id"`
	SnippetID   int64  `json:"-"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Message     string `json:"message"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
}

type SnippetFile struct {
	SnippetID int64  `json:"-"`
	Filename  string `json:"filename"`
	Language  string `json:"language"`
	Contents  string `json:"contents,omitempty"`
}

type SnippetComments []SnippetComment
type SnippetFiles []SnippetFile

type Snippet struct {
	ID          int64            `json:"id"`
	Username    string           `json:"username"`
	DisplayName string           `json:"display_name"`
	Description string           `json:"description"`
	Created     int64            `json:"created"`
	Updated     int64            `json:"updated"`
	Files       *SnippetFiles    `json:"files,omitempty"`
	Comments    *SnippetComments `json:"comments,omitempty"`
	Revisions   []string         `json:"revisions,omitempty"`
}

type Snippets []Snippet

func SnippetsUnread(db *sql.DB, username string) (*Snippets, error) {
	var snippets Snippets

	rows, err := db.Query(
		"SELECT snippet_id,username,display_name,description,created,updated "+
			"FROM snippet JOIN user USING (username) WHERE snippet_id NOT IN "+
			"(SELECT snippet_id FROM snippet_view WHERE username=?)",
		username,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var snippet Snippet

		rows.Scan(
			&snippet.ID,
			&snippet.Username,
			&snippet.DisplayName,
			&snippet.Description,
			&snippet.Created,
			&snippet.Updated,
		)

		snippets = append(snippets, snippet)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return &snippets, nil
}

func SnippetFetch(db *sql.DB, id int64) (*Snippet, error) {
	var snippet Snippet

	row := db.QueryRow(
		"SELECT snippet_id,username,display_name,description,created,updated "+
			"FROM snippet JOIN user USING (username) WHERE snippet_id=?",
		id,
	)

	err := row.Scan(
		&snippet.ID,
		&snippet.Username,
		&snippet.DisplayName,
		&snippet.Description,
		&snippet.Created,
		&snippet.Updated,
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	snippet.Files, err = SnippetFetchFiles(db, id)
	if err != nil {
		return nil, err
	}

	snippet.Comments, err = SnippetFetchComments(db, id)
	if err != nil {
		return nil, err
	}

	return &snippet, nil
}

func SnippetFetchComments(db *sql.DB, id int64) (*SnippetComments, error) {
	var comments SnippetComments

	rows, err := db.Query(
		"SELECT comment_id,snippet_id,username,display_name,message,created,updated FROM "+
			"snippet_comment JOIN user USING (username) WHERE snippet_id=? ORDER BY created",
		id,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var comment SnippetComment

		rows.Scan(
			&comment.ID,
			&comment.SnippetID,
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

func SnippetFetchFiles(db *sql.DB, id int64) (*SnippetFiles, error) {
	var files SnippetFiles

	rows, err := db.Query(
		"SELECT snippet_id,filename,language FROM "+
			"snippet_file WHERE snippet_id=? ORDER BY filename",
		id,
	)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var file SnippetFile

		rows.Scan(
			&file.SnippetID,
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
