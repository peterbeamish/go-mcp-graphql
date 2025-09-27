package main

import (
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server/graphql"
	"github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server/resolver"
	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

var (
	mcpServerOnce sync.Once
	mcpServer     *graphqlmcp.MCPGraphQLServer
	mcpServerErr  error
)

func getOrCreateMCPServer() (*graphqlmcp.MCPGraphQLServer, error) {
	mcpServerOnce.Do(func() {
		log.Println("Creating MCP server...")
		graphqlURL := "http://localhost:8080/query"
		mcpServer, mcpServerErr = graphqlmcp.NewMCPGraphQLServer(graphqlURL)
		if mcpServerErr != nil {
			log.Printf("Failed to create MCP server: %v", mcpServerErr)
		} else {
			log.Println("MCP server created successfully")
		}
	})
	return mcpServer, mcpServerErr
}

func main() {
	// Create GraphQL server
	graphqlResolver := resolver.NewResolver()
	graphqlServer := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: graphqlResolver}))

	// Start GraphQL server in background
	go func() {
		log.Println("Starting GraphQL server on :8080...")

		// Add playground for testing
		http.Handle("/", playground.Handler("GraphQL playground", "/query"))
		http.Handle("/query", graphqlServer)
		http.Handle("/graphql", graphqlServer)

		log.Println("📊 GraphQL Playground: http://localhost:8080")
		log.Println("🔍 GraphQL Endpoint: http://localhost:8080/query")
		log.Println("📋 Introspection: http://localhost:8080/graphql")

		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatalf("GraphQL server failed: %v", err)
		}
	}()

	// Wait a moment for GraphQL server to start
	time.Sleep(2 * time.Second)

	// Create MCP server
	server, err := getOrCreateMCPServer()
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Start MCP HTTP server
	log.Println("🚀 Starting MCP server on :8081")
	log.Println("🤖 MCP Server: http://localhost:8081/mcp")
	log.Println("🛠️  Tools List: http://localhost:8081/tools")
	log.Println("📋 Schema Info: http://localhost:8081/schema")
	log.Println("❤️  Health Check: http://localhost:8081/health")
	log.Println("")
	log.Println("Available MCP tools:")
	log.Println("  - MCP tools are available at http://localhost:8081/mcp")
	log.Println("  - Uses Server-Sent Events (SSE) for streaming")
	log.Println("  - Uses http_server.go StartHTTPServer pattern")
	log.Println("")

	// Start MCP server
	if err := graphqlmcp.StartHTTPServer(server, ":8081"); err != nil {
		log.Fatalf("MCP server failed: %v", err)
	}
}
