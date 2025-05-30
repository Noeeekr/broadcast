package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/Noeeekr/broadcast_server/internal/client"
	Panic "github.com/Noeeekr/broadcast_server/pkg/panic"
	"github.com/gorilla/websocket"
)

type Client struct {
	upgrader websocket.Upgrader
}

func (c *Client) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

func main() {
	var port int
	var debug bool
	flag.IntVar(&port, "port", 0, "The port of the server to connect to")
	flag.BoolVar(&debug, "debug", true, "Enable debug mechanisms, like panicking when finding a non implemented mechanism")

	flag.Parse()

	var p Panic.Panic
	p.EnableDebug(debug)

	if port == 0 {
		fmt.Println("Port is necessary")
		return
	}

	fmt.Println("Trying to start connection")

	// Starts connection to send messages
	var c client.Client = client.New("ws://localhost", port)
	if err := c.Connect(); err != nil {
		fmt.Println("Connection finished", err.Error())
	}
}
