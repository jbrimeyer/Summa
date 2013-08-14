package summa

import (
	"database/sql"
	_ "go-sqlite3"
)

type snippetComment struct {
	ID          int64  `json:"id"`
	SnippetID   string `json:"-"`
	Username    string `json:"username"`
	DisplayName string `json:"display_name"`
	Message     string `json:"message"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
}

type snippetFile struct {
	SnippetID string `json:"-"`
	Filename  string `json:"filename"`
	Language  string `json:"language"`
	Contents  string `json:"contents,omitempty"`
}

type snippetComments []snippetComment
type snippetFiles []snippetFile

type snippet struct {
	ID          string          `json:"id"`
	Username    string          `json:"username"`
	DisplayName string          `json:"display_name"`
	Description string          `json:"description"`
	Created     int64           `json:"created"`
	Updated     int64           `json:"updated"`
	Files       snippetFiles    `json:"files,omitempty"`
	Comments    snippetComments `json:"comments,omitempty"`
	Revisions   []string        `json:"revisions,omitempty"`
}

type snippets []snippet

// snippetCreate will create a new snippet and return it's id
func snippetCreate(db *sql.DB, snip *snippet) (string, error) {
	var err error

	tx, err := db.Begin()
	if err != nil {
		return "", err
	}

	defer (func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	})()

	ms := UnixMilliseconds()
	var id string
	for {
		id = Reverse(ToBase36(ms))
		var count int64
		row := db.QueryRow("SELECT COUNT(*) FROM snippet WHERE snippet_id=?", id)
		err = row.Scan(&count)

		if err != nil {
			return "", err
		}

		if count == 0 {
			break
		}

		ms--
	}

	_, err = db.Exec(
		"INSERT INTO snippet VALUES (?,?,?,?,0)",
		id,
		snip.Username,
		snip.Description,
		ms,
	)
	if err != nil {
		return "", err
	}

	for _, file := range snip.Files {
		_, err = db.Exec(
			"INSERT INTO snippet_file VALUES (?,?,?)",
			id,
			file.Filename,
			file.Language,
		)
		if err != nil {
			return "", err
		}
	}

	err = repoCreate(id, nil, snip.Files)

	return id, nil
}

// snippetDelete permanently removes a snippet
func snippetDelete(db *sql.DB, id string) error {
	queries := []string{
		"DELETE FROM snippet WHERE snippet_id=?",
		"DELETE FROM snippet_comment WHERE snippet_id=?",
		"DELETE FROM snippet_file WHERE snippet_id=?",
		"DELETE FROM snippet_view WHERE snippet_id=?",
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	for _, q := range queries {
		_, err = tx.Exec(q, id)
		if err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()

	// TODO: Remove git repository

	return nil
}

// snippetIsOwnedBy returns true if the snippet with the given id is
// owned by the given username
func snippetIsOwnedBy(db *sql.DB, id, username string) (bool, error) {
	var count int64
	row := db.QueryRow(
		"SELECT COUNT(*) FROM snippet WHERE snippet_id=? AND username=?",
		id,
		username,
	)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

// snippetsUnread will return unread snippets for a specific user
func snippetsUnread(db *sql.DB, username string) (*snippets, error) {
	var snips snippets

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
		var snip snippet

		rows.Scan(
			&snip.ID,
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

// snippetFetch will fetch an individual snippet by ID
func snippetFetch(db *sql.DB, id string) (*snippet, error) {
	var snip snippet

	row := db.QueryRow(
		"SELECT snippet_id,username,display_name,description,created,updated "+
			"FROM snippet JOIN user USING (username) WHERE snippet_id=?",
		id,
	)

	err := row.Scan(
		&snip.ID,
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

	infoLog.Printf("%+v", snip.Comments)

	return &snip, nil
}

// snippetFetchComments will fetch the comments for a specific snippet
func snippetFetchComments(db *sql.DB, id string) (snippetComments, error) {
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

	return comments, nil
}

// snippetFetchFiles will fetch the files for a sepcific snippet
func snippetFetchFiles(db *sql.DB, id string) (snippetFiles, error) {
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

	return files, nil
}
