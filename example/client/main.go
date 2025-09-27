package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

func main() {
	// Create an HTTP client to communicate with the MCP server
	client := graphqlmcp.CreateHTTPClient("http://localhost:8080")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// List available tools
	fmt.Println("Fetching available tools...")
	tools, err := client.ListTools(ctx)
	if err != nil {
		log.Fatalf("Failed to list tools: %v", err)
	}

	fmt.Printf("Found %d available tools:\n", len(tools))
	for _, tool := range tools {
		fmt.Printf("  - %s (%s): %s\n",
			tool["name"],
			tool["type"],
			tool["description"])
	}

	// Example: Call a query tool (this depends on the GraphQL schema)
	// For the countries API, we might have a "countries" query
	if len(tools) > 0 {
		firstTool := tools[0]
		toolName := firstTool["name"].(string)

		fmt.Printf("\nCalling tool: %s\n", toolName)

		// Call the tool with some arguments
		// Note: The actual arguments depend on the GraphQL schema
		response, err := client.CallTool(ctx, toolName, map[string]interface{}{
			// Add appropriate arguments based on the tool
		})
		if err != nil {
			log.Printf("Failed to call tool %s: %v", toolName, err)
		} else {
			fmt.Printf("Tool response: %+v\n", response)
		}
	}
}
