package client

import (
	"net"

	model "github.com/socketChat/models"
)

const (
	defaultID = "123"
)

type Client struct {
	openID  string
	conn    net.Conn
	msgChan chan<- model.Message
}

func (c *Client) Listen() {
	for {
		select {}
	}
}

func NewClient(conn net.Conn, mc chan<- model.Message) *Client {
	return &Client{
		openID:  defaultID,
		conn:    conn,
		msgChan: mc,
	}
}
