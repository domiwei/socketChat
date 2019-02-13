package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/socketChat/chatserver/server"
)

var (
	host = flag.String("h", "localhost", "server host name")
	port = flag.String("p", "1024", "port")
)

func main() {
	flag.Parse()
	server, err := server.NewWebSocketServer(*port)
	if err != nil {
		panic(err)
	}
	// Notify shutdown
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		server.ShutDown()
	}()
	// Serve
	server.Serve()
}
