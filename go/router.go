package main

import (
	"fmt"
	"github.com/garyburd/redigo/redis"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
)

func router() {
	c, err := redis.Dial("tcp", ":6379")
	if err != nil {
		panic(err)
	}

	defer c.Close()

	// hashedMessage, err := bcrypt.GenerateFromPassword([]byte("eko"), bcrypt.DefaultCost)
	// if err != nil {
	// 	panic(err)
	// }
	// c.Do("SET", "message1", hashedMessage)

	world, err := redis.String(c.Do("GET", "message1"))
	if err != nil {
		fmt.Println("key not found")
	}

	fmt.Println(world)

	router := gin.Default()
	router.LoadHTMLGlob("../templates/*")
	router.Static("/src", "../src")

	router.GET("/r/:r", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/j/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/form_post", func(c *gin.Context) {
		name := c.PostForm("name")
		age, _ := strconv.Atoi(c.PostForm("age"))

		user1 := User{name, age}

		err = bcrypt.CompareHashAndPassword([]byte(world), []byte(name))

		if err == nil {
			fmt.Println("chelsea")
		} else {
			fmt.Println("City")
		}

		c.JSON(200, gin.H{
			"user": user1,
			"err":  err,
		})
	})

	router.Run(":8080")
}
