package server

import (
	"fmt"
	"net/http"
	"time"

	Broadcast "github.com/Noeeekr/broadcast_server/internal"
	"github.com/Noeeekr/broadcast_server/pkg/instance"
	"github.com/gorilla/websocket"
)

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Turn connection into websocket connections
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Failed to enable websocket connection."))
	}
	defer conn.Close()

	// Incoming messages channel
	var messages chan []byte = make(chan []byte, 10)
	defer close(messages)

	id, err := s.Add(messages)
	defer s.Remove(id)

	s.SimpleLog("Connection", "Client connected "+conn.RemoteAddr().String(), "")
	fmt.Println()
	go func() {
		for {
			_, msg, err := conn.ReadMessage()

			if err != nil {
				s.SimpleLog("Connection", "Closed connection to "+conn.RemoteAddr().String(), "")
				break
			}
			s.SimpleLog("Broadcast", "Message request recieved from "+conn.LocalAddr().String(), "")

			s.listener.SendMessage(string(msg))
		}
	}()

	// Shutdown mechanism
	var Shutdown instance.Shutdown
	Proceed := Shutdown.Wait()
	defer Proceed()

	// Multi state machine handling incoming/outcoming messages and program termination messages
outerloop:
	for {
		select {
		// Close all connections on interrupt signal
		case <-instance.InterruptContext.Done():
			err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server requested close"), time.Now().Add(time.Second*2))
			if err != nil {
				s.ErrorLog("Failed to close connection properly", err.Error())
			}

			break outerloop
		// Send message to all connections
		case message, ok := <-messages:
			if !ok {
				break outerloop
			}
			// If message contains specific command, terminate.
			if string(message) == Broadcast.CommandsCloseConnection {
				message = []byte("Ending connection soon")
				conn.WriteMessage(websocket.TextMessage, message)

				Shutdown.Terminate()
				break outerloop
			}

			conn.WriteMessage(websocket.TextMessage, message)
			s.SimpleLog("BROADCAST", "Message sent to "+conn.RemoteAddr().String(), "")
			break
		}
	}
}
