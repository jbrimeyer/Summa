package summa

import (
	"database/sql"
	_ "go-sqlite3"
	"strings"
)

func apiProfile(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	username, _ := req.Data["username"].(string)

	if username == "" {
		username = req.Username
	}

	u, err := userFetch(db, username)
	if err != nil {
		return &internalServerError{"Could not fetch user", err}
	}

	if u == nil {
		return &notFoundError{"User does not exist"}
	}

	resp["user"] = u

	return nil
}

func apiProfileUpdate(db *sql.DB, req apiRequest, resp apiResponseData) apiError {
	var u User

	name, _ := req.Data["displayName"].(string)
	email, _ := req.Data["email"].(string)

	u.Username = req.Username
	u.DisplayName = strings.TrimSpace(name)
	u.Email = strings.TrimSpace(email)

	if u.DisplayName == "" {
		return &conflictError{apiResponseData{"field": "displayName"}}
	}

	if u.Email == "" {
		return &conflictError{apiResponseData{"field": "email"}}
	}

	err := userUpdate(db, &u)
	if err != nil {
		return &internalServerError{"Could not update user", err}
	}

	return nil
}
