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
	username, password := c.PostForm("username"), c.PostForm("password")

	if len(username) < 3 || len(password) < 6 {
		c.JSON(http.StatusBadRequest, errors.BadRequest)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	var user db.User
	if exists := user.Exists(username); exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username already exists"})
		return
	}

	if _, err = db.Exec(
		"INSERT INTO users(username, password) VALUES(?, ?)",
		username, hashedPassword); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	setSession(c, "user", username, 0)
	c.JSON(http.StatusOK, gin.H{"success": true, "user": username})
}

func login(c *gin.Context) {
	username, password := c.PostForm("username"), c.PostForm("password")

	if len(username) < 3 || len(password) < 6 {
		c.JSON(http.StatusBadRequest, errors.BadRequest)
		return
	}

	var databaseUsername, databasePassword string

	// Search the database for the username provided
	// If it exists grab the password for validation
	smt := fmt.Sprintf("SELECT username, password FROM users WHERE username = '%s'", username)

	// If not then redirect to the login page
	if err := db.QueryRowScan(smt, &databaseUsername, &databasePassword); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": err})
		return
	}

	// Validate the password
	if err := bcrypt.CompareHashAndPassword(
		[]byte(databasePassword), []byte(password)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Passwords does not match"})
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
	title, src, description, by := c.PostForm("title"), c.PostForm("src"), c.PostForm("description"), c.PostForm("posted_by")

	if len(title) < 1 || len(description) < 5 || len(src) < 5 {
		c.JSON(http.StatusBadRequest, errors.BadRequest)
		return
	}
	res, err := db.Exec(
		"INSERT INTO posts(title, description, src, posted_by) VALUES(?, ?, ?, ?)", title, description, src, by)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	id, _ := res.LastInsertId()
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func editPost(c *gin.Context) {
	var post db.Post
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&post); err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest)
		return
	}
	defer c.Request.Body.Close()

	id := c.Param("id")

	if _, err := db.Exec("update posts set title = (?), description = (?), src = (?) where id=?", post.Title, post.Description, post.Src, id); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"done": true})
}

func getPosts(c *gin.Context) {
	var posts []db.Post
	var post db.Post

	rows, err := db.Query("select id, title, src, description, likes from posts")
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&post.ID, &post.Title, &post.Src, &post.Description, &post.Likes); err != nil {
			c.JSON(http.StatusInternalServerError, errors.Internal)
			return
		}
		posts = append(posts, post)
	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func getPostByID(c *gin.Context) {
	var post db.Post
	id := c.Param("id")
	err := db.QueryRowScan("select id, title, src, description, likes from posts where id = "+id, &post.ID, &post.Title, &post.Src, &post.Description, &post.Likes)

	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusBadRequest, gin.H{"post": nil})
			return
		}
		log.Println(err)
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"post": post})
}

func getPostsByUsername(c *gin.Context) {
	var id, title string
	name := c.Param("name")
	posts := map[string]string{}

	rows, err := db.Query("select id, title from posts where posted_by=?", name)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&id, &title); err != nil {
			c.JSON(http.StatusInternalServerError, errors.Internal)
			return
		}
		posts[id] = title
	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"posts": posts})
}

func postComment(c *gin.Context) {
	text, postID, sender := c.PostForm("text"), c.PostForm("post_id"), c.PostForm("sender")
	time := time.Now().Unix()

	if _, err := db.Exec("INSERT INTO comments(text, sender, post_id, time) VALUES(?, ?, ?, ?)", text, sender, postID, time); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
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
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&comment.Text, &comment.Sender, &comment.PostID, &comment.Time); err != nil {
			c.JSON(http.StatusInternalServerError, errors.Internal)
			return
		}
		comments = append(comments, comment)
	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"comments": comments})
}

func uploadFile(c *gin.Context) {
	file, handler, err := c.Request.FormFile("photo")
	if err != nil {
		c.JSON(http.StatusBadRequest, errors.BadRequest)
		return
	}
	defer file.Close()
	f, err := os.OpenFile(uploadPath+"/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}
	defer f.Close()
	io.Copy(f, file)
}

func deletePostByID(c *gin.Context) {
	postID, user := c.PostForm("id"), c.PostForm("user")
	sessionUser := sessions.Default(c).Get("user")

	if user != sessionUser {
		c.JSON(http.StatusUnauthorized, errors.Unauthorized)
		return
	}

	if _, err := db.Exec("delete from posts where id=? and posted_by=? limit 1", postID, user); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"deleted": true})
}

func changePassword(c *gin.Context) {
	current, newPassword := c.PostForm("current"), c.PostForm("newPassword")
	username := sessions.Default(c).Get("user")

	if username == nil || len(current) < 6 || len(newPassword) < 6 {
		c.JSON(http.StatusBadRequest, errors.BadRequest)
		return
	}

	var password string

	smt := fmt.Sprintf("SELECT password FROM users WHERE username= '%s'", username.(string))

	if err := db.QueryRowScan(smt, &password); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(password), []byte(current)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"err": "Passwords does not match"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	if _, err := db.Exec("update users set password=(?) where username=(?)", hashedPassword, username); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"done": true})
}

func postLikes(c *gin.Context) {
	var like db.Like
	decoder := json.NewDecoder(c.Request.Body)
	if err := decoder.Decode(&like); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
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
			c.JSON(http.StatusInternalServerError, errors.Internal)
			return
		}

		c.JSON(http.StatusOK, gin.H{"success": true})
		return
	case err != nil:
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	default:
		if _, err = db.Exec("delete from post_likes where post_id=? and user=? limit 1", like.PostID, like.User); err != nil {
			c.JSON(http.StatusInternalServerError, errors.Internal)
			return
		}
	}
}

func getLikes(c *gin.Context) {
	var user string
	var users []string

	id := c.Param("id")

	rows, err := db.Query("select user from post_likes where post_id=?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&user); err != nil {
			c.JSON(http.StatusInternalServerError, errors.Internal)
			return
		}
		users = append(users, user)
	}
	if err = rows.Err(); err != nil {
		c.JSON(http.StatusInternalServerError, errors.Internal)
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": users})
}
