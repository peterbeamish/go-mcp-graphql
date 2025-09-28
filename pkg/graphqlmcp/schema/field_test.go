package schema

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestField_GenerateQueryStringWithSchema(t *testing.T) {
	tests := []struct {
		name        string
		field       *Field
		schema      *Schema
		expectedErr bool
		contains    []string
	}{
		{
			name: "nil schema",
			field: &Field{
				Name: "getUser",
				Type: &TypeRef{
					Name: "User",
					Kind: "OBJECT",
				},
			},
			schema:      nil,
			expectedErr: true,
		},
		{
			name: "simple field without arguments",
			field: &Field{
				Name: "getUser",
				Type: &TypeRef{
					Name: "User",
					Kind: "OBJECT",
				},
				ASTType: &ast.Type{
					NamedType: "User",
				},
			},
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"User": {
						Name: "User",
						Kind: ast.Object,
						Fields: []*ast.FieldDefinition{
							{
								Name: "id",
								Type: &ast.Type{
									NamedType: "ID",
								},
							},
							{
								Name: "name",
								Type: &ast.Type{
									NamedType: "String",
								},
							},
						},
					},
				},
			},
			expectedErr: false,
			contains:    []string{"query", "getUser", "id", "name"},
		},
		{
			name: "field with arguments",
			field: &Field{
				Name: "getUser",
				Type: &TypeRef{
					Name: "User",
					Kind: "OBJECT",
				},
				Args: []*Argument{
					{
						Name: "id",
						Type: &TypeRef{
							Kind: "NON_NULL",
							OfType: &TypeRef{
								Name: "ID",
								Kind: "SCALAR",
							},
						},
					},
				},
				ASTType: &ast.Type{
					NamedType: "User",
				},
			},
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"User": {
						Name: "User",
						Kind: ast.Object,
						Fields: []*ast.FieldDefinition{
							{
								Name: "id",
								Type: &ast.Type{
									NamedType: "ID",
								},
							},
							{
								Name: "name",
								Type: &ast.Type{
									NamedType: "String",
								},
							},
						},
					},
				},
			},
			expectedErr: false,
			contains:    []string{"query", "getUser", "$id: ID!", "id: $id", "id", "name"},
		},
		{
			name: "field with list return type",
			field: &Field{
				Name: "getUsers",
				Type: &TypeRef{
					Kind: "LIST",
					OfType: &TypeRef{
						Name: "User",
						Kind: "OBJECT",
					},
				},
				ASTType: &ast.Type{
					Elem: &ast.Type{
						NamedType: "User",
					},
				},
			},
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"User": {
						Name: "User",
						Kind: ast.Object,
						Fields: []*ast.FieldDefinition{
							{
								Name: "id",
								Type: &ast.Type{
									NamedType: "ID",
								},
							},
							{
								Name: "name",
								Type: &ast.Type{
									NamedType: "String",
								},
							},
						},
					},
				},
			},
			expectedErr: false,
			contains:    []string{"query", "getUsers", "id", "name"},
		},
		{
			name: "field with scalar return type",
			field: &Field{
				Name: "getCount",
				Type: &TypeRef{
					Name: "Int",
					Kind: "SCALAR",
				},
				ASTType: &ast.Type{
					NamedType: "Int",
				},
			},
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"Int": {
						Name: "Int",
						Kind: ast.Scalar,
					},
				},
			},
			expectedErr: false,
			contains:    []string{"query", "getCount"},
		},
		{
			name: "field with missing type definition",
			field: &Field{
				Name: "getUser",
				Type: &TypeRef{
					Name: "User",
					Kind: "OBJECT",
				},
				ASTType: &ast.Type{
					NamedType: "User",
				},
			},
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{},
			},
			expectedErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.field.GenerateQueryStringWithSchema(tt.schema)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("GenerateQueryStringWithSchema() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateQueryStringWithSchema() unexpected error: %v", err)
				return
			}

			for _, expected := range tt.contains {
				if !containsString(result, expected) {
					t.Errorf("GenerateQueryStringWithSchema() result does not contain expected string: %s", expected)
				}
			}
		})
	}
}

