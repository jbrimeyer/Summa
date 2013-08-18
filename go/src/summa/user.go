package summa

import (
	"database/sql"
	_ "go-sqlite3"
)

type User struct {
	Username    string `json:"username"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
}

func userExists(db *sql.DB, username string) (bool, error) {
	var count int64
	row := db.QueryRow(
		"SELECT COUNT(*) FROM user WHERE username=?",
		username,
	)
	err := row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count == 1, nil
}

func userFetch(db *sql.DB, username string) (*User, error) {
	var u User

	row := db.QueryRow(
		"SELECT username,display_name,email FROM user WHERE username=?",
		username,
	)

	err := row.Scan(
		&u.Username,
		&u.DisplayName,
		&u.Email,
	)

	switch {
	case err == sql.ErrNoRows:
		return nil, nil
	case err != nil:
		return nil, err
	}

	return &u, nil
}

func userCreate(db *sql.DB, u *User) error {
	_, err := db.Exec(
		"INSERT INTO user VALUES (?,?,?)",
		u.Username,
		u.DisplayName,
		u.Email,
	)

	return err
}

func userUpdate(db *sql.DB, u *User) error {
	_, err := db.Exec(
		"UPDATE user SET display_name=?,email=? WHERE username=?",
		u.DisplayName,
		u.Email,
		u.Username,
	)

	return err
}
