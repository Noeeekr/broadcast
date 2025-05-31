package client

import (
	"bufio"
	"fmt"
	"os"

	"github.com/Noeeekr/broadcast_server/pkg/instance"
)

// Client estabilishes a connection to broadcast server to send messages (to broadcast to all servers) and recieve broadcasted messages
type Client struct {
	instance.Logger

	// Necessary to check which broadcasted messages are from itself so it won't show in terminal, messages are removed after that
	MessageHistory []string

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
		MessageHistory: []string{},
		Logger:         *instance.NewLogger(),
		MessagesRcvd:   make(chan string, 10),
		MessagesToSend: make(chan string, 10),
		Url:            fmt.Sprintf("%s:%d", url, port),
	}
}

// Run estabilishes a connection to broadcast server.
// It output all recieved messages to command line interface and sends all inputs from cli to server to broadcast.
func (c *Client) Run() error {
	go c.HandleClientMessages()
	err := c.listen(c.Url)

	return err
}

func (c *Client) HandleClientMessages() {
	fmt.Println("Listening for messages. Press [ENTER] to send message.")

	reader := bufio.NewReader(os.Stdout)
	for {
		if msg, err := reader.ReadString('\n'); err == nil {
			c.SendMessage(msg)
		} else {
			fmt.Println("Error happened capturing message from terminal: " + err.Error())
			var instance instance.Shutdown
			instance.Terminate()
		}
	}
}

// Sends a message to server to broadcast
func (c *Client) SendMessage(message string) {
	c.MessagesToSend <- message
}
