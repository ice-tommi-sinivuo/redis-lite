package server

import (
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/tsinivuo/redis-lite/pkg/commands"
)

// Server represents the Redis-Lite TCP server
type Server struct {
	address        string
	port           int
	listener       net.Listener
	commandHandler *commands.CommandHandler
	connections    map[net.Conn]*Connection
	mutex          sync.RWMutex
	shutdown       chan struct{}
	running        bool
}

// NewServer creates a new Redis-Lite server
func NewServer(address string, port int) *Server {
	server := &Server{
		address:        address,
		port:           port,
		commandHandler: commands.NewCommandHandler(),
		connections:    make(map[net.Conn]*Connection),
		shutdown:       make(chan struct{}),
	}

	// Register built-in commands
	server.commandHandler.Register(commands.NewPingCommand())
	server.commandHandler.Register(commands.NewEchoCommand())

	return server
}

// Start starts the server and begins accepting connections
func (s *Server) Start() error {
	addr := fmt.Sprintf("%s:%d", s.address, s.port)
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("failed to listen on %s: %w", addr, err)
	}

	s.listener = listener
	s.running = true

	log.Printf("Redis-Lite server started on %s", addr)

	// Accept connections in a goroutine
	go s.acceptConnections()

	// Wait for shutdown signal
	<-s.shutdown

	return nil
}

// Stop gracefully stops the server
func (s *Server) Stop() error {
	if !s.running {
		return nil
	}

	s.running = false

	// Close all active connections
	s.mutex.Lock()
	for conn := range s.connections {
		conn.Close()
	}
	s.mutex.Unlock()

	// Close the listener
	if s.listener != nil {
		if err := s.listener.Close(); err != nil {
			return fmt.Errorf("failed to close listener: %w", err)
		}
	}

	// Signal shutdown
	close(s.shutdown)

	log.Println("Redis-Lite server stopped")
	return nil
}

// acceptConnections accepts incoming connections and handles them
func (s *Server) acceptConnections() {
	for s.running {
		conn, err := s.listener.Accept()
		if err != nil {
			if s.running {
				log.Printf("Error accepting connection: %v", err)
			}
			continue
		}

		// Create connection handler
		connection := NewConnection(conn, s.commandHandler)

		// Track the connection
		s.mutex.Lock()
		s.connections[conn] = connection
		s.mutex.Unlock()

		// Handle the connection in a separate goroutine
		go func() {
			defer func() {
				// Remove connection from tracking
				s.mutex.Lock()
				delete(s.connections, conn)
				s.mutex.Unlock()

				// Close the connection
				conn.Close()
			}()

			connection.Handle()
		}()
	}
}

// GetCommandHandler returns the command handler for testing purposes
func (s *Server) GetCommandHandler() *commands.CommandHandler {
	return s.commandHandler
}
