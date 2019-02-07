package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	model "github.com/socketChat/models"
)

type Client struct {
	conn      net.Conn
	tcpaddr   *net.TCPAddr
	InputChan chan []byte
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
		InputChan: make(chan []byte, 1024),
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
		println("Received from:", c.conn.RemoteAddr(), string(buffer[:n]))
		//channel <- buffer[:n]
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
				Text:      string(data),
			}
			b, _ := json.Marshal(&msg)
			_, err := c.conn.Write(b)
			if err != nil {
				fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
			}
			println("Write to:", c.conn.RemoteAddr(), string(data))
		}
	}
}

func main() {
	server := "localhost:1024"
	client := NewClient(server)

	for {
		var s string
		fmt.Scan(&s)
		client.InputChan <- []byte(s)
	}

	//go readServer()
	//go writeServer()
}
