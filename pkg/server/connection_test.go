package server

import (
	"bytes"
	"net"
	"testing"
	"time"

	"github.com/tsinivuo/redis-lite/pkg/commands"
	"github.com/tsinivuo/redis-lite/pkg/resp"
)

// mockConn implements net.Conn for testing
type mockConn struct {
	readBuffer  *bytes.Buffer
	writeBuffer *bytes.Buffer
}

func newMockConn(input string) *mockConn {
	return &mockConn{
		readBuffer:  bytes.NewBufferString(input),
		writeBuffer: &bytes.Buffer{},
	}
}

func (m *mockConn) Read(b []byte) (n int, err error) {
	return m.readBuffer.Read(b)
}

func (m *mockConn) Write(b []byte) (n int, err error) {
	return m.writeBuffer.Write(b)
}

func (m *mockConn) Close() error {
	return nil
}

func (m *mockConn) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 6379}
}

func (m *mockConn) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 12345}
}

func (m *mockConn) SetDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) SetWriteDeadline(t time.Time) error {
	return nil
}

func (m *mockConn) getWrittenData() string {
	return m.writeBuffer.String()
}

func TestNewConnection(t *testing.T) {
	conn := newMockConn("")
	handler := commands.NewCommandHandler()

	connection := NewConnection(conn, handler)

	if connection == nil {
		t.Fatal("NewConnection returned nil")
	}

	if connection.conn != conn {
		t.Error("Connection does not store the correct net.Conn")
	}

	if connection.commandHandler != handler {
		t.Error("Connection does not store the correct command handler")
	}
}

func TestConnection_processCommand_ValidPing(t *testing.T) {
	conn := newMockConn("")
	handler := commands.NewCommandHandler()
	handler.Register(commands.NewPingCommand())

	connection := NewConnection(conn, handler)

	// Create a PING command message: *1\r\n$4\r\nPING\r\n
	pingArray := resp.NewArray([]*resp.Message{
		resp.NewBulkString("PING"),
	})

	response := connection.processCommand(pingArray)

	if response.Type != resp.SimpleString {
		t.Errorf("Expected SimpleString response, got %s", response.Type)
	}

	if response.Value.(string) != "PONG" {
		t.Errorf("Expected 'PONG' response, got '%s'", response.Value.(string))
	}
}

func TestConnection_processCommand_ValidEcho(t *testing.T) {
	conn := newMockConn("")
	handler := commands.NewCommandHandler()
	handler.Register(commands.NewEchoCommand())

	connection := NewConnection(conn, handler)

	// Create an ECHO command message: *2\r\n$4\r\nECHO\r\n$5\r\nhello\r\n
	echoArray := resp.NewArray([]*resp.Message{
		resp.NewBulkString("ECHO"),
		resp.NewBulkString("hello"),
	})

	response := connection.processCommand(echoArray)

	if response.Type != resp.BulkString {
		t.Errorf("Expected BulkString response, got %s", response.Type)
	}

	if response.Value.(string) != "hello" {
		t.Errorf("Expected 'hello' response, got '%s'", response.Value.(string))
	}
}

func TestConnection_processCommand_CaseInsensitive(t *testing.T) {
	conn := newMockConn("")
	handler := commands.NewCommandHandler()
	handler.Register(commands.NewPingCommand())

	connection := NewConnection(conn, handler)

	// Test lowercase command
	pingArray := resp.NewArray([]*resp.Message{
		resp.NewBulkString("ping"),
	})

	response := connection.processCommand(pingArray)

	if response.Type != resp.SimpleString {
		t.Errorf("Expected SimpleString response, got %s", response.Type)
	}

	if response.Value.(string) != "PONG" {
		t.Errorf("Expected 'PONG' response, got '%s'", response.Value.(string))
	}
}

func TestConnection_processCommand_UnknownCommand(t *testing.T) {
	conn := newMockConn("")
	handler := commands.NewCommandHandler()

	connection := NewConnection(conn, handler)

	// Create an unknown command message
	unknownArray := resp.NewArray([]*resp.Message{
		resp.NewBulkString("UNKNOWN"),
	})

	response := connection.processCommand(unknownArray)

	if response.Type != resp.Error {
		t.Errorf("Expected Error response, got %s", response.Type)
	}

	expectedMsg := "ERR unknown command 'UNKNOWN'"
	if response.Value.(string) != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, response.Value.(string))
	}
}

func TestConnection_processCommand_InvalidInput(t *testing.T) {
	conn := newMockConn("")
	handler := commands.NewCommandHandler()

	connection := NewConnection(conn, handler)

	testCases := []struct {
		name     string
		message  *resp.Message
		expected string
	}{
		{
			name:     "non-array message",
			message:  resp.NewSimpleString("NOT_ARRAY"),
			expected: "ERR Protocol error: expected array",
		},
		{
			name:     "null array",
			message:  resp.NewNullArray(),
			expected: "ERR Protocol error: null array",
		},
		{
			name:     "empty array",
			message:  resp.NewArray([]*resp.Message{}),
			expected: "ERR Protocol error: empty array",
		},
		{
			name: "non-string command name",
			message: resp.NewArray([]*resp.Message{
				resp.NewInteger(123),
			}),
			expected: "ERR Protocol error: command name must be a string",
		},
		{
			name: "null command name",
			message: resp.NewArray([]*resp.Message{
				resp.NewNullBulkString(),
			}),
			expected: "ERR Protocol error: null command name",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response := connection.processCommand(tc.message)

			if response.Type != resp.Error {
				t.Errorf("Expected Error response, got %s", response.Type)
			}

			if response.Value.(string) != tc.expected {
				t.Errorf("Expected error message '%s', got '%s'", tc.expected, response.Value.(string))
			}
		})
	}
}
