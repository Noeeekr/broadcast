package instance

// Instance implements general features for this program
type Instance struct {
	Shutdown *Shutdown
}

func New() Instance {
	return Instance{
		Shutdown: &Shutdown{
			shutdownCallbacks: []ShutdownFunc{},
		},
	}
}
