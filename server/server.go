package server

import (
	"net/http"

	"github.com/jeonjonghyeok/chat-mqtt/api"
)

type Config struct {
	Address string
	Url     string
}

func ListenAndServe(c Config) error {
	return http.ListenAndServe(c.Address, api.API())
}
