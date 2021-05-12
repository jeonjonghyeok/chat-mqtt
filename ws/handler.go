package ws

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     func(r *http.Request) bool { return true },
}
var Conn *websocket.Conn
var NewConn *conn

func ChatHandler(roomID, userID int) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		Conn, err = upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
			return
		}
		defer Conn.Close()

		NewConn = newConn(Conn, roomID, userID)
		if err := NewConn.run(); err != nil {
			log.Println(err)
		}
	})
}
