package client

import (
	"fmt"
)

// Client estabilishes a connection to broadcast server to send messages (to broadcast to all servers) and recieve broadcasted messages
type Client struct {
	// Contains messages recieved from Client.Connect()
	MessagesRcvd   chan string
	MessagesToSend chan string

	// Broadcast server info
	Url  string
	Port int
}

// Client estabilishes a connection to broadcast server to send messages (to broadcast to all servers) and recieve broadcasted messages
func New(url string, port int) Client {
	return Client{
		MessagesRcvd:   make(chan string, 10),
		MessagesToSend: make(chan string, 10),
		Url:            fmt.Sprintf("%s:%d", url, port),
	}
}

// Connect estabilishes a connection to broadcast server
func (c *Client) Connect() error {
	go c.HandleClientMessages()
	err := c.listen(c.Url)

	return err
}

func (c *Client) HandleClientMessages() {
	fmt.Println("Listening for messages. Press [ENTER] to send message.")

	var message string
	for {
		fmt.Scan(&message)

		c.SendMessage(message)
	}
}

// Sends a message to server to broadcast
func (c *Client) SendMessage(message string) {
	c.MessagesToSend <- message
}
