package commands

import (
	"fmt"

	"github.com/tsinivuo/redis-lite/pkg/resp"
)

// EchoCommand implements the ECHO command
type EchoCommand struct{}

// NewEchoCommand creates a new ECHO command
func NewEchoCommand() *EchoCommand {
	return &EchoCommand{}
}

// Name returns the command name
func (c *EchoCommand) Name() string {
	return "ECHO"
}

// Validate checks if the ECHO command arguments are valid
func (c *EchoCommand) Validate(args []*resp.Message) error {
	// ECHO requires exactly 1 argument
	if len(args) != 1 {
		return fmt.Errorf("wrong number of arguments for 'echo' command")
	}
	return nil
}

// Execute processes the ECHO command
func (c *EchoCommand) Execute(args []*resp.Message) (*resp.Message, error) {
	arg := args[0]

	switch arg.Type {
	case resp.BulkString:
		if arg.Value == nil {
			return resp.NewNullBulkString(), nil
		}
		return resp.NewBulkString(arg.Value.(string)), nil
	case resp.SimpleString:
		return resp.NewBulkString(arg.Value.(string)), nil
	case resp.Integer:
		return resp.NewBulkString(fmt.Sprintf("%d", arg.Value.(int64))), nil
	default:
		return resp.NewError("ERR invalid argument type for ECHO"), nil
	}
}
