package summa

import (
	"database/sql"
	_ "go-sqlite3"
	"strings"
)

func apiCommentCreate(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	var comment snippetComment

	comment.SnippetID, _ = req.Data["snippet_id"].(string)
	if comment.SnippetID == "" {
		return &badRequestError{"The 'snippet_id' field must be a string"}
	}

	exists, err := snippetExists(db, comment.SnippetID)
	if err != nil {
		return &internalServerError{"Could not check if snippet exists", err}
	}

	if !exists {
		return &badRequestError{"No such snippet"}
	}

	comment.Markdown, _ = req.Data["message"].(string)
	if strings.TrimSpace(comment.Markdown) == "" {
		return &conflictError{apiResponseData{"field": "message"}}
	}

	comment.Username = req.Username
	err = snippetCommentCreate(db, &comment)
	if err != nil {
		return &internalServerError{"Could not create comment", err}
	}

	comment.DisplayName = req.User.DisplayName
	resp["comment"] = comment

	snippetMarkUnread(db, comment.SnippetID)
	snippetMarkReadBy(db, comment.SnippetID, comment.Username)

	return nil
}

func apiCommentUpdate(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	id, _ := req.Data["id"].(string)
	if id == "" {
		return &badRequestError{"The 'id' field must be a string"}
	}

	owned, err := snippetCommentIsOwnedBy(db, id, req.Username)
	if err != nil {
		return &internalServerError{"Could not check comment ownership", err}
	}

	if !owned {
		return &forbiddenError{"You do not have permission to delete this comment"}
	}

	message, _ := req.Data["message"].(string)
	if strings.TrimSpace(message) == "" {
		return &conflictError{apiResponseData{"field": "message"}}
	}

	comment, err := snippetCommentFetch(db, id)
	if err != nil {
		return &internalServerError{"Could not fetch comment", err}
	}

	comment.Markdown = message

	err = snippetCommentUpdate(db, comment)
	if err != nil {
		return &internalServerError{"Could not update comment", err}
	}

	resp["comment"] = comment

	snippetMarkUnread(db, comment.SnippetID)
	snippetMarkReadBy(db, comment.SnippetID, comment.Username)

	return nil
}

func apiCommentDelete(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	id, ok := req.Data["id"].(string)

	if !ok {
		return &badRequestError{"The 'id' field must be a string"}
	}

	owned, err := snippetCommentIsOwnedBy(db, id, req.Username)
	if err != nil {
		return &internalServerError{"Could not check comment ownership", err}
	}

	if !owned {
		return &forbiddenError{"You do not have permission to delete this comment"}
	}

	err = snippetCommentDelete(db, id)
	if err != nil {
		return &internalServerError{"Could not delete comment", err}
	}

	return nil
}
