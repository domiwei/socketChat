package main

import (
	"bufio"
	"flag"
	"os"

	"github.com/socketChat/chatclient/client"
)

var (
	host = flag.String("h", "localhost", "server host name")
	port = flag.String("p", "1024", "port")
	name = flag.String("n", "", "openID in chat room")
)

func main() {
	flag.Parse()
	if *name == "" {
		panic("Need an open name")
		return
	}
	server := *host + ":" + *port
	client := client.NewClient(server, *name, os.Stdout)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		client.InputChan <- text
	}
}
