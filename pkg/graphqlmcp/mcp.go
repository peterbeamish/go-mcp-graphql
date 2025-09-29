package graphqlmcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/go-logr/logr"
	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp/schema"
)

// MCPGraphQLServer represents an MCP server that provides GraphQL tools
type MCPGraphQLServer struct {
	executor  GraphQLExecutor
	Schema    *schema.Schema
	mcpServer *mcp.Server
	logger    logr.Logger
	options   *MCPGraphQLServerOptions
	testMode  bool
}

// NewMCPGraphQLServer creates a new MCP GraphQL server
// endpoint is the URL of the GraphQL server
func NewMCPGraphQLServer(endpoint string, opts ...MCPGraphQLServerOption) (*MCPGraphQLServer, error) {
	client := NewGraphQLClient(endpoint)
	return NewMCPGraphQLServerWithExecutor(client, opts...)
}

// NewMCPGraphQLServerWithExecutor creates a new MCP GraphQL server with a custom executor
func NewMCPGraphQLServerWithExecutor(executor GraphQLExecutor, opts ...MCPGraphQLServerOption) (*MCPGraphQLServer, error) {
	// Apply options
	options := NewMCPGraphQLServerOptions()
	for _, opt := range opts {
		opt(options)
	}

	// Set logger - use provided logger or default
	logger := options.Logger
	if logger.GetSink() == nil {
		logger = logr.Discard()
	}

	// Introspect the schema
	ctx := context.Background()
	schema, err := executor.IntrospectSchema(ctx)
	if err != nil {
		logger.Info("Failed to introspect GraphQL schema, continuing with empty schema", "error", err)
	}

	// Apply max depth configuration to schema if it was introspected successfully
	if schema != nil {
		schema.MaxDepth = options.MaxDepth
	}

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "graphql-mcp-server",
		Version: "1.0.0",
	}, nil)

	server := &MCPGraphQLServer{
		executor:  executor,
		Schema:    schema,
		mcpServer: mcpServer,
		logger:    logger,
		options:   options,
	}

	// Add tools for queries and mutations
	if server.Schema != nil {
		if err := server.addGraphQLTools(); err != nil {
			return nil, fmt.Errorf("failed to add GraphQL tools: %w", err)
		}
	} else {
		logger.Info("No schema introspected, skipping tool creation")
	}

	return server, nil
}

// addGraphQLTools adds MCP tools for all GraphQL queries and mutations
func (s *MCPGraphQLServer) addGraphQLTools() error {
	// Add query tools
	queries := s.Schema.GetQueries()
	for _, query := range queries {
		// Check if this query is allowed based on masking options
		if !s.options.isOperationAllowed(query.Name) {
			s.logger.V(1).Info("Skipping query due to masking rules", "query_name", query.Name)
			continue
		}

		if err := s.addQueryTool(query); err != nil {
			return fmt.Errorf("failed to add query tool for %s: %w", query.Name, err)
		}
	}

	// Add mutation tools
	mutations := s.Schema.GetMutations()
	for _, mutation := range mutations {
		// Check if this mutation is allowed based on masking options
		if !s.options.isOperationAllowed(mutation.Name) {
			s.logger.V(1).Info("Skipping mutation due to masking rules", "mutation_name", mutation.Name)
			continue
		}

		if err := s.addMutationTool(mutation); err != nil {
			return fmt.Errorf("failed to add mutation tool for %s: %w", mutation.Name, err)
		}
	}

	return nil
}

// addQueryTool adds an MCP tool for a GraphQL query
func (s *MCPGraphQLServer) addQueryTool(query *schema.Field) error {
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
		// Add passthru headers to context if available
		if passthruHeaders := GetPassthruHeaders(ctx); passthruHeaders != nil {
			ctx = AddPassthruHeadersToContext(ctx, passthruHeaders)
		}
		result, err := s.executeGraphQLOperation(ctx, query, input, "query")
		return result, nil, err
	}

	mcp.AddTool(s.mcpServer, tool, handler)
	return nil
}

// addMutationTool adds an MCP tool for a GraphQL mutation
func (s *MCPGraphQLServer) addMutationTool(mutation *schema.Field) error {
	toolName := "mutation_" + mutation.Name
	toolDescription := mutation.Description
	if toolDescription == "" {
		toolDescription = fmt.Sprintf("Execute GraphQL mutation: %s", mutation.Name)
	}

	// Create input schema for the tool
	inputSchema := s.createInputSchema(mutation)

	// Enhance description with input information
	if len(mutation.Args) > 0 {
		argNames := make([]string, len(mutation.Args))
		for i, arg := range mutation.Args {
			argNames[i] = arg.Name
		}
		toolDescription += fmt.Sprintf(" (Inputs: %s)", strings.Join(argNames, ", "))
	}

	tool := &mcp.Tool{
		Name:        toolName,
		Description: toolDescription,
		InputSchema: inputSchema,
	}

	// Create the handler function
	handler := func(ctx context.Context, req *mcp.CallToolRequest, input map[string]interface{}) (*mcp.CallToolResult, any, error) {
		// Add passthru headers to context if available
		if passthruHeaders := GetPassthruHeaders(ctx); passthruHeaders != nil {
			ctx = AddPassthruHeadersToContext(ctx, passthruHeaders)
		}
		result, err := s.executeGraphQLOperation(ctx, mutation, input, "mutation")
		return result, nil, err
	}

	mcp.AddTool(s.mcpServer, tool, handler)
	return nil
}

