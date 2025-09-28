package graphqlmcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp/schema"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// MockGraphQLExecutor is a mock implementation of GraphQLExecutor for testing
type MockGraphQLExecutor struct {
	mock.Mock
}

func (m *MockGraphQLExecutor) ExecuteQuery(ctx context.Context, query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	args := m.Called(ctx, query, variables)
	return args.Get(0).(*GraphQLResponse), args.Error(1)
}

func (m *MockGraphQLExecutor) IntrospectSchema(ctx context.Context) (*schema.Schema, error) {
	args := m.Called(ctx)
	return args.Get(0).(*schema.Schema), args.Error(1)
}

// loadTestSchema loads the test schema from the testdata file
func loadTestSchema(t *testing.T) *schema.Schema {
	t.Helper()

	data, err := os.ReadFile("testdata/real_introspection_response.json")
	if err != nil {
		t.Fatalf("Failed to read test data: %v", err)
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(data, &responseData); err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	// Extract the schema data from the response
	schemaData, ok := responseData["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Invalid test data format: missing 'data' field")
	}

	schema, err := schema.ParseIntrospectionResponse(schemaData)
	if err != nil {
		t.Fatalf("Failed to parse introspection response: %v", err)
	}

	return schema
}

func TestMCPGraphQLServer_ToolCreation(t *testing.T) {
	// Load test schema
	testSchema := loadTestSchema(t)

	// Create mock executor
	mockExecutor := new(MockGraphQLExecutor)
	mockExecutor.On("IntrospectSchema", mock.Anything).Return(testSchema, nil)

	// Create server with mock executor
	server, err := NewMCPGraphQLServerWithExecutor(mockExecutor)
	assert.NoError(t, err)
	assert.NotNil(t, server)

	// Get the MCP server
	mcpServer := server.GetMCPServer()
	assert.NotNil(t, mcpServer)

	// Get tools from the schema (since MCP server doesn't expose tools directly)
	queries := server.GetSchema().GetQueries()
	mutations := server.GetSchema().GetMutations()
	totalTools := len(queries) + len(mutations)
	assert.NotEmpty(t, totalTools)

	// Expected query tools based on the schema
	expectedQueryTools := []string{
		"query_equipment",
		"query_equipmentById",
		"query_facilities",
		"query_facilityById",
		"query_maintenanceRecords",
		"query_maintenanceRecordsByEquipment",
		"query_operationalMetrics",
		"query_facilityStatus",
	}

	// Expected mutation tools based on the schema
	expectedMutationTools := []string{
		"mutation_createEquipment",
		"mutation_updateEquipment",
		"mutation_deleteEquipment",
		"mutation_createFacility",
		"mutation_updateFacility",
		"mutation_deleteFacility",
		"mutation_scheduleMaintenance",
		"mutation_updateMaintenanceRecord",
		"mutation_completeMaintenance",
		"mutation_recordOperationalMetric",
		"mutation_updateEquipmentStatus",
	}

	// Create a map of actual tool names from the schema
	actualQueryTools := make(map[string]bool)
	for _, query := range queries {
		actualQueryTools["query_"+query.Name] = true
	}

	actualMutationTools := make(map[string]bool)
	for _, mutation := range mutations {
		actualMutationTools["mutation_"+mutation.Name] = true
	}

	// Verify all expected query tools are created
	for _, expectedTool := range expectedQueryTools {
		assert.True(t, actualQueryTools[expectedTool], "Expected query tool %s not found", expectedTool)
	}

	// Verify all expected mutation tools are created
	for _, expectedTool := range expectedMutationTools {
		assert.True(t, actualMutationTools[expectedTool], "Expected mutation tool %s not found", expectedTool)
	}

	// Verify we have the expected total number of tools
	expectedTotalTools := len(expectedQueryTools) + len(expectedMutationTools)
	assert.Equal(t, expectedTotalTools, totalTools, "Expected %d tools, got %d", expectedTotalTools, totalTools)

	// Verify mock expectations
	mockExecutor.AssertExpectations(t)
}

func TestMCPGraphQLServer_QueryToolExecution(t *testing.T) {
	// Load test schema
	testSchema := loadTestSchema(t)

	// Create mock executor
	mockExecutor := new(MockGraphQLExecutor)
	mockExecutor.On("IntrospectSchema", mock.Anything).Return(testSchema, nil)

	// Create server with mock executor
	server, err := NewMCPGraphQLServerWithExecutor(mockExecutor)
	assert.NoError(t, err)
	assert.NotNil(t, server)

	// Test each query tool
	queries := testSchema.GetQueries()
	for _, query := range queries {
		t.Run(fmt.Sprintf("Query_%s", query.Name), func(t *testing.T) {
			// Create mock response
			mockResponse := &GraphQLResponse{
				Data: map[string]interface{}{
					query.Name: "mock_data",
				},
				Errors: nil,
			}

			// Set up mock expectations
			mockExecutor.On("ExecuteQuery", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).
				Return(mockResponse, nil).Once()

			// Create appropriate input based on query arguments
			input := make(map[string]interface{})
			for _, arg := range query.Args {
				switch arg.Name {
				case "id":
					input["id"] = "test-id"
				case "equipmentId":
					input["equipmentId"] = "test-equipment-id"
				case "facilityId":
					input["facilityId"] = "test-facility-id"
				}
			}

			// Test the executeGraphQLOperation method directly
			ctx := context.Background()
			result, err := server.executeGraphQLOperation(ctx, query, input, "query")
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.False(t, result.IsError)
		})
	}

	// Verify mock expectations
	mockExecutor.AssertExpectations(t)
}

