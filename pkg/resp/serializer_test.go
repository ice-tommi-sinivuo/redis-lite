package resp

import (
	"testing"
)

func TestSerializeSimpleString(t *testing.T) {
	tests := []struct {
		name     string
		message  *Message
		expected string
	}{
		{
			name:     "OK response",
			message:  NewSimpleString("OK"),
			expected: "+OK\r\n",
		},
		{
			name:     "PONG response",
			message:  NewSimpleString("PONG"),
			expected: "+PONG\r\n",
		},
		{
			name:     "Empty string",
			message:  NewSimpleString(""),
			expected: "+\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := SerializeToString(test.message)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestSerializeError(t *testing.T) {
	tests := []struct {
		name     string
		message  *Message
		expected string
	}{
		{
			name:     "Error message",
			message:  NewError("Error message"),
			expected: "-Error message\r\n",
		},
		{
			name:     "ERR wrong type",
			message:  NewError("ERR wrong type"),
			expected: "-ERR wrong type\r\n",
		},
		{
			name:     "Empty error",
			message:  NewError(""),
			expected: "-\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := SerializeToString(test.message)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestSerializeInteger(t *testing.T) {
	tests := []struct {
		name     string
		message  *Message
		expected string
	}{
		{
			name:     "Positive integer",
			message:  NewInteger(1000),
			expected: ":1000\r\n",
		},
		{
			name:     "Negative integer",
			message:  NewInteger(-1),
			expected: ":-1\r\n",
		},
		{
			name:     "Zero",
			message:  NewInteger(0),
			expected: ":0\r\n",
		},
		{
			name:     "Large integer",
			message:  NewInteger(9223372036854775807),
			expected: ":9223372036854775807\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := SerializeToString(test.message)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestSerializeBulkString(t *testing.T) {
	tests := []struct {
		name     string
		message  *Message
		expected string
	}{
		{
			name:     "Normal bulk string",
			message:  NewBulkString("foobar"),
			expected: "$6\r\nfoobar\r\n",
		},
		{
			name:     "Empty bulk string",
			message:  NewBulkString(""),
			expected: "$0\r\n\r\n",
		},
		{
			name:     "Null bulk string",
			message:  NewNullBulkString(),
			expected: "$-1\r\n",
		},
		{
			name:     "Bulk string with special chars",
			message:  NewBulkString("hello\r\nworld"),
			expected: "$12\r\nhello\r\nworld\r\n",
		},
		{
			name:     "Bulk string with binary data",
			message:  NewBulkString("foo\x00bar"),
			expected: "$7\r\nfoo\x00bar\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := SerializeToString(test.message)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestSerializeArray(t *testing.T) {
	tests := []struct {
		name     string
		message  *Message
		expected string
	}{
		{
			name:     "Empty array",
			message:  NewArray([]*Message{}),
			expected: "*0\r\n",
		},
		{
			name:     "Array with one element",
			message:  NewArray([]*Message{NewBulkString("ping")}),
			expected: "*1\r\n$4\r\nping\r\n",
		},
		{
			name: "Array with two elements",
			message: NewArray([]*Message{
				NewBulkString("foo"),
				NewBulkString("bar"),
			}),
			expected: "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
		},
		{
			name:     "Null array",
			message:  NewNullArray(),
			expected: "*-1\r\n",
		},
		{
			name: "Mixed array",
			message: NewArray([]*Message{
				NewInteger(1),
				NewInteger(2),
				NewInteger(3),
			}),
			expected: "*3\r\n:1\r\n:2\r\n:3\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := SerializeToString(test.message)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

// Test specific cases mentioned in the Jira issue
func TestSerializeSpecificJiraTestCases(t *testing.T) {
	tests := []struct {
		name     string
		message  *Message
		expected string
	}{
		{
			name:     "Null bulk string: $-1\\r\\n",
			message:  NewNullBulkString(),
			expected: "$-1\r\n",
		},
		{
			name:     "Array with ping command: *1\\r\\n$4\\r\\nping\\r\\n",
			message:  NewArray([]*Message{NewBulkString("ping")}),
			expected: "*1\r\n$4\r\nping\r\n",
		},
		{
			name:     "Simple OK response: +OK\\r\\n",
			message:  NewSimpleString("OK"),
			expected: "+OK\r\n",
		},
		{
			name:     "Error message: -Error message\\r\\n",
			message:  NewError("Error message"),
			expected: "-Error message\r\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := SerializeToString(test.message)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if result != test.expected {
				t.Errorf("Expected %q, got %q", test.expected, result)
			}
		})
	}
}

func TestSerializeInvalidMessages(t *testing.T) {
	tests := []struct {
		name    string
		message *Message
	}{
		{
			name:    "Simple string with CR",
			message: NewSimpleString("Hello\rWorld"),
		},
		{
			name:    "Simple string with LF",
			message: NewSimpleString("Hello\nWorld"),
		},
		{
			name:    "Error with CR",
			message: NewError("Error\rmessage"),
		},
		{
			name:    "Error with LF",
			message: NewError("Error\nmessage"),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := SerializeToString(test.message)
			if err == nil {
				t.Error("Expected error for invalid message, got nil")
			}
		})
	}
}

func TestSerializeComplexNestedArray(t *testing.T) {
	// Test serializing a complex nested array structure
	message := NewArray([]*Message{
		NewArray([]*Message{
			NewInteger(1),
			NewInteger(2),
			NewInteger(3),
		}),
		NewArray([]*Message{
			NewSimpleString("Hello"),
			NewError("World"),
		}),
	})

	expected := "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n"

	result, err := SerializeToString(message)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestSerializeToBytes(t *testing.T) {
	message := NewSimpleString("OK")
	expected := []byte("+OK\r\n")

	result, err := SerializeToBytes(message)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if string(result) != string(expected) {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

// Test round-trip serialization and parsing
func TestRoundTripSerialization(t *testing.T) {
	tests := []struct {
		name    string
		message *Message
	}{
		{
			name:    "SimpleString",
			message: NewSimpleString("OK"),
		},
		{
			name:    "Error",
			message: NewError("Error message"),
		},
		{
			name:    "Integer",
			message: NewInteger(1000),
		},
		{
			name:    "BulkString",
			message: NewBulkString("foobar"),
		},
		{
			name:    "NullBulkString",
			message: NewNullBulkString(),
		},
		{
			name:    "Array",
			message: NewArray([]*Message{NewBulkString("foo"), NewBulkString("bar")}),
		},
		{
			name:    "NullArray",
			message: NewNullArray(),
		},
		{
			name: "ComplexArray",
			message: NewArray([]*Message{
				NewSimpleString("OK"),
				NewInteger(42),
				NewBulkString("hello"),
				NewArray([]*Message{NewInteger(1), NewInteger(2)}),
			}),
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Serialize the message
			serialized, err := SerializeToString(test.message)
			if err != nil {
				t.Fatalf("Serialization error: %v", err)
			}

			// Parse it back
			parsed, err := ParseString(serialized)
			if err != nil {
				t.Fatalf("Parsing error: %v", err)
			}

			// Compare the types
			if parsed.Type != test.message.Type {
				t.Errorf("Type mismatch: expected %s, got %s", test.message.Type, parsed.Type)
			}

			// Compare the values (this is a simplified comparison)
			if test.message.IsNull() {
				if !parsed.IsNull() {
					t.Error("Expected parsed message to be null")
				}
			} else {
				compareMessages(t, test.message, parsed)
			}
		})
	}
}

// Helper function to compare messages recursively
func compareMessages(t *testing.T, expected, actual *Message) {
	if expected.Type != actual.Type {
		t.Errorf("Type mismatch: expected %s, got %s", expected.Type, actual.Type)
		return
	}

	switch expected.Type {
	case SimpleString, Error:
		if expected.Value != actual.Value {
			t.Errorf("Value mismatch: expected %v, got %v", expected.Value, actual.Value)
		}
	case Integer:
		if expected.Value != actual.Value {
			t.Errorf("Value mismatch: expected %v, got %v", expected.Value, actual.Value)
		}
	case BulkString:
		if expected.IsNull() != actual.IsNull() {
			t.Errorf("Null mismatch: expected null=%v, got null=%v", expected.IsNull(), actual.IsNull())
		} else if !expected.IsNull() && expected.Value != actual.Value {
			t.Errorf("Value mismatch: expected %v, got %v", expected.Value, actual.Value)
		}
	case Array:
		if expected.IsNull() != actual.IsNull() {
			t.Errorf("Null mismatch: expected null=%v, got null=%v", expected.IsNull(), actual.IsNull())
		} else if !expected.IsNull() {
			expectedArr := expected.Value.([]*Message)
			actualArr := actual.Value.([]*Message)
			if len(expectedArr) != len(actualArr) {
				t.Errorf("Array length mismatch: expected %d, got %d", len(expectedArr), len(actualArr))
			} else {
				for i := range expectedArr {
					compareMessages(t, expectedArr[i], actualArr[i])
				}
			}
		}
	}
}
