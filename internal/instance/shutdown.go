package instance

import (
	"fmt"
	"os"
	"os/signal"
	"sync"

	Panic "github.com/Noeeekr/broadcast_server/pkg/panic"
)

var once sync.Once

// A callback to finish the task correctly
type ShutdownFunc func()

// Shutdown holds methods that helps finishing program smoothlier
type Shutdown struct {
	shutdownCallbacks []ShutdownFunc
}

func (s *Shutdown) onInterrupt(cb func()) {
	once.Do(
		func() {
			var signals chan os.Signal = make(chan os.Signal, 1)
			signal.Notify(signals, os.Interrupt, os.Kill)

			go func() {
				sig := <-signals
				fmt.Println("Finishing program - Signal detected", sig.String())

				cb()
				os.Exit(0)
			}()
		},
	)
}

// Enables gracefull shutdown
func (s *Shutdown) EnableGracefull() {
	var p Panic.Panic
	p.NotImplemented("Not implemented - mutex in callbacks for adding without shutting down program before")
	// quando a função é executada não tem porque lockar
	// mas quando ela é adiciona é importante lockar essa parte pra adicionar a nova

	s.onInterrupt(func() {
		for _, shutdownFunc := range s.shutdownCallbacks {
			shutdownFunc()
		}
	})
}

func (s *Shutdown) AddCallback(shutdownCallback func()) {
	s.shutdownCallbacks = append(s.shutdownCallbacks, shutdownCallback)
}
