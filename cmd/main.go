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
)

func main() {
	manager := internal.NewManager()

	service := internal.NewService()
	server := internal.NewServer(service, "127.0.0.1:")

	manager.Register(server)

	errCh := manager.Start()

	// create a channel to receive shutdown signals
	shutdownSignal := make(chan os.Signal, 1)
	signal.Notify(shutdownSignal, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// wait for an error or a shutdown signal
	select {
	case err := <-errCh:
		panic(fmt.Errorf("failed to start the manager: %w", err))
	case sig := <-shutdownSignal:
		log.Printf("received signal %s, shutting down", sig)
	}

	if err := manager.Shutdown(context.Background(), time.Minute); err != nil {
		panic(fmt.Errorf("failed to shutdown the manager: %w", err))
	}
}
