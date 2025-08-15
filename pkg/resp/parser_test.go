package resp

import (
	"testing"
)

func TestParseSimpleString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "OK response",
			input:    "+OK\r\n",
			expected: "OK",
		},
		{
			name:     "Empty string",
			input:    "+\r\n",
			expected: "",
		},
		{
			name:     "PONG response",
			input:    "+PONG\r\n",
			expected: "PONG",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if msg.Type != SimpleString {
				t.Errorf("Expected type SimpleString, got %s", msg.Type)
			}

			if msg.Value != test.expected {
				t.Errorf("Expected value %q, got %v", test.expected, msg.Value)
			}
		})
	}
}

func TestParseError(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Error message",
			input:    "-Error message\r\n",
			expected: "Error message",
		},
		{
			name:     "ERR wrong type",
			input:    "-ERR wrong type\r\n",
			expected: "ERR wrong type",
		},
		{
			name:     "Empty error",
			input:    "-\r\n",
			expected: "",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if msg.Type != Error {
				t.Errorf("Expected type Error, got %s", msg.Type)
			}

			if msg.Value != test.expected {
				t.Errorf("Expected value %q, got %v", test.expected, msg.Value)
			}
		})
	}
}

func TestParseInteger(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected int64
	}{
		{
			name:     "Positive integer",
			input:    ":1000\r\n",
			expected: 1000,
		},
		{
			name:     "Negative integer",
			input:    ":-1\r\n",
			expected: -1,
		},
		{
			name:     "Zero",
			input:    ":0\r\n",
			expected: 0,
		},
		{
			name:     "Large integer",
			input:    ":9223372036854775807\r\n",
			expected: 9223372036854775807,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if msg.Type != Integer {
				t.Errorf("Expected type Integer, got %s", msg.Type)
			}

			if msg.Value != test.expected {
				t.Errorf("Expected value %d, got %v", test.expected, msg.Value)
			}
		})
	}
}

func TestParseBulkString(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
		isNull   bool
	}{
		{
			name:     "Normal bulk string",
			input:    "$6\r\nfoobar\r\n",
			expected: "foobar",
		},
		{
			name:     "Empty bulk string",
			input:    "$0\r\n\r\n",
			expected: "",
		},
		{
			name:     "Null bulk string",
			input:    "$-1\r\n",
			expected: nil,
			isNull:   true,
		},
		{
			name:     "Bulk string with special chars",
			input:    "$12\r\nhello\r\nworld\r\n",
			expected: "hello\r\nworld",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if msg.Type != BulkString {
				t.Errorf("Expected type BulkString, got %s", msg.Type)
			}

			if test.isNull {
				if !msg.IsNull() {
					t.Error("Expected null bulk string")
				}
			} else {
				if msg.Value != test.expected {
					t.Errorf("Expected value %q, got %v", test.expected, msg.Value)
				}
			}
		})
	}
}

func TestParseArray(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectedLen int
		isNull      bool
	}{
		{
			name:        "Empty array",
			input:       "*0\r\n",
			expectedLen: 0,
		},
		{
			name:        "Array with one element",
			input:       "*1\r\n$4\r\nping\r\n",
			expectedLen: 1,
		},
		{
			name:        "Array with two elements",
			input:       "*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n",
			expectedLen: 2,
		},
		{
			name:   "Null array",
			input:  "*-1\r\n",
			isNull: true,
		},
		{
			name:        "Mixed array",
			input:       "*3\r\n:1\r\n:2\r\n:3\r\n",
			expectedLen: 3,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}

			if msg.Type != Array {
				t.Errorf("Expected type Array, got %s", msg.Type)
			}

			if test.isNull {
				if !msg.IsNull() {
					t.Error("Expected null array")
				}
			} else {
				arr := msg.Value.([]*Message)
				if len(arr) != test.expectedLen {
					t.Errorf("Expected array length %d, got %d", test.expectedLen, len(arr))
				}
			}
		})
	}
}

