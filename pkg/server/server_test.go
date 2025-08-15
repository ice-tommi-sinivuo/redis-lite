package server

import (
	"testing"
)

func TestNewServer(t *testing.T) {
	address := "127.0.0.1"
	port := 6379

	server := NewServer(address, port)

	if server == nil {
		t.Fatal("NewServer returned nil")
	}

	if server.address != address {
		t.Errorf("Expected address '%s', got '%s'", address, server.address)
	}

	if server.port != port {
		t.Errorf("Expected port %d, got %d", port, server.port)
	}

	if server.commandHandler == nil {
		t.Error("Command handler is nil")
	}

	if server.connections == nil {
		t.Error("Connections map is nil")
	}

	if server.shutdown == nil {
		t.Error("Shutdown channel is nil")
	}
}

func TestServer_CommandsRegistered(t *testing.T) {
	server := NewServer("127.0.0.1", 6379)

	// Check that PING command is registered
	cmd, exists := server.commandHandler.GetCommand("PING")
	if !exists {
		t.Error("PING command not registered")
	}

	if cmd.Name() != "PING" {
		t.Errorf("Expected PING command, got %s", cmd.Name())
	}

	// Check that ECHO command is registered
	cmd, exists = server.commandHandler.GetCommand("ECHO")
	if !exists {
		t.Error("ECHO command not registered")
	}

	if cmd.Name() != "ECHO" {
		t.Errorf("Expected ECHO command, got %s", cmd.Name())
	}
}

func TestServer_GetCommandHandler(t *testing.T) {
	server := NewServer("127.0.0.1", 6379)

	handler := server.GetCommandHandler()
	if handler == nil {
		t.Error("GetCommandHandler returned nil")
	}

	if handler != server.commandHandler {
		t.Error("GetCommandHandler returned different handler than expected")
	}
}

func TestServer_Stop_WhenNotRunning(t *testing.T) {
	server := NewServer("127.0.0.1", 6379)

	// Server is not running, Stop should not return error
	err := server.Stop()
	if err != nil {
		t.Errorf("Stop() returned error when server not running: %v", err)
	}
}
