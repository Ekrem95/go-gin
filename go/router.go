package main

import (
	"database/sql"
	_ "fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

func router() {
	// Redis()
	MySQL()
	router := gin.Default()
	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("session", store))
	router.LoadHTMLGlob("../templates/*")
	router.Static("/src", "../src")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	router.GET("/signup", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	router.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.html", gin.H{
			"title": "Main website",
		})
	})
	router.GET("/user", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	})

	router.POST("/signup", func(c *gin.Context) {
		username := c.PostForm("username")
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}

		var user string

		err = db.QueryRow("SELECT username FROM users WHERE username=?", username).Scan(&user)

		switch {
		// Username is available
		case err == sql.ErrNoRows:
			_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
			if err != nil {
				c.JSON(500, gin.H{
					"error": "Unable to Sign up.",
				})
				return
			}

			session := sessions.Default(c)
			user := username
			session.Set("user", username)
			session.Save()

			c.JSON(200, gin.H{
				"success": true,
				"user":    user,
			})
			return
		case err != nil:
			c.JSON(500, gin.H{
				"error": "An error occured.",
			})
			return
		default:
			c.JSON(200, gin.H{
				"error": "Username already exist.",
			})
		}
	})

	router.POST("/login", func(c *gin.Context) {

		username := c.PostForm("username")
		password := c.PostForm("password")

		var databaseUsername string
		var databasePassword string

		// Search the database for the username provided
		// If it exists grab the password for validation
		err = db.QueryRow("SELECT username, password FROM users WHERE username=?", username).Scan(&databaseUsername, &databasePassword)
		// If not then redirect to the login page
		if err != nil {
			c.JSON(200, gin.H{
				"err": err,
			})
			return
		}

		// Validate the password
		err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
		if err != nil {
			c.JSON(200, gin.H{
				"err":  err,
				"desc": "Passwords do not match",
			})
			return
		}

		session := sessions.Default(c)
		session.Set("user", databaseUsername)
		user := session.Get("user")
		session.Save()

		c.JSON(200, gin.H{
			"message": "hello " + databaseUsername,
			"user":    user,
		})
	})

	router.Run(":8080")
}
