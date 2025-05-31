package instance

var (
	enabled bool = true
)

type Debugger struct{}

func NewDebugger() *Debugger {
	return &Debugger{}
}

// Permits to NotImplement to panic
func (p *Debugger) EnableDebug(enable bool) {
	enabled = enable
}

func (p *Debugger) NotImplemented(message string) {
	if enabled {
		panic(message)
	}
	return
}
