package pkg

import (
	"context"
	"time"
)

// Component is an interface that all components must implement.
type Component interface {
	// Name returns the identifier of the component (mainly for logging purposes)
	Name() string

	// Start is called in order to start the component.
	// The initReady channel should be closed when the component's initialisation is started
	// and the manager can move onto the next component.
	Start(initReady chan struct{}) error

	// Shutdown is called to stop and cleanup the component.
	// There is a limited time to do so which is controlled by the ctx.
	Shutdown(ctx context.Context) error

	// StartTimeout returns timeout within which the component is expected to be start
	StartTimeout() time.Duration
}
