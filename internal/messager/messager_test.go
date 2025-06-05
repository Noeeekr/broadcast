package messager_test

import (
	"testing"

	"github.com/Noeeekr/broadcast_server/internal/messager"
	"github.com/godbus/dbus/v5"
)

func TestPingMessageListener(t *testing.T) {
	messagerInstance := messager.NewMessager()

	tests := []struct {
		Name              string
		ShouldStartServer bool
		ShouldPingSuccess bool
	}{
		{"Ping should reply success", false, false},
		{"Ping shouldn't reply success", true, true},
	}

	for _, test := range tests {
		var conn *dbus.Conn
		var detailedError messager.DetailedErrors

		if test.ShouldStartServer {
			conn, detailedError = startMessageListener(messagerInstance)
			if detailedError.Type != messager.ErrorNil {
				t.Fatal("Unable to execute test: couldn't start message listener -", detailedError.Description)
			}
		}

		t.Run(test.Name, func(t *testing.T) {
			success, err := messagerInstance.PingMessageListener()
			if err != nil {
				t.Fatal("Test failed: Error in PingMessageListener -", err.Error())
			}

			if success != test.ShouldPingSuccess {
				t.Fatal("Test failed: Ping returned", success, "when it should've returned", test.ShouldPingSuccess)
			}
		})

		if test.ShouldStartServer {
			conn.Close()
		}
	}
}

func startMessageListener(messagerInstance messager.IpcMessager) (*dbus.Conn, messager.DetailedErrors) {
	var cha chan []byte = make(chan []byte, 0)
	conn, detailedError := messagerInstance.Listener.StartMessageListener(cha)
	if detailedError.Type != messager.ErrorNil {
		return conn, detailedError
	}

	return conn, detailedError
}