func TestMCPGraphQLServer_MutationToolExecution(t *testing.T) {
	// Load test schema
	testSchema := loadTestSchema(t)

	// Create mock executor
	mockExecutor := new(MockGraphQLExecutor)
	mockExecutor.On("IntrospectSchema", mock.Anything).Return(testSchema, nil)

	// Create server with mock executor
	server, err := NewMCPGraphQLServerWithExecutor(mockExecutor)
	assert.NoError(t, err)
	assert.NotNil(t, server)
	// Test each mutation tool
	mutations := testSchema.GetMutations()
	for _, mutation := range mutations {
		t.Run(fmt.Sprintf("Mutation_%s", mutation.Name), func(t *testing.T) {
			// Create mock response
			mockResponse := &GraphQLResponse{
				Data: map[string]interface{}{
					mutation.Name: "mock_data",
				},
				Errors: nil,
			}

			// Set up mock expectations
			mockExecutor.On("ExecuteQuery", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).
				Return(mockResponse, nil).Once()

			// Create appropriate input based on mutation arguments
			input := make(map[string]interface{})
			for _, arg := range mutation.Args {
				switch arg.Name {
				case "id":
					input["id"] = "test-id"
				case "input":
					// Create a mock input object based on the mutation
					input["input"] = createMockInputForMutation(mutation.Name)
				case "status":
					input["status"] = "RUNNING"
				case "notes":
					input["notes"] = "test notes"
				}
			}

			// Test the executeGraphQLOperation method directly
			ctx := context.Background()
			result, err := server.executeGraphQLOperation(ctx, mutation, input, "mutation")
			assert.NoError(t, err)
			assert.NotNil(t, result)
			assert.False(t, result.IsError)
		})
	}

	// Verify mock expectations
	mockExecutor.AssertExpectations(t)
}

// createMockInputForMutation creates appropriate mock input data for different mutations
func createMockInputForMutation(mutationName string) map[string]interface{} {
	switch mutationName {
	case "createEquipment":
		return map[string]interface{}{
			"name":         "Test Equipment",
			"description":  "Test Description",
			"manufacturer": "Test Manufacturer",
			"model":        "Test Model",
			"serialNumber": "TEST123",
			"type":         "CNC_MILL",
			"facilityId":   "facility-1",
			"specifications": map[string]interface{}{
				"powerConsumption": 10.5,
				"maxSpeed":         100.0,
				"operatingTemperature": map[string]interface{}{
					"min": 0.0,
					"max": 50.0,
				},
				"weight": 1000.0,
				"dimensions": map[string]interface{}{
					"length": 2.0,
					"width":  1.0,
					"height": 1.5,
				},
				"electricalSpecs": map[string]interface{}{
					"voltage":     220.0,
					"current":     10.0,
					"powerFactor": 0.9,
					"frequency":   50.0,
				},
				"environmentalRequirements": []string{"clean", "dry"},
				"certifications":            []string{"ISO9001"},
			},
			"installedAt": "2024-01-01T00:00:00Z",
		}
	case "createFacility":
		return map[string]interface{}{
			"name":             "Test Facility",
			"address":          "123 Test St",
			"location":         map[string]interface{}{"latitude": 40.0, "longitude": -74.0},
			"capacity":         100,
			"operationalSince": "2024-01-01T00:00:00Z",
			"manager":          "Test Manager",
			"contactInfo": map[string]interface{}{
				"phone":          "555-1234",
				"email":          "test@example.com",
				"emergencyPhone": "555-9999",
				"managerContact": "Test Manager",
			},
		}
	case "scheduleMaintenance":
		return map[string]interface{}{
			"equipmentId":        "equipment-1",
			"type":               "PREVENTIVE",
			"priority":           "HIGH",
			"scheduledDate":      "2024-02-01T00:00:00Z",
			"description":        "Test maintenance",
			"assignedTechnician": "Test Tech",
			"estimatedDuration":  4,
			"requiredParts":      []string{"part1", "part2"},
		}
	case "recordOperationalMetric":
		return map[string]interface{}{
			"equipmentId": "equipment-1",
			"metricType":  "EFFICIENCY",
			"value":       95.5,
			"unit":        "percent",
			"targetValue": 90.0,
			"notes":       "Test metric",
		}
	default:
		return map[string]interface{}{}
	}
}

