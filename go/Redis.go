package main

import (
	"encoding/json"
	"fmt"
	"github.com/garyburd/redigo/redis"
	_ "log"
	_ "reflect"
)

// Redis func
func Redis() {
	conn, error := redis.Dial("tcp", ":6379")
	if error != nil {
		panic(error.Error())
	}
	defer conn.Close()

	// Send Redis a ping command and wait for a pong
	// user1 := User{"John", "22"}
	//
	// encoded, _ := json.Marshal(user1)
	//
	// conn.Do("LPUSH", "userstest", encoded)
	// conn.Do("LPOP", "userstest")

	var unencoded *User

	// Grabs the entire users list into an []string named users
	users, _ := redis.Strings(conn.Do("LRANGE", "userstest", 0, -1))
	// Grab one string value and convert it to type byte
	// Then decode the data into unencoded
	json.Unmarshal([]byte(users[1]), &unencoded)
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

	for _, v := range users {
		// log.Printf("value at [%d]=%v", i, v)
		json.Unmarshal([]byte(v), &unencoded)
		// fmt.Println(reflect.TypeOf(unencoded))
		// fmt.Println(unencoded.Password)
		if unencoded.Password == "26" {
			fmt.Println("123333333333331222222222222222")
		}
	}
}
