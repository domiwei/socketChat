package channel

import (
	"fmt"
	"net"

	model "github.com/socketChat/models"
)

type User struct {
	Conn      net.Conn
	MsgBuffer []model.Message
}

type Channel struct {
	ChannelID string
	MsgChan   chan model.Message
	Clients   map[string]User
}

func (c *Channel) Serve() {
	for {
		select {
		case msg := <-c.MsgChan:
			fmt.Println(msg)
		}
	}
}

func (c *Channel) Join(openID string, conn net.Conn) error {
	return nil
}

func (c *Channel) Leave(openID string) error {
	return nil
}

func NewChannel(cID string) *Channel {
	return &Channel{
		ChannelID: cID,
		MsgChan:   make(chan model.Message, 100000),
		Clients:   map[string]User{},
	}
}
