package graphqlmcp

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp/schema"
	"github.com/vektah/gqlparser/v2/ast"
)

// MCPGraphQLServer represents an MCP server that provides GraphQL tools
type MCPGraphQLServer struct {
	client    *GraphQLClient
	Schema    *schema.Schema
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

	// Schema is already parsed from introspection, no need to re-parse from SDL

	// Create MCP server
	mcpServer := mcp.NewServer(&mcp.Implementation{
		Name:    "graphql-mcp-server",
		Version: "1.0.0",
	}, nil)

	server := &MCPGraphQLServer{
		client:    client,
		Schema:    schema,
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
	queries := s.Schema.GetQueries()
	for _, query := range queries {
		if err := s.addQueryTool(query); err != nil {
			return fmt.Errorf("failed to add query tool for %s: %w", query.Name, err)
		}
	}

	// Add mutation tools
	mutations := s.Schema.GetMutations()
	for _, mutation := range mutations {
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
		result, err := s.executeGraphQLOperation(ctx, mutation, input, "mutation")
		return result, nil, err
	}

	mcp.AddTool(s.mcpServer, tool, handler)
	return nil
}

// createInputSchema creates a JSON schema for the tool input
func (s *MCPGraphQLServer) createInputSchema(field *schema.Field) map[string]interface{} {
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
func (s *MCPGraphQLServer) createArgumentSchema(arg *schema.Argument) map[string]interface{} {
	schema := map[string]interface{}{
		"type": arg.Type.ToJSONSchemaType(),
	}

	// Add description if available
	if arg.Description != "" {
		schema["description"] = arg.Description
	}

	// Handle list types
	if arg.Type.IsList() {
		schema["type"] = "array"
		// Create schema for list items
		itemSchema := s.createTypeRefSchema(arg.Type.OfType, arg.Description)
		schema["items"] = itemSchema
		return schema
	}

	// Handle input object types - resolve the actual input object definition
	if arg.Type.GetTypeName() != "" && !s.isBuiltinType(arg.Type.GetTypeName()) {
		if inputObjectSchema := s.createInputObjectSchema(arg.Type.GetTypeName()); inputObjectSchema != nil {
			// Merge the input object schema with the current schema
			for key, value := range inputObjectSchema {
				schema[key] = value
			}
		}
	}

	// Add default value if available
	if arg.DefaultValue != "" {
		schema["default"] = arg.DefaultValue
	}

	return schema
}

// executeGraphQLOperation executes a GraphQL query or mutation
func (s *MCPGraphQLServer) executeGraphQLOperation(ctx context.Context, field *schema.Field, input map[string]interface{}, operationType string) (*mcp.CallToolResult, error) {
	// Generate the GraphQL query/mutation string
	var queryString string
	var err error
	if operationType == "query" {
		queryString, err = field.GenerateQueryStringWithSchema(s.Schema)
		if err != nil {
			return nil, fmt.Errorf("failed to generate query string: %w", err)
		}
	} else {
		queryString, err = field.GenerateMutationStringWithSchema(s.Schema)
		if err != nil {
			return nil, fmt.Errorf("failed to generate mutation string: %w", err)
		}
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

	s.Schema = schema

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
func (s *MCPGraphQLServer) GetSchema() *schema.Schema {
	return s.Schema
}

// GetClient returns the GraphQL client
func (s *MCPGraphQLServer) GetClient() *GraphQLClient {
	return s.client
}

// TestCreateInputSchema creates a JSON schema for testing purposes
func (s *MCPGraphQLServer) TestCreateInputSchema(field *schema.Field) map[string]interface{} {
	return s.createInputSchema(field)
}

// TestCreateArgumentSchema creates a JSON schema for an argument for testing purposes
func (s *MCPGraphQLServer) TestCreateArgumentSchema(arg *schema.Argument) map[string]interface{} {
	return s.createArgumentSchema(arg)
}

// TestCreateInputObjectSchema creates a JSON schema for an input object type for testing purposes
func (s *MCPGraphQLServer) TestCreateInputObjectSchema(typeName string) map[string]interface{} {
	return s.createInputObjectSchema(typeName)
}

// createInputObjectSchema creates a detailed JSON schema for an input object type
func (s *MCPGraphQLServer) createInputObjectSchema(typeName string) map[string]interface{} {
	// Get the type definition from the schema
	typeDef := s.Schema.GetTypeDefinition(typeName)
	if typeDef == nil {
		return nil
	}

	// Only handle input object types
	if typeDef.Kind != "INPUT_OBJECT" {
		return nil
	}

	properties := make(map[string]interface{})
	required := []string{}

	// Process each field in the input object
	for _, field := range typeDef.Fields {
		fieldSchema := s.createInputFieldSchemaFromAST(field)
		properties[field.Name] = fieldSchema

		// Add to required if it's non-null and has no default value
		if s.isASTTypeNonNull(field.Type) && (field.DefaultValue == nil || field.DefaultValue.Raw == "") {
			required = append(required, field.Name)
		}
	}

	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	// Add description if available
	if typeDef.Description != "" {
		schema["description"] = typeDef.Description
	}

	return schema
}

// createInputFieldSchemaFromAST creates a JSON schema for an input field from AST
func (s *MCPGraphQLServer) createInputFieldSchemaFromAST(field *ast.FieldDefinition) map[string]interface{} {
	schema := map[string]interface{}{
		"type": s.astTypeToJSONSchemaType(field.Type),
	}

	// Add description if available
	if field.Description != "" {
		schema["description"] = field.Description
	}

	// Handle list types
	if s.isASTTypeList(field.Type) {
		schema["type"] = "array"
		itemSchema := s.createInputFieldSchemaFromAST(&ast.FieldDefinition{
			Name:        field.Name,
			Description: field.Description,
			Type:        field.Type.Elem,
		})
		schema["items"] = itemSchema
		return schema
	}

	// Handle nested input object types
	typeName := s.getASTTypeName(field.Type)
	if typeName != "" && !s.isBuiltinType(typeName) {
		if inputObjectSchema := s.createInputObjectSchema(typeName); inputObjectSchema != nil {
			// Merge the input object schema with the current schema
			for key, value := range inputObjectSchema {
				schema[key] = value
			}
		}
	}

	// Add default value if available
	if field.DefaultValue != nil {
		schema["default"] = field.DefaultValue.Raw
	}

	return schema
}

// isBuiltinType checks if a type name is a built-in GraphQL type
func (s *MCPGraphQLServer) isBuiltinType(typeName string) bool {
	builtinTypes := map[string]bool{
		"String":       true,
		"Int":          true,
		"Float":        true,
		"Boolean":      true,
		"ID":           true,
		"__Schema":     true,
		"__Type":       true,
		"__Field":      true,
		"__InputValue": true,
		"__EnumValue":  true,
		"__Directive":  true,
	}
	return builtinTypes[typeName]
}

// astTypeToJSONSchemaType converts AST type to JSON Schema type
func (s *MCPGraphQLServer) astTypeToJSONSchemaType(astType *ast.Type) string {
	if astType == nil {
		return "string"
	}

	// Get the base type name
	baseType := s.getASTTypeName(astType)

	// Convert GraphQL types to JSON Schema types
	switch baseType {
	case "String", "ID":
		return "string"
	case "Int":
		return "integer"
	case "Float":
		return "number"
	case "Boolean":
		return "boolean"
	default:
		return "object"
	}
}

// getASTTypeName extracts the type name from an AST type
func (s *MCPGraphQLServer) getASTTypeName(astType *ast.Type) string {
	if astType == nil {
		return ""
	}

	currentType := astType
	for currentType != nil {
		if currentType.NamedType != "" {
			return currentType.NamedType
		}
		currentType = currentType.Elem
	}
	return ""
}

// isASTTypeList checks if an AST type is a list
func (s *MCPGraphQLServer) isASTTypeList(astType *ast.Type) bool {
	if astType == nil {
		return false
	}

	// Check if this is a list type
	if astType.Elem != nil && !astType.NonNull {
		return true
	}

	// Check if wrapped in non-null
	if astType.NonNull && astType.Elem != nil {
		return s.isASTTypeList(astType.Elem)
	}

	return false
}

// isASTTypeNonNull checks if an AST type is non-null
func (s *MCPGraphQLServer) isASTTypeNonNull(astType *ast.Type) bool {
	if astType == nil {
		return false
	}

	// Check if this is a non-null type
	if astType.NonNull {
		return true
	}

	// For other types (including LIST), they are nullable unless wrapped in NON_NULL
	return false
}

// createTypeRefSchema creates a JSON schema for a TypeRef
func (s *MCPGraphQLServer) createTypeRefSchema(typeRef *schema.TypeRef, description string) map[string]interface{} {
	schema := map[string]interface{}{
		"type": typeRef.ToJSONSchemaType(),
	}

	// Add description if available
	if description != "" {
		schema["description"] = description
	}

	// Handle list types
	if typeRef.IsList() {
		schema["type"] = "array"
		itemSchema := s.createTypeRefSchema(typeRef.OfType, description)
		schema["items"] = itemSchema
		return schema
	}

	// Handle input object types
	if typeRef.GetTypeName() != "" && !s.isBuiltinType(typeRef.GetTypeName()) {
		if inputObjectSchema := s.createInputObjectSchema(typeRef.GetTypeName()); inputObjectSchema != nil {
			// Merge the input object schema with the current schema
			for key, value := range inputObjectSchema {
				schema[key] = value
			}
		}
	}

	return schema
}
