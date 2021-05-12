package ws

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/gorilla/websocket"
	"github.com/jeonjonghyeok/chat-mqtt/db"
	"github.com/jeonjonghyeok/chat-mqtt/vo"
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
	mutex  sync.Mutex
}

func newConn(wsConn *websocket.Conn, roomID, userID int) *conn {
	return &conn{
		wsConn: wsConn,
		roomID: roomID,
		userID: userID,
	}
}
func (c *conn) run() error {
	c.wg.Add(1)
	go c.readPump()
	c.wg.Wait()
	return nil
}

func (c *conn) pub(client mqtt.Client, msg vo.Message) {
	log.Println("Pub call")
	user, err := db.GetUser(c.userID)
	if err != nil {
		log.Println(err)
		return
	}

	msg.Sender = user
	msg.SenderID = c.userID
	msg.SentOn = time.Now()

	m, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	token := client.Publish("topic/test", 1, false, m)
	token.Wait()
	if token.Error() != nil {
		log.Println(token.Error())
	}
	time.Sleep(time.Second)
}

func (c *conn) sub(client mqtt.Client) {
	log.Println("subscribe call")
	topic := "topic/test"
	token := client.Subscribe(topic, 2, nil)
	token.Wait()
	fmt.Printf("Subscribed to topic: %s\n", topic)
}

func (c *conn) readPump() {
	defer c.wg.Done()
	client := NewBroker(broker, port, c.userID)
	c.sub(client)

	c.wsConn.SetReadLimit(maxMessageSize)
	c.wsConn.SetReadDeadline(time.Now().Add(readTimeout))
	c.wsConn.SetPongHandler(func(string) error {
		c.wsConn.SetReadDeadline(time.Now().Add(readTimeout))
		return nil
	})

	for {
		var msg vo.Message
		if err := c.wsConn.ReadJSON(&msg); err != nil {
			log.Println("err reading:", err)
			return
		}
		c.pub(client, msg)
	}
	client.Disconnect(250)
}
