package instance

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var once sync.Once
var InterruptContext context.Context
var interrupt context.CancelFunc
var waiter sync.WaitGroup
var isTerminating bool

func init() {
	InterruptContext, interrupt = signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
}

type Shutdown struct{}

func NewShutdown() *Shutdown {
	return &Shutdown{}
}

// Tells if the program is trying to terminate at the moment
func (s *Shutdown) IsTerminating() bool {
	return isTerminating
}

// Enables gracefull shutdown on SIGINT and SIGTERM
func (s *Shutdown) EnableGracefull() {
	var once sync.Once

	once.Do(func() {
		var signals chan os.Signal = make(chan os.Signal, 1)
		signal.Notify(signals, syscall.SIGINT, syscall.SIGTERM)

		go func() {
			sig := <-signals
			fmt.Println("Signal detected (" + sig.String() + ") Starting process shutdown.")
			s.Terminate()
		}()
	})
}

// Terminate sets IsTerminating return to true when called to tell new jobs to not proceed.
// After that, it for all Wait() to be eliminated with a Proceed() then finishes the program entirely.
// Needs logs
func (s *Shutdown) Terminate() {
	once.Do(func() {
		go func() {
			interrupt()
			isTerminating = true
			waiter.Wait()

			fmt.Println("All jobs finished. Shutting down.")
			os.Exit(0)
		}()
	})
}

// Wait stops Terminate() from finishing the process until Proceed() is called
func (s *Shutdown) Wait() (Proceed func()) {
	waiter.Add(1)

	return func() {
		waiter.Done()
	}
}
