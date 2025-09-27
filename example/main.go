package main

import (
	"context"
	"fmt"
	"log"

	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
	// Example GraphQL endpoint (you can replace this with any GraphQL server)
	graphqlEndpoint := "https://api.github.com/graphql"

	// Create the MCP GraphQL server
	server, err := graphqlmcp.NewMCPGraphQLServer(graphqlEndpoint)
	if err != nil {
		log.Fatalf("Failed to create MCP GraphQL server: %v", err)
	}

	// Set authentication header for GitHub GraphQL API
	// Note: You'll need to provide a valid GitHub token
	server.GetClient().SetHeader("Authorization", "Bearer YOUR_GITHUB_TOKEN")

	// Start the HTTP server
	fmt.Println("Starting MCP GraphQL server on :8080")
	fmt.Println("Available endpoints:")
	fmt.Println("  - POST /mcp - MCP protocol endpoint")
	fmt.Println("  - GET /health - Health check")
	fmt.Println("  - GET /schema - View GraphQL schema")
	fmt.Println("  - GET /tools - List available MCP tools")

	if err := graphqlmcp.StartHTTPServer(server, ":8080"); err != nil {
		log.Fatalf("Failed to start HTTP server: %v", err)
	}
}

// Example of how to use the HTTP client
func exampleClientUsage() {
	// Create an HTTP client
	client := graphqlmcp.CreateHTTPClient("http://localhost:8080")

	ctx := context.Background()

	// List available tools
	tools, err := client.ListTools(ctx)
	if err != nil {
		log.Printf("Failed to list tools: %v", err)
		return
	}

	fmt.Printf("Available tools: %+v\n", tools)

	// Call a tool (example with GitHub GraphQL API)
	// Note: This is just an example - actual tool names depend on the GraphQL schema
	response, err := client.CallTool(ctx, "query_viewer", map[string]interface{}{
		"login": "octocat",
	})
	if err != nil {
		log.Printf("Failed to call tool: %v", err)
		return
	}

	fmt.Printf("Tool response: %+v\n", response)
}
