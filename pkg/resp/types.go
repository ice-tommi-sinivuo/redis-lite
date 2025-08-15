package resp

import (
	"fmt"
)

// MessageType represents the different RESP data types
type MessageType int

const (
	// SimpleString represents a simple string (+)
	SimpleString MessageType = iota
	// Error represents an error message (-)
	Error
	// Integer represents an integer (:)
	Integer
	// BulkString represents a bulk string ($)
	BulkString
	// Array represents an array (*)
	Array
)

// String returns the string representation of the MessageType
func (mt MessageType) String() string {
	switch mt {
	case SimpleString:
		return "SimpleString"
	case Error:
		return "Error"
	case Integer:
		return "Integer"
	case BulkString:
		return "BulkString"
	case Array:
		return "Array"
	default:
		return "Unknown"
	}
}

// Message represents a RESP message with its type and value
type Message struct {
	Type  MessageType
	Value interface{}
}

// NewSimpleString creates a new simple string message
func NewSimpleString(value string) *Message {
	return &Message{
		Type:  SimpleString,
		Value: value,
	}
}

// NewError creates a new error message
func NewError(value string) *Message {
	return &Message{
		Type:  Error,
		Value: value,
	}
}

// NewInteger creates a new integer message
func NewInteger(value int64) *Message {
	return &Message{
		Type:  Integer,
		Value: value,
	}
}

// NewBulkString creates a new bulk string message
func NewBulkString(value string) *Message {
	return &Message{
		Type:  BulkString,
		Value: value,
	}
}

// NewNullBulkString creates a new null bulk string message
func NewNullBulkString() *Message {
	return &Message{
		Type:  BulkString,
		Value: nil,
	}
}

// NewArray creates a new array message
func NewArray(value []*Message) *Message {
	return &Message{
		Type:  Array,
		Value: value,
	}
}

// NewNullArray creates a new null array message
func NewNullArray() *Message {
	return &Message{
		Type:  Array,
		Value: nil,
	}
}

// String returns a string representation of the message for debugging
func (m *Message) String() string {
	switch m.Type {
	case SimpleString:
		return fmt.Sprintf("SimpleString(%q)", m.Value)
	case Error:
		return fmt.Sprintf("Error(%q)", m.Value)
	case Integer:
		return fmt.Sprintf("Integer(%d)", m.Value)
	case BulkString:
		if m.Value == nil {
			return "BulkString(null)"
		}
		return fmt.Sprintf("BulkString(%q)", m.Value)
	case Array:
		if m.Value == nil {
			return "Array(null)"
		}
		arr := m.Value.([]*Message)
		return fmt.Sprintf("Array(%d elements)", len(arr))
	default:
		return fmt.Sprintf("Unknown(%v)", m.Value)
	}
}

// IsNull returns true if the message represents a null value
func (m *Message) IsNull() bool {
	return m.Value == nil
}

// AsString returns the string value of the message or an error if not a string type
func (m *Message) AsString() (string, error) {
	switch m.Type {
	case SimpleString, Error:
		return m.Value.(string), nil
	case BulkString:
		if m.Value == nil {
			return "", fmt.Errorf("null bulk string")
		}
		return m.Value.(string), nil
	default:
		return "", fmt.Errorf("message type %s cannot be converted to string", m.Type)
	}
}

// AsInteger returns the integer value of the message or an error if not an integer
func (m *Message) AsInteger() (int64, error) {
	if m.Type != Integer {
		return 0, fmt.Errorf("message type %s cannot be converted to integer", m.Type)
	}
	return m.Value.(int64), nil
}

// AsArray returns the array value of the message or an error if not an array
func (m *Message) AsArray() ([]*Message, error) {
	if m.Type != Array {
		return nil, fmt.Errorf("message type %s cannot be converted to array", m.Type)
	}
	if m.Value == nil {
		return nil, fmt.Errorf("null array")
	}
	return m.Value.([]*Message), nil
}
