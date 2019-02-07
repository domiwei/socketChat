package clientconn

import (
	"encoding/json"
	"fmt"
	"net"

	"github.com/socketChat/chatserver/server"
	model "github.com/socketChat/models"
)

const (
	defaultID  = "123"
	bufferSize = 4096
)

type ClientConn struct {
	openID  string
	conn    net.Conn
	chanMgr *server.ChanMgr
}

func (c *ClientConn) Listen() {
	defer c.conn.Close()

	buffer := make([]byte, bufferSize)
	for {
		_, err := c.conn.Read(buffer)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		msg := model.Message{}
		if err := json.Unmarshal(buffer, &msg); err != nil {
			fmt.Println(err.Error())
			continue
		}
		ch, ok := c.chanMgr.Channels[msg.ChannelID]
		if !ok {
			continue
		}

		switch msg.Type {
		case model.Join:
			// Init openID and join
			c.openID = msg.OpenID
			if err := ch.Join(c.openID, c.conn); err != nil {
				//TODO
			}
		case model.Leave:
			if err := ch.Leave(c.openID); err != nil {
				//TODO
			}
		case model.Text:
			select {
			case ch.MsgChan <- msg:
				fmt.Println(msg)
			default:
				//TODO
			}
		}
	}
}

func NewClient(conn net.Conn, chanMgr *server.ChanMgr) *ClientConn {
	return &ClientConn{
		openID:  defaultID,
		conn:    conn,
		chanMgr: chanMgr,
	}
}
