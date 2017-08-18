package main

import (
	"database/sql"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
)

func common(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title": "Main",
	})
}

func getUser(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")
	// fmt.Println(user)
	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func signupPOST(c *gin.Context) {
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
}

func loginPOST(c *gin.Context) {

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
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("user", nil)
	user := session.Get("user")
	session.Save()

	c.JSON(200, gin.H{
		"logged Out": true,
		"user":       user,
	})
}

func addPost(c *gin.Context) {
	var post Post
	// post := Article{"a","b","c", [], 6}
	post.title = c.PostForm("title")
	post.description = c.PostForm("desc")
	post.src = c.PostForm("src")
	log.Println(post)
}

func getPosts(c *gin.Context) {
	var posts []Post
	var post Post
	var id, likes int
	var title, description string

	rows, err := db.Query("select id, title, description, likes from posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		error := rows.Scan(&id, &title, &description, &likes)
		if error != nil {
			log.Fatal(error)
		}
		post.id = id
		post.title = title
		post.description = description
		post.likes = likes

		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}
	log.Println(posts)

	c.JSON(200, gin.H{
		"posts": posts,
	})
}
