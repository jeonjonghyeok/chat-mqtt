package server

import (
	"net/http"

	"github.com/jeonjonghyeok/chat-mqtt/mqtt"
	"github.com/jeonjonghyeok/chat-mqtt/ws"
)

const (
	broker = "broker.emqx.io"
	port   = 1883
)

type Config struct {
	Address string
}

func ListenAndServe(c Config) error {
	mqtt.NewBroker(broker, port)
	return http.ListenAndServe(c.Address, ws.ChatHandler())
}
