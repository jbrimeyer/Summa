package summa

import (
	"database/sql"
	_ "go-sqlite3"
	"io/ioutil"
	"path"
)

type snippetFile struct {
	SnippetID string `json:"-"`
	Filename  string `json:"filename"`
	Language  string `json:"language"`
	Contents  string `json:"contents,omitempty"`
}

type snippetFiles []snippetFile

type snippet struct {
	ID          string          `json:"id"`
	Username    string          `json:"username"`
	DisplayName string          `json:"displayName"`
	Description string          `json:"description"`
	Created     int64           `json:"created"`
	Updated     int64           `json:"updated"`
	Files       snippetFiles    `json:"files,omitempty"`
	Comments    snippetComments `json:"comments,omitempty"`
	Revisions   []string        `json:"revisions,omitempty"`
}

// snippetExists checks is a snippet with the given ID exists
func snippetExists(db *sql.DB, id string) (bool, error) {
	var count int64
	row := db.QueryRow("SELECT COUNT(*) FROM snippet WHERE snippet_id=?", id)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	if count == 0 {
		return false, nil
	}

	return true, nil
}

// snippetCreate will create a new snippet and return it's id
func snippetCreate(db *sql.DB, snip *snippet, u *User) (string, error) {
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
		exists, err := snippetExists(db, id)
		if err != nil {
			return "", err
		}

		if !exists {
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

	err = repoCreate(id, u, snip.Files)

	return id, nil
}

func snippetUpdate(db *sql.DB, oldSnip, newSnip *snippet, u *User) error {
	var err error

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	defer (func() {
		if err == nil {
			tx.Commit()
		} else {
			tx.Rollback()
		}
	})()

	oldSnip.Updated = UnixMilliseconds()
	oldSnip.Description = newSnip.Description

	_, err = db.Exec(
		"UPDATE snippet SET description=?,updated=? WHERE snippet_id=?",
		oldSnip.Description,
		oldSnip.Updated,
		oldSnip.ID,
	)
	if err != nil {
		return err
	}

	_, err = db.Exec("DELETE FROM snippet_file WHERE snippet_id=?", oldSnip.ID)
	if err != nil {
		return err
	}

	for _, file := range newSnip.Files {
		_, err = db.Exec(
			"INSERT INTO snippet_file VALUES (?,?,?)",
			oldSnip.ID,
			file.Filename,
			file.Language,
		)
		if err != nil {
			return err
		}
	}

	err = repoUpdate(oldSnip.ID, u, oldSnip.Files, newSnip.Files)

	oldSnip.Files = newSnip.Files

	return nil
}

// snippetMarkReadBy will mark a snippet with a specified id as read
// by a specific user
func snippetMarkReadBy(db *sql.DB, id, username string) error {
	_, err := db.Exec(
		"REPLACE INTO snippet_view VALUES (?,?)",
		id,
		username,
	)

	return err
}

// snippetMarkUnread will mark a snippet as unread by all users
func snippetMarkUnread(db *sql.DB, id string) error {
	_, err := db.Exec(
		"DELETE FROM snippet_view WHERE snippet_id=?",
		id,
	)

	return err
}

// snippetMarkUnread will mark a snippet as unread by all users
func snippetMarkUnreadBy(db *sql.DB, id, username string) error {
	_, err := db.Exec(
		"DELETE FROM snippet_view WHERE snippet_id=? AND username=?",
		id,
		username,
	)

	return err
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

	repoDelete(id)

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

	return &snip, nil
}

// snippetFetch will fetch an individual snippet by ID, including it's comments
func snippetFetchAll(db *sql.DB, id string) (*snippet, error) {
	snip, err := snippetFetch(db, id)
	if err != nil {
		return nil, err
	}

	if snip == nil {
		return nil, nil
	}

	snip.Comments, err = snippetFetchComments(db, id)
	if err != nil {
		return nil, err
	}

	return snip, nil
}

// snippetFetchComments will fetch the comments for a specific snippet
func snippetFetchComments(db *sql.DB, id string) (snippetComments, error) {
	var comments snippetComments

	rows, err := db.Query(
		"SELECT comment_id,username,display_name,markdown,html,created,updated FROM "+
			"snippet_comment JOIN user USING (username) WHERE snippet_id=? ORDER BY created",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment snippetComment

		rows.Scan(
			&comment.ID,
			&comment.Username,
			&comment.DisplayName,
			&comment.Markdown,
			&comment.HTML,
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
	defer rows.Close()

	fsPath := repoPath(id)

	for rows.Next() {
		var file snippetFile

		rows.Scan(
			&file.Filename,
			&file.Language,
		)

		filePath := path.Join(fsPath, file.Filename)
		contents, err := ioutil.ReadFile(filePath)
		if err != nil {
			return nil, err
		}

		file.Contents = string(contents)

		files = append(files, file)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return files, nil
}
