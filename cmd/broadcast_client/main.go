package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Noeeekr/broadcast_server/internal/client"
	"github.com/Noeeekr/broadcast_server/internal/instance"
	"github.com/gorilla/websocket"
)

type Client struct {
	upgrader websocket.Upgrader
}

func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func main() {
	var Port int
	var EnableDebug bool
	flag.IntVar(&Port, "port", 0, "The port of the server to connect to")
	flag.BoolVar(&EnableDebug, "debug", true, "Enable debug mechanisms, like panicking when finding a non implemented mechanism")

	flag.Parse()

	var instance instance.Instance = instance.New()
	instance.EnableDebug(EnableDebug)
	instance.EnableGracefull()

	if Port == 0 {
		fmt.Println("Port is necessary")
		return
	}

	fmt.Println("Trying to start connection")

	// Starts connection to send messages
	var c client.Client = client.New("ws://localhost", Port)
	if err := c.Run(); err != nil {
		fmt.Println("Connection finished", err.Error())
	}
}
