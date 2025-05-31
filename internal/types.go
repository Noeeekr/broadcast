package Broadcast

type JSONMessage struct {
	Message string
	Sender  string
}

const (
	CommandsCloseConnection = "program:end"
)
