# Go MCP GraphQL Library

A Go library that introspects a GraphQL server and automatically generates MCP (Model Context Protocol) tools, allowing the MCP server to be hosted via HTTP.

## Features

- **GraphQL Introspection**: Automatically introspects any GraphQL server to understand its schema
- **MCP Tool Generation**: Converts GraphQL queries and mutations into MCP tools
- **HTTP Server**: Hosts the MCP server over HTTP for easy integration
- **Type Safety**: Leverages Go's type system for safe GraphQL operations
- **Flexible Transport**: Supports stdio, HTTP, and custom transport implementations

## Installation

```bash
go get github.com/peterbeamish/go-mcp-graphql
```

## Quick Start

### HTTP Server Example

```go
package main

import (
    "log"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // Create a GraphQL MCP server
    server, err := graphqlmcp.NewMCPGraphQLServer("https://countries.trevorblades.com/graphql")
    if err != nil {
        log.Fatal(err)
    }
    
    // Start HTTP server
    log.Println("Starting MCP GraphQL server on :8080")
    log.Fatal(graphqlmcp.StartHTTPServer(server, ":8080"))
}
```

### Stdio Server Example

```go
package main

import (
    "context"
    "log"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // Create a GraphQL MCP server
    server, err := graphqlmcp.NewMCPGraphQLServer("https://countries.trevorblades.com/graphql")
    if err != nil {
        log.Fatal(err)
    }
    
    // Run over stdio
    ctx := context.Background()
    log.Fatal(server.GetMCPServer().Run(ctx, &graphqlmcp.StdioTransport{}))
}
```

### HTTP Client Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // Create HTTP client
    client := graphqlmcp.CreateHTTPClient("http://localhost:8080")
    
    ctx := context.Background()
    
    // List available tools
    tools, err := client.ListTools(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Available tools: %+v\n", tools)
    
    // Call a tool
    response, err := client.CallTool(ctx, "query_countries", map[string]interface{}{
        "filter": map[string]interface{}{
            "continent": map[string]interface{}{
                "eq": "AF",
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Response: %+v\n", response)
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

- `example/main.go` - HTTP server with GitHub GraphQL API
- `example/simple/main.go` - Stdio server with countries API
- `example/client/main.go` - HTTP client example
- `example/gqlgen-server/` - Complete gqlgen-based GraphQL server
- `example/gqlgen-mcp-client/` - MCP client that introspects the gqlgen server
- `example/full-demo/` - Complete demo showing the full workflow
- `Makefile` - Comprehensive build and run commands

## Configuration

### Authentication

Set custom headers for GraphQL requests:

```go
server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql")
if err != nil {
    log.Fatal(err)
}

// Set authentication header
server.GetClient().SetHeader("Authorization", "Bearer YOUR_TOKEN")
```

### Custom Timeouts

```go
httpServer := graphqlmcp.NewHTTPServer(server.GetMCPServer())
httpServer.SetTimeout(60 * time.Second)
```

## Gqlgen Integration Example

The library includes a complete example using [gqlgen](https://github.com/99designs/gqlgen) to demonstrate the full workflow:

### Quick Start with Gqlgen

1. **Run the complete demo:**
   ```bash
   make demo
   ```

2. **Or run step by step:**
   ```bash
   # Terminal 1: Start GraphQL server
   make run-graphql
   
   # Terminal 2: Start MCP client
   make run-mcp
   
   # Terminal 3: Run full demo
   make run-full-demo
   ```

3. **Or use the Makefile for other operations:**
   ```bash
   make help              # Show all available commands
   make install           # Install dependencies
   make install-tools     # Install development tools
   make generate          # Generate gqlgen code
   make build             # Build all components
   make test              # Test all services
   make format            # Format all code
   make quality-all       # Run all quality checks
   make clean             # Clean build artifacts
   make status            # Check service status
   ```

### What the Gqlgen Example Includes

- **Complete GraphQL Server**: A blog API with users and posts using gqlgen
- **Schema Introspection**: Automatic discovery of queries and mutations
- **MCP Tool Generation**: Each GraphQL operation becomes an MCP tool
- **HTTP Integration**: MCP server accessible via HTTP endpoints
- **Full Demo**: End-to-end demonstration of the workflow

### Testing the Generated Tools

Once running, you can test the MCP tools:

```bash
# List available tools
curl http://localhost:8080/tools

# Test a query tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "query_posts", "arguments": {}}}'

# Test a mutation tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "mutation_createUser", "arguments": {"input": {"name": "Test User", "email": "test@example.com"}}}}'
```

## Development Tools

This project follows modern Go patterns for managing build tools and development dependencies:

### Tools Management
- **go.mod**: Contains all tool dependencies with proper versioning
- **tools.go**: Manages build tools for the gqlgen server example
- **Makefile**: Provides convenient commands for all development tasks

### Available Tools
- **Code Generation**: gqlgen for GraphQL code generation
- **Linting**: golangci-lint for comprehensive code analysis
- **Security**: gosec for security vulnerability scanning
- **Quality**: errcheck, ineffassign for code quality checks
- **Formatting**: goimports, gofmt, gci for code formatting
- **Testing**: testify for testing framework, mockery for mock generation
- **Development**: air for auto-reload during development

### Quick Development Setup
```bash
# Install all development tools
make install-tools

# Format and check code quality
make quality-all

# Run in development mode with auto-reload
make dev
```

## How It Works

1. **Introspection**: The library performs GraphQL introspection to understand the server's schema
2. **Tool Generation**: Each GraphQL query and mutation becomes an MCP tool
3. **Type Mapping**: GraphQL types are mapped to JSON Schema types for MCP tools
4. **HTTP Wrapper**: The MCP server is wrapped with HTTP transport for web access

## License

MIT