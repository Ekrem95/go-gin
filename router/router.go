package router

import (
	"log"
	"net/http"

	"github.com/ekrem95/go-gin/db"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	// UploadPath ...
	UploadPath = "./app/uploads"
)

// Default ...
func Default() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("session", store))
	r.LoadHTMLGlob("./app/templates/*")
	r.StaticFS("/src", http.Dir("./app/src"))
	r.StaticFile("/favicon.ico", "./app/templates/favicon.ico")

	// socketio
	server, err := websocket()
	if err != nil {
		log.Fatal(err)
	}

	r.GET("/", common)
	r.GET("/signup", common)
	r.GET("/login", common)
	r.GET("/add", common)
	r.GET("/upload", common)
	r.GET("/user", getUser)
	r.GET("/messages", db.RedisGetMsgs)
	r.GET("/api/posts", getPosts)
	r.GET("/api/postbyid/:id", getPostByID)
	r.GET("/api/commentsbyid/:id", getCommentsByID)
	r.GET("/p/*all", common)
	r.GET("/myposts", common)
	r.GET("/api/getpostbyusername/:name", getPostsByUsername)
	r.GET("/edit/:id", common)
	r.GET("/changepassword", common)
	r.GET("/get_likes/:id", getLikes)

	r.POST("/signup", signup)
	r.POST("/login", login)
	r.POST("/logout", logout)
	r.POST("/add", addPost)
	r.POST("/comment", postComment)
	r.POST("/upload", uploadFile)
	r.POST("/edit/:id", editPost)
	r.POST("/delete/:id", deletePostByID)
	r.POST("/changepassword", changePassword)
	r.POST("/post_likes", postLikes)

	// socketio
	r.GET("/socket.io/", gin.WrapH(server))
	r.POST("/socket.io/", gin.WrapH(server))

	return r
}
