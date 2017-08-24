package main

import (
	"database/sql"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
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
	post.Title = c.PostForm("title")
	post.Description = c.PostForm("desc")
	post.Src = c.PostForm("src")
	post.PostedBy = c.PostForm("postedBy")

	_, err = db.Exec("INSERT INTO posts(title, description, src, posted_by) VALUES(?, ?, ?, ?)", post.Title, post.Description, post.Src, post.PostedBy)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{
			"error": "Unable to add.",
		})
		return
	}
	c.JSON(200, gin.H{
		"done": true,
	})
}

func getPosts(c *gin.Context) {
	var posts []Post
	var post Post

	rows, error := db.Query("select id, title, src, description, likes from posts")
	if error != nil {
		log.Fatal(error)
	}
	defer rows.Close()

	for rows.Next() {
		error := rows.Scan(&post.ID, &post.Title, &post.Src, &post.Description, &post.Likes)
		if error != nil {
			log.Fatal(error)
		}

		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"posts": posts,
	})
}

func getPostByID(c *gin.Context) {
	id := c.Param("id")
	var post Post
	error := db.QueryRow("select id, title, src, description, likes from posts where id =?", id).Scan(&post.ID, &post.Title, &post.Src, &post.Description, &post.Likes)
	if error != nil {
		log.Fatal(error)
		c.JSON(200, gin.H{
			"post": nil,
		})
	}

	c.JSON(200, gin.H{
		"post": post,
	})
}

func getPostByUsername(c *gin.Context) {
	name := c.Param("name")
	var id, title string
	posts := map[string]string{}

	rows, error := db.Query("select id, title from posts where posted_by=?", name)
	if error != nil {
		log.Fatal(error)
	}
	defer rows.Close()

	for rows.Next() {
		error := rows.Scan(&id, &title)
		if error != nil {
			log.Fatal(error)
		}
		posts[id] = title
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"p": posts,
	})
}

func postComment(c *gin.Context) {
	var comment Comment

	comment.Sender = c.PostForm("sender")
	comment.PostID = c.PostForm("postId")
	comment.Text = c.PostForm("text")
	comment.Time = time.Now().Unix()

	_, err = db.Exec("INSERT INTO comments(text, sender, postId, time) VALUES(?, ?, ?, ?)", comment.Text, comment.Sender, comment.PostID, comment.Time)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{
			"error": "Unable to add comment.",
		})
		return
	}
	c.JSON(200, gin.H{
		"done": true,
	})
}

func getCommentsByID(c *gin.Context) {
	var comments []Comment
	var comment Comment
	id := c.Param("id")

	rows, error := db.Query("select text, sender, postId, time from comments where postId=?", id)
	if error != nil {
		log.Fatal(error)
	}
	defer rows.Close()

	for rows.Next() {
		error := rows.Scan(&comment.Text, &comment.Sender, &comment.PostID, &comment.Time)
		if error != nil {
			log.Fatal(error)
		}

		comments = append(comments, comment)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"comments": comments,
	})
}

func uploadFile(c *gin.Context) {
	file, handler, err := c.Request.FormFile("photo")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()
	f, err := os.OpenFile("./photos/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}
