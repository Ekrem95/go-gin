package main

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/googollee/go-socket.io"
	"log"
)

func main() {
	// Redis()
	MySQL()
	router := gin.Default()
	store, _ := sessions.NewRedisStore(10, "tcp", "localhost:6379", "", []byte("secret"))
	router.Use(sessions.Sessions("session", store))
	router.LoadHTMLGlob("../templates/*")
	router.Static("/src", "../src")

	// socketio
	server, socketErr := socketio.NewServer(nil)
	if socketErr != nil {
		log.Fatal(socketErr)
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Join("chat")

		so.On("msg", func(msg string) {
			so.BroadcastTo("chat", "dist", msg)
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
	router.GET("/user", getUser)

	router.POST("/signup", signupPOST)
	router.POST("/login", loginPOST)
	router.POST("/logout", logout)

	// socketio
	router.GET("/socket.io/", gin.WrapH(server))
	router.POST("/socket.io/", gin.WrapH(server))

	router.Run(":8080")
}
