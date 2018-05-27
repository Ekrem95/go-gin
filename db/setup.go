package db

import "os"

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
	`
	CREATE TABLE IF NOT EXISTS comments (
	id INT(11) NOT NULL AUTO_INCREMENT,
	text varchar(255),
	post_id varchar(11),
	time INT(22),
	sender varchar(255),
	primary key (id) )
	`,
	`
	CREATE TABLE IF NOT EXISTS post_likes (
	id INT(11) NOT NULL AUTO_INCREMENT,
	post_id varchar(11),
	user varchar(11),
	primary key (id) )
	`,
}

// TestSQLConnection ...
func TestSQLConnection() error {
	if os.Getenv("ENV") == "TEST" {
		dsn = "root:pass@/go_gin_test"
	} else {
		dsn = "root:pass@/go_gin"
	}

	db, err := open()
	if err != nil {
		return err
	}
	// sql.DB should be long lived "defer" closes it once this function ends
	defer db.Close()

	// Test the connection to the database
	if err = db.Ping(); err != nil {
		return err
	}

	for _, smt := range smts {
		if _, err = db.Exec(smt); err != nil {
			return err
		}
	}

	return nil
}
