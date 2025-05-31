package instance

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var once sync.Once
var shutdownCallbacks []ShutdownFunc

// A callback to finish the task correctly
type ShutdownFunc interface {
	Shutdown()
}

// Shutdown holds methods that helps finishing program smoothlier
type Shutdown struct{}

func (s *Shutdown) onInterrupt(shutdownFunc func()) {
	once.Do(
		func() {
			var signals chan os.Signal = make(chan os.Signal, 1)
			signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

			go func() {
				sig := <-signals
				fmt.Println("Finishing program - Signal detected", sig.String())
				signal.Stop(signals)

				shutdownFunc()

				os.Exit(0)
			}()
		},
	)
}

// Enables gracefull shutdown
func (s *Shutdown) EnableGracefull() {
	var debug Debugger
	debug.NotImplemented("Not implemented - mutex in callbacks for adding without shutting down program before")
	// quando a função é executada não tem porque lockar
	// mas quando ela é adiciona é importante lockar essa parte pra adicionar a nova

	s.onInterrupt(func() {
		for _, cb := range shutdownCallbacks {
			cb.Shutdown()
		}
	})
}

func (s *Shutdown) AddCallback(shutdownCallback ShutdownFunc) {
	shutdownCallbacks = append(shutdownCallbacks, shutdownCallback)
}

// Terminate emits SIGTERM to the current process. If EnableGracefull was called before, the process will terminate gracefully.
func (s *Shutdown) Terminate() {
	process, err := os.FindProcess(os.Getpid())
	if err != nil {
		fmt.Println("Failed to identify process from inside to shutdown gracefully")
		fmt.Println("Program will continue running until recieving an external signal to shutdown")
		fmt.Println(err.Error())
		return
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		fmt.Println("Failed to shutdown gracefully from inside the program")
		fmt.Println("Program will continue running until recieving an external signal to shutdown")
		fmt.Println(err.Error())
		return
	}
}

func NewShutdown() *Shutdown {
	return &Shutdown{}
}
