package ws

import (
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
)

const (
	writeTimeout   = 10 * time.Second
	readTimeout    = 60 * time.Second
	pingPeriod     = 10 * time.Second
	maxMessageSize = 512
)
const (
	broker = "broker.emqx.io"
	port   = 1883
)

type conn struct {
	wsConn *websocket.Conn
	wg     sync.WaitGroup
	roomID int
	userID int
}

func newConn(wsConn *websocket.Conn, roomID, userID int) *conn {
	return &conn{
		wsConn: wsConn,
		roomID: roomID,
		userID: userID,
	}
}
func (c *conn) run() error {
	/*sub, err := db.NewChatroomSubscription(c.chatroomID)
	if err != nil {
		return err
	}
	c.sub = sub

	c.wg.Add(2)

	c.wg.Wait()
	*/
	client := NewBroker(broker, port, c.userID)
	c.sub(client)
	c.pub(client)

	client.Disconnect(250)
	return nil
}

func (c *conn) pub(client mqtt.Client) {
	log.Println("publish call")
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func (c *conn) sub(client mqtt.Client) {
	log.Println("subscribe call")
	topic := "topic/test"
	token := client.Subscribe(topic, 2, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}
