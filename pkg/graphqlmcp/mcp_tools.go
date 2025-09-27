package graphqlmcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// MCPGraphQLServer represents an MCP server that provides GraphQL tools
type MCPGraphQLServer struct {
	client    *GraphQLClient
	schema    *Schema
	mcpServer *mcp.Server
}

// NewMCPGraphQLServer creates a new MCP GraphQL server
func NewMCPGraphQLServer(endpoint string) (*MCPGraphQLServer, error) {
	client := NewGraphQLClient(endpoint)

	// Introspect the schema
	ctx := context.Background()
	schema, err := client.IntrospectSchema(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to introspect GraphQL schema: %w", err)
	}

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "graphql-mcp-server",
		Version: "1.0.0",
	}, nil)

	server := &MCPGraphQLServer{
		client:    client,
		schema:    schema,
		mcpServer: mcpServer,
	}

	// Add tools for queries and mutations
	if err := server.addGraphQLTools(); err != nil {
		return nil, fmt.Errorf("failed to add GraphQL tools: %w", err)
	}

	return server, nil
}

// addGraphQLTools adds MCP tools for all GraphQL queries and mutations
func (s *MCPGraphQLServer) addGraphQLTools() error {
	// Add query tools
	queries := s.schema.GetQueries()
	for _, query := range queries {
		if err := s.addQueryTool(query); err != nil {
			return fmt.Errorf("failed to add query tool for %s: %w", query.Name, err)
		}
	}

	// Add mutation tools
	mutations := s.schema.GetMutations()
	for _, mutation := range mutations {
		if err := s.addMutationTool(mutation); err != nil {
			return fmt.Errorf("failed to add mutation tool for %s: %w", mutation.Name, err)
		}
	}

	return nil
}

// addQueryTool adds an MCP tool for a GraphQL query
func (s *MCPGraphQLServer) addQueryTool(query *Field) error {
	toolName := "query_" + query.Name
	toolDescription := query.Description
	if toolDescription == "" {
		toolDescription = fmt.Sprintf("Execute GraphQL query: %s", query.Name)
	}

	// Create input schema for the tool
	inputSchema := s.createInputSchema(query)

	tool := &mcp.Tool{
		Name:        toolName,
		Description: toolDescription,
		InputSchema: inputSchema,
	}

	// Create the handler function
	handler := func(ctx context.Context, req *mcp.CallToolRequest, input map[string]interface{}) (*mcp.CallToolResult, any, error) {
		result, err := s.executeGraphQLOperation(ctx, query, input, "query")
		return result, nil, err
	}

	mcp.AddTool(s.mcpServer, tool, handler)
	return nil
}

// addMutationTool adds an MCP tool for a GraphQL mutation
func (s *MCPGraphQLServer) addMutationTool(mutation *Field) error {
	toolName := "mutation_" + mutation.Name
	toolDescription := mutation.Description
	if toolDescription == "" {
		toolDescription = fmt.Sprintf("Execute GraphQL mutation: %s", mutation.Name)
	}

	// Create input schema for the tool
	inputSchema := s.createInputSchema(mutation)

	tool := &mcp.Tool{
		Name:        toolName,
		Description: toolDescription,
		InputSchema: inputSchema,
	}

	// Create the handler function
	handler := func(ctx context.Context, req *mcp.CallToolRequest, input map[string]interface{}) (*mcp.CallToolResult, any, error) {
		result, err := s.executeGraphQLOperation(ctx, mutation, input, "mutation")
		return result, nil, err
	}

	mcp.AddTool(s.mcpServer, tool, handler)
	return nil
}

// createInputSchema creates a JSON schema for the tool input
func (s *MCPGraphQLServer) createInputSchema(field *Field) map[string]interface{} {
	properties := make(map[string]interface{})
	required := []string{}

	// Add arguments as properties
	for _, arg := range field.Args {
		argSchema := s.createArgumentSchema(arg)
		properties[arg.Name] = argSchema

		// Add to required if it's non-null and has no default value
		if arg.Type.IsNonNull() && arg.DefaultValue == "" {
			required = append(required, arg.Name)
		}
	}

	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// createArgumentSchema creates a JSON schema for a GraphQL argument
func (s *MCPGraphQLServer) createArgumentSchema(arg *Argument) map[string]interface{} {
	schema := map[string]interface{}{
		"type": arg.Type.ToJSONSchemaType(),
	}

	// Add description if available (Arguments don't have descriptions in GraphQL introspection)
	// This is a placeholder for future enhancement

	// Handle list types
	if arg.Type.IsList() {
		schema["type"] = "array"
		schema["items"] = map[string]interface{}{
			"type": arg.Type.OfType.ToJSONSchemaType(),
		}
	}

	// Add default value if available
	if arg.DefaultValue != "" {
		schema["default"] = arg.DefaultValue
	}

	return schema
}

// executeGraphQLOperation executes a GraphQL query or mutation
func (s *MCPGraphQLServer) executeGraphQLOperation(ctx context.Context, field *Field, input map[string]interface{}, operationType string) (*mcp.CallToolResult, error) {
	// Generate the GraphQL query/mutation string
	var queryString string
	if operationType == "query" {
		queryString = field.GenerateQueryString()
	} else {
		queryString = field.GenerateMutationString()
	}

	// Execute the GraphQL operation
	resp, err := s.client.ExecuteQuery(ctx, queryString, input)
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("GraphQL %s failed: %v", operationType, err),
				},
			},
		}, nil
	}

	// Check for GraphQL errors
	if len(resp.Errors) > 0 {
		errorMessages := make([]string, len(resp.Errors))
		for i, err := range resp.Errors {
			errorMessages[i] = err.Message
		}
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("GraphQL %s errors: %s", operationType, strings.Join(errorMessages, "; ")),
				},
			},
		}, nil
	}

	// Convert response to JSON string
	jsonData, err := json.MarshalIndent(resp.Data, "", "  ")
	if err != nil {
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to marshal response: %v", err),
				},
			},
		}, nil
	}

	return &mcp.CallToolResult{
		IsError: false,
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: string(jsonData),
			},
		},
	}, nil
}

// GetMCPServer returns the underlying MCP server
func (s *MCPGraphQLServer) GetMCPServer() *mcp.Server {
	return s.mcpServer
}

// RefreshSchema re-introspects the GraphQL schema and updates tools
func (s *MCPGraphQLServer) RefreshSchema() error {
	ctx := context.Background()
	schema, err := s.client.IntrospectSchema(ctx)
	if err != nil {
		return fmt.Errorf("failed to refresh GraphQL schema: %w", err)
	}

	s.schema = schema

	// Recreate the MCP server with new tools
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "graphql-mcp-server",
		Version: "1.0.0",
	}, nil)

	s.mcpServer = mcpServer

	// Add tools for queries and mutations
	if err := s.addGraphQLTools(); err != nil {
		return fmt.Errorf("failed to add GraphQL tools after refresh: %w", err)
	}

	return nil
}

// GetSchema returns the current GraphQL schema
func (s *MCPGraphQLServer) GetSchema() *Schema {
	return s.schema
}

// GetClient returns the GraphQL client
func (s *MCPGraphQLServer) GetClient() *GraphQLClient {
	return s.client
}
