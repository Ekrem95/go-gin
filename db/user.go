package db

import (
	"database/sql"
	"fmt"
)

// Exists ...
func (u User) Exists(username string) bool {
	var user string

	smt := fmt.Sprintf("SELECT username FROM users WHERE username = '%s'", username)
	if err := QueryRowScan(smt, &user); err == sql.ErrNoRows {
		// username is available
		return false
	}
	return true
}
