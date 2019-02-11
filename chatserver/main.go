package main

import (
	"flag"

	"github.com/socketChat/chatserver/server"
)

var (
	host = flag.String("-h", "localhost", "server host name")
	port = flag.String("-p", "1024", "port")
)

func main() {
	flag.Parse()
	server, err := server.NewServer(*host, *port)
	if err != nil {
		panic(err)
	}
	server.Serve()
}
