package router

import (
	"log"

	"github.com/ekrem95/go-gin/db"
	"github.com/googollee/go-socket.io"
)

func websocket() (*socketio.Server, error) {
	server, err := socketio.NewServer(nil)
	if err != nil {
		return nil, err
	}
	server.On("connection", func(so socketio.Socket) {
		log.Println("on connection")

		so.Join("chat")

		so.On("msg", func(msg *db.Message) {
			so.BroadcastTo("chat", "dist", msg)
			db.RedisSaveMsg(msg)
		})
		so.On("disconnection", func() {
			log.Println("on disconnect")
		})
	})
	server.On("error", func(so socketio.Socket, err error) {
		log.Println("error:", err)
	})

	return server, nil
}
