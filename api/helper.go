package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jeonjonghyeok/chat-mqtt/db"
	"github.com/jeonjonghyeok/chat-mqtt/token"
	"github.com/jeonjonghyeok/chat-mqtt/vo"
)

//t-chat
//const apikey = "9ZEJVUK1FMVKGJEZ9YHJILAZ7IWKCHTM"
const apikey = "GSRHNFLVOCK0YPVX6NGFWWZPMI5G0EOK"

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	must(json.NewEncoder(w).Encode(v))
}

func parseJSON(r io.Reader, v interface{}) {
	if err := json.NewDecoder(r).Decode(v); err != nil {
		log.Println("parsing json body:", err)
		panic(malformedInputError)
	}
}

func must(err error) {
	if err == db.ErrUnauthorized {
		panic(unauthorizedError)
	} else if err != nil {
		log.Println("internal error:", err)
		panic(internalError)
	}
}

func parseIntParam(r *http.Request, key string) int {
	vars := mux.Vars(r)
	if v, ok := vars[key]; ok {
		i, err := strconv.Atoi(v)
		must(err)
		return i
	}
	panic(malformedInputError)

}

func userID(r *http.Request) int {
	t := r.URL.Query().Get("token")
	id, err := token.Parse(t)
	if err != nil {
		log.Println(err)
		panic(unauthorizedError)
	}
	return id
}
func userName(r *http.Request) (name string) {
	t := r.URL.Query().Get("token")
	id, err := token.Parse(t)
	if err != nil {
		log.Println(err)
		panic(unauthorizedError)
	}
	name, err = db.GetUser(id)
	must(err)

	return name
}
func parseTToken(r *http.Request) string {
	t := r.URL.Query().Get("token")
	//log.Println("token = ", t)
	ttoken, err := token.Tparse(t)
	if err != nil {
		log.Println(err)
		panic(unauthorizedError)
	}
	return ttoken
}

func postUser(w http.ResponseWriter, v interface{}, url string) string {
	var token vo.ResponseToken
	pbytes, err := json.Marshal(v)
	must(err)
	buff := bytes.NewBuffer(pbytes)
	req, err := http.NewRequest("POST", url, buff)
	must(err)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("apikey", apikey)

	client := &http.Client{}
	resp, err := client.Do(req)
	must(err)

	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)

	json.Unmarshal([]byte(respBody), &token)
	if token.Status == "error" {
		must(errors.New("sign up fail"))
	}

	return token.Data.Token

}
