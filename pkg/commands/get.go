package commands

import (
	"fmt"

	"github.com/tsinivuo/redis-lite/pkg/resp"
	"github.com/tsinivuo/redis-lite/pkg/storage"
)

// GetCommand implements the GET command
type GetCommand struct{}

// NewGetCommand creates a new GET command
func NewGetCommand() *GetCommand {
	return &GetCommand{}
}

// Name returns the command name
func (c *GetCommand) Name() string {
	return "GET"
}

// Validate checks if the GET command arguments are valid
func (c *GetCommand) Validate(args []*resp.Message) error {
	// GET requires exactly 1 argument: key
	if len(args) != 1 {
		return fmt.Errorf("wrong number of arguments for 'get' command")
	}
	return nil
}

// Execute processes the GET command
func (c *GetCommand) Execute(args []*resp.Message, store storage.Store) (*resp.Message, error) {
	// Extract key from arguments
	keyArg := args[0]

	// Convert key to string
	var key string
	switch keyArg.Type {
	case resp.BulkString:
		if keyArg.Value == nil {
			return resp.NewError("ERR key cannot be null"), nil
		}
		key = keyArg.Value.(string)
	case resp.SimpleString:
		key = keyArg.Value.(string)
	default:
		return resp.NewError("ERR invalid key type"), nil
	}

	// Retrieve the value from storage
	value, exists := store.Get(key)
	if !exists {
		// Return null bulk string for non-existent keys
		return resp.NewNullBulkString(), nil
	}

	// Return the value as a bulk string
	return resp.NewBulkString(value), nil
}
