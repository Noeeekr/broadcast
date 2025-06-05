package main

import (
	"flag"
	"fmt"

	"github.com/Noeeekr/broadcast_server/internal/instance"
	"github.com/Noeeekr/broadcast_server/internal/ipc"
	"github.com/Noeeekr/broadcast_server/internal/server"
)

// Starts a server if a --port flag is provided
// Sends a message to the server if already exists
//
//	Otherwise prints message to enable the server first
func main() {
	var port int
	var enableDebug bool

	flag.IntVar(&port, "port", 3332, "Defines the port the server will be listening to")
	flag.BoolVar(&enableDebug, "debug", true, "Defines if not implemented features will panic.")
	flag.Parse()

	// Enable production panic for unimplemented features
	var instance instance.Instance = instance.New()
	instance.EnableDebug(enableDebug)

	// Enable gracefull shutdown
	instance.EnableGracefull()

	var server *server.Server = server.New()

	// If port flag is present then it is a server start
	if port != 3332 {
		if port < 3332 {
			fmt.Println("[ERROR] Port must be bigger than 3332")
		}
		if Error := server.Serve(port); Error.Type != "" {
			fmt.Println("[ERROR] Failed to start server on port", port)
			fmt.Println("[ERROR]", Error.Description)
			instance.Terminate()
		}
	}

	// Otherwise, it is a message, act as CLI
	// Check if server is running
	communication := ipc.New()
	if communication.IsServerListening() {
		if message := instance.ArgsToString(); message != "" {
			if Error := communication.Sender.Send(message); Error.Type != ipc.ErrorNil {
				fmt.Println("Failed to send message to server - ", Error.Description)
			} else {
				fmt.Println("Message successfully sent..")
			}
			instance.Terminate()
		} else {
			fmt.Println("Provide a message to send to all clients")
			instance.Terminate()
		}
	} else {
		fmt.Println("Start a server first to send messages")
		instance.Terminate()
	}
}
