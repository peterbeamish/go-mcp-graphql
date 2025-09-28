# Examples and Use Cases

Comprehensive examples showing how to use the Go MCP GraphQL library in real-world scenarios.

## API Reference

### Core Types

#### `MCPGraphQLServer`
The main server type that provides GraphQL tools via MCP.

```go
type MCPGraphQLServer struct {
    executor  GraphQLExecutor
    Schema    *schema.Schema
    mcpServer *mcp.Server
    logger    *slog.Logger
    options   *MCPGraphQLServerOptions
}
```

#### `MCPGraphQLServerOptions`
Configuration options for the MCP server.

```go
type MCPGraphQLServerOptions struct {
    Logger          *slog.Logger
    Mask            *MaskConfig
    PassthruHeaders []string
}
```

### Constructor Functions

#### `NewMCPGraphQLServer(endpoint string, opts ...MCPGraphQLServerOption)`
Creates a new MCP GraphQL server with the given endpoint and options.

#### `NewMCPGraphQLServerWithExecutor(executor GraphQLExecutor, opts ...MCPGraphQLServerOption)`
Creates a new MCP GraphQL server with a custom executor and options.

### Option Functions

#### `WithLogger(logger *slog.Logger)`
Sets a custom logger for the MCP server.

#### `WithMask(allowList, blockList []string)`
Configures operation filtering by name or pattern.

#### `WithPassthruHeaders(headers []string)`
Configures which headers to pass through from MCP requests to GraphQL requests.

### HTTP Functions

#### `GetCompleteMux(server *MCPGraphQLServer)`
Returns an HTTP mux with all MCP endpoints.

#### `CreateHTTPClient(endpoint string)`
Creates an HTTP client for connecting to an MCP server.

## Basic Examples

### Simple HTTP Server

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

### HTTP Client

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

## Advanced Examples

### Authentication with Custom Headers

```go
package main

import (
    "log"
    "net/http"
    "os"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // Create authenticated GraphQL client
    client := graphqlmcp.NewGraphQLClient("https://api.github.com/graphql")
    client.SetHeader("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
    
    // Create MCP server with authenticated client
    server, err := graphqlmcp.NewMCPGraphQLServerWithExecutor(client)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create HTTP server
    mux := graphqlmcp.GetCompleteMux(server)
    
    log.Println("Starting authenticated MCP GraphQL server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### Dynamic Authentication with Passthru Headers

This example shows how to use passthru headers to dynamically pass authentication from MCP clients to the GraphQL server:

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // Create MCP server with passthru headers for dynamic authentication
    server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
        graphqlmcp.WithPassthruHeaders([]string{
            "Authorization",    // Bearer tokens from clients
            "X-User-ID",        // User identification
            "X-Request-ID",     // Request tracing
            "X-Tenant-ID",      // Multi-tenant support
        }),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create HTTP server
    mux := graphqlmcp.GetCompleteMux(server)
    
    log.Println("Starting MCP GraphQL server with dynamic authentication on :8080")
    log.Println("Clients can now send their own authentication headers!")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

**Client Usage:**
```bash
# Different clients can use different authentication methods
curl -H "Authorization: Bearer client1-token" \
     -H "X-User-ID: user123" \
     -X POST http://localhost:8080/mcp \
     -d '{"method":"tools/call","params":{"name":"query_profile","arguments":{}}}'

curl -H "Authorization: Bearer client2-token" \
     -H "X-User-ID: user456" \
     -H "X-Tenant-ID: tenant-abc" \
     -X POST http://localhost:8080/mcp \
     -d '{"method":"tools/call","params":{"name":"query_orders","arguments":{}}}'
```

### Custom Logging and Filtering

```go
package main

