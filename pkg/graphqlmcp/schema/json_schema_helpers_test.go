package schema

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestSchema_CreateInputObjectSchema_ListFields(t *testing.T) {
	tests := []struct {
		name     string
		schema   *Schema
		typeName string
		expected map[string]interface{}
	}{
		{
			name: "input object with list of strings",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"CreatePostInput": {
						Name: "CreatePostInput",
						Kind: ast.InputObject,
						Fields: []*ast.FieldDefinition{
							{
								Name:        "tags",
								Description: "Post tags",
								Type:        ast.ListType(ast.NonNullNamedType("String", nil), nil),
							},
						},
					},
				},
			},
			typeName: "CreatePostInput",
			expected: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"tags": map[string]interface{}{
						"type":        "array",
						"description": "Post tags",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
		{
			name: "input object with non-null list of strings",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"CreateUserInput": {
						Name: "CreateUserInput",
						Kind: ast.InputObject,
						Fields: []*ast.FieldDefinition{
							{
								Name:        "skills",
								Description: "User skills",
								Type:        ast.NonNullListType(ast.NonNullNamedType("String", nil), nil),
							},
						},
					},
				},
			},
			typeName: "CreateUserInput",
			expected: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"skills": map[string]interface{}{
						"type":        "array",
						"description": "User skills",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
		{
			name: "input object with list of integers",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"CreateOrderInput": {
						Name: "CreateOrderInput",
						Kind: ast.InputObject,
						Fields: []*ast.FieldDefinition{
							{
								Name:        "quantities",
								Description: "Item quantities",
								Type:        ast.ListType(ast.NonNullNamedType("Int", nil), nil),
							},
						},
					},
				},
			},
			typeName: "CreateOrderInput",
			expected: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"quantities": map[string]interface{}{
						"type":        "array",
						"description": "Item quantities",
						"items": map[string]interface{}{
							"type": "integer",
						},
					},
				},
			},
		},
		{
			name: "input object with mixed field types including lists",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"CreateProductInput": {
						Name: "CreateProductInput",
						Kind: ast.InputObject,
						Fields: []*ast.FieldDefinition{
							{
								Name:        "name",
								Description: "Product name",
								Type:        ast.NonNullNamedType("String", nil),
							},
							{
								Name:        "tags",
								Description: "Product tags",
								Type:        ast.ListType(ast.NonNullNamedType("String", nil), nil),
							},
							{
								Name:        "categories",
								Description: "Product categories",
								Type:        ast.NonNullListType(ast.NonNullNamedType("String", nil), nil),
							},
							{
								Name:        "price",
								Description: "Product price",
								Type:        ast.NamedType("Float", nil),
							},
						},
					},
				},
			},
			typeName: "CreateProductInput",
			expected: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "Product name",
					},
					"tags": map[string]interface{}{
						"type":        "array",
						"description": "Product tags",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"categories": map[string]interface{}{
						"type":        "array",
						"description": "Product categories",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
					"price": map[string]interface{}{
						"type":        "number",
						"description": "Product price",
					},
				},
				"required": []string{"name", "categories"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.CreateInputObjectSchema(tt.typeName)

			if result == nil {
				t.Errorf("CreateInputObjectSchema() = nil, want %v", tt.expected)
				return
			}

			// Check type
			if result["type"] != tt.expected["type"] {
				t.Errorf("CreateInputObjectSchema() type = %v, want %v", result["type"], tt.expected["type"])
			}

			// Check properties
			if expectedProps, ok := tt.expected["properties"].(map[string]interface{}); ok {
				if resultProps, ok := result["properties"].(map[string]interface{}); ok {
					for key, expectedValue := range expectedProps {
						if resultValue, exists := resultProps[key]; !exists {
							t.Errorf("CreateInputObjectSchema() missing property: %s", key)
						} else if !mapsEqual(resultValue, expectedValue) {
							t.Errorf("CreateInputObjectSchema() property %s = %v, want %v", key, resultValue, expectedValue)
						}
					}
				} else {
					t.Errorf("CreateInputObjectSchema() properties is not a map")
				}
			}

			// Check required fields
			if expectedRequired, ok := tt.expected["required"].([]string); ok {
				if resultRequired, ok := result["required"].([]string); ok {
					if len(resultRequired) != len(expectedRequired) {
						t.Errorf("CreateInputObjectSchema() required length = %d, want %d", len(resultRequired), len(expectedRequired))
					} else {
						for i, expected := range expectedRequired {
							if resultRequired[i] != expected {
								t.Errorf("CreateInputObjectSchema() required[%d] = %v, want %v", i, resultRequired[i], expected)
							}
						}
					}
				} else {
					t.Errorf("CreateInputObjectSchema() required is not a string slice")
				}
			}
		})
	}
}

