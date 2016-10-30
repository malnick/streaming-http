package main

import (
	"fmt"
	"os"

	"github.com/malnick/test-chunked-server/client"
	"github.com/malnick/test-chunked-server/server"
)

func main() {
	if os.Args[1] == "server" {
		fmt.Println("Server started...")
		server.StartServer()
	}

	if os.Args[1] == "client" {
		fmt.Println("Client started...")
		client.StartClient()
	}

	fmt.Println("Error, no selection")
}
