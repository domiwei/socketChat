package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	model "github.com/socketChat/models"
)

type Client struct {
	conn       net.Conn
	tcpaddr    *net.TCPAddr
	openID     string
	chatOutput io.Writer
	InputChan  chan string
}

func NewClient(server, openName string, chatOutput io.Writer) *Client {
	tcpAddr, err := net.ResolveTCPAddr("tcp4", server)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
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
		tcpaddr:    tcpAddr,
		openID:     openName,
		chatOutput: chatOutput,
		InputChan:  make(chan string, 1024),
	}
	go client.read()
	go client.write()
	fmt.Fprintln(client.chatOutput, "successfully connect...")
	return client
}

func (c *Client) read() {
	defer c.conn.Close()
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

func (c *Client) write() {
	defer c.conn.Close()
	for {
		select {
		case data := <-c.InputChan:
			msg := model.Message{
				Type:      model.Text,
				ChannelID: "happy-pig-year",
				OpenID:    c.openID,
				Text:      data,
			}
			b, err := json.Marshal(&msg)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
				return
			}
			_, err = c.conn.Write(b)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
				return
			}
		}
	}
}
