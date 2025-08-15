package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// Parser handles parsing RESP messages from byte streams
type Parser struct {
	reader *bufio.Reader
}

// NewParser creates a new RESP parser with the given reader
func NewParser(reader io.Reader) *Parser {
	return &Parser{
		reader: bufio.NewReader(reader),
	}
}

// Parse parses a single RESP message from the input stream
func (p *Parser) Parse() (*Message, error) {
	// Read the first byte to determine the message type
	typeByte, err := p.reader.ReadByte()
	if err != nil {
		return nil, fmt.Errorf("failed to read message type: %w", err)
	}

	switch typeByte {
	case '+':
		return p.parseSimpleString()
	case '-':
		return p.parseError()
	case ':':
		return p.parseInteger()
	case '$':
		return p.parseBulkString()
	case '*':
		return p.parseArray()
	default:
		return nil, fmt.Errorf("invalid message type: %c", typeByte)
	}
}

// parseSimpleString parses a simple string message (+OK\r\n)
func (p *Parser) parseSimpleString() (*Message, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, fmt.Errorf("failed to read simple string: %w", err)
	}
	return NewSimpleString(line), nil
}

// parseError parses an error message (-Error message\r\n)
func (p *Parser) parseError() (*Message, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, fmt.Errorf("failed to read error message: %w", err)
	}
	return NewError(line), nil
}

// parseInteger parses an integer message (:1000\r\n)
func (p *Parser) parseInteger() (*Message, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, fmt.Errorf("failed to read integer: %w", err)
	}

	value, err := strconv.ParseInt(line, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("invalid integer format: %s", line)
	}

	return NewInteger(value), nil
}

// parseBulkString parses a bulk string message ($6\r\nfoobar\r\n or $-1\r\n for null)
func (p *Parser) parseBulkString() (*Message, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, fmt.Errorf("failed to read bulk string length: %w", err)
	}

	length, err := strconv.Atoi(line)
	if err != nil {
		return nil, fmt.Errorf("invalid bulk string length: %s", line)
	}

	// Handle null bulk string
	if length == -1 {
		return NewNullBulkString(), nil
	}

	if length < 0 {
		return nil, fmt.Errorf("invalid bulk string length: %d", length)
	}

	// Read the string data plus the trailing \r\n
	data := make([]byte, length)
	_, err = io.ReadFull(p.reader, data)
	if err != nil {
		return nil, fmt.Errorf("failed to read bulk string data: %w", err)
	}

	// Read and verify the trailing \r\n
	if err := p.expectCRLF(); err != nil {
		return nil, fmt.Errorf("missing CRLF after bulk string: %w", err)
	}

	return NewBulkString(string(data)), nil
}

// parseArray parses an array message (*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n or *-1\r\n for null)
func (p *Parser) parseArray() (*Message, error) {
	line, err := p.readLine()
	if err != nil {
		return nil, fmt.Errorf("failed to read array length: %w", err)
	}

	length, err := strconv.Atoi(line)
	if err != nil {
		return nil, fmt.Errorf("invalid array length: %s", line)
	}

	// Handle null array
	if length == -1 {
		return NewNullArray(), nil
	}

	if length < 0 {
		return nil, fmt.Errorf("invalid array length: %d", length)
	}

	// Parse each element in the array
	elements := make([]*Message, length)
	for i := 0; i < length; i++ {
		element, err := p.Parse()
		if err != nil {
			return nil, fmt.Errorf("failed to parse array element %d: %w", i, err)
		}
		elements[i] = element
	}

	return NewArray(elements), nil
}

// readLine reads a line terminated by \r\n and returns the content without the terminator
func (p *Parser) readLine() (string, error) {
	line, err := p.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	// Remove \r\n terminator
	if len(line) < 2 || line[len(line)-2:] != "\r\n" {
		return "", fmt.Errorf("line not terminated with CRLF: %q", line)
	}

	return line[:len(line)-2], nil
}

// expectCRLF reads and verifies that the next two bytes are \r\n
func (p *Parser) expectCRLF() error {
	cr, err := p.reader.ReadByte()
	if err != nil {
		return err
	}
	if cr != '\r' {
		return fmt.Errorf("expected \\r, got %c", cr)
	}

	lf, err := p.reader.ReadByte()
	if err != nil {
		return err
	}
	if lf != '\n' {
		return fmt.Errorf("expected \\n, got %c", lf)
	}

	return nil
}

// ParseString is a convenience function to parse a RESP message from a string
func ParseString(input string) (*Message, error) {
	parser := NewParser(strings.NewReader(input))
	return parser.Parse()
}
