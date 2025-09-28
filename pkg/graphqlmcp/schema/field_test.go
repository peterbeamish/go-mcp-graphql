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
