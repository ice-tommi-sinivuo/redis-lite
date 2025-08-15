package commands

import (
	"fmt"

	"github.com/tsinivuo/redis-lite/pkg/resp"
)

// PingCommand implements the PING command
type PingCommand struct{}

// NewPingCommand creates a new PING command
func NewPingCommand() *PingCommand {
	return &PingCommand{}
}

// Name returns the command name
func (c *PingCommand) Name() string {
	return "PING"
}

// Validate checks if the PING command arguments are valid
func (c *PingCommand) Validate(args []*resp.Message) error {
	// PING can have 0 or 1 argument
	if len(args) > 1 {
		return fmt.Errorf("wrong number of arguments for 'ping' command")
	}
	return nil
}

// Execute processes the PING command
func (c *PingCommand) Execute(args []*resp.Message) (*resp.Message, error) {
	// If no arguments, return "PONG"
	if len(args) == 0 {
		return resp.NewSimpleString("PONG"), nil
	}

	// If one argument, echo it back
	arg := args[0]
	switch arg.Type {
	case resp.BulkString:
		if arg.Value == nil {
			return resp.NewNullBulkString(), nil
		}
		return resp.NewBulkString(arg.Value.(string)), nil
	case resp.SimpleString:
		return resp.NewSimpleString(arg.Value.(string)), nil
	default:
		return resp.NewError("ERR invalid argument type for PING"), nil
	}
}
