package api

import (
	"net/http"

	"github.com/jeonjonghyeok/chat-mqtt/ws"
)

func createRoom(w http.ResponseWriter, r *http.Request) {

}
func getRooms(w http.ResponseWriter, r *http.Request) {

}

func connectToRoom(w http.ResponseWriter, r *http.Request) {
	ws.ChatHandler().ServeHTTP(w, r)
}
