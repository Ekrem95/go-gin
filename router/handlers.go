package router

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ekrem95/go-gin/db"
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

	c.JSON(http.StatusOK, gin.H{
		"user": user,
	})
}

func signup(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) > 0 && len(password) > 0 {

		hashedPassword, error := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if error != nil {
			panic(error)
		}

		var user string

		smt := fmt.Sprintf("SELECT username FROM users WHERE username = '%s'", username)

		err := db.QueryRowScan(smt, &user)

		if err != nil {
			fmt.Println(err)
		}

		switch {
		// Username is available
		case err == sql.ErrNoRows:
			_, err = db.Exec("INSERT INTO users(username, password) VALUES(?, ?)", username, hashedPassword)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": "Unable to Sign up.",
				})
				return
			}

			session := sessions.Default(c)
			user := username
			session.Set("user", username)
			session.Save()

			c.JSON(http.StatusOK, gin.H{
				"success": true,
				"user":    user,
			})
			return
		case err != nil:
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "An error occured.",
			})
			return
		default:
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Username already exists.",
			})
		}
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both fields are required."})
	}
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) > 0 && len(password) > 0 {

		var databaseUsername string
		var databasePassword string

		// Search the database for the username provided
		// If it exists grab the password for validation
		smt := fmt.Sprintf("SELECT username, password FROM users WHERE username = '%s'", username)

		err := db.QueryRowScan(smt, &databaseUsername, &databasePassword)
		// If not then redirect to the login page
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"err": err,
			})
			return
		}

		// Validate the password
		err = bcrypt.CompareHashAndPassword([]byte(databasePassword), []byte(password))
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
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

		c.JSON(http.StatusOK, gin.H{
			"message": "hello " + databaseUsername,
			"user":    user,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Both fields are required."})
	}
}

func logout(c *gin.Context) {
	session := sessions.Default(c)
	session.Set("user", nil)
	user := session.Get("user")
	session.Save()

	c.JSON(http.StatusOK, gin.H{
		"logged Out": true,
		"user":       user,
	})
}

func addPost(c *gin.Context) {
	var post db.Post
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&post)
	if err != nil {
		panic(err)
	}
	defer c.Request.Body.Close()

	if len(post.Title) > 0 && len(post.Description) > 0 && len(post.Src) > 0 {
		res, err := db.Exec("INSERT INTO posts(title, description, src, posted_by) VALUES(?, ?, ?, ?)", post.Title, post.Description, post.Src, post.PostedBy)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to add.",
			})
			return
		}

		id, _ := res.LastInsertId()
		c.JSON(http.StatusOK, gin.H{
			"id": id,
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Fields for title, description and src are required.",
		})
	}

}

func editPost(c *gin.Context) {
	var post db.Post
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&post)
	if err != nil {
		panic(err)
	}
	defer c.Request.Body.Close()

	id := c.Param("id")

	_, err = db.Exec("update posts set title = (?), description = (?), src = (?) where id=?", post.Title, post.Description, post.Src, id)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to edit.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"done": true,
	})
}

func getPosts(c *gin.Context) {
	var posts []db.Post
	var post db.Post

	rows, err := db.Query("select id, title, src, description, likes from posts")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&post.ID, &post.Title, &post.Src, &post.Description, &post.Likes)
		if err != nil {
			log.Fatal(err)
		}

		posts = append(posts, post)
	}
	err = rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"posts": posts,
	})
}

func getPostByID(c *gin.Context) {
	id := c.Param("id")
	var post db.Post
	error := db.QueryRowScan("select id, title, src, description, likes from posts where id = "+id, &post.ID, &post.Title, &post.Src, &post.Description, &post.Likes)

	if error != nil {
		if error == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{
				"post": nil,
			})
			return
		}
		log.Fatal(error)
		c.JSON(http.StatusInternalServerError, gin.H{
			"post": nil,
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"post": post,
	})
}

func getPostsByUsername(c *gin.Context) {
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
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"p": posts,
	})
}

func postComment(c *gin.Context) {
	var comment db.Comment
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&comment)
	if err != nil {
		panic(err)
	}
	defer c.Request.Body.Close()

	comment.Time = time.Now().Unix()

	_, err = db.Exec("INSERT INTO comments(text, sender, post_id, time) VALUES(?, ?, ?, ?)", comment.Text, comment.Sender, comment.PostID, comment.Time)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to add comment.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"done": true,
	})
}

func getCommentsByID(c *gin.Context) {
	var comments []db.Comment
	var comment db.Comment
	id := c.Param("id")

	rows, error := db.Query("select text, sender, post_id, time from comments where post_id=?", id)
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
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
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
	sessionUser := sessions.Default(c).Get("user")

	if user != sessionUser {
		c.JSON(http.StatusOK, gin.H{
			"err": "Unable to delete post.",
		})
		return
	}

	_, err := db.Exec("delete from posts where id=? and posted_by=? limit 1", id, user)
	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Unable to delete post.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"deleted": true,
	})
}

func changePassword(c *gin.Context) {
	current := c.PostForm("current")
	newPassword := c.PostForm("newPassword")
	username := sessions.Default(c).Get("user")

	if username == nil || len(current) < 6 || len(newPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{
			"err": "Bad Request",
		})
		return
	}

	var password string

	smt := fmt.Sprintf("SELECT password FROM users WHERE username= '%s'", username.(string))

	err := db.QueryRowScan(smt, &password)

	if err != nil {
		log.Fatal(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Internal Server Error",
		})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(password), []byte(current))
	if err != nil {
		c.JSON(http.StatusOK, gin.H{
			"err": "Passwords do not match",
		})
		return
	}

	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if error != nil {
		panic(error)
	}

	_, err = db.Exec("update users set password=(?) where username=(?)", hashedPassword, username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to change password.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"done": true,
	})
}

func postLikes(c *gin.Context) {
	var like db.Like
	decoder := json.NewDecoder(c.Request.Body)
	error := decoder.Decode(&like)
	if error != nil {
		panic(error)
	}
	defer c.Request.Body.Close()

	var id int
	var username string

	smt := fmt.Sprintf("SELECT post_id, user FROM post_likes WHERE post_id= '%s' and user= '%s'", like.PostID, like.User)

	err := db.QueryRowScan(smt, &id, &username)

	switch {
	case err == sql.ErrNoRows:
		_, err = db.Exec("INSERT INTO post_likes (post_id, user) VALUES(?, ?)", like.PostID, like.User)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to like.",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"success": true,
		})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "An error occured.",
		})
		return
	default:
		_, err = db.Exec("delete from post_likes where post_id=? and user=? limit 1", like.PostID, like.User)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Unable to dislike.",
			})
			return
		}
	}

	// ****************************************************************

	// _,err = db.Exec("insert into post_likes (post_id, user)  select * from (select " + postID + ", '" + user + "') as tmp where not exists ( select post_id, user from post_likes where post_id = " + postID + "  and user = '" + user + "' ) limit 1")
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
	err := rows.Err()
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, gin.H{
		"users": users,
	})

}
