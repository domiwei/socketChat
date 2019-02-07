package client

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

type Client struct {
	openID string
	conn   net.Conn
	server *server.Server
}

func (c *Client) Listen() {
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
		ch, ok := c.server.Channels[msg.ChannelID]
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

func NewClient(conn net.Conn, server *server.Server) *Client {
	return &Client{
		openID: defaultID,
		conn:   conn,
		server: server,
	}
}
