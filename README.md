# Redis-Lite

Redis-Lite is a simplified Redis server clone implemented in Go. The project aims to provide a lightweight, educational implementation of Redis core functionality while maintaining compatibility with the Redis protocol.

## Features

### Currently Implemented

- **RESP Protocol Support**: Full implementation of Redis Serialization Protocol (RESP) with support for all data types:
  - Simple Strings (`+OK\r\n`)
  - Errors (`-Error message\r\n`)
  - Integers (`:1000\r\n`)
  - Bulk Strings (`$6\r\nfoobar\r\n`)
  - Arrays (`*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n`)
  - Null values for both Bulk Strings and Arrays

### Planned Features

- TCP Server with concurrent client handling
- Core Redis commands (GET, SET, PING, ECHO, etc.)
- Key expiration (TTL) support
- In-memory data storage with thread safety
- Persistence to disk
- Background key cleanup

## Architecture

The project follows a layered architecture with clear separation of concerns:

- **Protocol Layer** (`pkg/resp/`): RESP message parsing and serialization
- **Network Layer** (`pkg/server/`): TCP server and connection handling
- **Command Layer** (`pkg/commands/`): Command routing and execution
- **Storage Layer** (`pkg/storage/`): In-memory data storage
- **Persistence Layer** (`pkg/persistence/`): Disk serialization

For detailed architecture information, see [docs/architecture.md](docs/architecture.md).

## Getting Started

### Prerequisites

- Go 1.21 or later

### Installation

```bash
git clone https://github.com/tsinivuo/redis-lite.git
cd redis-lite
go mod tidy
```

### Running Tests

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test ./... -v

# Run tests for specific package
go test ./pkg/resp/... -v
```

### Usage

*Note: The server implementation is still in progress. Currently, only the RESP protocol implementation is available.*

```go
package main

import (
    "fmt"
    "github.com/tsinivuo/redis-lite/pkg/resp"
)

func main() {
    // Parse a RESP message
    message, err := resp.ParseString("+OK\r\n")
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Parsed: %s\n", message.String())
    
    // Serialize a message
    msg := resp.NewSimpleString("PONG")
    serialized, err := resp.SerializeToString(msg)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Serialized: %q\n", serialized)
}
```

## RESP Protocol Implementation

The RESP (Redis Serialization Protocol) implementation is fully compatible with Redis protocol specification:

### Supported Data Types

1. **Simple Strings**: Single-line strings prefixed with `+`
2. **Errors**: Error messages prefixed with `-`
3. **Integers**: 64-bit signed integers prefixed with `:`
4. **Bulk Strings**: Binary-safe strings with length prefix `$`
5. **Arrays**: Collections of other RESP types prefixed with `*`

### Example Messages

```
+OK\r\n                    # Simple string "OK"
-Error message\r\n         # Error message
:1000\r\n                  # Integer 1000
$6\r\nfoobar\r\n          # Bulk string "foobar"
$-1\r\n                    # Null bulk string
*2\r\n$3\r\nfoo\r\n$3\r\nbar\r\n  # Array ["foo", "bar"]
*-1\r\n                    # Null array
```

### Error Handling

The implementation gracefully handles malformed messages and provides detailed error information:

- Invalid message type indicators
- Malformed length specifications
- Missing CRLF terminators
- Length mismatches
- Nested parsing errors

## Testing

The project includes comprehensive test coverage for all RESP protocol functionality:

- Unit tests for all message types
- Edge case testing (null values, empty strings, etc.)
- Malformed message handling
- Round-trip serialization/parsing verification
- Complex nested structure support

## Development

### Code Quality

The project follows Go best practices and coding standards:

- Simple, readable code with clear naming
- Comprehensive error handling
- Extensive unit test coverage
- Documentation for all public APIs
- Consistent code formatting with `gofmt`

### Contributing

1. Fork the repository
2. Create a feature branch
3. Implement your changes with tests
4. Ensure all tests pass
5. Submit a pull request

## License

This project is for educational purposes. See the LICENSE file for details.

## Status

🚧 **In Development**: This project is actively being developed. The RESP protocol implementation is complete and fully tested. Server implementation and Redis commands are planned for future releases.

### Current Implementation Status

- ✅ RESP Protocol (Complete)
- ⏳ TCP Server (Planned)
- ⏳ Core Commands (Planned)
- ⏳ Storage Layer (Planned)
- ⏳ Persistence (Planned)
- ⏳ Expiry Management (Planned)