// Test specific cases mentioned in the Jira issue
func TestSpecificJiraTestCases(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		testFunc func(*testing.T, *Message)
	}{
		{
			name:  "Null bulk string: $-1\\r\\n",
			input: "$-1\r\n",
			testFunc: func(t *testing.T, msg *Message) {
				if msg.Type != BulkString {
					t.Errorf("Expected BulkString, got %s", msg.Type)
				}
				if !msg.IsNull() {
					t.Error("Expected null bulk string")
				}
			},
		},
		{
			name:  "Array with ping command: *1\\r\\n$4\\r\\nping\\r\\n",
			input: "*1\r\n$4\r\nping\r\n",
			testFunc: func(t *testing.T, msg *Message) {
				if msg.Type != Array {
					t.Errorf("Expected Array, got %s", msg.Type)
				}
				arr := msg.Value.([]*Message)
				if len(arr) != 1 {
					t.Errorf("Expected array length 1, got %d", len(arr))
				}
				if arr[0].Type != BulkString {
					t.Errorf("Expected first element to be BulkString, got %s", arr[0].Type)
				}
				if arr[0].Value != "ping" {
					t.Errorf("Expected first element value 'ping', got %v", arr[0].Value)
				}
			},
		},
		{
			name:  "Simple OK response: +OK\\r\\n",
			input: "+OK\r\n",
			testFunc: func(t *testing.T, msg *Message) {
				if msg.Type != SimpleString {
					t.Errorf("Expected SimpleString, got %s", msg.Type)
				}
				if msg.Value != "OK" {
					t.Errorf("Expected value 'OK', got %v", msg.Value)
				}
			},
		},
		{
			name:  "Error message: -Error message\\r\\n",
			input: "-Error message\r\n",
			testFunc: func(t *testing.T, msg *Message) {
				if msg.Type != Error {
					t.Errorf("Expected Error, got %s", msg.Type)
				}
				if msg.Value != "Error message" {
					t.Errorf("Expected value 'Error message', got %v", msg.Value)
				}
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			msg, err := ParseString(test.input)
			if err != nil {
				t.Fatalf("Unexpected error: %v", err)
			}
			test.testFunc(t, msg)
		})
	}
}

func TestParseInvalidMessages(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{
			name:  "Invalid type character",
			input: "?invalid\r\n",
		},
		{
			name:  "Missing CRLF",
			input: "+OK",
		},
		{
			name:  "Invalid integer",
			input: ":abc\r\n",
		},
		{
			name:  "Invalid bulk string length",
			input: "$abc\r\n",
		},
		{
			name:  "Negative bulk string length (not -1)",
			input: "$-2\r\n",
		},
		{
			name:  "Bulk string length mismatch",
			input: "$5\r\nfoo\r\n",
		},
		{
			name:  "Invalid array length",
			input: "*abc\r\n",
		},
		{
			name:  "Negative array length (not -1)",
			input: "*-2\r\n",
		},
		{
			name:  "Incomplete array",
			input: "*2\r\n$3\r\nfoo\r\n",
		},
		{
			name:  "Line with only CR",
			input: "+OK\r",
		},
		{
			name:  "Line with only LF",
			input: "+OK\n",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := ParseString(test.input)
			if err == nil {
				t.Error("Expected error for invalid input, got nil")
			}
		})
	}
}

func TestComplexNestedArray(t *testing.T) {
	// Test parsing a complex nested array structure
	input := "*2\r\n*3\r\n:1\r\n:2\r\n:3\r\n*2\r\n+Hello\r\n-World\r\n"

	msg, err := ParseString(input)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	if msg.Type != Array {
		t.Errorf("Expected Array, got %s", msg.Type)
	}

	arr := msg.Value.([]*Message)
	if len(arr) != 2 {
		t.Errorf("Expected array length 2, got %d", len(arr))
	}

	// Check first sub-array
	if arr[0].Type != Array {
		t.Errorf("Expected first element to be Array, got %s", arr[0].Type)
	}
	subArr1 := arr[0].Value.([]*Message)
	if len(subArr1) != 3 {
		t.Errorf("Expected first sub-array length 3, got %d", len(subArr1))
	}

	// Check second sub-array
	if arr[1].Type != Array {
		t.Errorf("Expected second element to be Array, got %s", arr[1].Type)
	}
	subArr2 := arr[1].Value.([]*Message)
	if len(subArr2) != 2 {
		t.Errorf("Expected second sub-array length 2, got %d", len(subArr2))
	}

	// Check that elements have correct types and values
	if subArr1[0].Value != int64(1) {
		t.Errorf("Expected first element of first sub-array to be 1, got %v", subArr1[0].Value)
	}
	if subArr2[0].Value != "Hello" {
		t.Errorf("Expected first element of second sub-array to be 'Hello', got %v", subArr2[0].Value)
	}
}
