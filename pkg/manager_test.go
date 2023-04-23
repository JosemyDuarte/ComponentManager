package pkg_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/JosemyDuarte/ComponentManager/pkg"
)

// Test the start and shutdown of the manager on a happy path
func TestComponentManager_StartAndShutdown(t *testing.T) {
	fakeComponents := []*fakeComponent{
		{
			name:             "health-check",
			startDuration:    5 * time.Millisecond,
			startTimeout:     10 * time.Millisecond,
			shutdownDuration: 5 * time.Millisecond,
		},
		{
			name:             "grpc-server",
			startDuration:    5 * time.Millisecond,
			startTimeout:     10 * time.Millisecond,
			shutdownDuration: 5 * time.Millisecond,
		},
		{
			name:             "kafka-consumer",
			startDuration:    5 * time.Millisecond,
			startTimeout:     10 * time.Millisecond,
			shutdownDuration: 5 * time.Millisecond,
		},
	}

	// Create a new manager
	m := pkg.NewManager()

	for _, fc := range fakeComponents {
		m.Register(fc)
	}

	// Start the manager
	errCh := m.Start()
	select {
	case err := <-errCh:
		require.NoError(t, err, "failed to start the manager")
	default:
		// No error, continue
	}

	// Check that all the components are started
	for _, fc := range fakeComponents {
		assert.True(t, fc.IsStarted(), "component %s is not started", fc.Name())
	}

	// Shutdown the manager
	require.NoError(t, m.Shutdown(context.Background(), time.Second), "failed to shutdown the manager")

	// Check that all the components are shutdown
	for _, fc := range fakeComponents {
		assert.True(t, fc.IsShutdown(), "component %s is not shutdown", fc.Name())
	}
}

// Test scenario where a component fails to start due to a timeout
func TestComponentManager_StartTimeout(t *testing.T) {
	fakeComponents := []*fakeComponent{
		{
			name:          "health-check",
			startDuration: 5 * time.Millisecond,
			startTimeout:  10 * time.Millisecond,
		},
		{
			name:          "grpc-server",
			startDuration: 5 * time.Millisecond,
			startTimeout:  3 * time.Millisecond, // This component will timeout
		},
		{
			name:          "kafka-consumer",
			startDuration: 5 * time.Millisecond,
			startTimeout:  10 * time.Millisecond,
		},
	}

	// Create a new manager
	m := pkg.NewManager()

	for _, fc := range fakeComponents {
		m.Register(fc)
	}

	// Start the manager
	errCh := m.Start()

	// Check that the manager failed to start
	select {
	case err := <-errCh:
		require.Error(t, err, "manager should have failed to start")
		require.IsType(t, pkg.ErrStartTimeout{}, err, "error should be a timeout error")
	default:
		t.Fatal("manager should have failed to start")
	}
}

func TestComponentManager_StartErr(t *testing.T) {
	fakeComponents := []*fakeComponent{
		{
			name:          "health-check",
			startDuration: 5 * time.Millisecond,
			startTimeout:  10 * time.Millisecond,
		},
		{
			name:          "grpc-server",
			startDuration: 5 * time.Millisecond,
			startTimeout:  10 * time.Millisecond,
			startError:    assert.AnError, // This component will fail to start
		},
		{
			name:          "kafka-consumer",
			startDuration: 5 * time.Millisecond,
			startTimeout:  10 * time.Millisecond,
		},
	}

	// Create a new manager
	m := pkg.NewManager()

	for _, fc := range fakeComponents {
		m.Register(fc)
	}

	// Start the manager
	errCh := m.Start()

	// Check that the manager failed to start
	select {
	case err := <-errCh:
		require.Error(t, err, "manager should have failed to start")
	default:
		t.Fatal("manager should have failed to start")
	}

	// Check that grpc-server is not started
	for _, fc := range fakeComponents {
		if fc.Name() == "grpc-server" {
			assert.False(t, fc.IsStarted(), "component %s should not be started", fc.Name())
		}
	}
}

// Test when a component returns an error after it has signaled that it has started
func TestComponentManager_StartErrAfterStarted(t *testing.T) {
	// Create a new manager
	m := pkg.NewManager()

	m.Register(&fakeComponent{
		name:          "health-check",
		startDuration: 5 * time.Millisecond,
		startTimeout:  10 * time.Millisecond,
	})
	m.Register(&fakeComponent{
		name:          "kafka-consumer",
		startDuration: 5 * time.Millisecond,
		startTimeout:  10 * time.Millisecond,
	})
	m.Register(&fakeComponentWithErrOnStart{
		fakeComponent: fakeComponent{
			name:          "grpc-server",
			startDuration: 5 * time.Millisecond,
			startTimeout:  10 * time.Millisecond,
			startError:    assert.AnError, // This component will fail to start after it has signaled that it has started
		},
	})

	// Start the manager
	errCh := m.Start()

	select {
	case err := <-errCh:
		require.NoError(t, err, "failed to start the manager")
	default:
		// No error, continue
	}

	// Wait for components to fail after they signaled that they have started
	time.Sleep(20 * time.Millisecond)

	// Check that the manager error channel contains an error
	select {
	case err := <-errCh:
		require.Error(t, err, "manager should have failed to start")
	default:
		t.Fatal("manager start should have failed")
	}
}

