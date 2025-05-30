package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Turn connection into websocket connections
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to enable continuous connection."))
	}
	defer conn.Close()

	fmt.Println("client connected")

	// Request channel
	var requests chan []byte = make(chan []byte, 1)
	go func() {
		for {
			_, msg, err := conn.ReadMessage()

			if err != nil {
				fmt.Println("Websockets conn closed")
				break
			}

			requests <- msg
		}
		close(requests)
	}()

	// Broadcast channel
	var messages chan string = make(chan string, 10)
	id, err := s.Add(messages)
	defer s.Remove(id)

outerloop:
	for {
		select {
		case message, ok := <-messages:
			if !ok {
				break outerloop
			}
			conn.WriteMessage(websocket.TextMessage, []byte(message))
			fmt.Println("message sent to clients")
			break
		case request, ok := <-requests:
			if !ok {
				break outerloop
			}
			fmt.Println("message recieved from clients")
			conn.WriteMessage(websocket.TextMessage, []byte("ok"))
			s.listener.SendMessage(string(request))
			break
		}
	}
}
