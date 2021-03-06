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
		panic("Must need an open name")
		return
	}
	client := client.NewClient(*host, *port, *name, os.Stdout, client.WebSocket)
	defer client.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		text, _ := reader.ReadString('\n')
		if err := client.Send(text); err != nil {
			break
		}
	}
}
