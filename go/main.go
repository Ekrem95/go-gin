package main

import (
	"log"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
)

func main() {
	MySQL()
	r := gin.Default()
	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	r.Use(sessions.Sessions("session", store))
	r.LoadHTMLGlob("../templates/*")
	r.Static("/src", "../src")
	r.StaticFile("/favicon.ico", "../templates/favicon.ico")

	// socketio
	server, socketErr := socketio.NewServer(nil)
	if socketErr != nil {
		log.Fatal(socketErr)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Join("chat")

		so.On("msg", func(msg *Message) {
			so.BroadcastTo("chat", "dist", msg)
			RedisSaveMsg(msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	r.GET("/", common)
	r.GET("/signup", common)
	r.GET("/login", common)
	r.GET("/add", common)
	r.GET("/upload", common)
	r.GET("/user", getUser)
	r.GET("/messages", RedisGetMsgs)
	r.GET("/api/posts", getPosts)
	r.GET("/api/postbyid/:id", getPostByID)
	r.GET("/api/commentsbyid/:id", getCommentsByID)
	r.GET("/p/*all", common)
	r.GET("/myposts", common)
	r.GET("/api/getpostbyusername/:name", getPostByUsername)
	r.GET("/edit/:id", common)

	r.POST("/signup", signupPOST)
	r.POST("/login", loginPOST)
	r.POST("/logout", logout)
	r.POST("/add", addPost)
	r.POST("/comment", postComment)
	r.POST("/upload", uploadFile)
	r.POST("/edit/:id", editPost)

	// socketio
	r.GET("/socket.io/", gin.WrapH(server))
	r.POST("/socket.io/", gin.WrapH(server))

	r.Run(":8080")
}
