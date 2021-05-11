package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/jeonjonghyeok/chat-mqtt/vo"
)

var messagePubHandler mqtt.MessageHandler = func(client mqtt.Client, msg mqtt.Message) {
	log.Println("message pub handler call")
	Conn.SetWriteDeadline(time.Now().Add(writeTimeout))
	var m vo.Message
	json.Unmarshal(msg.Payload(), &m)

	log.Println("Text:", m)
	if err := Conn.WriteJSON(m); err != nil {
		log.Println(err)
		return
	}

}

var connectHandler mqtt.OnConnectHandler = func(client mqtt.Client) {
	fmt.Println("Connected")
}

var connectLostHandler mqtt.ConnectionLostHandler = func(client mqtt.Client, err error) {
	fmt.Printf("Connect lost: %v\n", err)
}

func NewBroker(broker string, port, userID int) (Client mqtt.Client) {
	opts := mqtt.NewClientOptions()
	opts.AddBroker(fmt.Sprintf("tcp://%s:%d", broker, port))
	opts.SetClientID(fmt.Sprint(userID))
	opts.SetUsername("emqx")
	opts.SetPassword("public")
	opts.SetDefaultPublishHandler(messagePubHandler)
	opts.OnConnect = connectHandler
	opts.OnConnectionLost = connectLostHandler
	Client = mqtt.NewClient(opts)
	if token := Client.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}
	return
}
