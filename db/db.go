package db

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
)

var dsn string

// RedisAddress ...
var RedisAddress string

func init() {
	redis := os.Getenv("REDIS_ADDRESS")
	if redis != "" {
		RedisAddress = redis
	} else {
		RedisAddress = "localhost:6379"
	}
}

func open() (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}

// Exec ...
func Exec(smt string, args ...interface{}) (sql.Result, error) {
	db, err := open()
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
	db, err := open()
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
	db, err := open()
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

// RedisGetMsgs func
func RedisGetMsgs(c *gin.Context) {
	conn, err := redis.Dial("tcp", RedisAddress)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()
	// Grabs the entire messages list into an []string named messages
	messages, _ := redis.Strings(conn.Do("LRANGE", "messages", 0, -1))

	c.JSON(200, gin.H{"messages": messages})
}

// RedisSaveMsg func
func RedisSaveMsg(msg *Message) {
	conn, err := redis.Dial("tcp", RedisAddress)
	if err != nil {
		panic(err.Error())
	}
	defer conn.Close()

	encoded, _ := json.Marshal(msg)

	conn.Do("LPUSH", "messages", encoded)
	conn.Do("LTRIM", "messages", 0, 99)
}
