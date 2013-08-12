package summa

import (
	"bytes"
	"crypto/sha1"
	"database/sql"
	"encoding/binary"
	"fmt"
	_ "go-sqlite3"
	"math/rand"
	"time"
)

// isValidSession checks to determine if a given username and token
// combine to make a valid session
func isValidSession(db *sql.DB, username, token string) (bool, error) {
	expired := UnixMilliseconds() - config.SessionExpire

	// Remove expired sessions
	_, err := db.Exec(
		"DELETE FROM user_session WHERE timestamp <= ?",
		expired,
	)
	if err != nil {
		return false, err
	}

	row := db.QueryRow(
		"SELECT COUNT(*) FROM user_session WHERE username=? AND token=?",
		username,
		token,
	)

	var count int64
	err = row.Scan(&count)
	if err != nil {
		return false, err
	}

	return count != 0, nil
}

// createSession generates a random session token for
// a given username and stores it in the database
func createSession(db *sql.DB, username string) (string, error) {
	rand.Seed(time.Now().UnixNano())
	buf := new(bytes.Buffer)
	binary.Write(buf, binary.LittleEndian, rand.Int63())

	hasher := sha1.New()
	hasher.Write(buf.Bytes())
	hasher.Write([]byte(username))

	token := fmt.Sprintf("%x", hasher.Sum(nil))

	_, err := db.Exec(
		"REPLACE INTO user_session VALUES (?,?,?)",
		username,
		token,
		UnixMilliseconds(),
	)

	return token, err
}

// removeSesion removes a session with a given username
// and token from the database
func removeSession(db *sql.DB, username, token string) error {
	_, err := db.Exec(
		"DELETE FROM user_session WHERE username=? AND token=?",
		username,
		token,
	)

	return err
}
