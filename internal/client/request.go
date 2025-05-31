package client

import (
	"fmt"

	Broadcast "github.com/Noeeekr/broadcast_server/internal"
	"github.com/Noeeekr/broadcast_server/pkg/instance"
	"github.com/gorilla/websocket"
)

// Listen estabilishes a connection to broadcast server and listen to all messages
func (c *Client) listen(url string) error {
	var debug instance.Debugger
	debug.NotImplemented("internal/client/request.go - Not implemented - Fmt the response message from server 'conn estabilished. welcome from server'")
	Conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	c.SimpleLog("Connection", "Estabilished", "")

	var instance instance.Shutdown

	go func() {
		for {
			msgType, msg, err := Conn.ReadMessage()
			debug.NotImplemented("internal/client/request.go - Not implemented - Handle error case for EOF | Conn Closed")

			if err != nil {
				fmt.Println("server requested to close connection")
				instance.Terminate()
				return
			}

			if msgType == websocket.TextMessage {
				c.MessagesRcvd <- string(msg)
			} else {
				c.MessagesRcvd <- "Unknown message recieved"
			}
		}
	}()

	for {
		select {
		case msg := <-c.MessagesRcvd:
			if msg == Broadcast.CommandsCloseConnection {
				fmt.Println("Server request to close connection. Terminating")
				instance.Terminate()
				break
			}
			fmt.Println("Message recieved:", msg)
			break
		case msg := <-c.MessagesToSend:
			err := Conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Failed to send message")
			}
			break
		}
	}
}
