package pkg

import "fmt"

// ErrTimeout is returned when a component times out during initialization.
type ErrTimeout struct {
	ComponentName string
}

func (e ErrTimeout) Error() string {
	if e.ComponentName == "" {
		return "process timed out"
	}

	return fmt.Sprintf("component %s timed out", e.ComponentName)
}
