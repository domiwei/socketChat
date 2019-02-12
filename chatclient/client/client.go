package client

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"sync"
	"time"

	model "github.com/socketChat/models"
	"golang.org/x/net/websocket"
)

type ConnType int

const (
	Socket ConnType = iota
	WebSocket
)

type Client struct {
	conn       io.ReadWriteCloser
	openID     string
	chatOutput io.Writer
	InputChan  chan string
	wg         sync.WaitGroup
}

func NewClient(host, port, openName string, chatOutput io.Writer, conntype ConnType) *Client {
	var conn io.ReadWriteCloser
	switch conntype {
	case Socket:
		server := host + ":" + port
		tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
		if err != nil {
			log.Fatal(err)
		}

		socketconn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Fatal(err)
		}
		conn = socketconn
	case WebSocket:
		origin := fmt.Sprintf("http://%s/", host)
		url := fmt.Sprintf("ws://%s:%s/chat", host, port)
		websocketconn, err := websocket.Dial(url, "", origin)
		if err != nil {
			log.Fatal(err)
		}
		conn = websocketconn
	}

	// Join channel
	msg := model.Message{
		Type:      model.Join,
		ChannelID: "happy-pig-year",
		OpenID:    openName,
	}
	b, _ := json.Marshal(&msg)
	if _, err := conn.Write(b); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	client := &Client{
		conn:       conn,
		openID:     openName,
		chatOutput: chatOutput,
		InputChan:  make(chan string, 1024),
	}
	go client.read()
	fmt.Fprintln(client.chatOutput, "successfully connect...")
	return client
}

func (c *Client) read() {
	defer func() {
		fmt.Fprintf(c.chatOutput, "Connection lost...")
		c.wg.Done()
	}()
	c.wg.Add(1)
	buffer := make([]byte, 65536)
	for {
		n, err := c.conn.Read(buffer)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			return
		}
		msgs := []model.Message{}
		if err := json.Unmarshal(buffer[:n], &msgs); err != nil {
			fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			return
		}
		for _, msg := range msgs {
			if msg.OpenID == c.openID {
				msg.OpenID = "->" + msg.OpenID
			}
			t := time.Unix(msg.Timestamp, int64(0)).Format("Mon Jan _2 15:04:05 2006")
			fmt.Fprintf(c.chatOutput, "%s (%s): %s", msg.OpenID, t, msg.Text)
		}
	}
}

func (c *Client) Send(text string) error {
	msg := model.Message{
		Type:      model.Text,
		ChannelID: "happy-pig-year",
		OpenID:    c.openID,
		Text:      text,
	}
	b, err := json.Marshal(&msg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return err
	}
	_, err = c.conn.Write(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		return err
	}
	return nil
}

func (c *Client) Close() {
	c.conn.Close()
	c.wg.Wait()
}
