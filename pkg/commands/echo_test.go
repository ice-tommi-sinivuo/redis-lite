package commands

import (
	"testing"

	"github.com/tsinivuo/redis-lite/pkg/resp"
)

func TestEchoCommand_Name(t *testing.T) {
	cmd := NewEchoCommand()
	if cmd.Name() != "ECHO" {
		t.Errorf("Expected command name 'ECHO', got '%s'", cmd.Name())
	}
}

func TestEchoCommand_Validate(t *testing.T) {
	cmd := NewEchoCommand()

	testCases := []struct {
		name string
		args []*resp.Message
		want bool
	}{
		{
			name: "exactly one argument",
			args: []*resp.Message{resp.NewBulkString("hello")},
			want: true,
		},
		{
			name: "no arguments",
			args: []*resp.Message{},
			want: false,
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

func TestEchoCommand_Execute(t *testing.T) {
	cmd := NewEchoCommand()

	testCases := []struct {
		name         string
		args         []*resp.Message
		expectedType resp.MessageType
		expectedVal  interface{}
	}{
		{
			name:         "bulk string argument",
			args:         []*resp.Message{resp.NewBulkString("hello world")},
			expectedType: resp.BulkString,
			expectedVal:  "hello world",
		},
		{
			name:         "simple string argument",
			args:         []*resp.Message{resp.NewSimpleString("test")},
			expectedType: resp.BulkString,
			expectedVal:  "test",
		},
		{
			name:         "integer argument",
			args:         []*resp.Message{resp.NewInteger(42)},
			expectedType: resp.BulkString,
			expectedVal:  "42",
		},
		{
			name:         "null bulk string argument",
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

func TestEchoCommand_Execute_InvalidType(t *testing.T) {
	cmd := NewEchoCommand()

	// Test with array argument (should return error)
	args := []*resp.Message{resp.NewArray([]*resp.Message{resp.NewBulkString("test")})}
	response, err := cmd.Execute(args)

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if response.Type != resp.Error {
		t.Errorf("Expected Error response, got %s", response.Type)
	}

	expectedMsg := "ERR invalid argument type for ECHO"
	if response.Value.(string) != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, response.Value.(string))
	}
}
