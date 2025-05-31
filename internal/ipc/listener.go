package ipc

import (
	"fmt"

	"github.com/godbus/dbus/v5"
)

// Empty string if no error happened
type DetailedErrors struct {
	Type        ConnectionErrors
	Description string
}

// Listens for messages sent from Sender.
type Listener struct {
	// A channel to hold recieved messages
	messages chan string
}

// Start a dbus connection to listen for command line messages.
// Throws an error if another connection is already started.
// Stores recieved messages in messageChannel.
func (s *Listener) StartMessageListener(messageChannel chan string) (conn *dbus.Conn, errors DetailedErrors) {
	conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return conn, DetailedErrors{
			Type:        ErrorFailedConnection,
			Description: fmt.Sprintf("%s - %s", "Failed to start CLI message listener", err.Error()),
		}
	}

	reply, err := conn.RequestName(name, dbus.NameFlagDoNotQueue)
	if err != nil {
		return conn, DetailedErrors{
			Type:        ErrorFailedConnection,
			Description: fmt.Sprintf("%s - %s", "Failed to start CLI message listener", err.Error()),
		}
	}

	if reply != dbus.RequestNameReplyPrimaryOwner {
		return conn, DetailedErrors{
			Type:        ErrorConnAlreadyExists,
			Description: fmt.Sprintf("%s - %s", "Failed to start CLI message listener", "Server already started"),
		}
	}

	err = conn.AddMatchSignal(
		dbus.WithMatchObjectPath(dbus.ObjectPath(object_path)),
		dbus.WithMatchInterface(interface_name),
		dbus.WithMatchMember(signal_name),
	)
	defer conn.RemoveMatchSignal(
		dbus.WithMatchObjectPath(dbus.ObjectPath(object_path)),
		dbus.WithMatchInterface(interface_name),
		dbus.WithMatchMember(signal_name),
	)

	if err != nil {
		return conn, DetailedErrors{
			Type:        ErrorFailedConnection,
			Description: fmt.Sprintf("%s - %s", "Failed to start CLI message listener", err.Error()),
		}
	}

	s.messages = messageChannel

	err = conn.Export(s, dbus.ObjectPath(object_path), interface_name)
	if err != nil {
		return conn, DetailedErrors{
			Type:        ErrorFailedConnection,
			Description: fmt.Sprintf("%s - %s", "Failed to start CLI message listener", err.Error()),
		}
	}

	return conn, DetailedErrors{
		Type:        ErrorNil,
		Description: "",
	}
}

// Called internally to store a message send via CLI sender
func (s *Listener) SendMessage(message string) (reply string, err *dbus.Error) {
	s.messages <- message
	return "Acknoledge", nil
}
