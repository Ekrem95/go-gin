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

func setSession(c *gin.Context, k, v string, MaxAge int) {
	session := sessions.Default(c)

	if MaxAge != 0 {
		session.Options(sessions.Options{MaxAge: 604800})
	}

	session.Set(k, v)
	session.Save()
}

func common(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", gin.H{"title": "Main"})
}

func getUser(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("user")

	c.JSON(http.StatusOK, gin.H{"user": user})
}

func signup(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) < 3 || len(password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	var user db.User
	if exists := user.Exists(username); exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Username already exists"})
		return
	}

	if _, err = db.Exec(
		"INSERT INTO users(username, password) VALUES(?, ?)",
		username, hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
		return
	}

	setSession(c, "user", username, 0)
	c.JSON(http.StatusOK, gin.H{"success": true, "user": username})
}

func login(c *gin.Context) {
	username := c.PostForm("username")
	password := c.PostForm("password")

	if len(username) < 3 || len(password) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}

	var databaseUsername string
	var databasePassword string

	// Search the database for the username provided
	// If it exists grab the password for validation
	smt := fmt.Sprintf("SELECT username, password FROM users WHERE username = '%s'", username)

	// If not then redirect to the login page
	if err := db.QueryRowScan(smt, &databaseUsername, &databasePassword); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": err})
		return
	}

	// Validate the password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(databasePassword), []byte(password)); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"err": "Passwords do not match"})
		return
	}

	setSession(c, "user", databaseUsername, 604800)
	c.JSON(http.StatusOK, gin.H{"user": databaseUsername})
}

func logout(c *gin.Context) {
	setSession(c, "user", "", 0)
	c.JSON(http.StatusOK, gin.H{"success": true})
}

func addPost(c *gin.Context) {
	var post db.Post
	decoder := json.NewDecoder(c.Request.Body)
	err := decoder.Decode(&post)
	if err != nil {
		panic(err)
	}
	defer c.Request.Body.Close()

	if len(post.Title) < 1 || len(post.Description) < 5 || len(post.Src) < 5 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Bad Request"})
		return
	}
	res, err := db.Exec(
		"INSERT INTO posts(title, description, src, posted_by) VALUES(?, ?, ?, ?)", post.Title, post.Description, post.Src, post.PostedBy)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to add."})
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func editPost(c *gin.Context) {
	var post db.Post
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&post); err != nil {
		panic(err)
	}
	defer c.Request.Body.Close()

	id := c.Param("id")

	if _, err := db.Exec("update posts set title = (?), description = (?), src = (?) where id=?", post.Title, post.Description, post.Src, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to edit."})
		return
	}

	c.JSON(http.StatusOK, gin.H{"done": true})
}

func getPosts(c *gin.Context) {
	var posts []db.Post
	var post db.Post

	rows, err := db.Query("select id, title, src, description, likes from posts")
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&post.ID, &post.Title, &post.Src, &post.Description, &post.Likes); err != nil {
			panic(err)
		}

		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func getPostByID(c *gin.Context) {
	id := c.Param("id")
	var post db.Post
	err := db.QueryRowScan("select id, title, src, description, likes from posts where id = "+id, &post.ID, &post.Title, &post.Src, &post.Description, &post.Likes)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"post": nil})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

func getPostsByUsername(c *gin.Context) {
	name := c.Param("name")
	var id, title string
	posts := map[string]string{}

	rows, err := db.Query("select id, title from posts where posted_by=?", name)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&id, &title); err != nil {
			panic(err)
		}

		posts[id] = title
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func postComment(c *gin.Context) {
	var comment db.Comment
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&comment); err != nil {
		panic(err)
	}
	defer c.Request.Body.Close()

	comment.Time = time.Now().Unix()

	if _, err := db.Exec("INSERT INTO comments(text, sender, post_id, time) VALUES(?, ?, ?, ?)", comment.Text, comment.Sender, comment.PostID, comment.Time); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Unable to add comment.",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{"done": true})
}

func getCommentsByID(c *gin.Context) {
	var comments []db.Comment
	var comment db.Comment
	id := c.Param("id")

	rows, err := db.Query("select text, sender, post_id, time from comments where post_id=?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&comment.Text, &comment.Sender, &comment.PostID, &comment.Time); err != nil {
			panic(err)
		}

		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
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

func deletePostByID(c *gin.Context) {
	id := c.PostForm("id")
	user := c.PostForm("user")
	sessionUser := sessions.Default(c).Get("user")

	if user != sessionUser {
		c.JSON(http.StatusOK, gin.H{"err": "Unable to delete post."})
		return
	}

	if _, err := db.Exec("delete from posts where id=? and posted_by=? limit 1", id, user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Unable to delete post.",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

func changePassword(c *gin.Context) {
	current := c.PostForm("current")
	newPassword := c.PostForm("newPassword")
	username := sessions.Default(c).Get("user")

	if username == nil || len(current) < 6 || len(newPassword) < 6 {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Bad Request"})
		return
	}

	var password string

	smt := fmt.Sprintf("SELECT password FROM users WHERE username= '%s'", username.(string))

	if err := db.QueryRowScan(smt, &password); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"err": "Internal Server Error",
		})
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(password), []byte(current)); err != nil {
		c.JSON(http.StatusOK, gin.H{"err": "Passwords do not match"})
		return
	}

	hashedPassword, error := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if error != nil {
		panic(error)
	}

	if _, err := db.Exec("update users set password=(?) where username=(?)", hashedPassword, username); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"done": true})
}

func postLikes(c *gin.Context) {
	var like db.Like
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&like); err != nil {
		panic(err)
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
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal Error"})
		return
	default:
		if _, err = db.Exec("delete from post_likes where post_id=? and user=? limit 1", like.PostID, like.User); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Internal Error",
			})
			return
		}
	}
}

func getLikes(c *gin.Context) {
	id := c.Param("id")

	var user string
	var users []string

	rows, err := db.Query("select user from post_likes where post_id=?", id)
	if err != nil {
		panic(err)
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&user); err != nil {
			panic(err)
		}

		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		panic(err)
	}

	c.JSON(http.StatusOK, gin.H{"users": users})

}
