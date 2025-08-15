package commands

import (
	"testing"

	"github.com/tsinivuo/redis-lite/pkg/resp"
)

// mockCommand is a test implementation of the Command interface
type mockCommand struct {
	name string
}

func (m *mockCommand) Name() string {
	return m.name
}

func (m *mockCommand) Execute(args []*resp.Message) (*resp.Message, error) {
	return resp.NewSimpleString("OK"), nil
}

func (m *mockCommand) Validate(args []*resp.Message) error {
	return nil
}

func TestNewCommandHandler(t *testing.T) {
	handler := NewCommandHandler()
	if handler == nil {
		t.Fatal("NewCommandHandler returned nil")
	}

	if handler.commands == nil {
		t.Fatal("CommandHandler.commands map is nil")
	}
}

func TestCommandHandler_Register(t *testing.T) {
	handler := NewCommandHandler()
	command := &mockCommand{name: "TEST"}

	handler.Register(command)

	// Check if command was registered
	cmd, exists := handler.GetCommand("TEST")
	if !exists {
		t.Fatal("Command was not registered")
	}

	if cmd.Name() != "TEST" {
		t.Errorf("Expected command name 'TEST', got '%s'", cmd.Name())
	}
}

func TestCommandHandler_Execute(t *testing.T) {
	handler := NewCommandHandler()
	command := &mockCommand{name: "TEST"}
	handler.Register(command)

	// Test successful execution
	response, err := handler.Execute("TEST", []*resp.Message{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.Type != resp.SimpleString {
		t.Errorf("Expected SimpleString response, got %s", response.Type)
	}

	if response.Value.(string) != "OK" {
		t.Errorf("Expected 'OK' response, got '%s'", response.Value.(string))
	}
}

func TestCommandHandler_Execute_UnknownCommand(t *testing.T) {
	handler := NewCommandHandler()

	// Test unknown command
	response, err := handler.Execute("UNKNOWN", []*resp.Message{})
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.Type != resp.Error {
		t.Errorf("Expected Error response, got %s", response.Type)
	}

	expectedMsg := "ERR unknown command 'UNKNOWN'"
	if response.Value.(string) != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, response.Value.(string))
	}
}

func TestCommandHandler_GetCommand(t *testing.T) {
	handler := NewCommandHandler()
	command := &mockCommand{name: "TEST"}
	handler.Register(command)

	// Test existing command
	cmd, exists := handler.GetCommand("TEST")
	if !exists {
		t.Fatal("Expected command to exist")
	}

	if cmd.Name() != "TEST" {
		t.Errorf("Expected command name 'TEST', got '%s'", cmd.Name())
	}

	// Test non-existing command
	_, exists = handler.GetCommand("NONEXISTENT")
	if exists {
		t.Fatal("Expected command to not exist")
	}
}
