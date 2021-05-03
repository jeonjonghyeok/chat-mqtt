package main

import (
	"net/http"

	"github.com/jeonjonghyeok/chat-mqtt/api"
)

func main() {
	http.ListenAndServe(":5000", api.API())
}
