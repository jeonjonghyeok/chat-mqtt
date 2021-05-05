package db

import (
	"errors"
	"log"

	"github.com/jeonjonghyeok/chat-mqtt/vo"
)

// CreateRoom 방 만들기
func CreateRoom(name string) (id int, err error) {
	err = db.QueryRow(`INSERT INTO chatrooms (name)
		VALUES ($1)
		RETURNING id`, name).Scan(&id)
	return
}

// RoomExists 방 존재 여부
func RoomExists(id int) (exists bool, err error) {
	err = db.QueryRow(`SELECT EXISTS(
			SELECT 1 FROM chatrooms WHERE id = $1)`, id).
		Scan(&exists)
	return
}

// GetRooms 모든 채팅방 리턴
func GetRooms(uid int) ([]vo.Room, error) {
	rooms := []vo.Room{}
	rows, err := db.Query(`SELECT id, name, s_id, b_id FROM chatrooms`)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var room vo.Room
		var sid, bid interface{}
		if err := rows.Scan(&room.ID, &room.Name, &sid, &bid); err != nil {
			log.Println(err)
			return nil, err
		}
		s, sOK := sid.(int64)
		b, bOK := bid.(int64)

		if sOK && bOK {
			if int(s) != uid && int(b) != uid {
				continue
			}
		}

		rooms = append(rooms, room)
	}

	return rooms, nil
}

func GetMessages(roomid int) ([]vo.Message, error) {
	var messages []vo.Message
	rows, err := db.Query(`SELECT m.id, u.username, m.sender_id, m.text, m.sent_on FROM messages m
	JOIN users u ON
	u.id = m.sender_id
	WHERE chatroom_id = $1
	ORDER BY m.id`, roomid)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var message vo.Message
		if err := rows.Scan(&message.ID, &message.Sender, &message.SenderID, &message.Text, &message.SentOn); err != nil {
			log.Println(err)
			return nil, err
		}
		messages = append(messages, message)

	}
	return messages, nil
}
func GetAnotherUser(roomID, id int) (int, error) {
	var sid, bid int
	err := db.QueryRow(`SELECT s_id, b_id 
		FROM chatrooms
		WHERE id=$1`, roomID).Scan(&sid, &bid)
	log.Println("sid, bid=", sid, bid)
	if err != nil {
		log.Println(err)
		return 0, err
	}
	if id == sid {
		return bid, nil
	} else if id == bid {
		return sid, nil
	}
	return 0, errors.New("user not pull join channel")

}

func ConnectToRoom(uid, roomID int) error {
	sid, bid, err := GetChatUser(roomID)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(sid, bid)
	s, sidOk := sid.(int64)
	_, bidOk := bid.(int64)
	if !sidOk && !bidOk {
		res, err := db.Exec(`UPDATE chatrooms SET s_id = $1 WHERE id=$2`, uid, roomID)
		if err != nil {
			log.Println(err)
			return err
		}
		if r, err := res.RowsAffected(); r == 0 || err != nil {
			log.Println(err)
			return err
		}
	} else if !bidOk {
		if uid != int(s) {
			res, err := db.Exec(`UPDATE chatrooms SET b_id = $1 WHERE id=$2`, uid, roomID)
			if err != nil {
				log.Println(err)
				return err
			}
			if r, err := res.RowsAffected(); r == 0 || err != nil {
				log.Println(err)
				return err
			}
		}

	}
	return nil
}

func GetChatUser(roomID int) (sid, bid interface{}, err error) {
	err = db.QueryRow(`SELECT s_id, b_id 
		FROM chatrooms
		WHERE id=$1`, roomID).Scan(&sid, &bid)
	if err != nil {
		log.Println(err)
		return
	}
	return
}