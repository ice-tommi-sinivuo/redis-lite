package commands

import (
	"github.com/tsinivuo/redis-lite/pkg/resp"
	"github.com/tsinivuo/redis-lite/pkg/storage"
)

// Command defines the interface that all Redis commands must implement
type Command interface {
	// Name returns the command name (e.g., "PING", "ECHO")
	Name() string

	// Execute processes the command with given arguments and returns a RESP message
	Execute(args []*resp.Message, store storage.Store) (*resp.Message, error)

	// Validate checks if the command arguments are valid
	Validate(args []*resp.Message) error
}

// CommandHandler manages command registration and execution
type CommandHandler struct {
	commands map[string]Command
}

// NewCommandHandler creates a new command handler
func NewCommandHandler() *CommandHandler {
	return &CommandHandler{
		commands: make(map[string]Command),
	}
}

// Register registers a command with the handler
func (h *CommandHandler) Register(command Command) {
	h.commands[command.Name()] = command
}

// Execute executes a command by name with the given arguments
func (h *CommandHandler) Execute(commandName string, args []*resp.Message, store storage.Store) (*resp.Message, error) {
	command, exists := h.commands[commandName]
	if !exists {
		return resp.NewError("ERR unknown command '" + commandName + "'"), nil
	}

	if err := command.Validate(args); err != nil {
		return resp.NewError("ERR " + err.Error()), nil
	}

	return command.Execute(args, store)
}

// GetCommand returns a command by name
func (h *CommandHandler) GetCommand(name string) (Command, bool) {
	cmd, exists := h.commands[name]
	return cmd, exists
}
