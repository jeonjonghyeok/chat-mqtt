package ws

import (
	"sync"
	"time"

	"github.com/gorilla/websocket"
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
	//sub        db.ChatroomSubscription
	chatroomID int
	senderID   int
}

func newConn(wsConn *websocket.Conn) *conn {
	return &conn{
		wsConn: wsConn,
		//chatroomID: chatroomID,
		//senderID:   senderID,
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
	return nil
}

func (c *conn) publish() {

}

func (c *conn) subscribe() {

}
