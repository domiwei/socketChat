package chat

import (
	"fmt"
	"net"

	model "github.com/socketChat/models"
)

type User struct {
	Conn      net.Conn
	MsgBuffer []model.Message
}

type ChatRoom struct {
	RoomID  string
	MsgChan chan model.Message
	Clients map[string]User
}

func (c *ChatRoom) Serve() {
	for {
		select {
		case msg := <-c.MsgChan:
			fmt.Println(msg)
		}
	}
}

func NewChatRoom(roomID string) *ChatRoom {
	return &ChatRoom{
		RoomID:  roomID,
		MsgChan: make(chan model.Message, 100000),
		Clients: map[string]User{},
	}
}
