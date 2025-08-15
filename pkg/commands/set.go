package commands

import (
	"fmt"

	"github.com/tsinivuo/redis-lite/pkg/resp"
	"github.com/tsinivuo/redis-lite/pkg/storage"
)

// SetCommand implements the SET command
type SetCommand struct{}

// NewSetCommand creates a new SET command
func NewSetCommand() *SetCommand {
	return &SetCommand{}
}

// Name returns the command name
func (c *SetCommand) Name() string {
	return "SET"
}

// Validate checks if the SET command arguments are valid
func (c *SetCommand) Validate(args []*resp.Message) error {
	// SET requires exactly 2 arguments: key and value
	if len(args) != 2 {
		return fmt.Errorf("wrong number of arguments for 'set' command")
	}
	return nil
}

// Execute processes the SET command
func (c *SetCommand) Execute(args []*resp.Message, store storage.Store) (*resp.Message, error) {
	// Extract key and value from arguments
	keyArg := args[0]
	valueArg := args[1]

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

	// Convert value to string
	var value string
	switch valueArg.Type {
	case resp.BulkString:
		if valueArg.Value == nil {
			value = ""
		} else {
			value = valueArg.Value.(string)
		}
	case resp.SimpleString:
		value = valueArg.Value.(string)
	case resp.Integer:
		value = fmt.Sprintf("%d", valueArg.Value.(int64))
	default:
		return resp.NewError("ERR invalid value type"), nil
	}

	// Store the key-value pair
	if err := store.Set(key, value); err != nil {
		return resp.NewError("ERR " + err.Error()), nil
	}

	// Return OK response
	return resp.NewSimpleString("OK"), nil
}
