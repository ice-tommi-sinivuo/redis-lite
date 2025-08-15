package server

import (
	"log"
	"net"
	"strings"

	"github.com/tsinivuo/redis-lite/pkg/commands"
	"github.com/tsinivuo/redis-lite/pkg/resp"
)

// Connection represents a client connection to the server
type Connection struct {
	conn           net.Conn
	parser         *resp.Parser
	serializer     *resp.Serializer
	commandHandler *commands.CommandHandler
}

// NewConnection creates a new connection handler
func NewConnection(conn net.Conn, commandHandler *commands.CommandHandler) *Connection {
	return &Connection{
		conn:           conn,
		parser:         resp.NewParser(conn),
		serializer:     resp.NewSerializer(conn),
		commandHandler: commandHandler,
	}
}

// Handle processes incoming messages from the client
func (c *Connection) Handle() {
	log.Printf("New client connected: %s", c.conn.RemoteAddr())
	defer log.Printf("Client disconnected: %s", c.conn.RemoteAddr())

	for {
		// Parse incoming RESP message
		message, err := c.parser.Parse()
		if err != nil {
			log.Printf("Error parsing message: %v", err)
			return
		}

		// Process the command
		response := c.processCommand(message)

		// Send response back to client
		if err := c.serializer.Serialize(response); err != nil {
			log.Printf("Error serializing response: %v", err)
			return
		}
	}
}

// processCommand processes a command message and returns a response
func (c *Connection) processCommand(message *resp.Message) *resp.Message {
	// Commands should be arrays in RESP protocol
	if message.Type != resp.Array {
		return resp.NewError("ERR Protocol error: expected array")
	}

	// Handle null array
	if message.Value == nil {
		return resp.NewError("ERR Protocol error: null array")
	}

	args := message.Value.([]*resp.Message)

	// Commands need at least one element (the command name)
	if len(args) == 0 {
		return resp.NewError("ERR Protocol error: empty array")
	}

	// First element should be the command name
	commandNameMsg := args[0]
	if commandNameMsg.Type != resp.BulkString && commandNameMsg.Type != resp.SimpleString {
		return resp.NewError("ERR Protocol error: command name must be a string")
	}

	var commandName string
	switch commandNameMsg.Type {
	case resp.BulkString:
		if commandNameMsg.Value == nil {
			return resp.NewError("ERR Protocol error: null command name")
		}
		commandName = commandNameMsg.Value.(string)
	case resp.SimpleString:
		commandName = commandNameMsg.Value.(string)
	}

	// Convert command name to uppercase (Redis commands are case-insensitive)
	commandName = strings.ToUpper(commandName)

	// Get command arguments (everything after the command name)
	commandArgs := args[1:]

	// Execute the command
	response, err := c.commandHandler.Execute(commandName, commandArgs)
	if err != nil {
		return resp.NewError("ERR " + err.Error())
	}

	return response
}