func TestComponentManager_ShutdownErr(t *testing.T) {
	fakeComponents := []*fakeComponent{
		{
			name:             "health-check",
			startDuration:    5 * time.Millisecond,
			startTimeout:     10 * time.Millisecond,
			shutdownDuration: 5 * time.Millisecond,
		},
		{
			name:             "grpc-server",
			startDuration:    5 * time.Millisecond,
			startTimeout:     10 * time.Millisecond,
			shutdownDuration: 5 * time.Millisecond,
			shutdownErr:      assert.AnError, // This component will fail to shutdown
		},
		{
			name:             "kafka-consumer",
			startDuration:    5 * time.Millisecond,
			startTimeout:     10 * time.Millisecond,
			shutdownDuration: 5 * time.Millisecond,
		},
	}

	// Create a new manager
	m := pkg.NewManager()

	for _, fc := range fakeComponents {
		m.Register(fc)
	}

	// Start the manager
	errCh := m.Start()

	// Check that the manager started
	select {
	case err := <-errCh:
		require.NoError(t, err, "manager should have started")
	default:
		// No error, continue
	}

	// Shutdown the manager
	err := m.Shutdown(context.Background(), time.Second)

	// Check that the manager failed to shutdown
	require.Error(t, err, "manager should have failed to shutdown")
}

func TestComponentManager_ShutdownGracePeriod(t *testing.T) {
	// The sum of the shutdown duration of all the components is greater than the grace period
	fakeComponents := []*fakeComponent{
		{
			name:             "health-check",
			startDuration:    5 * time.Millisecond,
			startTimeout:     10 * time.Millisecond,
			shutdownDuration: 3 * time.Millisecond,
		},
		{
			name:             "grpc-server",
			startDuration:    5 * time.Millisecond,
			startTimeout:     10 * time.Millisecond,
			shutdownDuration: 2 * time.Millisecond,
		},
		{
			name:             "kafka-consumer",
			startDuration:    5 * time.Millisecond,
			startTimeout:     10 * time.Millisecond,
			shutdownDuration: 2 * time.Millisecond,
		},
	}

	// Create a new manager
	m := pkg.NewManager()

	for _, fc := range fakeComponents {
		m.Register(fc)
	}

	// Start the manager
	errCh := m.Start()

	// Check that the manager started
	select {
	case err := <-errCh:
		require.NoError(t, err, "manager should have started")
	default:
		// No error, continue
	}

	// Shutdown the manager
	err := m.Shutdown(context.Background(), 5*time.Millisecond)

	// Check that the manager failed to shutdown
	require.Error(t, err, "manager should have failed to shutdown")

	// Check that the error is a timeout error
	require.IsType(t, pkg.ErrShutdownTimeout{}, err, "error should be a timeout error")
}

// Define a fake component that takes some time to start and shutdown
type fakeComponent struct {
	name             string
	startDuration    time.Duration
	startTimeout     time.Duration
	isStarted        bool
	startError       error
	shutdownDuration time.Duration
	isShutdown       bool
	shutdownErr      error
}

func (c *fakeComponent) Name() string {
	return c.name
}

func (c *fakeComponent) Start(initReady chan struct{}) error {
	// Simulate an initialization that takes some time
	time.Sleep(c.startDuration)
	if c.startError != nil {
		return c.startError
	}

	c.isStarted = true
	close(initReady)
	return nil
}

func (c *fakeComponent) IsStarted() bool {
	return c.isStarted
}

func (c *fakeComponent) Shutdown(_ context.Context) error {
	// Simulate a shutdown that takes some time
	time.Sleep(c.shutdownDuration)
	if c.shutdownErr != nil {
		return c.shutdownErr
	}

	c.isShutdown = true
	return nil
}

func (c *fakeComponent) IsShutdown() bool {
	return c.isShutdown
}

func (c *fakeComponent) StartTimeout() time.Duration {
	return c.startTimeout
}

// Define a fake component that fails to start after it has signaled that it is ready
type fakeComponentWithErrOnStart struct {
	fakeComponent
}

func (c *fakeComponentWithErrOnStart) Start(initReady chan struct{}) error {
	c.isStarted = true
	close(initReady)
	time.Sleep(c.startDuration)

	return c.startError
}
