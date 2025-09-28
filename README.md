# Go MCP GraphQL Library

A Go library that introspects a GraphQL server and automatically generates MCP (Model Context Protocol) tools, allowing the MCP server to be hosted via HTTP.

## Table of Contents

- [Quick Start](docs/quickstart.md) - Get up and running in minutes
- [Configuration](docs/config.md) - Authentication, logging, and filtering options
- [Query Generation](docs/query-generation.md) - How GraphQL operations become MCP tools
- [Examples](docs/examples.md) - Real-world usage patterns and use cases

## What is this?

This library bridges GraphQL APIs with AI chat sessions by creating an MCP (Model Context Protocol) server that exposes GraphQL operations as tools that AI assistants can use.

### How it works in practice
1. **You** run the MCP server pointing to your GraphQL API
2. **AI assistants** (like Claude, ChatGPT, etc.) connect to your MCP server
3. **Users** can ask the AI to "get all users" or "create a new post" 
4. **The AI** automatically calls the right GraphQL operations through the MCP tools
5. **Your GraphQL API** receives the requests and returns data
6. **The AI** presents the results to the user in a conversational way

## Features

- **GraphQL Introspection**: Automatically introspects any GraphQL server to understand its schema
- **MCP Tool Generation**: Converts GraphQL queries and mutations into MCP tools
- **HTTP Server**: Hosts the MCP server over HTTP for easy integration
- **Type Safety**: Leverages Go's type system for safe GraphQL operations
- **Options Pattern**: Configurable with `WithLogger()`, `WithMask()`, and `WithPassthruHeaders()` options
- **Header Passthrough**: Automatically pass authentication and tracing headers from MCP requests to GraphQL requests
- **Advanced GraphQL Support**: Unions, interfaces, complex types, and full schema introspection

## Installation

```bash
go get github.com/peterbeamish/go-mcp-graphql
```

## Quick Start

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // Create a GraphQL MCP server
    server, err := graphqlmcp.NewMCPGraphQLServer("https://countries.trevorblades.com/graphql")
    if err != nil {
        log.Fatal(err)
    }
    
    // Create HTTP server with all MCP endpoints
    mux := graphqlmcp.GetCompleteMux(server)
    
    // Start HTTP server
    log.Println("Starting MCP GraphQL server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### With Authentication (Passthru Headers)

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // Create a GraphQL MCP server with authentication passthrough
    server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
        graphqlmcp.WithPassthruHeaders([]string{
            "Authorization",  // Bearer tokens
            "X-User-ID",      // User identification  
            "X-Request-ID",   // Request tracing
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create HTTP server with all MCP endpoints
    mux := graphqlmcp.GetCompleteMux(server)
    
    // Start HTTP server
    log.Println("Starting authenticated MCP GraphQL server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## API Endpoints

When running as an HTTP server, the following endpoints are available:

- `POST /mcp` - Main MCP protocol endpoint
- `GET /health` - Health check endpoint
- `GET /schema` - View the introspected GraphQL schema
- `GET /tools` - List all available MCP tools

## Examples

See the `example/` directory for complete working examples:

- `example/client/main.go` - HTTP client example
- `example/gqlgen-server/` - Complete gqlgen-based GraphQL server with advanced features
- `example/full-demo/` - Complete demo showing the full workflow

### Running the Complete Demo

```bash
# Run the complete demo (standalone application)
make run-demo
```

The `run-full-demo` application is a standalone demo that runs everything in one process - you don't need to run separate services.

## Development

```bash
# Install dependencies
make install

# Format and check code quality
make format
make lint

# Run tests
make test
```

## License

MIT