func TestSchema_CreateInputObjectSchema_RealData(t *testing.T) {
	// Test with real introspection data to ensure the fix works
	data, err := os.ReadFile("testdata/real_introspection_response.json")
	if err != nil {
		t.Skip("Test data not available")
	}

	var responseData map[string]interface{}
	if err := json.Unmarshal(data, &responseData); err != nil {
		t.Fatalf("Failed to unmarshal test data: %v", err)
	}

	// Extract the schema data from the response
	schemaData, ok := responseData["data"].(map[string]interface{})
	if !ok {
		t.Fatal("Invalid test data format: missing 'data' field")
	}

	parsedSchema, err := ParseIntrospectionResponse(schemaData)
	if err != nil {
		t.Fatalf("Failed to parse introspection response: %v", err)
	}

	// Test EquipmentSpecificationsInput schema
	inputSchema := parsedSchema.CreateInputObjectSchema("EquipmentSpecificationsInput")
	if inputSchema == nil {
		t.Fatal("Failed to create EquipmentSpecificationsInput schema")
	}

	// Check that certifications is a simple array of strings
	if props, ok := inputSchema["properties"].(map[string]interface{}); ok {
		if certs, ok := props["certifications"].(map[string]interface{}); ok {
			// Should be a simple array
			if certs["type"] != "array" {
				t.Errorf("certifications type = %v, want 'array'", certs["type"])
			}

			// Items should be simple strings, not nested arrays
			if items, ok := certs["items"].(map[string]interface{}); ok {
				if items["type"] != "string" {
					t.Errorf("certifications items type = %v, want 'string'", items["type"])
				}

				// Should NOT have nested items (this would indicate double-nested arrays)
				if _, hasNestedItems := items["items"]; hasNestedItems {
					t.Error("certifications items should not have nested 'items' field (indicates double-nested arrays)")
				}
			} else {
				t.Error("certifications items should be a map")
			}
		} else {
			t.Error("certifications field not found in EquipmentSpecificationsInput")
		}

		// Check that environmentalRequirements is also a simple array of strings
		if envReqs, ok := props["environmentalRequirements"].(map[string]interface{}); ok {
			// Should be a simple array
			if envReqs["type"] != "array" {
				t.Errorf("environmentalRequirements type = %v, want 'array'", envReqs["type"])
			}

			// Items should be simple strings, not nested arrays
			if items, ok := envReqs["items"].(map[string]interface{}); ok {
				if items["type"] != "string" {
					t.Errorf("environmentalRequirements items type = %v, want 'string'", items["type"])
				}

				// Should NOT have nested items
				if _, hasNestedItems := items["items"]; hasNestedItems {
					t.Error("environmentalRequirements items should not have nested 'items' field (indicates double-nested arrays)")
				}
			} else {
				t.Error("environmentalRequirements items should be a map")
			}
		} else {
			t.Error("environmentalRequirements field not found in EquipmentSpecificationsInput")
		}
	} else {
		t.Error("EquipmentSpecificationsInput properties should be a map")
	}
}

