package graphqlmcp

import (
	"encoding/json"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Note: We now use the MCP SDK's built-in StreamableHTTPHandler
// which handles all the HTTP transport details including SSE support

// StartHTTPServer starts an HTTP server with the MCP GraphQL server
func GetMux(server *MCPGraphQLServer) *http.ServeMux {
	// Create a mux for routing
	mux := http.NewServeMux()

	// MCP endpoint using the streamable handler
	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return server.GetMCPServer()
		},
		nil,
	)
	mux.Handle("/mcp", mcpHandler)

	// Add a health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "graphql-mcp-server",
		})
	})

	// Add a schema endpoint to view the GraphQL schema
	mux.HandleFunc("/schema", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		schema := server.GetSchema()
		response := map[string]interface{}{
			"schema": schema,
			"sdl":    schema.GetSchemaSDL(),
		}
		json.NewEncoder(w).Encode(response)
	})

	// Add a tools endpoint to list available MCP tools
	mux.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
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
	})

	return mux
}
