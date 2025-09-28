package schema

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestSchema_GetQueries(t *testing.T) {
	tests := []struct {
		name     string
		schema   *Schema
		expected []*Field
	}{
		{
			name: "nil QueryType",
			schema: &Schema{
				QueryType: nil,
			},
			expected: nil,
		},
		{
			name: "empty QueryType",
			schema: &Schema{
				QueryType: &Type{
					Name:   "Query",
					Kind:   "OBJECT",
					Fields: []*Field{},
				},
			},
			expected: []*Field{},
		},
		{
			name: "QueryType with fields",
			schema: &Schema{
				QueryType: &Type{
					Name: "Query",
					Kind: "OBJECT",
					Fields: []*Field{
						{
							Name: "getUser",
							Type: &TypeRef{
								Name: "User",
								Kind: "OBJECT",
							},
						},
						{
							Name: "getUsers",
							Type: &TypeRef{
								Kind: "LIST",
								OfType: &TypeRef{
									Name: "User",
									Kind: "OBJECT",
								},
							},
						},
					},
				},
			},
			expected: []*Field{
				{
					Name: "getUser",
					Type: &TypeRef{
						Name: "User",
						Kind: "OBJECT",
					},
				},
				{
					Name: "getUsers",
					Type: &TypeRef{
						Kind: "LIST",
						OfType: &TypeRef{
							Name: "User",
							Kind: "OBJECT",
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.GetQueries()
			if len(result) != len(tt.expected) {
				t.Errorf("GetQueries() returned %d fields, want %d", len(result), len(tt.expected))
				return
			}
			for i, field := range result {
				if field.Name != tt.expected[i].Name {
					t.Errorf("GetQueries()[%d].Name = %v, want %v", i, field.Name, tt.expected[i].Name)
				}
			}
		})
	}
}

func TestSchema_GetMutations(t *testing.T) {
	tests := []struct {
		name     string
		schema   *Schema
		expected []*Field
	}{
		{
			name: "nil MutationType",
			schema: &Schema{
				MutationType: nil,
			},
			expected: nil,
		},
		{
			name: "empty MutationType",
			schema: &Schema{
				MutationType: &Type{
					Name:   "Mutation",
					Kind:   "OBJECT",
					Fields: []*Field{},
				},
			},
			expected: []*Field{},
		},
		{
			name: "MutationType with fields",
			schema: &Schema{
				MutationType: &Type{
					Name: "Mutation",
					Kind: "OBJECT",
					Fields: []*Field{
						{
							Name: "createUser",
							Type: &TypeRef{
								Name: "User",
								Kind: "OBJECT",
							},
						},
						{
							Name: "updateUser",
							Type: &TypeRef{
								Name: "User",
								Kind: "OBJECT",
							},
						},
					},
				},
			},
			expected: []*Field{
				{
					Name: "createUser",
					Type: &TypeRef{
						Name: "User",
						Kind: "OBJECT",
					},
				},
				{
					Name: "updateUser",
					Type: &TypeRef{
						Name: "User",
						Kind: "OBJECT",
					},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.GetMutations()
			if len(result) != len(tt.expected) {
				t.Errorf("GetMutations() returned %d fields, want %d", len(result), len(tt.expected))
				return
			}
			for i, field := range result {
				if field.Name != tt.expected[i].Name {
					t.Errorf("GetMutations()[%d].Name = %v, want %v", i, field.Name, tt.expected[i].Name)
				}
			}
		})
	}
}

func TestSchema_GetTypeDefinition(t *testing.T) {
	userDef := &ast.Definition{
		Name: "User",
		Kind: ast.Object,
	}
	postDef := &ast.Definition{
		Name: "Post",
		Kind: ast.Object,
	}

	tests := []struct {
		name     string
		schema   *Schema
		typeName string
		expected *ast.Definition
	}{
		{
			name: "nil typeRegistry",
			schema: &Schema{
				typeRegistry: nil,
			},
			typeName: "User",
			expected: nil,
		},
		{
			name: "empty typeRegistry",
			schema: &Schema{
				typeRegistry: make(map[string]*ast.Definition),
			},
			typeName: "User",
			expected: nil,
		},
		{
			name: "type found",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"User": userDef,
					"Post": postDef,
				},
			},
			typeName: "User",
			expected: userDef,
		},
		{
			name: "type not found",
			schema: &Schema{
				typeRegistry: map[string]*ast.Definition{
					"User": userDef,
					"Post": postDef,
				},
			},
			typeName: "Comment",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.GetTypeDefinition(tt.typeName)
			if result != tt.expected {
				t.Errorf("GetTypeDefinition(%s) = %v, want %v", tt.typeName, result, tt.expected)
			}
		})
	}
}

func TestSchema_GetSchemaSDL(t *testing.T) {
	tests := []struct {
		name     string
		schema   *Schema
		expected string
	}{
		{
			name: "nil parsedSchema",
			schema: &Schema{
				parsedSchema: nil,
			},
			expected: "",
		},
		{
			name: "empty parsedSchema",
			schema: &Schema{
				parsedSchema: &ast.Schema{},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.GetSchemaSDL()
			if result != tt.expected {
				t.Errorf("GetSchemaSDL() = %v, want %v", result, tt.expected)
			}
		})
	}
}
