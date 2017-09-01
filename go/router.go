package main

import (
	"database/sql"
	"encoding/json"
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
	session.Options(sessions.Options{MaxAge: 604800})
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
	decoder := json.NewDecoder(c.Request.Body)
	error := decoder.Decode(&post)
	if error != nil {
		panic(error)
	}
	defer c.Request.Body.Close()

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

func editPost(c *gin.Context) {
	id := c.Param("id")
	title := c.PostForm("title")
	description := c.PostForm("description")
	src := c.PostForm("src")

	_, err = db.Exec("update posts set title = (?), description = (?), src = (?) where id=?", title, description, src, id)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{
			"error": "Unable to edit.",
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
	file, handler, errFile := c.Request.FormFile("photo")
	if errFile != nil {
		fmt.Println(errFile)
		return
	}
	defer file.Close()
	f, errFile := os.OpenFile("./photos/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if errFile != nil {
		fmt.Println(errFile)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func deletePostByID(c *gin.Context) {
	id := c.PostForm("id")
	user := c.PostForm("user")
	log.Println("---------------------------------")
	sessionUser := sessions.Default(c).Get("user")

	if user != sessionUser {
		c.JSON(200, gin.H{
			"err": "Unable to delete post.",
		})
		return
	}

	_, err = db.Exec("delete from posts where id=? and posted_by=? limit 1", id, user)
	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{
			"err": "Unable to delete post.",
		})
		return
	}

	c.JSON(200, gin.H{
		"deleted": true,
	})
}
func changePassword(c *gin.Context) {
	current := c.PostForm("current")
	new := c.PostForm("newPassword")
	username := sessions.Default(c).Get("user")

	var password string

	err = db.QueryRow("SELECT password FROM users WHERE username=?", username).Scan(&password)

	if err != nil {
		log.Fatal(err)
		c.JSON(500, gin.H{
			"err": "Internal Server Error",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(current))
	if err != nil {
		c.JSON(200, gin.H{
			"err": "Passwords do not match",
		})
		return
	}

	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(new), bcrypt.DefaultCost)
	if error != nil {
		panic(error)
	}

	_, err = db.Exec("update users set password=(?) where username=(?)", hashedPassword, username)
	if err != nil {
		c.JSON(500, gin.H{
			"error": "Unable to change password.",
		})
		return
	}

	c.JSON(200, gin.H{
		"done": true,
	})
}

func postLikes(c *gin.Context) {
	postID := c.PostForm("postID")
	user := c.PostForm("user")

	var id int
	var username string

	err = db.QueryRow("SELECT post_id, user FROM post_likes WHERE post_id=? and user=?", postID, user).Scan(&id, &username)

	switch {
	case err == sql.ErrNoRows:
		_, err = db.Exec("INSERT INTO post_likes (post_id, user) VALUES(?, ?)", postID, user)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Unable to add.",
			})
			return
		}

		c.JSON(200, gin.H{
			"success": true,
		})
		return
	case err != nil:
		c.JSON(500, gin.H{
			"error": "An error occured.",
		})
		return
	default:
		_, err = db.Exec("delete from post_likes where post_id=? and user=? limit 1", postID, user)
		if err != nil {
			c.JSON(500, gin.H{
				"error": "Unable to delete.",
			})
			return
		}
	}

	// _, err = db.Exec("insert into post_likes (post_id, user)  select * from (select " + postID + ", '" + user + "') as tmp where not exists ( select post_id, user from post_likes where post_id = " + postID + "  and user = '" + user + "' ) limit 1")
	// if err != nil {
	// 	log.Fatal(err)
	// 	c.JSON(500, gin.H{
	// 		"error": "Unable to add.",
	// 	})
	// 	return
	// }
	// c.JSON(200, gin.H{
	// 	"done": true,
	// })
}

func getLikes(c *gin.Context) {
	id := c.Param("id")

	var user string
	var users []string

	rows, error := db.Query("select user from post_likes where post_id=?", id)
	if error != nil {
		log.Fatal(error)
	}
	defer rows.Close()

	for rows.Next() {
		error := rows.Scan(&user)
		if error != nil {
			log.Fatal(error)
		}

		users = append(users, user)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(200, gin.H{
		"users": users,
	})

}
