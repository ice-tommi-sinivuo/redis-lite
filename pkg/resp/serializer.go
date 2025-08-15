package resp

import (
	"fmt"
	"io"
	"strings"
)

// Serializer handles serializing RESP messages to byte streams
type Serializer struct {
	writer io.Writer
}

// NewSerializer creates a new RESP serializer with the given writer
func NewSerializer(writer io.Writer) *Serializer {
	return &Serializer{
		writer: writer,
	}
}

// Serialize serializes a RESP message to the output stream
func (s *Serializer) Serialize(message *Message) error {
	switch message.Type {
	case SimpleString:
		return s.serializeSimpleString(message.Value.(string))
	case Error:
		return s.serializeError(message.Value.(string))
	case Integer:
		return s.serializeInteger(message.Value.(int64))
	case BulkString:
		return s.serializeBulkString(message.Value)
	case Array:
		return s.serializeArray(message.Value)
	default:
		return fmt.Errorf("unsupported message type: %s", message.Type)
	}
}

// serializeSimpleString serializes a simple string (+OK\r\n)
func (s *Serializer) serializeSimpleString(value string) error {
	// Simple strings cannot contain CR or LF
	if strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("simple string cannot contain CR or LF characters")
	}

	_, err := fmt.Fprintf(s.writer, "+%s\r\n", value)
	return err
}

// serializeError serializes an error message (-Error message\r\n)
func (s *Serializer) serializeError(value string) error {
	// Error messages cannot contain CR or LF
	if strings.ContainsAny(value, "\r\n") {
		return fmt.Errorf("error message cannot contain CR or LF characters")
	}

	_, err := fmt.Fprintf(s.writer, "-%s\r\n", value)
	return err
}

// serializeInteger serializes an integer (:1000\r\n)
func (s *Serializer) serializeInteger(value int64) error {
	_, err := fmt.Fprintf(s.writer, ":%d\r\n", value)
	return err
}

// serializeBulkString serializes a bulk string ($6\r\nfoobar\r\n or $-1\r\n for null)
func (s *Serializer) serializeBulkString(value interface{}) error {
	if value == nil {
		// Null bulk string
		_, err := s.writer.Write([]byte("$-1\r\n"))
		return err
	}

	str := value.(string)
	length := len(str)

	// Write length header
	_, err := fmt.Fprintf(s.writer, "$%d\r\n", length)
	if err != nil {
		return err
	}

	// Write the string data
	_, err = s.writer.Write([]byte(str))
	if err != nil {
		return err
	}

	// Write trailing CRLF
	_, err = s.writer.Write([]byte("\r\n"))
	return err
}

// serializeArray serializes an array (*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n or *-1\r\n for null)
func (s *Serializer) serializeArray(value interface{}) error {
	if value == nil {
		// Null array
		_, err := s.writer.Write([]byte("*-1\r\n"))
		return err
	}

	array := value.([]*Message)
	length := len(array)

	// Write length header
	_, err := fmt.Fprintf(s.writer, "*%d\r\n", length)
	if err != nil {
		return err
	}

	// Serialize each element
	for i, element := range array {
		if err := s.Serialize(element); err != nil {
			return fmt.Errorf("failed to serialize array element %d: %w", i, err)
		}
	}

	return nil
}

// SerializeToString is a convenience function to serialize a RESP message to a string
func SerializeToString(message *Message) (string, error) {
	var builder strings.Builder
	serializer := NewSerializer(&builder)

	err := serializer.Serialize(message)
	if err != nil {
		return "", err
	}

	return builder.String(), nil
}

// SerializeToBytes is a convenience function to serialize a RESP message to bytes
func SerializeToBytes(message *Message) ([]byte, error) {
	str, err := SerializeToString(message)
	if err != nil {
		return nil, err
	}
	return []byte(str), nil
}
