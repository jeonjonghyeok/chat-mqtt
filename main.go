package main

import (
	"fmt"
	"log"

	"github.com/jeonjonghyeok/chat-mqtt/server"
)

const (
	DB_USER     = "jjh"
	DB_PASSWORD = "zxcvasdf123"
	DB_NAME     = "tt"
)

func main() {
	if err := server.ListenAndServe(server.Config{
		Address: ":5000",
		Url: fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable",
			DB_USER, DB_PASSWORD, DB_NAME),
	}); err != nil {
		log.Fatal(err)
	}

}
