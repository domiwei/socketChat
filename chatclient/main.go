package main

import (
	"bufio"
	"os"

	"github.com/socketChat/chatclient/client"
)

func main() {
	server := "localhost:1024"
	client := client.NewClient(server)

	for {
		reader := bufio.NewReader(os.Stdin)
		text, _ := reader.ReadString('\n')
		client.InputChan <- text
	}
}
