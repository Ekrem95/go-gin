package main

import (
	"database/sql"
	"encoding/json"
	"os"

	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
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
	`
	CREATE TABLE IF NOT EXISTS comments (
	id INT(11) NOT NULL AUTO_INCREMENT,
	text varchar(255),
	postId varchar(11),
	time INT(22),
	sender varchar(255),
	primary key (id) )
	`,
}

// MySQL func
func MySQL() {
	// database := os.Getenv("mysql")

	var dsn string

	if os.Getenv("ENV") == "TEST" {
		dsn = "root:secret@/go_gin_test"
	} else {
		dsn = "root:secret@/go_gin"
	}

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		panic(err.Error())
	}

	for _, smt := range smts {
		if _, err = db.Exec(smt); err != nil {
			panic(err.Error())
		}
	}

	// sql.DB should be long lived "defer" closes it once this function ends
	defer db.Close()

	// Test the connection to the database
	err = db.Ping()
	if err != nil {
		panic(err.Error())
	}

}

// RedisGetMsgs func
func RedisGetMsgs(c *gin.Context) {
	conn, error := redis.Dial("tcp", ":6379")
	if error != nil {
		panic(error.Error())
	}
	defer conn.Close()

	// var message *Message

	// Grabs the entire users list into an []string named users
	messages, _ := redis.Strings(conn.Do("LRANGE", "messagestest", 0, -1))
	// Grab one string value and convert it to type byte
	// Then decode the data into unencoded
	// json.Unmarshal([]byte(messages[0]), &message)
	// fmt.Println(unencoded.Name)

	// var user User
	//
	// json.Unmarshal([]byte(users[0]), &user)
	// if err == nil {
	// 	fmt.Printf("%+v\n", user)
	// 	fmt.Println(reflect.TypeOf(user))
	// } else {
	// 	fmt.Println(err)
	// 	fmt.Printf("%+v\n", user)
	// }
	//
	// fmt.Println(users[0])
	// fmt.Println(reflect.TypeOf(users[0]))

	// for _, v := range messages {
	// 	// log.Printf("value at [%d]=%v", i, v)
	// 	json.Unmarshal([]byte(v), &message)
	// 	// fmt.Println(reflect.TypeOf(unencoded))
	// 	// fmt.Println(unencoded.Password)
	// 	fmt.Println("*****************************")
	// 	fmt.Println(message)
	// 	fmt.Println(reflect.TypeOf(message))
	// }

	c.JSON(200, gin.H{
		"messages": messages,
	})
}

// RedisSaveMsg func
func RedisSaveMsg(msg *Message) {
	conn, error := redis.Dial("tcp", ":6379")
	if error != nil {
		panic(error.Error())
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
	//conn.Do("LPOP", "messagestest")

	// var message *Message

	// Grabs the entire users list into an []string named users
	// messages, _ := redis.Strings(conn.Do("LRANGE", "messagestest", 0, -1))
	// Grab one string value and convert it to type byte
	// Then decode the data into unencoded
	// json.Unmarshal([]byte(messages[0]), &message)
	// fmt.Println(unencoded.Name)

	// var user User
	//
	// json.Unmarshal([]byte(users[0]), &user)
	// if err == nil {
	// 	fmt.Printf("%+v\n", user)
	// 	fmt.Println(reflect.TypeOf(user))
	// } else {
	// 	fmt.Println(err)
	// 	fmt.Printf("%+v\n", user)
	// }
	//
	// fmt.Println(users[0])
	// fmt.Println(reflect.TypeOf(users[0]))

	// for _, v := range messages {
	// 	// log.Printf("value at [%d]=%v", i, v)
	// 	json.Unmarshal([]byte(v), &message)
	// 	// fmt.Println(reflect.TypeOf(unencoded))
	// 	// fmt.Println(unencoded.Password)
	// 	fmt.Println("*****************************")
	// 	fmt.Println(message)
	// }
}
