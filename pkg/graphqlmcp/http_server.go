package graphqlmcp

import (
	"encoding/json"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// HTTP handler functions for MCP GraphQL server endpoints
// These functions provide individual handlers that can be registered on any http.ServeMux

// GetMCPHandler returns the MCP endpoint handler using the MCP SDK's StreamableHTTPHandler
// which handles all the HTTP transport details including SSE support
func GetMCPHandler(server *MCPGraphQLServer) http.Handler {
	return mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return server.GetMCPServer()
		},
		nil,
	)
}

// GetHealthHandler returns a health check endpoint handler
func GetHealthHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "graphql-mcp-server",
		})
	}
}

// GetSchemaHandler returns a schema endpoint handler to view the GraphQL schema
func GetSchemaHandler(server *MCPGraphQLServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		schema := server.GetSchema()
		response := map[string]interface{}{
			"schema": schema,
			"sdl":    schema.GetSchemaSDL(),
		}
		json.NewEncoder(w).Encode(response)
	}
}

// GetToolsHandler returns a tools endpoint handler to list available MCP tools
func GetToolsHandler(server *MCPGraphQLServer) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get available tools from the schema
		queries := server.GetSchema().GetQueries()
		mutations := server.GetSchema().GetMutations()

		tools := make([]map[string]interface{}, 0, len(queries)+len(mutations))

		// Add query tools
		for _, query := range queries {
			inputSchema := server.createInputSchema(query)
			tools = append(tools, map[string]interface{}{
				"name":        "query_" + query.Name,
				"description": query.Description,
				"type":        "query",
				"inputSchema": inputSchema,
			})
		}

		// Add mutation tools
		for _, mutation := range mutations {
			inputSchema := server.createInputSchema(mutation)
			tools = append(tools, map[string]interface{}{
				"name":        "mutation_" + mutation.Name,
				"description": mutation.Description,
				"type":        "mutation",
				"inputSchema": inputSchema,
			})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"tools": tools,
			"count": len(tools),
		})
	}
}

// GetCompleteMux returns a complete http.ServeMux with all MCP GraphQL server endpoints
// This is a convenience function that registers all handlers on a new mux
func GetCompleteMux(server *MCPGraphQLServer) *http.ServeMux {
	mux := http.NewServeMux()

	// Register all handlers
	mux.Handle("/mcp", GetMCPHandler(server))
	mux.HandleFunc("/health", GetHealthHandler())
	mux.HandleFunc("/schema", GetSchemaHandler(server))
	mux.HandleFunc("/tools", GetToolsHandler(server))

	return mux
}