// createInputSchema creates a JSON schema for the tool input
func (s *MCPGraphQLServer) createInputSchema(field *schema.Field) map[string]interface{} {
	return s.Schema.CreateInputSchema(field)
}

// executeGraphQLOperation executes a GraphQL query or mutation
func (s *MCPGraphQLServer) executeGraphQLOperation(ctx context.Context, field *schema.Field, input map[string]interface{}, operationType string) (*mcp.CallToolResult, error) {
	// Generate a request ID for tracking
	requestID := fmt.Sprintf("req_%d", time.Now().UnixNano())

	// Log tool call initiation
	s.logger.Info("Tool call initiated",
		"request_id", requestID,
		"operation_type", operationType,
		"field_name", field.Name,
		"input_args", len(field.Args),
		"input_values", input,
	)

	// Generate the GraphQL query/mutation string
	var queryString string
	var err error
	if operationType == "query" {
		queryString, err = field.GenerateQueryStringWithSchema(s.Schema)
		if err != nil {
			s.logger.Error(err, "Failed to generate query string",
				"request_id", requestID,
				"field_name", field.Name,
			)
			return nil, fmt.Errorf("failed to generate query string: %w", err)
		}
	} else {
		queryString, err = field.GenerateMutationStringWithSchema(s.Schema)
		if err != nil {
			s.logger.Error(err, "Failed to generate mutation string",
				"request_id", requestID,
				"field_name", field.Name,
			)
			return nil, fmt.Errorf("failed to generate mutation string: %w", err)
		}
	}

	// Log the generated query/mutation
	s.logger.V(1).Info("Generated GraphQL operation",
		"request_id", requestID,
		"operation_type", operationType,
		"field_name", field.Name,
		"query", queryString,
	)

	// Execute the GraphQL operation
	startTime := time.Now()
	resp, err := s.executor.ExecuteQuery(ctx, queryString, input)
	duration := time.Since(startTime)

	if err != nil {
		s.logger.Error(err, "GraphQL execution failed",
			"request_id", requestID,
			"operation_type", operationType,
			"field_name", field.Name,
			"duration_ms", duration.Milliseconds(),
		)
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
		s.logger.Info("GraphQL operation returned errors",
			"request_id", requestID,
			"operation_type", operationType,
			"field_name", field.Name,
			"duration_ms", duration.Milliseconds(),
			"error_count", len(resp.Errors),
			"errors", errorMessages,
		)
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
		s.logger.Error(err, "Failed to marshal GraphQL response",
			"request_id", requestID,
			"operation_type", operationType,
			"field_name", field.Name,
			"duration_ms", duration.Milliseconds(),
		)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{
				&mcp.TextContent{
					Text: fmt.Sprintf("Failed to marshal response: %v", err),
				},
			},
		}, nil
	}

	// Log successful completion
	s.logger.Info("Tool call completed successfully",
		"request_id", requestID,
		"operation_type", operationType,
		"field_name", field.Name,
		"duration_ms", duration.Milliseconds(),
		"response_size_bytes", len(jsonData),
	)

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

// SetLogger sets a custom logger for the server
func (s *MCPGraphQLServer) SetLogger(logger logr.Logger) {
	s.logger = logger
}

// RefreshSchema re-introspects the GraphQL schema and updates tools
func (s *MCPGraphQLServer) RefreshSchema() error {
	ctx := context.Background()
	schema, err := s.executor.IntrospectSchema(ctx)
	if err != nil {
		return fmt.Errorf("failed to refresh GraphQL schema: %w", err)
	}

	// Apply max depth configuration to the refreshed schema
	if schema != nil {
		schema.MaxDepth = s.options.MaxDepth
	}

	s.Schema = schema

	// Recreate the MCP server with new tools
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "graphql-mcp-server",
		Version: "1.0.0",
	}, nil)

	s.mcpServer = mcpServer

	// Add tools for queries and mutations (respecting masking options)
	if err := s.addGraphQLTools(); err != nil {
		return fmt.Errorf("failed to add GraphQL tools after refresh: %w", err)
	}

	return nil
}

// GetSchema returns the current GraphQL schema
func (s *MCPGraphQLServer) GetSchema() *schema.Schema {
	return s.Schema
}

// GetExecutor returns the GraphQL executor
func (s *MCPGraphQLServer) GetExecutor() GraphQLExecutor {
	return s.executor
}

// ExtractPassthruHeaders extracts the configured passthru headers from the request
func (s *MCPGraphQLServer) ExtractPassthruHeaders(r *http.Request) map[string]string {
	if len(s.options.PassthruHeaders) == 0 {
		return nil
	}

	headers := make(map[string]string)
	for _, headerName := range s.options.PassthruHeaders {
		if value := r.Header.Get(headerName); value != "" {
			headers[headerName] = value
		}
	}

	return headers
}
