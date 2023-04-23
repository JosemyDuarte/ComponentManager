package pkg

import (
	"fmt"
	"time"
)

// ErrStartTimeout is returned when a component times out during initialization.
type ErrStartTimeout struct {
	ComponentName string
}

func (e ErrStartTimeout) Error() string {
	return fmt.Sprintf("component %s timed out", e.ComponentName)
}

// ErrShutdownTimeout is returned when shutdown times out
type ErrShutdownTimeout struct {
	TimeOut time.Duration
}

func (e ErrShutdownTimeout) Error() string {
	return fmt.Sprintf("shutdown timed out after %s", e.TimeOut)
}
