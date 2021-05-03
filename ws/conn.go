package ws

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/jeonjonghyeok/chat-mqtt/mqtt"
)

const (
	writeTimeout   = 10 * time.Second
	readTimeout    = 60 * time.Second
	pingPeriod     = 10 * time.Second
	maxMessageSize = 512
)

type conn struct {
	wsConn *websocket.Conn
	wg     sync.WaitGroup
}

func newConn(wsConn *websocket.Conn) *conn {
	return &conn{
		wsConn: wsConn,
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
	c.sub()
	c.pub()
	return nil
}

func (c *conn) pub() {
	log.Println("publish call")
	num := 10
	for i := 0; i < num; i++ {
		text := fmt.Sprintf("Message %d", i)
		token := mqtt.Client.Publish("topic/test", 0, false, text)
		token.Wait()
		time.Sleep(time.Second)
	}
}

func (c *conn) sub() {
	log.Println("subscribe call")
	topic := "topic/test"
	token := mqtt.Client.Subscribe(topic, 2, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s", topic)
}
