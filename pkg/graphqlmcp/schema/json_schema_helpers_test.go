package schema

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestSchema_CreateInputObjectSchema(t *testing.T) {
	tests := []struct {
		name     string
		schema   *Schema
		typeName string
		expected map[string]interface{}
	}{
		{
			name: "nil type definition",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{},
			},
			typeName: "NonExistentType",
			expected: nil,
		},
		{
			name: "non-input object type",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"User": {
						Name: "User",
						Kind: ast.Object,
					},
				},
			},
			typeName: "User",
			expected: nil,
		},
		{
			name: "simple input object",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"CreateUserInput": {
						Name:        "CreateUserInput",
						Kind:        ast.InputObject,
						Description: "Input for creating a user",
						Fields: []*ast.FieldDefinition{
							{
								Name:        "name",
								Description: "User name",
								Type: &ast.Type{
									NonNull: true,
									Elem: &ast.Type{
										NamedType: "String",
									},
								},
							},
							{
								Name:        "email",
								Description: "User email",
								Type: &ast.Type{
									NamedType: "String",
								},
							},
						},
					},
				},
			},
			typeName: "CreateUserInput",
			expected: map[string]interface{}{
				"type":        "object",
				"description": "Input for creating a user",
				"properties": map[string]interface{}{
					"name": map[string]interface{}{
						"type":        "string",
						"description": "User name",
					},
					"email": map[string]interface{}{
						"type":        "string",
						"description": "User email",
					},
				},
				"required": []string{"name"},
			},
		},
		{
			name: "input object with list field",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"CreatePostInput": {
						Name: "CreatePostInput",
						Kind: ast.InputObject,
						Fields: []*ast.FieldDefinition{
							{
								Name: "tags",
								Type: &ast.Type{
									Elem: &ast.Type{
										NonNull: true,
										Elem: &ast.Type{
											NamedType: "String",
										},
									},
								},
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
						"type": "array",
						"items": map[string]interface{}{
							"type": "string",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.CreateInputObjectSchema(tt.typeName)

			if tt.expected == nil {
				if result != nil {
					t.Errorf("CreateInputObjectSchema() = %v, want nil", result)
				}
				return
			}

			if result == nil {
				t.Errorf("CreateInputObjectSchema() = nil, want %v", tt.expected)
				return
			}

			// Check type
			if result["type"] != tt.expected["type"] {
				t.Errorf("CreateInputObjectSchema() type = %v, want %v", result["type"], tt.expected["type"])
			}

			// Check description if expected
			if expectedDesc, ok := tt.expected["description"]; ok {
				if result["description"] != expectedDesc {
					t.Errorf("CreateInputObjectSchema() description = %v, want %v", result["description"], expectedDesc)
				}
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

func TestSchema_CreateArgumentSchema(t *testing.T) {
	tests := []struct {
		name     string
		schema   *Schema
		arg      *Argument
		expected map[string]interface{}
	}{
		{
			name: "simple string argument",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{},
			},
			arg: &Argument{
				Name:        "name",
				Description: "User name",
				Type: &TypeRef{
					Name: "String",
					Kind: "SCALAR",
				},
			},
			expected: map[string]interface{}{
				"type":        "string",
				"description": "User name",
			},
		},
		{
			name: "non-null string argument",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{},
			},
			arg: &Argument{
				Name:        "id",
				Description: "User ID",
				Type: &TypeRef{
					Kind: "NON_NULL",
					OfType: &TypeRef{
						Name: "ID",
						Kind: "SCALAR",
					},
				},
			},
			expected: map[string]interface{}{
				"type":        "string",
				"description": "User ID",
			},
		},
		{
			name: "list argument",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{},
			},
			arg: &Argument{
				Name:        "tags",
				Description: "User tags",
				Type: &TypeRef{
					Kind: "LIST",
					OfType: &TypeRef{
						Name: "String",
						Kind: "SCALAR",
					},
				},
			},
			expected: map[string]interface{}{
				"type":        "array",
				"description": "User tags",
				"items": map[string]interface{}{
					"type":        "string",
					"description": "User tags",
				},
			},
		},
		{
			name: "argument with default value",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{},
			},
			arg: &Argument{
				Name:         "limit",
				Description:  "Number of items to return",
				DefaultValue: "10",
				Type: &TypeRef{
					Name: "Int",
					Kind: "SCALAR",
				},
			},
			expected: map[string]interface{}{
				"type":        "integer",
				"description": "Number of items to return",
				"default":     "10",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.CreateArgumentSchema(tt.arg)

			if !mapsEqual(result, tt.expected) {
				t.Errorf("CreateArgumentSchema() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestSchema_CreateInputSchema(t *testing.T) {
	tests := []struct {
		name     string
		schema   *Schema
		field    *Field
		expected map[string]interface{}
	}{
		{
			name: "field without arguments",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{},
			},
			field: &Field{
				Name: "getUser",
				Type: &TypeRef{
					Name: "User",
					Kind: "OBJECT",
				},
				Args: []*Argument{},
			},
			expected: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			name: "field with arguments",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{},
			},
			field: &Field{
				Name: "getUser",
				Type: &TypeRef{
					Name: "User",
					Kind: "OBJECT",
				},
				Args: []*Argument{
					{
						Name:        "id",
						Description: "User ID",
						Type: &TypeRef{
							Kind: "NON_NULL",
							OfType: &TypeRef{
								Name: "ID",
								Kind: "SCALAR",
							},
						},
					},
					{
						Name:        "includeDeleted",
						Description: "Include deleted users",
						Type: &TypeRef{
							Name: "Boolean",
							Kind: "SCALAR",
						},
						DefaultValue: "false",
					},
				},
			},
			expected: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"id": map[string]interface{}{
						"type":        "string",
						"description": "User ID",
					},
					"includeDeleted": map[string]interface{}{
						"type":        "boolean",
						"description": "Include deleted users",
						"default":     "false",
					},
				},
				"required": []string{"id"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.CreateInputSchema(tt.field)

			if !mapsEqual(result, tt.expected) {
				t.Errorf("CreateInputSchema() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Helper function to compare maps recursively
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
