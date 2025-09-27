package main

import (
	"fmt"
	"log"
	"time"

	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
	// Wait a moment for the GraphQL server to start
	fmt.Println("Waiting for GraphQL server to start...")
	time.Sleep(2 * time.Second)

	// Create the MCP GraphQL server that connects to our gqlgen server
	server, err := graphqlmcp.NewMCPGraphQLServer("http://localhost:8081/graphql")
	if err != nil {
		log.Fatalf("Failed to create MCP GraphQL server: %v", err)
	}

	// Display the introspected schema
	fmt.Println("\n=== Introspected GraphQL Schema ===")
	schema := server.GetSchema()

	fmt.Printf("Query Type: %s\n", schema.QueryType.Name)
	fmt.Printf("Mutation Type: %s\n", schema.MutationType.Name)

	fmt.Println("\nAvailable Queries:")
	for _, query := range schema.GetQueries() {
		fmt.Printf("  - %s: %s\n", query.Name, query.Description)
	}

	fmt.Println("\nAvailable Mutations:")
	for _, mutation := range schema.GetMutations() {
		fmt.Printf("  - %s: %s\n", mutation.Name, mutation.Description)
	}

	// Start the MCP server over HTTP
	fmt.Println("\n=== Starting MCP Server ===")
	fmt.Println("MCP server will be available at: http://localhost:8080/mcp")
	fmt.Println("Health check: http://localhost:8080/health")
	fmt.Println("Schema endpoint: http://localhost:8080/schema")
	fmt.Println("Tools endpoint: http://localhost:8080/tools")

	if err := graphqlmcp.StartHTTPServer(server, ":8080"); err != nil {
		log.Fatalf("Failed to start MCP server: %v", err)
	}
}
