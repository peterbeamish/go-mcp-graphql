# Quick Start Guide

Get up and running with the Go MCP GraphQL library in minutes.

## Installation

```bash
go get github.com/peterbeamish/go-mcp-graphql
```

## Basic HTTP Server

The simplest way to get started is with an HTTP server:

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

## HTTP Client Example

Connect to an MCP server and use the tools:

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

## Running the Complete Demo

The project includes a complete working example with a GraphQL server:

```bash
# Run the complete demo (starts GraphQL server, MCP server, and client)
make demo

# Or run step by step:
make run-graphql  # Terminal 1: Start GraphQL server
make run-mcp      # Terminal 2: Start MCP server  
make run-full-demo # Terminal 3: Run client demo
```

## Testing Your Setup

Once running, test the MCP tools:

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

## Next Steps

- [Configuration Options](config.md) - Learn about authentication, timeouts, and masking
- [Query Generation](query-generation.md) - Understand how GraphQL operations are converted to tools
- [Examples](examples.md) - See advanced usage patterns and real-world examples
