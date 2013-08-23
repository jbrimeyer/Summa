package summa

import (
	"database/sql"
	_ "go-sqlite3"
)

type snippetComment struct {
	ID          int64  `json:"id"`
	SnippetID   string `json:"-"`
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Markdown    string `json:"markdown"`
	HTML        string `json:"html"`
	Created     int64  `json:"created"`
	Updated     int64  `json:"updated"`
}

type snippetComments []snippetComment

// snippetCommentExists returns true if the comment with the given id exists in the database
func snippetCommentExists(db *sql.DB, id string) (bool, error) {
	var count int64
	row := db.QueryRow(
		"SELECT COUNT(*) FROM snippet_comment WHERE comment_id=?",
		id,
	)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

// snippetCommentIsOwnedBy returns true if the comment with the given id is
// owned by the given username
func snippetCommentIsOwnedBy(db *sql.DB, id, username string) (bool, error) {
	var count int64
	row := db.QueryRow(
		"SELECT COUNT(*) FROM snippet_comment WHERE comment_id=? AND username=?",
		id,
		username,
	)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

// snippetCommentFetch will fetch an individual comment by ID
func snippetCommentFetch(db *sql.DB, id string) (*snippetComment, error) {
	var comment snippetComment

	row := db.QueryRow(
		"SELECT comment_id,snippet_id,username,display_name,markdown,html,created,updated "+
			"FROM snippet_comment JOIN user USING (username) WHERE comment_id=?",
		id,
	)

	err := row.Scan(
		&comment.ID,
		&comment.SnippetID,
		&comment.Username,
		&comment.DisplayName,
		&comment.Markdown,
		&comment.HTML,
		&comment.Created,
		&comment.Updated,
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	return &comment, nil
}

// snippetCommentCreate will create a new comment in the database
func snippetCommentCreate(db *sql.DB, comment *snippetComment) error {
	var err error

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	comment.Created = UnixMilliseconds()
	comment.HTML = markdownParse(comment.Markdown)

	result, err := db.Exec(
		"INSERT INTO snippet_comment VALUES (NULL,?,?,?,?,?,0)",
		comment.SnippetID,
		comment.Username,
		comment.Markdown,
		comment.HTML,
		comment.Created,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	comment.ID, err = result.LastInsertId()
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// snippetCommentUpdate will update an existing comment in the database
func snippetCommentUpdate(db *sql.DB, comment *snippetComment) error {
	var err error

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	comment.Updated = UnixMilliseconds()

	_, err = db.Exec(
		"UPDATE snippet_comment SET markdown=?,html=?,updated=? WHERE comment_id=?",
		comment.Markdown,
		markdownParse(comment.Markdown),
		comment.Updated,
		comment.ID,
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()

	return nil
}

// snippetCommentDelete permanently removes a comment from the database
func snippetCommentDelete(db *sql.DB, id string) error {
	queries := []string{
		"DELETE FROM snippet_comment WHERE comment_id=?",
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

	return nil
}
