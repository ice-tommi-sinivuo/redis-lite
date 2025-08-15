package commands

import (
	"testing"

	"github.com/tsinivuo/redis-lite/pkg/resp"
	"github.com/tsinivuo/redis-lite/pkg/storage"
)

func TestSetCommand_Name(t *testing.T) {
	cmd := NewSetCommand()
	if cmd.Name() != "SET" {
		t.Errorf("Expected command name 'SET', got '%s'", cmd.Name())
	}
}

func TestSetCommand_Validate(t *testing.T) {
	cmd := NewSetCommand()

	tests := []struct {
		name    string
		args    []*resp.Message
		wantErr bool
	}{
		{
			name:    "valid args",
			args:    []*resp.Message{resp.NewBulkString("key"), resp.NewBulkString("value")},
			wantErr: false,
		},
		{
			name:    "no args",
			args:    []*resp.Message{},
			wantErr: true,
		},
		{
			name:    "one arg",
			args:    []*resp.Message{resp.NewBulkString("key")},
			wantErr: true,
		},
		{
			name:    "three args",
			args:    []*resp.Message{resp.NewBulkString("key"), resp.NewBulkString("value"), resp.NewBulkString("extra")},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.Validate(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSetCommand_Execute(t *testing.T) {
	cmd := NewSetCommand()
	store := storage.NewMemoryStore()

	tests := []struct {
		name         string
		args         []*resp.Message
		wantResponse *resp.Message
		wantErr      bool
	}{
		{
			name:         "set bulk string key and value",
			args:         []*resp.Message{resp.NewBulkString("key1"), resp.NewBulkString("value1")},
			wantResponse: resp.NewSimpleString("OK"),
			wantErr:      false,
		},
		{
			name:         "set simple string key and bulk string value",
			args:         []*resp.Message{resp.NewSimpleString("key2"), resp.NewBulkString("value2")},
			wantResponse: resp.NewSimpleString("OK"),
			wantErr:      false,
		},
		{
			name:         "set bulk string key and simple string value",
			args:         []*resp.Message{resp.NewBulkString("key3"), resp.NewSimpleString("value3")},
			wantResponse: resp.NewSimpleString("OK"),
			wantErr:      false,
		},
		{
			name:         "set with integer value",
			args:         []*resp.Message{resp.NewBulkString("key4"), resp.NewInteger(42)},
			wantResponse: resp.NewSimpleString("OK"),
			wantErr:      false,
		},
		{
			name:         "set with empty string value",
			args:         []*resp.Message{resp.NewBulkString("key5"), resp.NewBulkString("")},
			wantResponse: resp.NewSimpleString("OK"),
			wantErr:      false,
		},
		{
			name:         "set with null bulk string value",
			args:         []*resp.Message{resp.NewBulkString("key6"), resp.NewNullBulkString()},
			wantResponse: resp.NewSimpleString("OK"),
			wantErr:      false,
		},
		{
			name:    "invalid key type (null)",
			args:    []*resp.Message{resp.NewNullBulkString(), resp.NewBulkString("value")},
			wantErr: false, // Should return error response, not Go error
		},
		{
			name:    "invalid key type (integer)",
			args:    []*resp.Message{resp.NewInteger(123), resp.NewBulkString("value")},
			wantErr: false, // Should return error response, not Go error
		},
		{
			name:    "invalid value type (array)",
			args:    []*resp.Message{resp.NewBulkString("key"), resp.NewArray([]*resp.Message{})},
			wantErr: false, // Should return error response, not Go error
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response, err := cmd.Execute(tt.args, store)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr && tt.wantResponse != nil {
				if response.Type != tt.wantResponse.Type {
					t.Errorf("Execute() response type = %v, want %v", response.Type, tt.wantResponse.Type)
					return
				}

				if response.Value != tt.wantResponse.Value {
					t.Errorf("Execute() response value = %v, want %v", response.Value, tt.wantResponse.Value)
				}
			}
		})
	}
}

func TestSetCommand_Integration(t *testing.T) {
	cmd := NewSetCommand()
	store := storage.NewMemoryStore()

	// Test successful SET operation
	args := []*resp.Message{resp.NewBulkString("testkey"), resp.NewBulkString("testvalue")}
	response, err := cmd.Execute(args, store)

	if err != nil {
		t.Errorf("Execute() returned error: %v", err)
	}

	if response.Type != resp.SimpleString || response.Value != "OK" {
		t.Errorf("Expected OK response, got %v", response)
	}

	// Verify the value was actually stored
	storedValue, exists := store.Get("testkey")
	if !exists {
		t.Error("Key was not stored in the storage")
	}
	if storedValue != "testvalue" {
		t.Errorf("Expected stored value 'testvalue', got '%s'", storedValue)
	}
}

func TestSetCommand_OverwriteValue(t *testing.T) {
	cmd := NewSetCommand()
	store := storage.NewMemoryStore()

	// Set initial value
	args1 := []*resp.Message{resp.NewBulkString("key"), resp.NewBulkString("value1")}
	cmd.Execute(args1, store)

	// Overwrite with new value
	args2 := []*resp.Message{resp.NewBulkString("key"), resp.NewBulkString("value2")}
	response, err := cmd.Execute(args2, store)

	if err != nil {
		t.Errorf("Execute() returned error: %v", err)
	}

	if response.Type != resp.SimpleString || response.Value != "OK" {
		t.Errorf("Expected OK response, got %v", response)
	}

	// Verify the new value
	storedValue, exists := store.Get("key")
	if !exists {
		t.Error("Key was not found in storage")
	}
	if storedValue != "value2" {
		t.Errorf("Expected stored value 'value2', got '%s'", storedValue)
	}
}
