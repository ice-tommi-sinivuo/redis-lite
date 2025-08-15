package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/tsinivuo/redis-lite/pkg/server"
)

const (
	// DefaultAddress is the default address the server listens on
	DefaultAddress = "127.0.0.1"
	// DefaultPort is the default port the server listens on (Redis standard)
	DefaultPort = 6379
)

func main() {
	// Create server
	srv := server.NewServer(DefaultAddress, DefaultPort)

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		if err := srv.Start(); err != nil {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, stopping server...")

	// Stop server gracefully
	if err := srv.Stop(); err != nil {
		log.Printf("Error stopping server: %v", err)
	}

	log.Println("Server stopped successfully")
}
