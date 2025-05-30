package Panic

var (
	enabled bool = true
)

type Panic struct{}

// Permits to NotImplement to pani
func (p *Panic) EnableDebug(enable bool) {
	enabled = enable
}

func (p *Panic) NotImplemented(message string) {
	if enabled {
		panic(message)
	}
	return
}
