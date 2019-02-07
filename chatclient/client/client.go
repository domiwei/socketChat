package client

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"time"

	model "github.com/socketChat/models"
)

type Client struct {
	conn      net.Conn
	tcpaddr   *net.TCPAddr
	InputChan chan string
}

func NewClient(server string) *Client {
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
		OpenID:    "keweigg",
	}
	b, _ := json.Marshal(&msg)
	if _, err := conn.Write(b); err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}

	client := &Client{
		conn:      conn,
		tcpaddr:   tcpAddr,
		InputChan: make(chan string, 1024),
	}
	go client.read()
	go client.write()
	fmt.Println("successfully connect...")
	return client
}

func (c *Client) read() {
	buffer := make([]byte, 2048)
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
			t := time.Unix(msg.Timestamp, int64(0)).Format("Mon Jan _2 15:04:05 2006")
			fmt.Printf("%s (%s): %s", msg.OpenID, t, msg.Text)
		}
		//println("Received from:", c.conn.RemoteAddr(), string(buffer[:n]))
	}
}

func (c *Client) write() {
	for {
		select {
		case data := <-c.InputChan:
			msg := model.Message{
				Type:      model.Text,
				ChannelID: "happy-pig-year",
				OpenID:    "keweigg",
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
			}
			//println("Write to:", c.conn.RemoteAddr(), data)
		}
	}
}
