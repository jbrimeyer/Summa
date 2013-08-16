package summa

import (
	"database/sql"
	"fmt"
	_ "go-sqlite3"
	"regexp"
	"strings"
)

func apiSnippet(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	id, ok := req.Data["id"].(string)

	if !ok {
		return &badRequestError{"The 'id' field must be a string"}
	}

	snippet, err := snippetFetch(db, id)
	if err != nil {
		return &internalServerError{"Could not fetch snippet", err}
	}

	if snippet == nil {
		return &notFoundError{"No such snippet"}
	}

	var markRead bool
	switch req.Data["markRead"].(type) {
	case bool:
		markRead = req.Data["markRead"].(bool)

	case string:
		markRead = req.Data["markRead"].(string) != ""
	}

	if markRead {
		err := snippetMarkReadBy(db, id, req.Username)
		if err != nil {
			return &internalServerError{"Could not mark snippet read", err}
		}
	}

	resp["snippet"] = snippet

	return nil
}

func apiSnippetCreate(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	reqForFiles := []string{"filename", "language", "contents"}
	var snip snippet

	snip.Description, _ = req.Data["description"].(string)
	if snip.Description == "" {
		return &conflictError{apiResponseData{"field": "description"}}
	}

	fileRegex := regexp.MustCompile("(?i)^[a-z0-9_.-]+$")
	var files snippetFiles

	switch req.Data["files"].(type) {
	case []interface{}:
		filenames := make(map[string]bool)
		for i, v := range req.Data["files"].([]interface{}) {
			switch v.(type) {
			case map[string]interface{}:
				vmap := v.(map[string]interface{})
				fields := make(map[string]string)
				var file snippetFile

				for _, required := range reqForFiles {
					strVal, ok := vmap[required].(string)
					if !ok || strings.TrimSpace(strVal) == "" {
						return &conflictError{apiResponseData{"field": fmt.Sprintf("file[%d].%s", i, required)}}
					}
					fields[required] = strings.TrimSpace(strVal)
				}

				lcFilename := strings.ToLower(fields["filename"])
				_, ok := filenames[lcFilename]
				if !fileRegex.MatchString(fields["filename"]) || ok {
					return &conflictError{apiResponseData{"field": fmt.Sprintf("file[%d].filename", i)}}
				}

				filenames[lcFilename] = true

				file.Filename = fields["filename"]
				file.Language = fields["language"]
				file.Contents = fields["contents"]

				files = append(files, file)
			default:
				return &badRequestError{"'files' field is malformed"}
			}
		}
	default:
		return &conflictError{apiResponseData{"field": "files"}}
	}

	snip.Files = files
	snip.Username = req.Username

	id, err := snippetCreate(db, &snip, req.User)
	if err != nil {
		return &internalServerError{"Could not create snippet", err}
	}

	snippetMarkReadBy(db, id, req.Username)

	resp["id"] = id

	return nil
}

func apiSnippetUpdate(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	// TODO: apiSnippetUpdate
	return nil
}

func apiSnippetDelete(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	id, ok := req.Data["id"].(string)

	if !ok {
		return &badRequestError{"The 'id' field must be a string"}
	}

	owned, err := snippetIsOwnedBy(db, id, req.Username)
	if err != nil {
		return &internalServerError{"Could not check snippet ownership", err}
	}

	if !owned {
		return &forbiddenError{"You do not have permission to delete this snippet"}
	}

	err = snippetDelete(db, id)
	if err != nil {
		return &internalServerError{"Could not delete snippet", err}
	}

	return nil
}