import (
    "log"
    "log/slog"
    "net/http"
    "os"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // Create custom logger
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelDebug,
    }))
    
    // Create server with custom logger and operation filtering
    server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
        graphqlmcp.WithLogger(logger),
        graphqlmcp.WithMask(
            []string{"^get.*", "^list.*"},  // Only allow get and list operations
            []string{".*admin.*", ".*delete.*"}, // Block admin and delete operations
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create HTTP server
    mux := graphqlmcp.GetCompleteMux(server)
    
    log.Println("Starting filtered MCP GraphQL server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

### Multiple GraphQL Endpoints

```go
package main

import (
    "log"
    "net/http"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // Create servers for different GraphQL endpoints
    userServer, err := graphqlmcp.NewMCPGraphQLServer("https://users.api.com/graphql",
        graphqlmcp.WithMask([]string{"^get.*", "^list.*"}, nil),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    productServer, err := graphqlmcp.NewMCPGraphQLServer("https://products.api.com/graphql",
        graphqlmcp.WithMask([]string{"^query.*"}, []string{".*admin.*"}),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create combined HTTP server
    mux := http.NewServeMux()
    
    // User endpoints
    userMux := graphqlmcp.GetCompleteMux(userServer)
    mux.Handle("/users/", http.StripPrefix("/users", userMux))
    
    // Product endpoints
    productMux := graphqlmcp.GetCompleteMux(productServer)
    mux.Handle("/products/", http.StripPrefix("/products", productMux))
    
    log.Println("Starting multi-endpoint MCP server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

## Real-World Use Cases

### GitHub API Integration

```go
package main

import (
    "context"
    "fmt"
    "log"
    "os"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // GitHub API with authentication
    client := graphqlmcp.NewGraphQLClient("https://api.github.com/graphql")
    client.SetHeader("Authorization", "Bearer "+os.Getenv("GITHUB_TOKEN"))
    
    server, err := graphqlmcp.NewMCPGraphQLServerWithExecutor(client)
    if err != nil {
        log.Fatal(err)
    }
    
    // Create HTTP client to test
    mcpClient := graphqlmcp.CreateHTTPClient("http://localhost:8080")
    
    ctx := context.Background()
    
    // List available GitHub tools
    tools, err := mcpClient.ListTools(ctx)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("GitHub tools available: %d\n", len(tools))
    
    // Example: Get user information
    response, err := mcpClient.CallTool(ctx, "query_user", map[string]interface{}{
        "login": "octocat",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("User data: %+v\n", response)
}
```

### E-commerce API Integration

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // E-commerce API with filtering
    server, err := graphqlmcp.NewMCPGraphQLServer("https://shop.example.com/graphql",
        graphqlmcp.WithMask(
            []string{"^get.*", "^list.*", "^search.*"}, // Allow read operations
            []string{".*admin.*", ".*delete.*", ".*update.*"}, // Block admin and write operations
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create HTTP server
    mux := graphqlmcp.GetCompleteMux(server)
    
    go func() {
        log.Println("Starting e-commerce MCP server on :8080")
        log.Fatal(http.ListenAndServe(":8080", mux))
    }()
    
    // Test the integration
    mcpClient := graphqlmcp.CreateHTTPClient("http://localhost:8080")
    ctx := context.Background()
    
    // Search for products
    response, err := mcpClient.CallTool(ctx, "query_searchProducts", map[string]interface{}{
        "query": "laptop",
        "filters": map[string]interface{}{
            "category": "electronics",
            "priceRange": map[string]interface{}{
                "min": 500,
                "max": 2000,
            },
        },
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Search results: %+v\n", response)
}
```

### Content Management System

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // CMS API with role-based filtering
    server, err := graphqlmcp.NewMCPGraphQLServer("https://cms.example.com/graphql",
        graphqlmcp.WithMask(
            []string{"^get.*", "^list.*", "^query.*"}, // Allow all read operations
            []string{".*admin.*", ".*delete.*"}, // Block admin and delete operations
        ),
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create HTTP server
    mux := graphqlmcp.GetCompleteMux(server)
    
    go func() {
        log.Println("Starting CMS MCP server on :8080")
        log.Fatal(http.ListenAndServe(":8080", mux))
    }()
    
    // Test content retrieval
    mcpClient := graphqlmcp.CreateHTTPClient("http://localhost:8080")
    ctx := context.Background()
    
    // Get published articles
    response, err := mcpClient.CallTool(ctx, "query_articles", map[string]interface{}{
        "status": "PUBLISHED",
        "limit": 10,
        "orderBy": "publishedAt",
    })
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Articles: %+v\n", response)
}
```

## Testing Examples

### Unit Testing

```go
package main

import (
    "context"
    "testing"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func TestMCPClient(t *testing.T) {
    // Create test server
    server, err := graphqlmcp.NewMCPGraphQLServer("https://countries.trevorblades.com/graphql")
    if err != nil {
        t.Fatal(err)
    }
    
    // Create test client
    client := graphqlmcp.CreateHTTPClient("http://localhost:8080")
    ctx := context.Background()
    
    // Test listing tools
    tools, err := client.ListTools(ctx)
    if err != nil {
        t.Fatal(err)
    }
    
    if len(tools) == 0 {
        t.Error("Expected tools to be available")
    }
    
    // Test calling a tool
    response, err := client.CallTool(ctx, "query_countries", map[string]interface{}{})
    if err != nil {
        t.Fatal(err)
    }
    
    if response.IsError {
        t.Error("Expected successful response")
    }
}
```

### Integration Testing

```go
package main

import (
    "context"
    "net/http"
    "testing"
    "time"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func TestIntegration(t *testing.T) {
    // Start test server
    server, err := graphqlmcp.NewMCPGraphQLServer("https://countries.trevorblades.com/graphql")
    if err != nil {
        t.Fatal(err)
    }
    
    mux := graphqlmcp.GetCompleteMux(server)
    httpServer := &http.Server{
        Addr:    ":8081",
        Handler: mux,
    }
    
    go func() {
        httpServer.ListenAndServe()
    }()
    
    // Wait for server to start
    time.Sleep(100 * time.Millisecond)
    
    // Test client
    client := graphqlmcp.CreateHTTPClient("http://localhost:8081")
    ctx := context.Background()
    
    // Test health endpoint
    resp, err := http.Get("http://localhost:8081/health")
    if err != nil {
        t.Fatal(err)
    }
    resp.Body.Close()
    
    if resp.StatusCode != 200 {
        t.Error("Expected health check to return 200")
    }
    
    // Test MCP tools
    tools, err := client.ListTools(ctx)
    if err != nil {
        t.Fatal(err)
    }
    
    if len(tools) == 0 {
        t.Error("Expected tools to be available")
    }
    
    // Cleanup
    httpServer.Shutdown(context.Background())
}
```

## Performance Examples

### High-Throughput Server

```go
package main

import (
    "log"
    "net/http"
    "time"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    // Create server with custom configuration
    server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
        graphqlmcp.WithLogger(nil), // Use default logger for performance
    )
    if err != nil {
        log.Fatal(err)
    }
    
    // Create HTTP server with performance tuning
    mux := graphqlmcp.GetCompleteMux(server)
    
    httpServer := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    
    log.Println("Starting high-performance MCP server on :8080")
    log.Fatal(httpServer.ListenAndServe())
}
```

### Load Testing

```go
package main

import (
    "context"
    "fmt"
    "log"
    "sync"
    "time"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    client := graphqlmcp.CreateHTTPClient("http://localhost:8080")
    ctx := context.Background()
    
    // Load test parameters
    numGoroutines := 10
    requestsPerGoroutine := 100
    
    var wg sync.WaitGroup
    start := time.Now()
    
    for i := 0; i < numGoroutines; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            
            for j := 0; j < requestsPerGoroutine; j++ {
                _, err := client.CallTool(ctx, "query_countries", map[string]interface{}{})
                if err != nil {
                    log.Printf("Request failed: %v", err)
                }
            }
        }()
    }
    
    wg.Wait()
    duration := time.Since(start)
    
    totalRequests := numGoroutines * requestsPerGoroutine
    fmt.Printf("Completed %d requests in %v\n", totalRequests, duration)
    fmt.Printf("Requests per second: %.2f\n", float64(totalRequests)/duration.Seconds())
}
```

## Error Handling Examples

### Graceful Error Handling

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
    client := graphqlmcp.CreateHTTPClient("http://localhost:8080")
    ctx := context.Background()
    
    // List tools with error handling
    tools, err := client.ListTools(ctx)
    if err != nil {
        log.Printf("Failed to list tools: %v", err)
        return
    }
    
    fmt.Printf("Available tools: %d\n", len(tools))
    
    // Call tool with error handling
    response, err := client.CallTool(ctx, "query_countries", map[string]interface{}{})
    if err != nil {
        log.Printf("Failed to call tool: %v", err)
        return
    }
    
    if response.IsError {
        log.Printf("Tool returned error: %v", response.Content)
        return
    }
    
    fmt.Printf("Success: %v\n", response.Content)
}
```

### Retry Logic

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"
    
    "github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func callToolWithRetry(client *graphqlmcp.HTTPMCPClient, ctx context.Context, toolName string, args map[string]interface{}, maxRetries int) (*mcp.CallToolResult, error) {
    var lastErr error
    
    for i := 0; i < maxRetries; i++ {
        response, err := client.CallTool(ctx, toolName, args)
        if err == nil {
            return response, nil
        }
        
        lastErr = err
        log.Printf("Attempt %d failed: %v", i+1, err)
        
        if i < maxRetries-1 {
            time.Sleep(time.Duration(i+1) * time.Second)
        }
    }
    
    return nil, fmt.Errorf("failed after %d attempts: %w", maxRetries, lastErr)
}

func main() {
    client := graphqlmcp.CreateHTTPClient("http://localhost:8080")
    ctx := context.Background()
    
    response, err := callToolWithRetry(client, ctx, "query_countries", map[string]interface{}{}, 3)
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Response: %+v\n", response)
}
```
