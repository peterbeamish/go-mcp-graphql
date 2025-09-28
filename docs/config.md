# Configuration Guide

Learn how to configure the MCP GraphQL server with authentication, timeouts, logging, and operation filtering.

## Basic Configuration

### Using Options Pattern

The library supports an options pattern for configuration:

```go
// With custom logger
server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
    graphqlmcp.WithLogger(customLogger),
)

// With operation filtering
server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
    graphqlmcp.WithMask(
        []string{"^get.*", "^list.*"},  // Allow list
        []string{".*admin.*", ".*delete.*"}, // Block list
    ),
)

// Combined options
server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
    graphqlmcp.WithLogger(logger),
    graphqlmcp.WithMask(allowList, blockList),
)
```

## Authentication

### Custom Headers

Set custom headers for GraphQL requests:

```go
// Create a custom GraphQL client with authentication
client := graphqlmcp.NewGraphQLClient("https://api.example.com/graphql")
client.SetHeader("Authorization", "Bearer YOUR_TOKEN")
client.SetHeader("X-API-Key", "your-api-key")

// Create MCP server with the authenticated client
server, err := graphqlmcp.NewMCPGraphQLServerWithExecutor(client)
if err != nil {
    log.Fatal(err)
}
```

### Multiple Headers

```go
client := graphqlmcp.NewGraphQLClient("https://api.example.com/graphql")
client.SetHeader("Authorization", "Bearer YOUR_TOKEN")
client.SetHeader("Content-Type", "application/json")
client.SetHeader("User-Agent", "MyApp/1.0")
```

## Logging

### Custom Logger

Use a custom structured logger:

```go
import (
    "log/slog"
    "os"
)

// Create custom logger
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
    graphqlmcp.WithLogger(logger),
)
```

### Log Levels

The logger supports different levels:
- `Debug`: Detailed operation information
- `Info`: General operation status
- `Error`: Error conditions

## Operation Filtering (Masking)

Control which GraphQL operations are exposed as MCP tools:

### Allow List Only

Only expose operations matching specific patterns:

```go
server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
    graphqlmcp.WithMask([]string{"^get.*", "^list.*"}, nil),
)
```

### Block List Only

Block specific operations while allowing others:

```go
server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
    graphqlmcp.WithMask(nil, []string{".*admin.*", ".*delete.*"}),
)
```

### Combined Filtering

Use both allow and block lists (block list takes precedence):

```go
server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
    graphqlmcp.WithMask(
        []string{"^get.*", "^list.*"},  // Allow list
        []string{".*admin.*", ".*delete.*"}, // Block list
    ),
)
```

### Pattern Examples

- `^get.*` - Operations starting with "get"
- `.*admin.*` - Operations containing "admin"
- `^list[A-Z].*` - Operations starting with "list" followed by uppercase letter
- `.*_internal$` - Operations ending with "_internal"

## Timeouts

### HTTP Client Timeouts

```go
import (
    "net/http"
    "time"
)

// Create HTTP client with custom timeout
client := &http.Client{
    Timeout: 60 * time.Second,
}

// Use with HTTPMCPClient
mcpClient := graphqlmcp.CreateHTTPClient("http://localhost:8081")
// Note: HTTPMCPClient uses a 30-second default timeout
```

### GraphQL Request Timeouts

The GraphQL client respects HTTP client timeouts:

```go
client := graphqlmcp.NewGraphQLClient("https://api.example.com/graphql")
// Timeout is controlled by the underlying HTTP client
```

## Environment Variables

### Common Configuration

```bash
# GraphQL endpoint
export GRAPHQL_ENDPOINT="https://api.example.com/graphql"

# Authentication
export GRAPHQL_TOKEN="your-bearer-token"
export GRAPHQL_API_KEY="your-api-key"

# Server configuration
export MCP_PORT="8080"
export LOG_LEVEL="debug"
```

### Using Environment Variables

```go
import "os"

endpoint := os.Getenv("GRAPHQL_ENDPOINT")
if endpoint == "" {
    endpoint = "https://countries.trevorblades.com/graphql"
}

client := graphqlmcp.NewGraphQLClient(endpoint)
if token := os.Getenv("GRAPHQL_TOKEN"); token != "" {
    client.SetHeader("Authorization", "Bearer "+token)
}

server, err := graphqlmcp.NewMCPGraphQLServer(endpoint)
```

## Advanced Configuration

### Custom GraphQL Client

For advanced configuration, create a custom GraphQL client:

```go
import (
    "net/http"
    "time"
)

// Custom HTTP client
httpClient := &http.Client{
    Timeout: 30 * time.Second,
    Transport: &http.Transport{
        MaxIdleConns:        100,
        IdleConnTimeout:     90 * time.Second,
        DisableCompression:  true,
    },
}

// Custom GraphQL client
client := graphqlmcp.NewGraphQLClientWithHTTPClient("https://api.example.com/graphql", httpClient)
client.SetHeader("Authorization", "Bearer YOUR_TOKEN")

// Create MCP server
server, err := graphqlmcp.NewMCPGraphQLServerWithExecutor(client)
```

### Multiple Endpoints

Handle multiple GraphQL endpoints:

```go
// Create multiple servers for different endpoints
userServer, err := graphqlmcp.NewMCPGraphQLServer("https://users.api.com/graphql",
    graphqlmcp.WithMask([]string{"^get.*", "^list.*"}, nil),
)

productServer, err := graphqlmcp.NewMCPGraphQLServer("https://products.api.com/graphql",
    graphqlmcp.WithMask([]string{"^query.*"}, []string{".*admin.*"}),
)
```

## Error Handling

### Graceful Degradation

```go
server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql")
if err != nil {
    log.Printf("Failed to create MCP server: %v", err)
    // Handle error appropriately
    return
}

// Server will log errors but continue operating
```

### Custom Error Handling

```go
// Custom logger with error handling
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelError,
}))

server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
    graphqlmcp.WithLogger(logger),
)
```

## Best Practices

1. **Use structured logging** for better debugging
2. **Set appropriate timeouts** for your use case
3. **Filter operations** to only expose what's needed
4. **Handle authentication** securely
5. **Monitor performance** and adjust timeouts accordingly
6. **Use environment variables** for configuration
7. **Test with different GraphQL schemas** to ensure compatibility
