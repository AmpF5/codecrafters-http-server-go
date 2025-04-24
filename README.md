[![progress-banner](https://backend.codecrafters.io/progress/http-server/a646075d-b990-46cb-8ac9-3c2299709139)](https://app.codecrafters.io/users/codecrafters-bot?r=2qF)

# HTTP Server Implementation

A lightweight but feature-rich HTTP/1.1 server written in Go. This server was developed as part of the CodeCrafters "Build Your Own HTTP Server" challenge.

## Features

- Handles multiple concurrent client connections using goroutines
- Supports GET and POST HTTP methods
- Echo functionality that returns the provided content
- User-agent detection
- File operations (reading and writing)
- HTTP status code handling (200 OK, 201 Created, 404 Not Found)
- Content-type support (text/plain, application/octet-stream)

## Implementation Details

The server uses Go's native networking capabilities to create a TCP server that listens on port 4221. When a connection is received, it parses the HTTP request, extracts method, URL, headers, and request body when applicable.

Key components include:
- Request parsing
- Header management
- File operations
- Response building with appropriate status codes
- Concurrent connection handling

## Running the Server

```sh
./your_program.sh
```

This will start the server on port 4221, ready to accept incoming HTTP requests.

## Development

The main implementation is in `app/main.go`. To make changes:

1. Ensure you have Go 1.24 installed
2. Modify the code in `app/main.go`
3. Run the server using the provided script
4. Test with any HTTP client (curl, browser, etc.)

## License

This project is part of the CodeCrafters challenge and is intended for educational purposes.
