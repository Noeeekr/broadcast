package client

import (
	"fmt"
	"slices"

	"github.com/Noeeekr/broadcast_server/pkg/instance"
	"github.com/gorilla/websocket"
)

// Listen estabilishes a connection to broadcast server and listen to all messages
func (c *Client) listen(url string) error {
	c.MessagesRcvd = make(chan string, 10)
	defer close(c.MessagesRcvd)
	c.MessagesToSend = make(chan string, 10)
	defer close(c.MessagesToSend)

	Conn, _, err := websocket.DefaultDialer.Dial(url, nil)
	if err != nil {
		return err
	}
	c.SimpleLog("Connection", "Estabilished", "")

	var waiter instance.Shutdown

	go func() {
		for {
			msgType, msg, err := Conn.ReadMessage()

			if err != nil {
				fmt.Println("Server requested to close connection")
				Conn.Close()
				waiter.Terminate()
				return
			}

			if msgType == websocket.TextMessage {
				c.MessagesRcvd <- string(msg)
			} else {
				c.MessagesRcvd <- "Unknown message recieved"
			}
		}
	}()

outerloop:
	for {
	selectLoop:
		select {
		case <-instance.InterruptContext.Done():
			fmt.Println("Signal recieved. Finishing program")
			break outerloop
		case msg := <-c.MessagesRcvd:
			for i, message := range c.MessageHistory {
				if message == msg {
					c.MessageHistory = slices.Delete(c.MessageHistory, i, i)
					break selectLoop
				}
			}

			fmt.Printf("S:%s", msg)
			break
		case msg := <-c.MessagesToSend:
			err := Conn.WriteMessage(websocket.TextMessage, []byte(msg))
			if err != nil {
				fmt.Println("Failed to send message")
			}
			c.MessageHistory = append(c.MessageHistory, msg)
			break
		}
	}

	return nil
}
