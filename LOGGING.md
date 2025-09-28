# Logging in MCP GraphQL Server

This document describes the comprehensive logging system implemented in the MCP GraphQL server.

## Overview

The MCP GraphQL server now includes structured logging using Go's built-in `slog` package. This provides:

- **Structured logging** with key-value pairs for easy parsing and analysis
- **Request ID tracking** for correlating logs across different components
- **Performance metrics** including request duration and response sizes
- **Error context** with detailed error information
- **Configurable log levels** for different environments

## Logging Components

### 1. Tool Call Logging

When MCP tools are called, the following events are logged:

- **Tool call initiation**: When a tool call starts
- **Query generation**: When GraphQL queries/mutations are generated
- **Execution timing**: How long the GraphQL operation takes
- **Success/failure**: Whether the operation succeeded or failed
- **Response details**: Size and error information

### 2. HTTP Client Logging

For HTTP-based tool calls, additional logging includes:

- **HTTP request details**: URL, method, request size
- **HTTP response details**: Status code, response size, duration
- **Network errors**: Connection failures, timeouts, etc.

### 3. GraphQL Client Logging

The underlying GraphQL client logs:

- **Query execution**: What queries are being executed
- **Introspection**: Schema introspection operations
- **Network communication**: HTTP requests to GraphQL endpoints
- **Response parsing**: Success/failure of response parsing

## Configuration

### Basic Setup

```go
import (
    "log/slog"
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

// Configure verbose logging for development
logger := graphqlmcp.ConfigureVerboseLogging()

// Configure production logging (JSON format)
logger := graphqlmcp.ConfigureProductionLogging()

// Configure custom logging
logger := graphqlmcp.ConfigureLogging(slog.LevelInfo, true) // JSON format
logger := graphqlmcp.ConfigureLogging(slog.LevelDebug, false) // Text format
```

### Setting Loggers on Components

```go
// Create MCP server
server, err := graphqlmcp.NewMCPGraphQLServer("http://localhost:8080/query")
if err != nil {
    log.Fatal(err)
}

// Set custom logger
server.SetLogger(logger)

// For HTTP clients
client := graphqlmcp.CreateHTTPClient("http://localhost:8081")
client.SetLogger(logger)
```

## Log Levels

- **DEBUG**: Detailed information for debugging (query details, request/response bodies)
- **INFO**: General information about operations (tool calls, completions)
- **ERROR**: Error conditions (failures, network errors, parsing errors)

## Log Format Examples

### Text Format (Development)

```
time=2024-01-01T12:00:00Z level=INFO msg="Tool call initiated" request_id=req_1234567890 operation_type=query field_name=getUser input_args=2 input_values=map[id:123]
time=2024-01-01T12:00:00Z level=DEBUG msg="Generated GraphQL operation" request_id=req_1234567890 operation_type=query field_name=getUser query="query GetUser($id: ID!) { getUser(id: $id) { id name email } }"
time=2024-01-01T12:00:01Z level=INFO msg="Tool call completed successfully" request_id=req_1234567890 operation_type=query field_name=getUser duration_ms=150 response_size_bytes=245
```

### JSON Format (Production)

```json
{
  "time": "2024-01-01T12:00:00Z",
  "level": "INFO",
  "msg": "Tool call initiated",
  "request_id": "req_1234567890",
  "operation_type": "query",
  "field_name": "getUser",
  "input_args": 2,
  "input_values": {"id": "123"}
}
{
  "time": "2024-01-01T12:00:01Z",
  "level": "INFO",
  "msg": "Tool call completed successfully",
  "request_id": "req_1234567890",
  "operation_type": "query",
  "field_name": "getUser",
  "duration_ms": 150,
  "response_size_bytes": 245
}
```

## Request ID Tracking

Every operation gets a unique request ID that allows you to:

- Track a single tool call through all components
- Correlate logs across different services
- Debug issues by following the request flow

Request IDs are generated using timestamps and are unique across the application.

## Performance Metrics

The logging system tracks:

- **Duration**: How long operations take (in milliseconds)
- **Response sizes**: Size of responses in bytes
- **Request sizes**: Size of requests in bytes
- **Error counts**: Number of errors in responses

## Error Logging

Errors are logged with full context:

- **Error type**: What kind of error occurred
- **Error message**: The actual error message
- **Request context**: Which operation failed
- **Timing information**: When the error occurred
- **Additional context**: Relevant variables and state

## Example Usage

See the `example/full-demo/main.go` file for a complete example of how to set up logging in your application.

## Testing

Run the logging tests to verify functionality:

```bash
go test ./pkg/graphqlmcp -v -run TestConfigure
```

## Best Practices

1. **Use appropriate log levels**: DEBUG for development, INFO for production
2. **Include request IDs**: Always include request IDs in your own logs
3. **Log errors with context**: Include relevant context when logging errors
4. **Use structured logging**: Prefer key-value pairs over formatted strings
5. **Configure for environment**: Use JSON in production, text in development
