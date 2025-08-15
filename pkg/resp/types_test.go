package resp

import (
	"testing"
)

func TestMessageType_String(t *testing.T) {
	tests := []struct {
		msgType  MessageType
		expected string
	}{
		{SimpleString, "SimpleString"},
		{Error, "Error"},
		{Integer, "Integer"},
		{BulkString, "BulkString"},
		{Array, "Array"},
		{MessageType(999), "Unknown"},
	}

	for _, test := range tests {
		t.Run(test.expected, func(t *testing.T) {
			result := test.msgType.String()
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestNewSimpleString(t *testing.T) {
	msg := NewSimpleString("OK")

	if msg.Type != SimpleString {
		t.Errorf("Expected type SimpleString, got %s", msg.Type)
	}

	if msg.Value != "OK" {
		t.Errorf("Expected value OK, got %v", msg.Value)
	}
}

func TestNewError(t *testing.T) {
	msg := NewError("Error message")

	if msg.Type != Error {
		t.Errorf("Expected type Error, got %s", msg.Type)
	}

	if msg.Value != "Error message" {
		t.Errorf("Expected value 'Error message', got %v", msg.Value)
	}
}

func TestNewInteger(t *testing.T) {
	msg := NewInteger(1000)

	if msg.Type != Integer {
		t.Errorf("Expected type Integer, got %s", msg.Type)
	}

	if msg.Value != int64(1000) {
		t.Errorf("Expected value 1000, got %v", msg.Value)
	}
}

func TestNewBulkString(t *testing.T) {
	msg := NewBulkString("foobar")

	if msg.Type != BulkString {
		t.Errorf("Expected type BulkString, got %s", msg.Type)
	}

	if msg.Value != "foobar" {
		t.Errorf("Expected value 'foobar', got %v", msg.Value)
	}
}

func TestNewNullBulkString(t *testing.T) {
	msg := NewNullBulkString()

	if msg.Type != BulkString {
		t.Errorf("Expected type BulkString, got %s", msg.Type)
	}

	if msg.Value != nil {
		t.Errorf("Expected value nil, got %v", msg.Value)
	}

	if !msg.IsNull() {
		t.Error("Expected message to be null")
	}
}

func TestNewArray(t *testing.T) {
	elements := []*Message{
		NewBulkString("foo"),
		NewBulkString("bar"),
	}
	msg := NewArray(elements)

	if msg.Type != Array {
		t.Errorf("Expected type Array, got %s", msg.Type)
	}

	arr := msg.Value.([]*Message)
	if len(arr) != 2 {
		t.Errorf("Expected array length 2, got %d", len(arr))
	}

	if arr[0].Value != "foo" {
		t.Errorf("Expected first element 'foo', got %v", arr[0].Value)
	}

	if arr[1].Value != "bar" {
		t.Errorf("Expected second element 'bar', got %v", arr[1].Value)
	}
}

func TestNewNullArray(t *testing.T) {
	msg := NewNullArray()

	if msg.Type != Array {
		t.Errorf("Expected type Array, got %s", msg.Type)
	}

	if msg.Value != nil {
		t.Errorf("Expected value nil, got %v", msg.Value)
	}

	if !msg.IsNull() {
		t.Error("Expected message to be null")
	}
}

func TestMessage_String(t *testing.T) {
	tests := []struct {
		name     string
		message  *Message
		expected string
	}{
		{
			name:     "SimpleString",
			message:  NewSimpleString("OK"),
			expected: `SimpleString("OK")`,
		},
		{
			name:     "Error",
			message:  NewError("Error message"),
			expected: `Error("Error message")`,
		},
		{
			name:     "Integer",
			message:  NewInteger(1000),
			expected: "Integer(1000)",
		},
		{
			name:     "BulkString",
			message:  NewBulkString("foobar"),
			expected: `BulkString("foobar")`,
		},
		{
			name:     "NullBulkString",
			message:  NewNullBulkString(),
			expected: "BulkString(null)",
		},
		{
			name:     "Array",
			message:  NewArray([]*Message{NewBulkString("foo")}),
			expected: "Array(1 elements)",
		},
		{
			name:     "NullArray",
			message:  NewNullArray(),
			expected: "Array(null)",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := test.message.String()
			if result != test.expected {
				t.Errorf("Expected %s, got %s", test.expected, result)
			}
		})
	}
}

func TestMessage_AsString(t *testing.T) {
	tests := []struct {
		name        string
		message     *Message
		expected    string
		expectError bool
	}{
		{
			name:     "SimpleString",
			message:  NewSimpleString("OK"),
			expected: "OK",
		},
		{
			name:     "Error",
			message:  NewError("Error message"),
			expected: "Error message",
		},
		{
			name:     "BulkString",
			message:  NewBulkString("foobar"),
			expected: "foobar",
		},
		{
			name:        "NullBulkString",
			message:     NewNullBulkString(),
			expectError: true,
		},
		{
			name:        "Integer",
			message:     NewInteger(1000),
			expectError: true,
		},
		{
			name:        "Array",
			message:     NewArray([]*Message{}),
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.message.AsString()

			if test.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != test.expected {
					t.Errorf("Expected %s, got %s", test.expected, result)
				}
			}
		})
	}
}

func TestMessage_AsInteger(t *testing.T) {
	tests := []struct {
		name        string
		message     *Message
		expected    int64
		expectError bool
	}{
		{
			name:     "Integer",
			message:  NewInteger(1000),
			expected: 1000,
		},
		{
			name:        "SimpleString",
			message:     NewSimpleString("OK"),
			expectError: true,
		},
		{
			name:        "Error",
			message:     NewError("Error message"),
			expectError: true,
		},
		{
			name:        "BulkString",
			message:     NewBulkString("foobar"),
			expectError: true,
		},
		{
			name:        "Array",
			message:     NewArray([]*Message{}),
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.message.AsInteger()

			if test.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if result != test.expected {
					t.Errorf("Expected %d, got %d", test.expected, result)
				}
			}
		})
	}
}

func TestMessage_AsArray(t *testing.T) {
	elements := []*Message{
		NewBulkString("foo"),
		NewBulkString("bar"),
	}

	tests := []struct {
		name        string
		message     *Message
		expected    []*Message
		expectError bool
	}{
		{
			name:     "Array",
			message:  NewArray(elements),
			expected: elements,
		},
		{
			name:        "NullArray",
			message:     NewNullArray(),
			expectError: true,
		},
		{
			name:        "SimpleString",
			message:     NewSimpleString("OK"),
			expectError: true,
		},
		{
			name:        "Integer",
			message:     NewInteger(1000),
			expectError: true,
		},
		{
			name:        "BulkString",
			message:     NewBulkString("foobar"),
			expectError: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := test.message.AsArray()

			if test.expectError {
				if err == nil {
					t.Error("Expected error, got nil")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error: %v", err)
				}
				if len(result) != len(test.expected) {
					t.Errorf("Expected array length %d, got %d", len(test.expected), len(result))
				}
				for i, elem := range test.expected {
					if result[i].Value != elem.Value {
						t.Errorf("Expected element %d to be %v, got %v", i, elem.Value, result[i].Value)
					}
				}
			}
		})
	}
}
