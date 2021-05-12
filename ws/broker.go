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
	NewConn.mutex.Lock()
	defer NewConn.mutex.Unlock()
	log.Println("message pub handler call")
	NewConn.wsConn.SetWriteDeadline(time.Now().Add(writeTimeout))
	var m vo.Message
	if err := json.Unmarshal(msg.Payload(), &m); err != nil {
		log.Println(err)
		return
	}

	log.Println("ID:", m.ID, " Sender:", m.Sender, " SenderID:", m.SenderID, " Msg:", m.Text)
	if err := NewConn.wsConn.WriteJSON(m); err != nil {
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
