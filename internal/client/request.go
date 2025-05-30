package client

import (
	"fmt"

	Panic "github.com/Noeeekr/broadcast_server/pkg/panic"
	"github.com/gorilla/websocket"
)

// Listen estabilishes a connection to broadcast server and listen to all messages
func (c *Client) listen(url string) error {
	var p Panic.Panic
	p.NotImplemented("internal/client/request.go - Not implemented - Fmt the response message from server 'conn estabilished. welcome from server'")
	Conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}

	fmt.Println("Connection estabilised")

	go func() {
		defer Conn.Close()
		for {
			msgType, msg, err := Conn.ReadMessage()
			p.NotImplemented("internal/client/request.go - Not implemented - Handle error case for EOF | Conn Closed")

			if err != nil {
				fmt.Println("Error happened in internal/client/request", err.Error())
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
			fmt.Println("Message recieved:", msg)
			break
		case msg := <-c.MessagesToSend:
			c.send(msg)
			break
		}
	}
}

func (c *Client) send(message string) {
	Conn, _, err := websocket.DefaultDialer.Dial(c.Url, nil)
	if err != nil {
		fmt.Println("internal/client/request.go", err.Error())
		return
	}

	err = Conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		fmt.Println("internal/client/request.go 2", err.Error())
		return
	}

	fmt.Println("Message send successfully")
}
