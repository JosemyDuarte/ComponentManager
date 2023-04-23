# Component Manager

This repository aims to be a showcase on an attempt to simplify the managing of the lifecycle of multiple background services.
It provides an easy-to-use interface to start and stop components in a particular order and handle errors gracefully.

## Usage

To use the Component Manager, you need to create an instance of the [Manager](pkg/manager.go) struct and add your components to it. For example:

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/JosemyDuarte/ComponentManager/internal"
	"github.com/JosemyDuarte/ComponentManager/pkg"
)

    func main() {
        manager := pkg.NewManager()
    
        service := internal.NewService()
        server := internal.NewServer(service, "127.0.0.1:") // Server implements the Component interface
    
        manager.Register(server)
    
        manager.Start()

		// Wait for a while to let the components run.
		time.Sleep(10 * time.Second)
    
        if err := manager.Shutdown(context.Background(), time.Minute); err != nil {
            panic(fmt.Errorf("failed to shutdown the manager: %w", err))
        }
    }
```

In the example above, we created a `Server` components and added it to the manager. 
Then we called `Start` on the manager to start the components. After waiting for 10 seconds,
we called `Shutdown` to stop the components in reverse order. 

The components will be started in the order we register them and stopped in the reverse order.

**Example**

```
manager.Register(server)
manager.Register(database)
manager.Register(anotherComponent)

// The order of starting the components will be:
// 1. Server
// 2. Database
// 3. AnotherComponent

// The order of stopping the components will be:
// 1. AnotherComponent
// 2. Database
// 3. Server
```
Check out the [tests](pkg/manager_test.go) or the [main](cmd/main.go) for more examples.

## Component Interface
To add a component to the manager, your component needs to implement the [Component](pkg/component.go) interface:

```go
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
```

## Contributing
Contributions are welcome! If you find any bugs or have suggestions for improvements, please submit a pull request.
