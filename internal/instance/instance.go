package instance

import (
	"github.com/Noeeekr/broadcast_server/pkg/instance"
)

// Instance implements general features for this program
type Instance struct {
	*instance.Shutdown
	*instance.Debugger
	*instance.Logger
}

func New() Instance {
	return Instance{
		Logger:   instance.NewLogger(),
		Shutdown: instance.NewShutdown(),
		Debugger: instance.NewDebugger(),
	}
}
