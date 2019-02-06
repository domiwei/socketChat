package main

import (
	"flag"
	"fmt"
	"net"

	"github.com/socketChat/server/chat"
	"github.com/socketChat/server/client"
)

var (
	host = flag.String("-h", "localhost", "server host name")
	port = flag.String("-p", "1024", "port")
)

const ()

func main() {
	flag.Parse()
	addr := *host + ":" + *port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	// Init a chat room and run
	chatroom := chat.NewChatRoom("happy-pig-year")
	go chatroom.Serve()
	// Accept new client
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		client := client.NewClient(conn, chatroom.MsgChan)
		go client.Listen()
	}
}
