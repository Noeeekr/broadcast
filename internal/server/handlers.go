package server

import (
	"net/http"
	"time"

	Broadcast "github.com/Noeeekr/broadcast_server/internal"
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

	s.SimpleLog("Connection", "Client connected "+conn.RemoteAddr().String(), "")

	// Request channel
	var requests chan []byte = make(chan []byte, 1)
	go func() {
		for {
			_, msg, err := conn.ReadMessage()

			if err != nil {
				s.SimpleLog("Connection", "Closed connection to "+conn.RemoteAddr().String(), "")
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
			if message == Broadcast.CommandsCloseConnection {
				err := conn.WriteControl(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server requested close"), time.Now().Add(time.Second*2))
				if err != nil {
					s.ErrorLog("Failed to close connection properly", err.Error())
				}
				break
			}

			conn.WriteMessage(websocket.TextMessage, []byte(message))
			s.SimpleLog("BROADCAST", "Message sent to "+conn.LocalAddr().String(), "")
			break
		case request, ok := <-requests:
			if !ok {
				break outerloop
			}
			s.SimpleLog("BROADCAST", "Message recieved from "+conn.LocalAddr().String(), "")
			s.listener.SendMessage(string(request))
			break
		}
	}
}
