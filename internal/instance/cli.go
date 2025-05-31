package instance

import (
	"fmt"
	"os"
)

// ArgsToString returns a single string made of all the command line arguments
func (i *Instance) ArgsToString() string {
	var message string

	if len(os.Args) > 1 {
		for _, arg := range os.Args[1:] {
			message = fmt.Sprint(message, arg)
		}
	}

	return message
}
