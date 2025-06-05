package messager

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

// IpcMessager can start a Inter Process Communication (IPC)
// It can either be used to start a broadcast server or send messages to one.
type IpcMessager struct {
	Sender   Sender
	Listener Listener
}

func NewMessager() IpcMessager {
	return IpcMessager{
		Sender: Sender{},
		Listener: Listener{
			messages: make(chan []byte, 10),
		},
	}
}

func (messager *IpcMessager) PingMessageListener() (Success bool, err error) {
	Conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return false, err
	}
	defer Conn.Close()

	reply, err := Conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return false, err
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return true, err
	} else {
		return false, err
	}
}
