package main

import (
	"flag"
	"log"
	"net"

	channel "github.com/socketChat/chatserver/channel"
	"github.com/socketChat/chatserver/clientconn"
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
		log.Println(err.Error())
		panic(err)
	}
	server := server.NewChanMgr(addr)
	// Init a chat room and run
	ch := channel.NewChannel("happy-pig-year")
	go ch.Serve()
	if err := server.AddChannel(ch); err != nil {
		log.Println(err.Error())
		panic(err)
	}
	// Accept new clientconn
	connID := int32(0)
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println(err.Error())
			break
		}
		clientconn := clientconn.NewClient(conn, connID, server)
		go clientconn.Listen()
		connID++
	}
}
