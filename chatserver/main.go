package main

import (
	"flag"
	"fmt"
	"net"

	channel "github.com/socketChat/chatserver/channel"
	"github.com/socketChat/chatserver/client"
	"github.com/socketChat/chatserver/server"
)

var (
	host = flag.String("-h", "localhost", "server host name")
	port = flag.String("-p", "1024", "port")
)

func main() {
	flag.Parse()
	addr := *host + ":" + *port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	server := server.NewServer(addr)
	// Init a chat room and run
	ch := channel.NewChannel("happy-pig-year")
	go ch.Serve()
	if err := server.AddChannel(ch); err != nil {
		fmt.Println(err.Error())
		panic(err)
	}
	// Accept new client
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err.Error())
			break
		}
		client := client.NewClient(conn, server)
		go client.Listen()
	}
}
