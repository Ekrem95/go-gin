package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"log"
)

func main() {
	MySQL()
	router := gin.Default()
	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("session", store))
	router.LoadHTMLGlob("../templates/*")
	router.Static("/src", "../src")
	router.StaticFile("/favicon.ico", "../templates/favicon.ico")

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

	router.GET("/", common)
	router.GET("/signup", common)
	router.GET("/login", common)
	router.GET("/add", common)
	router.GET("/user", getUser)
	router.GET("/messages", RedisGetMsgs)
	router.GET("/api/posts", getPosts)
	router.GET("/api/postbyid/:id", getPostByID)
	router.GET("/p/*all", common)

	router.POST("/signup", signupPOST)
	router.POST("/login", loginPOST)
	router.POST("/logout", logout)
	router.POST("/add", addPost)
	router.POST("/comment", postComment)

	// socketio
	router.GET("/socket.io/", gin.WrapH(server))
	router.POST("/socket.io/", gin.WrapH(server))

	router.Run(":8080")
}
