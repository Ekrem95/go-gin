package router

import (
	"log"
	"net/http"
	"os"

	"github.com/ekrem95/go-gin/db"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

const (
	uploadPath = "./app/uploads"
)

// Default ...
func Default() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()

	if _, err := os.Stat(uploadPath); os.IsNotExist(err) {
		os.Mkdir(uploadPath, 0700)
	}

	gopath := os.Getenv("GOPATH")
	public := gopath + "/src/github.com/ekrem95/go-gin/app"

	store, _ := sessions.NewRedisStore(10, "tcp", db.RedisAddress, "", []byte("secret"))
	r.Use(sessions.Sessions("session", store))
	r.LoadHTMLGlob(public + "/templates/*")
	r.StaticFS("/src", http.Dir(public+"/src"))
	r.StaticFile("/favicon.ico", public+"/templates/favicon.ico")

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
	ws, err := websocket()
	if err != nil {
		log.Fatal(err)
	}
	r.GET("/socket.io/", gin.WrapH(ws))
	r.POST("/socket.io/", gin.WrapH(ws))

	return r
}
