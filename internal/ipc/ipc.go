package ipc

import "github.com/godbus/dbus/v5"

const (
	name           string = "com.broadcastserver.Messages"
	interface_name string = name + ".MessageListenerInterface"
	object_path    string = "/com/broadcastserver/Messages/MessageListener"
	method         string = interface_name + ".SendMessage"
	signal_name    string = "SendMessage"
)

type ConnectionErrors string

const (
	ErrorConnAlreadyExists ConnectionErrors = "ConnAlreadyExists"
	ErrorFailedConnection  ConnectionErrors = "FailedConnection"
	ErrorNil               ConnectionErrors = ""
)

// Connection can start a Inter Process Communication (IPC)
// It can either be used to start a broadcast server or send messages to one.
type Connection struct {
	Sender   Sender
	Listener Listener
}

func New() Connection {
	return Connection{
		Sender: Sender{},
		Listener: Listener{
			messages: make(chan string, 30),
		},
	}
}

func (cli *Connection) IsServerListening() bool {
	Conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return true
	}
	defer Conn.Close()

	reply, err := Conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return true
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return true
	} else {
		return false
	}
}
