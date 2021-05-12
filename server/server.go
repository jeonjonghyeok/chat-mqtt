package server

import (
	"log"
	"net/http"

	"github.com/jeonjonghyeok/chat-mqtt/api"
	"github.com/jeonjonghyeok/chat-mqtt/db"
)

type Config struct {
	Address string
	Url     string
}

func ListenAndServe(c Config) error {
	if err := db.Connect(c.Url); err != nil {
		log.Fatal(err)
	}
	return http.ListenAndServe(c.Address, api.API())
}
