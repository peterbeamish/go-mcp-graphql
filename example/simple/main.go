package main

import (
	"context"
	"fmt"
	"log"

	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
	// Example with a simple GraphQL server
	// You can use any GraphQL server that supports introspection
	graphqlEndpoint := "https://countries.trevorblades.com/graphql"

	// Create the MCP GraphQL server
	server, err := graphqlmcp.NewMCPGraphQLServer(graphqlEndpoint)
	if err != nil {
		log.Fatalf("Failed to create MCP GraphQL server: %v", err)
	}

	// Get the underlying MCP server
	mcpServer := server.GetMCPServer()

	// Run the MCP server over stdio (for testing)
	fmt.Println("MCP GraphQL server running over stdio")
	fmt.Println("Available tools:")

	// List available tools
	queries := server.GetSchema().GetQueries()
	for _, query := range queries {
		fmt.Printf("  - query_%s: %s\n", query.Name, query.Description)
	}

	mutations := server.GetSchema().GetMutations()
	for _, mutation := range mutations {
		fmt.Printf("  - mutation_%s: %s\n", mutation.Name, mutation.Description)
	}

	// Connect to the server
	ctx := context.Background()
	conn, err := mcpServer.Connect(ctx, &graphqlmcp.StdioTransport{}, nil)
	if err != nil {
		log.Fatalf("Failed to connect MCP server: %v", err)
	}
	defer conn.Close()

	// Keep the connection alive
	select {}
}
