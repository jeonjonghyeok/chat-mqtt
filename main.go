package main

import (
	"log"

	"github.com/jeonjonghyeok/chat-mqtt/server"
)

func main() {
	if err := server.ListenAndServe(server.Config{
		Address: ":5000",
	}); err != nil {
		log.Fatal(err)
	}

}
