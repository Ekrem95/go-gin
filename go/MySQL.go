package main

import (
	"database/sql"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

var db *sql.DB
var err error

var smts = []string{
	`
	CREATE TABLE IF NOT EXISTS users (
	id INT(11) NOT NULL AUTO_INCREMENT,
	username varchar(255),
	password varchar(255),
	primary key (id) )
	`,
	`
	CREATE TABLE IF NOT EXISTS posts (
	id INT(11) NOT NULL AUTO_INCREMENT,
	title varchar(255),
	src varchar(255),
	description varchar(255),
	likes int(11) DEFAULT 0,
	posted_by varchar(255),
	primary key (id) )
	`,
}

// MySQL func
func MySQL() {
	envErr := godotenv.Load()
	if envErr != nil {
		log.Fatal("Error loading .env file")
	}

	database := os.Getenv("mysql")

	db, err = sql.Open("mysql", database)
	if err != nil {
		panic(err.Error())
	}

	for _, smt := range smts {
		if _, err = db.Exec(smt); err != nil {
			panic(err.Error())
		}
	}

	// sql.DB should be long lived "defer" closes it once this function ends
	// defer db.Close()

	// Test the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

}
