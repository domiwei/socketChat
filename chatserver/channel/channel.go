package channel

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	model "github.com/socketChat/models"
)

var (
	ErrOpenIDExist    = fmt.Errorf("OpenID for this user exists")
	ErrOpenIDNotExist = fmt.Errorf("OpenID for this user not exist")
	ErrMsgChanIsFull  = fmt.Errorf("Too many incoming messages")
)

type Channel struct {
	ChannelID     string
	msgChan       chan model.Message
	history       []model.Message
	users         map[model.ID]*user
	usersMutex    sync.Mutex
	broadcastChan chan struct{}
}

func NewChannel(cID string) *Channel {
	return &Channel{
		ChannelID:     cID,
		msgChan:       make(chan model.Message, 100000),
		history:       []model.Message{},
		users:         map[model.ID]*user{},
		broadcastChan: make(chan struct{}, 1),
	}
}

type user struct {
	Conn     net.Conn
	MsgIndex int32
}

func (c *Channel) Serve() {
	fmt.Println("Channel " + c.ChannelID + " is serving")
	go c.broadcast()
	for {
		select {
		case msg := <-c.msgChan:
			msg.Timestamp = time.Now().Unix()
			fmt.Println(msg)
			c.history = append(c.history, msg)
			c.notifyBroadcast()
		}
	}
}

func (c *Channel) Join(openID model.ID, conn net.Conn) error {
	c.usersMutex.Lock()
	defer c.usersMutex.Unlock()
	// Check if user exists or not
	if _, exist := c.users[openID]; exist {
		return ErrOpenIDExist
	}
	// Add a new user
	c.users[openID] = &user{
		Conn:     conn,
		MsgIndex: int32(0),
	}
	return nil
}

func (c *Channel) Leave(openID model.ID) error {
	c.usersMutex.Lock()
	defer c.usersMutex.Unlock()
	// Check if user exists or not
	if _, exist := c.users[openID]; !exist {
		return ErrOpenIDNotExist
	}
	// Remove
	delete(c.users, openID)
	return nil
}

func (c *Channel) SendMsg(msg model.Message) error {
	select {
	case c.msgChan <- msg:
		fmt.Println(msg)
	default:
		return ErrMsgChanIsFull
	}
	return nil
}

func (c *Channel) notifyBroadcast() {
	select {
	case c.broadcastChan <- struct{}{}:
	default:
	}
}

func (c *Channel) broadcast() error {
	for {
		select {
		case <-c.broadcastChan:
			// TODO: concurrently send
			historyEnd := int32(len(c.history))
			for id, user := range c.users {
				msgs := c.history[user.MsgIndex:historyEnd]
				b, _ := json.Marshal(msgs)
				// Send to client
				user.Conn.Write(b)
				// Update history indexing
				c.users[id].MsgIndex = historyEnd
			}
		}
	}
	return nil
}