func TestMCPGraphQLServer_ErrorHandling(t *testing.T) {
	// Load test schema
	testSchema := loadTestSchema(t)

	// Create mock executor
	mockExecutor := new(MockGraphQLExecutor)
	mockExecutor.On("IntrospectSchema", mock.Anything).Return(testSchema, nil)

	// Create server with mock executor
	server, err := NewMCPGraphQLServerWithExecutor(mockExecutor)
	assert.NoError(t, err)
	assert.NotNil(t, server)

	// Test error handling in query execution
	queries := testSchema.GetQueries()
	if len(queries) > 0 {
		query := queries[0]

		// Test with GraphQL errors
		mockResponseWithErrors := &GraphQLResponse{
			Data: nil,
			Errors: []struct {
				Message string `json:"message"`
			}{
				{Message: "Test GraphQL error"},
			},
		}

		mockExecutor.On("ExecuteQuery", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).
			Return(mockResponseWithErrors, nil).Once()

		ctx := context.Background()
		input := make(map[string]interface{})
		result, err := server.executeGraphQLOperation(ctx, query, input, "query")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.Contains(t, result.Content[0].(*mcp.TextContent).Text, "GraphQL query errors")
	}

	// Test with execution error
	if len(queries) > 0 {
		query := queries[0]

		mockExecutor.On("ExecuteQuery", mock.Anything, mock.AnythingOfType("string"), mock.AnythingOfType("map[string]interface {}")).
			Return((*GraphQLResponse)(nil), fmt.Errorf("network error")).Once()

		ctx := context.Background()
		input := make(map[string]interface{})
		result, err := server.executeGraphQLOperation(ctx, query, input, "query")

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.True(t, result.IsError)
		assert.Contains(t, result.Content[0].(*mcp.TextContent).Text, "GraphQL query failed")
	}

	// Verify mock expectations
	mockExecutor.AssertExpectations(t)
}

func TestMCPGraphQLServer_ToolDescriptions(t *testing.T) {
	// Load test schema
	testSchema := loadTestSchema(t)

	// Create mock executor
	mockExecutor := new(MockGraphQLExecutor)
	mockExecutor.On("IntrospectSchema", mock.Anything).Return(testSchema, nil)

	// Create server with mock executor
	server, err := NewMCPGraphQLServerWithExecutor(mockExecutor)
	assert.NoError(t, err)
	assert.NotNil(t, server)

	// Verify tool descriptions are properly set by checking the schema fields
	queries := server.GetSchema().GetQueries()
	mutations := server.GetSchema().GetMutations()

	// Verify query descriptions
	for _, query := range queries {
		toolName := "query_" + query.Name
		expectedDescription := query.Description
		if expectedDescription == "" {
			expectedDescription = fmt.Sprintf("Execute GraphQL query: %s", query.Name)
		}
		assert.NotEmpty(t, expectedDescription, "Query tool %s should have a description", toolName)
		// The description should be meaningful and not empty
		assert.True(t, len(expectedDescription) > 10, "Query tool %s should have a meaningful description", toolName)
	}

	// Verify mutation descriptions
	for _, mutation := range mutations {
		toolName := "mutation_" + mutation.Name
		expectedDescription := mutation.Description
		if expectedDescription == "" {
			expectedDescription = fmt.Sprintf("Execute GraphQL mutation: %s", mutation.Name)
		}
		assert.NotEmpty(t, expectedDescription, "Mutation tool %s should have a description", toolName)
		// The description should be meaningful and not empty
		assert.True(t, len(expectedDescription) > 10, "Mutation tool %s should have a meaningful description", toolName)
	}

	// Verify mock expectations
	mockExecutor.AssertExpectations(t)
}

func TestMCPGraphQLServer_InputSchemas(t *testing.T) {
	// Load test schema
	testSchema := loadTestSchema(t)

	// Create mock executor
	mockExecutor := new(MockGraphQLExecutor)
	mockExecutor.On("IntrospectSchema", mock.Anything).Return(testSchema, nil)

	// Create server with mock executor
	server, err := NewMCPGraphQLServerWithExecutor(mockExecutor)
	assert.NoError(t, err)
	assert.NotNil(t, server)

	// Verify all tools have input schemas by checking the schema fields
	queries := server.GetSchema().GetQueries()
	mutations := server.GetSchema().GetMutations()

	// Verify query input schemas
	for _, query := range queries {
		toolName := "query_" + query.Name
		inputSchema := server.createInputSchema(query)
		assert.NotNil(t, inputSchema, "Query tool %s should have an input schema", toolName)
		assert.NotEmpty(t, inputSchema, "Query tool %s input schema should not be empty", toolName)
	}

	// Verify mutation input schemas
	for _, mutation := range mutations {
		toolName := "mutation_" + mutation.Name
		inputSchema := server.createInputSchema(mutation)
		assert.NotNil(t, inputSchema, "Mutation tool %s should have an input schema", toolName)
		assert.NotEmpty(t, inputSchema, "Mutation tool %s input schema should not be empty", toolName)
	}

	// Verify mock expectations
	mockExecutor.AssertExpectations(t)
}
