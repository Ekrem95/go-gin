package main

import (
	"database/sql"
	_ "fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func common(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main",
	})
}

func router() {
	// Redis()
	MySQL()
	router := gin.Default()
	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("session", store))
	router.LoadHTMLGlob("../templates/*")
	router.Static("/src", "../src")

	// socketio
	server, err := socketio.NewServer(nil)
	if err != nil {
		log.Fatal(err)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Emit("some:event", "dataForClient")

		so.On("msg", func(msg string) {
			log.Println(msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	router.GET("/", common)
	router.GET("/signup", common)
	router.GET("/login", common)
	router.GET("/socket.io/", gin.WrapH(server))
	router.POST("/socket.io/", gin.WrapH(server))

	router.GET("/user", func(c *gin.Context) {
		session := sessions.Default(c)
		user := session.Get("user")
		// fmt.Println(user)
		c.JSON(http.StatusOK, gin.H{
			"user": user,
		})
	})
	router.GET("/cookie", func(c *gin.Context) {
		// fmt.Println(user)
		c.JSON(http.StatusOK, gin.H{
			"cookie": c.Request.Context(),
		})
	})

	router.POST("/signup", func(c *gin.Context) {
		username := c.PostForm("username")
		hashedPassword, error := bcrypt.GenerateFromPassword([]byte(c.PostForm("password")), bcrypt.DefaultCost)
		if error != nil {
			panic(error)
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
				"error": "Username already exists.",
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
	router.POST("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Set("user", nil)
		user := session.Get("user")
		session.Save()

		c.JSON(200, gin.H{
			"logged Out": true,
			"user":       user,
		})
	})

	router.Run(":8080")
}
