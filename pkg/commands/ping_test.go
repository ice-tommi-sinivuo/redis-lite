package commands

import (
	"testing"

	"github.com/tsinivuo/redis-lite/pkg/resp"
)

func TestPingCommand_Name(t *testing.T) {
	cmd := NewPingCommand()
	if cmd.Name() != "PING" {
		t.Errorf("Expected command name 'PING', got '%s'", cmd.Name())
	}
}

func TestPingCommand_Validate(t *testing.T) {
	cmd := NewPingCommand()

	// Test valid cases
	testCases := []struct {
		name string
		args []*resp.Message
		want bool
	}{
		{
			name: "no arguments",
			args: []*resp.Message{},
			want: true,
		},
		{
			name: "one argument",
			args: []*resp.Message{resp.NewBulkString("hello")},
			want: true,
		},
		{
			name: "too many arguments",
			args: []*resp.Message{
				resp.NewBulkString("hello"),
				resp.NewBulkString("world"),
			},
			want: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := cmd.Validate(tc.args)
			if tc.want && err != nil {
				t.Errorf("Expected validation to pass, but got error: %v", err)
			}
			if !tc.want && err == nil {
				t.Errorf("Expected validation to fail, but it passed")
			}
		})
	}
}

func TestPingCommand_Execute(t *testing.T) {
	cmd := NewPingCommand()

	testCases := []struct {
		name         string
		args         []*resp.Message
		expectedType resp.MessageType
		expectedVal  interface{}
	}{
		{
			name:         "no arguments returns PONG",
			args:         []*resp.Message{},
			expectedType: resp.SimpleString,
			expectedVal:  "PONG",
		},
		{
			name:         "bulk string argument echoed back",
			args:         []*resp.Message{resp.NewBulkString("hello")},
			expectedType: resp.BulkString,
			expectedVal:  "hello",
		},
		{
			name:         "simple string argument echoed back",
			args:         []*resp.Message{resp.NewSimpleString("world")},
			expectedType: resp.SimpleString,
			expectedVal:  "world",
		},
		{
			name:         "null bulk string echoed back",
			args:         []*resp.Message{resp.NewNullBulkString()},
			expectedType: resp.BulkString,
			expectedVal:  nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			response, err := cmd.Execute(tc.args)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if response.Type != tc.expectedType {
				t.Errorf("Expected response type %s, got %s", tc.expectedType, response.Type)
			}

			if response.Value != tc.expectedVal {
				t.Errorf("Expected response value %v, got %v", tc.expectedVal, response.Value)
			}
		})
	}
}

func TestPingCommand_Execute_InvalidType(t *testing.T) {
	cmd := NewPingCommand()

	// Test with integer argument (should return error)
	args := []*resp.Message{resp.NewInteger(42)}
	response, err := cmd.Execute(args)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.Type != resp.Error {
		t.Errorf("Expected Error response, got %s", response.Type)
	}

	expectedMsg := "ERR invalid argument type for PING"
	if response.Value.(string) != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, response.Value.(string))
	}
}
