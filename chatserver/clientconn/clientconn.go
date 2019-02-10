package clientconn

import (
	"encoding/json"
	"log"
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
	connID  int32
	conn    net.Conn
	chanMgr *server.ChanMgr
}

func (c *ClientConn) Listen() {
	defer func() {
		log.Printf("%s left chatroom", c.openID)
		// For now, use openID as key to find out all channels user belonging,
		// But actually need to lookup by uid to avoid change of openID.
		err := c.chanMgr.LeaveAllChannels(c.openID)
		if err != nil {
			log.Println(err.Error())
		}
		c.conn.Close()
	}()

	log.Println("Connecting...")
	buffer := make([]byte, bufferSize)
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			log.Println(err.Error())
			break
		}
		msg := model.Message{}
		if err := json.Unmarshal(buffer[:n], &msg); err != nil {
			log.Println(err.Error())
			continue
		}
		// Get target channel and do approrpiate action
		ch, err := c.chanMgr.GetChannel(msg.ChannelID)
		if err != nil {
			log.Printf(err.Error(), msg.ChannelID)
			continue
		}
		switch msg.Type {
		case model.Join:
			// Init openID and join
			c.openID = msg.OpenID
			if err := ch.Join(model.ID(c.openID), c.conn); err != nil {
				log.Println(err.Error())
				return
			}
		case model.Leave:
			if err := ch.Leave(model.ID(c.openID)); err != nil {
				log.Println(err.Error())
				return
			}
			break
		case model.Text:
			if err := ch.SendMsg(msg); err != nil {
				log.Println(err.Error())
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
