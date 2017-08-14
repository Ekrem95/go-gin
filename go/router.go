package main

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func router() {
	Redis()

	router := gin.Default()
	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("session", store))
	router.LoadHTMLGlob("../templates/*")
	router.Static("/src", "../src")

	router.GET("/r/:r", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})

	router.GET("/j/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	router.POST("/form_post", func(c *gin.Context) {
		session := sessions.Default(c)
		var count int
		v := session.Get("count")
		if v == nil {
			count = 0
		} else {
			count = v.(int)
			count++
		}
		session.Set("count", count)
		session.Save()

		name := c.PostForm("name")
		password, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		newUser := User{name, string(password)}
		fmt.Println(newUser)
		// encoded, _ := json.Marshal(newUser)

		c.JSON(200, gin.H{
			"user":  newUser,
			"err":   err,
			"count": count,
		})
	})

	router.Run(":8080")
}
