package ipc

import (
	"github.com/godbus/dbus/v5"
)

type Sender struct{}

// Starts an Inter Process Communication (IPC) service to send a message to a running server
func (s *Sender) Send(message string) DetailedErrors {
	Conn, err := dbus.ConnectSessionBus()
	if err != nil {
		return DetailedErrors{
			Type:        ErrorFailedConnection,
			Description: err.Error(),
		}
	}
	defer Conn.Close()

	call := Conn.Object(name, dbus.ObjectPath(object_path)).Call(method, 0, message)
	if call.Err != nil {
		return DetailedErrors{
			Type:        ErrorFailedConnection,
			Description: call.Err.Error(),
		}
	}

	var r1 string
	err = call.Store(&r1)
	if err != nil {
		return DetailedErrors{
			Type:        ErrorFailedConnection,
			Description: err.Error(),
		}
	}

	return DetailedErrors{
		Type:        ErrorNil,
		Description: "",
	}
}
