package commands

import (
	"testing"

	"github.com/tsinivuo/redis-lite/pkg/resp"
	"github.com/tsinivuo/redis-lite/pkg/storage"
)

func TestGetCommand_Name(t *testing.T) {
	cmd := NewGetCommand()
	if cmd.Name() != "GET" {
		t.Errorf("Expected command name 'GET', got '%s'", cmd.Name())
	}
}

func TestGetCommand_Validate(t *testing.T) {
	cmd := NewGetCommand()

	tests := []struct {
		name    string
		args    []*resp.Message
		wantErr bool
	}{
		{
			name:    "valid args",
			args:    []*resp.Message{resp.NewBulkString("key")},
			wantErr: false,
		},
		{
			name:    "no args",
			args:    []*resp.Message{},
			wantErr: true,
		},
		{
			name:    "two args",
			args:    []*resp.Message{resp.NewBulkString("key"), resp.NewBulkString("extra")},
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

func TestGetCommand_Execute(t *testing.T) {
	cmd := NewGetCommand()
	store := storage.NewMemoryStore()

	// Pre-populate store with test data
	store.Set("existing_key", "existing_value")
	store.Set("empty_value", "")

	tests := []struct {
		name         string
		args         []*resp.Message
		wantResponse *resp.Message
		wantErr      bool
	}{
		{
			name:         "get existing key with bulk string",
			args:         []*resp.Message{resp.NewBulkString("existing_key")},
			wantResponse: resp.NewBulkString("existing_value"),
			wantErr:      false,
		},
		{
			name:         "get existing key with simple string",
			args:         []*resp.Message{resp.NewSimpleString("existing_key")},
			wantResponse: resp.NewBulkString("existing_value"),
			wantErr:      false,
		},
		{
			name:         "get non-existent key",
			args:         []*resp.Message{resp.NewBulkString("nonexistent_key")},
			wantResponse: resp.NewNullBulkString(),
			wantErr:      false,
		},
		{
			name:         "get key with empty value",
			args:         []*resp.Message{resp.NewBulkString("empty_value")},
			wantResponse: resp.NewBulkString(""),
			wantErr:      false,
		},
		{
			name:    "invalid key type (null)",
			args:    []*resp.Message{resp.NewNullBulkString()},
			wantErr: false, // Should return error response, not Go error
		},
		{
			name:    "invalid key type (integer)",
			args:    []*resp.Message{resp.NewInteger(123)},
			wantErr: false, // Should return error response, not Go error
		},
		{
			name:    "invalid key type (array)",
			args:    []*resp.Message{resp.NewArray([]*resp.Message{})},
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

				// Handle null bulk string comparison
				if response.Type == resp.BulkString && tt.wantResponse.Type == resp.BulkString {
					if response.Value == nil && tt.wantResponse.Value == nil {
						// Both null, OK
					} else if response.Value != nil && tt.wantResponse.Value != nil {
						if response.Value != tt.wantResponse.Value {
							t.Errorf("Execute() response value = %v, want %v", response.Value, tt.wantResponse.Value)
						}
					} else {
						t.Errorf("Execute() response value = %v, want %v", response.Value, tt.wantResponse.Value)
					}
				} else {
					if response.Value != tt.wantResponse.Value {
						t.Errorf("Execute() response value = %v, want %v", response.Value, tt.wantResponse.Value)
					}
				}
			}
		})
	}
}

func TestGetCommand_Integration(t *testing.T) {
	getCmd := NewGetCommand()
	setCmd := NewSetCommand()
	store := storage.NewMemoryStore()

	// First, set a value using SET command
	setArgs := []*resp.Message{resp.NewBulkString("integrationkey"), resp.NewBulkString("integrationvalue")}
	setResponse, err := setCmd.Execute(setArgs, store)

	if err != nil {
		t.Errorf("SET Execute() returned error: %v", err)
	}

	if setResponse.Type != resp.SimpleString || setResponse.Value != "OK" {
		t.Errorf("Expected OK response from SET, got %v", setResponse)
	}

	// Now, get the value using GET command
	getArgs := []*resp.Message{resp.NewBulkString("integrationkey")}
	getResponse, err := getCmd.Execute(getArgs, store)

	if err != nil {
		t.Errorf("GET Execute() returned error: %v", err)
	}

	if getResponse.Type != resp.BulkString {
		t.Errorf("Expected BulkString response, got %v", getResponse.Type)
	}

	if getResponse.Value != "integrationvalue" {
		t.Errorf("Expected value 'integrationvalue', got '%v'", getResponse.Value)
	}
}

func TestGetCommand_NonExistentKey(t *testing.T) {
	cmd := NewGetCommand()
	store := storage.NewMemoryStore()

	// Try to get a key that doesn't exist
	args := []*resp.Message{resp.NewBulkString("nonexistent")}
	response, err := cmd.Execute(args, store)

	if err != nil {
		t.Errorf("Execute() returned error: %v", err)
	}

	if response.Type != resp.BulkString {
		t.Errorf("Expected BulkString response type, got %v", response.Type)
	}

	if response.Value != nil {
		t.Errorf("Expected null value for non-existent key, got %v", response.Value)
	}
}

func TestGetCommand_KeyTypes(t *testing.T) {
	cmd := NewGetCommand()
	store := storage.NewMemoryStore()

	// Test different key types that should work
	store.Set("test", "value")

	tests := []struct {
		name    string
		keyArg  *resp.Message
		wantErr bool
	}{
		{
			name:    "bulk string key",
			keyArg:  resp.NewBulkString("test"),
			wantErr: false,
		},
		{
			name:    "simple string key",
			keyArg:  resp.NewSimpleString("test"),
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			args := []*resp.Message{tt.keyArg}
			response, err := cmd.Execute(args, store)

			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if response.Type != resp.BulkString {
					t.Errorf("Expected BulkString response, got %v", response.Type)
				}
				if response.Value != "value" {
					t.Errorf("Expected value 'value', got %v", response.Value)
				}
			}
		})
	}
}
