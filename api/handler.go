package api

import (
	"log"
	"net/http"

	"github.com/jeonjonghyeok/chat-mqtt/db"
	"github.com/jeonjonghyeok/chat-mqtt/token"
	"github.com/jeonjonghyeok/chat-mqtt/ws"
)

func test(w http.ResponseWriter, r *http.Request) {
	log.Println("test")
}

type usernamePassword struct {
	Username string `json:"username"`
	Password string `json:"password"`
}
type User struct {
	UserID string `json:"userid"`
	Passwd string `json:"password"`
}

func signup(w http.ResponseWriter, r *http.Request) {
	url := "https://api.trustalk.co.kr/quorum/user/register"

	var req usernamePassword
	var tuser User

	//t-talk token
	parseJSON(r.Body, &req)
	tuser.UserID = req.Username
	tuser.Passwd = req.Password
	ttoken := postUser(w, tuser, url)

	//t-chat signup
	log.Println("username=", req.Username, "password=", req.Password)
	id, err := db.CreateUser(req.Username, req.Password)
	must(err)

	//t-chat token create
	t, err := token.New(id, ttoken)
	must(err)

	writeJSON(w, struct {
		Token string `json:"token"`
	}{t})

}

func signin(w http.ResponseWriter, r *http.Request) {
	url := "https://api.trustalk.co.kr/quorum/user/login"

	var puser User

	var req usernamePassword
	parseJSON(r.Body, &req)
	puser.UserID = req.Username
	puser.Passwd = req.Password
	ttoken := postUser(w, puser, url)

	id, err := db.FindUser(req.Username, req.Password)
	must(err)

	t, err := token.New(id, ttoken)
	must(err)
	writeJSON(w, struct {
		Token string `json:"token"`
	}{t})

}
func createRoom(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Name string `json:"name"`
	}
	parseJSON(r.Body, &req)

	id, err := db.CreateRoom(req.Name)
	must(err)
	writeJSON(w, struct {
		ID int `json:"id"`
	}{id})
}

func getRooms(w http.ResponseWriter, r *http.Request) {
	uid := userID(r)
	rooms, err := db.GetRooms(uid)
	must(err)
	writeJSON(w, rooms)
}

func connectToRoom(w http.ResponseWriter, r *http.Request) {
	uid := userID(r)
	roomID := parseIntParam(r, "id")
	exists, err := db.RoomExists(roomID)
	must(err)

	if !exists {
		panic(notFoundError)
	}
	must(db.ConnectToRoom(uid, roomID))

	ws.ChatHandler(roomID, uid).ServeHTTP(w, r)
}
