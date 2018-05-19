package db

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

// var db *sql.DB
// var err error
var dsn string

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

// Exec ...
func Exec(smt string, args ...interface{}) (sql.Result, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	res, err := db.Exec(smt, args...)
	if err != nil {
		return nil, err
	}
	return res, nil
}

// Query ...
func Query(smt string, args ...interface{}) (*sql.Rows, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(smt, args...)
	if err != nil {
		return nil, err
	}

	return rows, nil
}

// QueryRowScan ...
func QueryRowScan(smt string, dest ...interface{}) error {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	err = db.QueryRow(smt).Scan(dest...)
	if err != nil {
		return err
	}

	return nil
}

// TestSQLConnection ...
func TestSQLConnection() error {
	if os.Getenv("ENV") == "TEST" {
		dsn = "root:pass@/go_gin_test"
	} else {
		dsn = "root:pass@/go_gin"
	}

	db, err := sql.Open("mysql", dsn)
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

// RedisGetMsgs func
func RedisGetMsgs(c *gin.Context) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	// Grabs the entire users list into an []string named users
	messages, _ := redis.Strings(conn.Do("LRANGE", "messagestest", 0, -1))
	// Grab one string value and convert it to type byte
	// Then decode the data into unencoded
	// json.Unmarshal([]byte(messages[0]), &message)
	// fmt.Println(unencoded.Name)

	c.JSON(200, gin.H{
		"messages": messages,
	})
}

// RedisSaveMsg func
func RedisSaveMsg(msg *Message) {
	conn, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	// Send Redis a ping command and wait for a pong
	message1 := msg
	// message1 := Message{"Hello", "1502982731", "tormond"}
	//
	encoded, _ := json.Marshal(message1)
	//
	conn.Do("LPUSH", "messagestest", encoded)
	conn.Do("LTRIM", "messagestest", 0, 99)
}
