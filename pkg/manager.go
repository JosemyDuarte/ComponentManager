package pkg

import (
	"context"
	"fmt"
	"log"
	"time"
)

// Manager is responsible for managing all the component's lifecycle.
type Manager struct {
	components []Component
}

// NewManager returns a new Manager instance.
func NewManager() *Manager {
	return &Manager{}
}

// Register registers a new component.
func (m *Manager) Register(component Component) {
	m.components = append(m.components, component)
}

// Start starts all the components and returns an error channel.
//
// The returned error channel receives all the errors from the components, and it is the caller's responsibility
// to check for errors. Since there can be errors happening even after the components have signaled that
// they are ready.
//
// It publishes an error to the channel if any of the components fail to start.
// The components are started in the order they were registered.
// It won't continue to the next component until the previous one has finished its initialization.
// If a component fails to start, the manager will stop starting the remaining components and publish an error.
func (m *Manager) Start() chan error {
	log.Printf("Starting %d components...\n", len(m.components))
	errCh := make(chan error, len(m.components))

	start := time.Now()

	// Start all the components in order.
	for _, c := range m.components {
		log.Printf("Starting component %s...\n", c.Name())

		// Create a channel to signal that the component's initialization is done.
		initReady := make(chan struct{})

		// Start the component.
		go func() {
			errCh <- c.Start(initReady)
		}()

		// Wait for the component's initialization to complete.
		select {
		case <-initReady:
			log.Printf("Component %s initialization completed\n", c.Name())
		case <-time.After(c.StartTimeout()):
			errCh <- StartTimeoutError{ComponentName: c.Name()}

			return errCh
		}

		// Check for errors during startup.
		select {
		case err := <-errCh:
			if err != nil {
				errCh <- fmt.Errorf("failed to start component %s: %w", c.Name(), err)

				return errCh
			}
		default:
			log.Printf("Component %s started successfully\n", c.Name())
		}
	}

	log.Printf("All components started in %v", time.Since(start))

	return errCh
}

// Shutdown shuts down all the components in reverse order it was registered.
// It returns an error if overall shutdown takes longer than the specified grace period.
// It tries to shutdown all the components even if some of them fail to shutdown.
// It won't continue to the next component until the previous one has finished its shutdown.
func (m *Manager) Shutdown(ctx context.Context, gracePeriod time.Duration) error {
	log.Printf("Shutting down %d components...\n", len(m.components))

	start := time.Now()

	var errs []error

	shutdownDone := make(chan struct{})
	go func() {
		// Start shutting down all the components in reverse order.
		for i := len(m.components) - 1; i >= 0; i-- {
			c := m.components[i]
			log.Printf("Shutting down component %s...\n", c.Name())

			// Check for errors during shutdown.
			if err := c.Shutdown(ctx); err != nil {
				errs = append(errs, fmt.Errorf("failed to shutdown component %s: %w", c.Name(), err))
			} else {
				log.Printf("Component %s shutdown successfully\n", c.Name())
			}
		}

		close(shutdownDone)
	}()

	// Wait for the components shutdown to complete or timeout.
	select {
	case <-shutdownDone:
		log.Printf("Shutdown finished in %v\n", time.Since(start))
	case <-time.After(gracePeriod):
		return ShutdownTimeoutError{TimeOut: gracePeriod}
	}

	// Check if there were any errors during shutdown.
	if len(errs) > 0 {
		return ShutdownError{Errors: errs}
	}

	log.Printf("Shutdown completed successfully\n")

	return nil
}
