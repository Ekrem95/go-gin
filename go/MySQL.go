package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	_ "golang.org/x/crypto/bcrypt"
)

var db *sql.DB
var err error

// MySQL func
func MySQL() {
	db, err = sql.Open("mysql", "database:123456@/golang")
	if err != nil {
		panic(err.Error())
	}
	// sql.DB should be long lived "defer" closes it once this function ends
	// defer db.Close()

	// Test the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

}
