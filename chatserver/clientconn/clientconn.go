package clientconn

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/socketChat/chatserver/server"
	model "github.com/socketChat/models"
)

const (
	defaultID  = "123"
	bufferSize = 4096
)

type ClientConn struct {
	openID  string
	connID  int32
	conn    net.Conn
	chanMgr *server.ChanMgr
}

func (c *ClientConn) Listen() {
	defer func() {
		fmt.Fprintf(os.Stderr, "%s left chatroom", c.openID)
		c.conn.Close()
	}()

	fmt.Println("Connecting...")
	buffer := make([]byte, bufferSize)
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		msg := model.Message{}
		if err := json.Unmarshal(buffer[:n], &msg); err != nil {
			fmt.Println(err.Error())
			continue
		}
		ch, err := c.chanMgr.GetChannel(msg.ChannelID)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error(), msg.ChannelID)
			continue
		}

		switch msg.Type {
		case model.Join:
			// Init openID and join
			c.openID = msg.OpenID
			if err := ch.Join(model.ID(c.openID), c.conn); err != nil {
				fmt.Println(err.Error())
				return
			}
		case model.Leave:
			if err := ch.Leave(model.ID(c.openID)); err != nil {
				fmt.Println(err.Error())
				return
			}
			break
		case model.Text:
			if err := ch.SendMsg(msg); err != nil {
				fmt.Println(err.Error())
			}
		}
	}
}

func NewClient(conn net.Conn, connID int32, chanMgr *server.ChanMgr) *ClientConn {
	return &ClientConn{
		openID:  defaultID,
		conn:    conn,
		connID:  connID,
		chanMgr: chanMgr,
	}
}