func TestSchema_createItemSchemaFromAST(t *testing.T) {
	tests := []struct {
		name     string
		astType  *ast.Type
		expected map[string]interface{}
	}{
		{
			name:     "nil AST type",
			astType:  nil,
			expected: map[string]interface{}{"type": "string"},
		},
		{
			name:     "string type",
			astType:  ast.NamedType("String", nil),
			expected: map[string]interface{}{"type": "string"},
		},
		{
			name:     "non-null string type",
			astType:  ast.NonNullNamedType("String", nil),
			expected: map[string]interface{}{"type": "string"},
		},
		{
			name:     "integer type",
			astType:  ast.NamedType("Int", nil),
			expected: map[string]interface{}{"type": "integer"},
		},
		{
			name:     "float type",
			astType:  ast.NamedType("Float", nil),
			expected: map[string]interface{}{"type": "number"},
		},
		{
			name:     "boolean type",
			astType:  ast.NamedType("Boolean", nil),
			expected: map[string]interface{}{"type": "boolean"},
		},
		{
			name:     "ID type",
			astType:  ast.NamedType("ID", nil),
			expected: map[string]interface{}{"type": "string"},
		},
	}

	schema := &Schema{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := schema.createItemSchemaFromAST(tt.astType, nil)

			if !mapsEqual(result, tt.expected) {
				t.Errorf("createItemSchemaFromAST() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSchema_CreateInputObjectSchema_EnumFields(t *testing.T) {
	tests := []struct {
		name     string
		schema   *Schema
		typeName string
		expected map[string]interface{}
	}{
		{
			name: "input object with enum field",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"UserRole": {
						Name: "UserRole",
						Kind: ast.Enum,
						EnumValues: []*ast.EnumValueDefinition{
							{Name: "ADMIN"},
							{Name: "USER"},
							{Name: "GUEST"},
						},
					},
					"CreateUserInput": {
						Name: "CreateUserInput",
						Kind: ast.InputObject,
						Fields: []*ast.FieldDefinition{
							{
								Name:        "role",
								Description: "User role",
								Type:        ast.NonNullNamedType("UserRole", nil),
							},
						},
					},
				},
			},
			typeName: "CreateUserInput",
			expected: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"role": map[string]interface{}{
						"type":        "string",
						"description": "User role",
						"enum":        []interface{}{"ADMIN", "USER", "GUEST"},
					},
				},
				"required": []interface{}{"role"},
			},
		},
		{
			name: "input object with optional enum field",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"Status": {
						Name: "Status",
						Kind: ast.Enum,
						EnumValues: []*ast.EnumValueDefinition{
							{Name: "ACTIVE"},
							{Name: "INACTIVE"},
							{Name: "PENDING"},
						},
					},
					"UpdateUserInput": {
						Name: "UpdateUserInput",
						Kind: ast.InputObject,
						Fields: []*ast.FieldDefinition{
							{
								Name:        "status",
								Description: "User status",
								Type:        ast.NamedType("Status", nil),
							},
						},
					},
				},
			},
			typeName: "UpdateUserInput",
			expected: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"status": map[string]interface{}{
						"type":        "string",
						"description": "User status",
						"enum":        []interface{}{"ACTIVE", "INACTIVE", "PENDING"},
					},
				},
			},
		},
		{
			name: "input object with list of enums",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"Priority": {
						Name: "Priority",
						Kind: ast.Enum,
						EnumValues: []*ast.EnumValueDefinition{
							{Name: "LOW"},
							{Name: "MEDIUM"},
							{Name: "HIGH"},
							{Name: "CRITICAL"},
						},
					},
					"CreateTaskInput": {
						Name: "CreateTaskInput",
						Kind: ast.InputObject,
						Fields: []*ast.FieldDefinition{
							{
								Name:        "priorities",
								Description: "Task priorities",
								Type:        ast.ListType(ast.NonNullNamedType("Priority", nil), nil),
							},
						},
					},
				},
			},
			typeName: "CreateTaskInput",
			expected: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"priorities": map[string]interface{}{
						"type":        "array",
						"description": "Task priorities",
						"items": map[string]interface{}{
							"type": "string",
							"enum": []interface{}{"LOW", "MEDIUM", "HIGH", "CRITICAL"},
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.CreateInputObjectSchema(tt.typeName)

			// Check if the result has the expected structure
			if result["type"] != "object" {
				t.Errorf("Expected type to be 'object', got %v", result["type"])
			}

			properties, ok := result["properties"].(map[string]interface{})
			if !ok {
				t.Errorf("Expected properties to be a map")
				return
			}

			// Check each property
			for fieldName, expectedField := range tt.expected["properties"].(map[string]interface{}) {
				actualField, exists := properties[fieldName]
				if !exists {
					t.Errorf("Expected field %s to exist", fieldName)
					continue
				}

				actualFieldMap := actualField.(map[string]interface{})
				expectedFieldMap := expectedField.(map[string]interface{})

				// Check type
				if actualFieldMap["type"] != expectedFieldMap["type"] {
					t.Errorf("Field %s: expected type %v, got %v", fieldName, expectedFieldMap["type"], actualFieldMap["type"])
				}

				// Check description
				if actualFieldMap["description"] != expectedFieldMap["description"] {
					t.Errorf("Field %s: expected description %v, got %v", fieldName, expectedFieldMap["description"], actualFieldMap["description"])
				}

				// Check enum values
				if expectedEnum, hasEnum := expectedFieldMap["enum"]; hasEnum {
					actualEnum, hasActualEnum := actualFieldMap["enum"]
					if !hasActualEnum {
						t.Errorf("Field %s: expected enum values but none found", fieldName)
					} else {
						// Convert to string slices for comparison
						actualEnumSlice := actualEnum.([]string)
						expectedEnumSlice := make([]string, len(expectedEnum.([]interface{})))
						for i, v := range expectedEnum.([]interface{}) {
							expectedEnumSlice[i] = v.(string)
						}

						// Check if slices have same length and contain same values
						if len(actualEnumSlice) != len(expectedEnumSlice) {
							t.Errorf("Field %s: expected %d enum values, got %d", fieldName, len(expectedEnumSlice), len(actualEnumSlice))
						} else {
							// Create maps for comparison
							actualMap := make(map[string]bool)
							expectedMap := make(map[string]bool)
							for _, v := range actualEnumSlice {
								actualMap[v] = true
							}
							for _, v := range expectedEnumSlice {
								expectedMap[v] = true
							}

							for v := range expectedMap {
								if !actualMap[v] {
									t.Errorf("Field %s: expected enum value %s not found", fieldName, v)
								}
							}
							for v := range actualMap {
								if !expectedMap[v] {
									t.Errorf("Field %s: unexpected enum value %s found", fieldName, v)
								}
							}
						}
					}
				}
			}

			// Check required fields
			if expectedRequired, hasRequired := tt.expected["required"]; hasRequired {
				actualRequired, hasActualRequired := result["required"]
				if !hasActualRequired {
					t.Errorf("Expected required fields but none found")
				} else {
					actualRequiredSlice := actualRequired.([]string)
					expectedRequiredSlice := make([]string, len(expectedRequired.([]interface{})))
					for i, v := range expectedRequired.([]interface{}) {
						expectedRequiredSlice[i] = v.(string)
					}

					if len(actualRequiredSlice) != len(expectedRequiredSlice) {
						t.Errorf("Expected %d required fields, got %d", len(expectedRequiredSlice), len(actualRequiredSlice))
					} else {
						actualMap := make(map[string]bool)
						expectedMap := make(map[string]bool)
						for _, v := range actualRequiredSlice {
							actualMap[v] = true
						}
						for _, v := range expectedRequiredSlice {
							expectedMap[v] = true
						}

						for v := range expectedMap {
							if !actualMap[v] {
								t.Errorf("Expected required field %s not found", v)
							}
						}
						for v := range actualMap {
							if !expectedMap[v] {
								t.Errorf("Unexpected required field %s found", v)
							}
						}
					}
				}
			}
		})
	}
}

// Helper function to compare maps recursively (reused from existing test)
func mapsEqual(a, b interface{}) bool {
	switch aVal := a.(type) {
	case map[string]interface{}:
		if bMap, ok := b.(map[string]interface{}); ok {
			if len(aVal) != len(bMap) {
				return false
			}
			for key, aValue := range aVal {
				if bValue, exists := bMap[key]; !exists || !mapsEqual(aValue, bValue) {
					return false
				}
			}
			return true
		}
		return false
	case []string:
		if bSlice, ok := b.([]string); ok {
			if len(aVal) != len(bSlice) {
				return false
			}
			for i, aValue := range aVal {
				if bSlice[i] != aValue {
					return false
				}
			}
			return true
		}
		return false
	case []interface{}:
		if bSlice, ok := b.([]interface{}); ok {
			if len(aVal) != len(bSlice) {
				return false
			}
			for i, aValue := range aVal {
				if !mapsEqual(aValue, bSlice[i]) {
					return false
				}
			}
			return true
		}
		return false
	default:
		return a == b
	}
}
