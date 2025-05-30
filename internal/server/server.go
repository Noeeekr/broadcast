package server

import (
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/Noeeekr/broadcast_server/internal/ipc"
	"github.com/gorilla/websocket"
)

type ChanInfo struct {
	cha chan string
	id  int
}

type Connections struct {
	connMutex sync.Mutex
	channels  map[int]ChanInfo

	limit int
	last  int
}

// Server implements a websocket connection with dbus connection to listen for ipc messages.
type Server struct {
	*Connections

	// Turn HTTP connections into Websockets
	upgrader websocket.Upgrader

	// Inter process comunication (IPC) between server and Command Line Interface (CLI)
	listener ipc.Listener
	messages chan string

	// All open connections
	port int
}

func New() *Server {
	return &Server{
		Connections: &Connections{
			last:     0,
			channels: make(map[int]ChanInfo),
			limit:    10000,
		},
		listener: ipc.Listener{},
		messages: make(chan string, 10),
		upgrader: websocket.Upgrader{},
	}
}

// Adds a new channels to broadcast. Remove must be called when the connection is closed
func (s *Server) Add(cha chan string) (id int, err error) {
	conn := s.Connections

	conn.connMutex.Lock()
	defer conn.connMutex.Unlock()

	if conn.last >= s.limit {
		return id, errors.New("Connection limit reached")
	}

	// Add the channel to the next conn id
	conn.last++
	conn.channels[conn.last] = ChanInfo{
		id:  conn.last,
		cha: cha,
	}

	return id, nil
}

// Remove channels from broadcast
func (s *Server) Remove(id int) {
	conn := s.Connections
	conn.connMutex.Lock()
	defer conn.connMutex.Unlock()

	// Replace chann to be removed with last chan
	if conn.last != id {
		lastChannel := conn.channels[conn.last]
		lastChannel.id = id

		conn.channels[id] = lastChannel
	}

	delete(conn.channels, conn.last)
	conn.last--
}

func (s *Server) Serve(port int) ipc.DetailedErrors {
	if conn, Error := s.listener.StartMessageListener(s.messages); Error.Type != ipc.ErrorNil {
		conn.Close()
		return Error
	}

	fmt.Println("[LOG] Starting Server on port", port)

	go s.Broadcast()

	return ipc.DetailedErrors{
		Type:        ipc.ErrorFailedConnection,
		Description: http.ListenAndServe(fmt.Sprintf(":%d", port), s).Error(),
	}
}

func (s *Server) Broadcast() {
	message := <-s.messages

	for _, messageChannel := range s.channels {
		messageChannel.cha <- message
	}
}
