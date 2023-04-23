package pkg

import (
	"fmt"
	"time"
)

// StartTimeoutError is returned when a component times out during initialization.
type StartTimeoutError struct {
	ComponentName string
}

func (e StartTimeoutError) Error() string {
	return fmt.Sprintf("component %s timed out", e.ComponentName)
}

// ShutdownTimeoutError is returned when shutdown times out.
type ShutdownTimeoutError struct {
	TimeOut time.Duration
}

func (e ShutdownTimeoutError) Error() string {
	return fmt.Sprintf("shutdown timed out after %s", e.TimeOut)
}

// ShutdownError is returned when shutdown fails.
// It contains all the errors returned by the components.
type ShutdownError struct {
	Errors []error
}

func (e ShutdownError) Error() string {
	return fmt.Sprintf("failed to shutdown: %v", e.Errors)
}
