package clientconn

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"log"
	"net"

	"github.com/socketChat/chatserver/channel"
	model "github.com/socketChat/models"
)

const (
	defaultID  = "ConnectedUser"
	bufferSize = 4096
)

type ClientConn struct {
	openID  string
	userID  model.ID
	connID  int32
	conn    net.Conn
	chanMgr *channel.ChanMgr
}

func (c *ClientConn) Listen() {
	defer func() {
		log.Printf("%s, %s left chatroom", c.openID, c.userID)
		err := c.chanMgr.LeaveAllChannels(c.userID)
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
		// Fill in userID in msg model
		msg.UserID = c.userID
		// Get target channel and do approrpiate action
		ch, err := c.chanMgr.GetChannel(msg.ChannelID)
		if err != nil {
			log.Printf(err.Error(), msg.ChannelID)
			continue
		}
		switch msg.Type {
		case model.Join:
			// Join channel
			if err := ch.Join(c.userID, msg.OpenID, c.conn); err != nil {
				log.Println(err.Error())
				return
			}
		case model.Leave:
			if err := ch.Leave(c.userID); err != nil {
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

func NewClient(conn net.Conn, connID int32, chanMgr *channel.ChanMgr) *ClientConn {
	return &ClientConn{
		openID:  defaultID,
		userID:  pseudoUUID(),
		conn:    conn,
		connID:  connID,
		chanMgr: chanMgr,
	}
}

func pseudoUUID() (uuid model.ID) {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		log.Println("Error: ", err)
		return
	}
	uuid = model.ID(fmt.Sprintf("%X-%X-%X-%X-%X", b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))
	return
}
