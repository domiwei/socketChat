package channel

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	model "github.com/socketChat/models"
)

var (
	ErrOpenIDExist    = fmt.Errorf("OpenID exists")
	ErrUserIDExist    = fmt.Errorf("UserID exists")
	ErrOpenIDNotExist = fmt.Errorf("OpenID does not exist")
	ErrMsgChanIsFull  = fmt.Errorf("Too many incoming messages")
)

const (
	msgChanBufSize = 10000
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
		msgChan:       make(chan model.Message, msgChanBufSize),
		history:       []model.Message{},
		users:         map[model.ID]*user{},
		broadcastChan: make(chan struct{}, 1),
	}
}

type user struct {
	OpenID   string
	Conn     io.ReadWriteCloser
	MsgIndex int32
}

func (c *Channel) Serve() {
	log.Println("Channel " + c.ChannelID + " is serving")
	go c.broadcast()
	for {
		select {
		case msg := <-c.msgChan:
			msg.Timestamp = time.Now().Unix()
			if msg.UserID != "" { // Not a system message
				u, ok := c.users[msg.UserID]
				if ok {
					msg.OpenID = u.OpenID
				} else {
					log.Println("Unknown userID in chatroom ", c.ChannelID)
					break
				}
			}
			c.history = append(c.history, msg)
			c.notifyBroadcast()
		}
	}
}

func (c *Channel) Join(userID model.ID, openID string, conn io.ReadWriteCloser) error {
	c.usersMutex.Lock()
	defer c.usersMutex.Unlock()
	// Check if this user exists or not
	if _, exist := c.users[userID]; exist {
		return ErrOpenIDExist
	}
	// Check if this openID exists or not
	for _, user := range c.users {
		if user.OpenID == openID {
			return ErrUserIDExist
		}
	}
	// Add a new user
	c.users[userID] = &user{
		OpenID:   openID,
		Conn:     conn,
		MsgIndex: int32(0),
	}
	// Notify all clients
	c.SendMsg(model.Message{
		Type:      model.Text,
		OpenID:    "system",
		ChannelID: c.ChannelID,
		Text:      openID + " joins chatroom\n",
	})
	return nil
}

func (c *Channel) Leave(userID model.ID) error {
	c.usersMutex.Lock()
	defer c.usersMutex.Unlock()
	// Check if user exists or not
	u, exist := c.users[userID]
	if !exist {
		return ErrOpenIDNotExist
	}
	// Remove
	delete(c.users, userID)
	c.SendMsg(model.Message{
		Type:      model.Text,
		OpenID:    "system",
		ChannelID: c.ChannelID,
		Text:      u.OpenID + " left chatroom\n",
	})
	return nil
}

func (c *Channel) SendMsg(msg model.Message) error {
	select {
	case c.msgChan <- msg:
		log.Println("New message: ", msg.Text, msg.ChannelID, msg.OpenID)
	default:
		return ErrMsgChanIsFull
	}
	return nil
}

func (c *Channel) notifyBroadcast() {
	select {
	case c.broadcastChan <- struct{}{}: // Just notify
	default:
	}
}

func (c *Channel) broadcast() error {
	for {
		select {
		case <-c.broadcastChan:
			// Concurrently broadcast messages
			leftIDs := []model.ID{}
			historyEnd := int32(len(c.history))
			var wg sync.WaitGroup
			for id, u := range c.users {
				wg.Add(1)
				go func(id model.ID, u *user) {
					defer wg.Done()
					msgs := c.history[u.MsgIndex:historyEnd]
					b, _ := json.Marshal(msgs)
					// Send to client
					if _, err := u.Conn.Write(b); err != nil {
						log.Println("Error on " + string(id) + " : " + err.Error())
						leftIDs = append(leftIDs, id)
						return
					}
					// Update index of chat history for each client
					c.users[id].MsgIndex = historyEnd
				}(id, u)
			}
			wg.Wait()
			for _, id := range leftIDs {
				if err := c.Leave(id); err != nil {
					log.Println("Error: ", err.Error())
				}
			}
			log.Println(c.users)
		}
	}
	return nil
}
