# Go MCP GraphQL Library

A Go library that introspects a GraphQL server and automatically generates MCP (Model Context Protocol) tools, allowing the MCP server to be hosted via HTTP.

## What is this?

This library bridges GraphQL APIs with AI chat sessions by creating an MCP (Model Context Protocol) server that exposes GraphQL operations as tools that AI assistants can use.

### The MCP Server
The **MCP Server** acts as a bridge between AI chat sessions and your GraphQL API. It:
- Connects to any GraphQL server (like GitHub's API, your company's API, etc.)
- Automatically discovers all available queries and mutations through introspection
- Converts each GraphQL operation into an MCP "tool" that AI assistants can call
- Runs as an HTTP server that AI clients can connect to
- Handles authentication, parameter validation, and response formatting

### The MCP Client
The **MCP Client** is what AI assistants use to interact with your GraphQL API. It:
- Connects to the MCP server via HTTP
- Lists all available GraphQL operations as tools
- Calls specific tools (GraphQL queries/mutations) with parameters
- Receives structured responses that the AI can understand and work with

**When is the MCP Client used?**
The MCP Client is used by AI assistants (like Claude, ChatGPT, or other MCP-compatible tools) when they need to:
- **Query data**: "Show me all users", "Get the latest posts", "Find products by category"
- **Create resources**: "Create a new user", "Add a blog post", "Update a record"
- **Modify data**: "Update user profile", "Delete old records", "Change settings"
- **Analyze information**: "Summarize user activity", "Generate reports", "Find patterns"

The AI assistant automatically determines which GraphQL operations to call based on the user's natural language request, then uses the MCP Client to execute those operations against your GraphQL API.

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
    client := graphqlmcp.CreateHTTPClient("http://localhost:8081")
    
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

- `example/client/main.go` - HTTP client example
- `example/gqlgen-server/` - Complete gqlgen-based GraphQL server
- `example/full-demo/` - Complete demo showing the full workflow
- `Makefile` - Comprehensive build and run commands

## Configuration

### Authentication

Set custom headers for GraphQL requests by creating a custom GraphQL client:

```go
// Create a custom GraphQL client with authentication
client := graphqlmcp.NewGraphQLClient("https://api.example.com/graphql")
client.SetHeader("Authorization", "Bearer YOUR_TOKEN")

// Create MCP server with the authenticated client
server, err := graphqlmcp.NewMCPGraphQLServerWithExecutor(client)
if err != nil {
    log.Fatal(err)
}
```

### Custom Timeouts

```go
// Create HTTP client with custom timeout
client := &http.Client{
    Timeout: 60 * time.Second,
}

// Use with HTTPMCPClient
mcpClient := graphqlmcp.CreateHTTPClient("http://localhost:8081")
// Note: HTTPMCPClient uses a 30-second default timeout
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
   make generate          # Generate gqlgen code
   make build             # Build all components
   make test              # Test all services
   make format            # Format all code
   make lint              # Run linter
   make clean             # Clean build artifacts
   make rebuild           # Clean and rebuild all components
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
curl http://localhost:8081/tools

# Test a query tool
curl -X POST http://localhost:8081/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "query_posts", "arguments": {}}}'

# Test a mutation tool
curl -X POST http://localhost:8081/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "mutation_createUser", "arguments": {"input": {"name": "Test User", "email": "test@example.com"}}}}'
```

## Development Tools

This project follows modern Go patterns for managing build tools and development dependencies:

### Tools Management
- **go.mod**: Contains all tool dependencies with proper versioning
- **Makefile**: Provides convenient commands for all development tasks

### Available Tools
- **Code Generation**: gqlgen for GraphQL code generation
- **Testing**: testify for testing framework
- **Formatting**: gofmt for code formatting

### Quick Development Setup
```bash
# Install dependencies
make install

# Format and check code quality
make format
make lint

# Run tests
make test
```

## How It Works

1. **Introspection**: The library performs GraphQL introspection to understand the server's schema
2. **Tool Generation**: Each GraphQL query and mutation becomes an MCP tool
3. **Type Mapping**: GraphQL types are mapped to JSON Schema types for MCP tools
4. **HTTP Wrapper**: The MCP server is wrapped with HTTP transport for web access

## License

MIT