func TestField_GenerateMutationStringWithSchema(t *testing.T) {
	tests := []struct {
		name        string
		field       *Field
		schema      *Schema
		expectedErr bool
		contains    []string
	}{
		{
			name: "nil schema",
			field: &Field{
				Name: "createUser",
				Type: &TypeRef{
					Name: "User",
					Kind: "OBJECT",
				},
			},
			schema:      nil,
			expectedErr: true,
		},
		{
			name: "simple mutation without arguments",
			field: &Field{
				Name: "createUser",
				Type: &TypeRef{
					Name: "User",
					Kind: "OBJECT",
				},
				ASTType: &ast.Type{
					NamedType: "User",
				},
			},
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"User": {
						Name: "User",
						Kind: ast.Object,
						Fields: []*ast.FieldDefinition{
							{
								Name: "id",
								Type: &ast.Type{
									NamedType: "ID",
								},
							},
							{
								Name: "name",
								Type: &ast.Type{
									NamedType: "String",
								},
							},
						},
					},
				},
			},
			expectedErr: false,
			contains:    []string{"mutation", "createUser", "id", "name"},
		},
		{
			name: "mutation with arguments",
			field: &Field{
				Name: "createUser",
				Type: &TypeRef{
					Name: "User",
					Kind: "OBJECT",
				},
				Args: []*Argument{
					{
						Name: "input",
						Type: &TypeRef{
							Kind: "NON_NULL",
							OfType: &TypeRef{
								Name: "CreateUserInput",
								Kind: "INPUT_OBJECT",
							},
						},
					},
				},
				ASTType: &ast.Type{
					NamedType: "User",
				},
			},
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"User": {
						Name: "User",
						Kind: ast.Object,
						Fields: []*ast.FieldDefinition{
							{
								Name: "id",
								Type: &ast.Type{
									NamedType: "ID",
								},
							},
							{
								Name: "name",
								Type: &ast.Type{
									NamedType: "String",
								},
							},
						},
					},
				},
			},
			expectedErr: false,
			contains:    []string{"mutation", "createUser", "$input: CreateUserInput!", "input: $input", "id", "name"},
		},
		{
			name: "mutation with scalar return type",
			field: &Field{
				Name: "deleteUser",
				Type: &TypeRef{
					Name: "Boolean",
					Kind: "SCALAR",
				},
				ASTType: &ast.Type{
					NamedType: "Boolean",
				},
			},
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"Boolean": {
						Name: "Boolean",
						Kind: ast.Scalar,
					},
				},
			},
			expectedErr: false,
			contains:    []string{"mutation", "deleteUser"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tt.field.GenerateMutationStringWithSchema(tt.schema)

			if tt.expectedErr {
				if err == nil {
					t.Errorf("GenerateMutationStringWithSchema() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("GenerateMutationStringWithSchema() unexpected error: %v", err)
				return
			}

			for _, expected := range tt.contains {
				if !containsString(result, expected) {
					t.Errorf("GenerateMutationStringWithSchema() result does not contain expected string: %s", expected)
				}
			}
		})
	}
}

// Helper function to check if a string contains a substring
func containsString(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func TestField_GenerateQueryStringWithSchema_Personnel(t *testing.T) {
	// Create a schema with Personnel interface
	// Use the same definitions for both parsedSchema and typeRegistry
	personnelDef := &ast.Definition{
		Name: "Personnel",
		Kind: ast.Interface,
		Fields: []*ast.FieldDefinition{
			{Name: "id", Type: ast.NamedType("ID", nil)},
			{Name: "name", Type: ast.NamedType("String", nil)},
			{Name: "email", Type: ast.NamedType("String", nil)},
			{Name: "joinedAt", Type: ast.NamedType("String", nil)},
		},
	}

	managerDef := &ast.Definition{
		Name:       "Manager",
		Kind:       ast.Object,
		Interfaces: []string{"Personnel"},
		Fields: []*ast.FieldDefinition{
			{Name: "id", Type: ast.NamedType("ID", nil)},
			{Name: "name", Type: ast.NamedType("String", nil)},
			{Name: "email", Type: ast.NamedType("String", nil)},
			{Name: "joinedAt", Type: ast.NamedType("String", nil)},
			{Name: "department", Type: ast.NamedType("String", nil)},
			{Name: "level", Type: ast.NamedType("Int", nil)},
		},
	}

	associateDef := &ast.Definition{
		Name:       "Associate",
		Kind:       ast.Object,
		Interfaces: []string{"Personnel"},
		Fields: []*ast.FieldDefinition{
			{Name: "id", Type: ast.NamedType("ID", nil)},
			{Name: "name", Type: ast.NamedType("String", nil)},
			{Name: "email", Type: ast.NamedType("String", nil)},
			{Name: "joinedAt", Type: ast.NamedType("String", nil)},
			{Name: "jobTitle", Type: ast.NamedType("String", nil)},
			{Name: "reportsTo", Type: ast.NamedType("Manager", nil)},
		},
	}

	schema := &Schema{
		parsedSchema: &ast.Schema{
			Types: map[string]*ast.Definition{
				"Personnel": personnelDef,
				"Manager":   managerDef,
				"Associate": associateDef,
			},
		},
		typeRegistry: map[string]*ast.Definition{
			"Personnel": personnelDef,
			"Manager":   managerDef,
			"Associate": associateDef,
		},
	}

	// Create a field that returns Personnel array
	field := &Field{
		Name:    "personnel",
		ASTType: ast.ListType(ast.NamedType("Personnel", nil), nil),
	}

	// Generate the query
	query, err := field.GenerateQueryStringWithSchema(schema)
	if err != nil {
		t.Fatalf("Failed to generate query: %v", err)
	}

	// Print the generated query for demonstration
	t.Logf("Generated Personnel Interface Query:\n%s", query)

	// Check that the query contains the expected elements
	expectedElements := []string{
		"query {",
		"personnel {",
		"id",
		"name",
		"email",
		"joinedAt",
		"__typename",
		"... on Manager {",
		"department",
		"level",
		"... on Associate {",
		"jobTitle",
		"reportsTo",
	}

	for _, element := range expectedElements {
		if !containsString(query, element) {
			t.Errorf("Query should contain '%s', but got: %s", element, query)
		}
	}

	// Verify the structure is correct
	if !containsString(query, "query {\n  personnel {\n    ") {
		t.Errorf("Query should have proper structure, but got: %s", query)
	}
}
