package schema

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestIsBuiltinType(t *testing.T) {
	tests := []struct {
		name     string
		typeName string
		expected bool
	}{
		{
			name:     "String",
			typeName: "String",
			expected: true,
		},
		{
			name:     "Int",
			typeName: "Int",
			expected: true,
		},
		{
			name:     "Float",
			typeName: "Float",
			expected: true,
		},
		{
			name:     "Boolean",
			typeName: "Boolean",
			expected: true,
		},
		{
			name:     "ID",
			typeName: "ID",
			expected: true,
		},
		{
			name:     "__Schema",
			typeName: "__Schema",
			expected: true,
		},
		{
			name:     "__Type",
			typeName: "__Type",
			expected: true,
		},
		{
			name:     "__Field",
			typeName: "__Field",
			expected: true,
		},
		{
			name:     "__InputValue",
			typeName: "__InputValue",
			expected: true,
		},
		{
			name:     "__EnumValue",
			typeName: "__EnumValue",
			expected: true,
		},
		{
			name:     "__Directive",
			typeName: "__Directive",
			expected: true,
		},
		{
			name:     "User",
			typeName: "User",
			expected: false,
		},
		{
			name:     "CustomType",
			typeName: "CustomType",
			expected: false,
		},
		{
			name:     "empty string",
			typeName: "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsBuiltinType(tt.typeName)
			if result != tt.expected {
				t.Errorf("IsBuiltinType(%s) = %v, want %v", tt.typeName, result, tt.expected)
			}
		})
	}
}

func TestASTTypeToJSONSchemaType(t *testing.T) {
	tests := []struct {
		name     string
		astType  *ast.Type
		expected string
	}{
		{
			name:     "nil AST type",
			astType:  nil,
			expected: "string",
		},
		{
			name: "String type",
			astType: &ast.Type{
				NamedType: "String",
			},
			expected: "string",
		},
		{
			name: "ID type",
			astType: &ast.Type{
				NamedType: "ID",
			},
			expected: "string",
		},
		{
			name: "Int type",
			astType: &ast.Type{
				NamedType: "Int",
			},
			expected: "integer",
		},
		{
			name: "Float type",
			astType: &ast.Type{
				NamedType: "Float",
			},
			expected: "number",
		},
		{
			name: "Boolean type",
			astType: &ast.Type{
				NamedType: "Boolean",
			},
			expected: "boolean",
		},
		{
			name: "Object type",
			astType: &ast.Type{
				NamedType: "User",
			},
			expected: "object",
		},
		{
			name: "non-null String",
			astType: &ast.Type{
				NonNull: true,
				Elem: &ast.Type{
					NamedType: "String",
				},
			},
			expected: "string",
		},
		{
			name: "list of String",
			astType: &ast.Type{
				Elem: &ast.Type{
					NamedType: "String",
				},
			},
			expected: "string",
		},
		{
			name: "unknown type",
			astType: &ast.Type{
				NamedType: "CustomType",
			},
			expected: "object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ASTTypeToJSONSchemaType(tt.astType)
			if result != tt.expected {
				t.Errorf("ASTTypeToJSONSchemaType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestGetASTTypeName(t *testing.T) {
	tests := []struct {
		name     string
		astType  *ast.Type
		expected string
	}{
		{
			name:     "nil AST type",
			astType:  nil,
			expected: "",
		},
		{
			name: "simple type",
			astType: &ast.Type{
				NamedType: "User",
			},
			expected: "User",
		},
		{
			name: "non-null wrapper",
			astType: &ast.Type{
				NonNull: true,
				Elem: &ast.Type{
					NamedType: "String",
				},
			},
			expected: "String",
		},
		{
			name: "list wrapper",
			astType: &ast.Type{
				Elem: &ast.Type{
					NamedType: "String",
				},
			},
			expected: "String",
		},
		{
			name: "non-null list wrapper",
			astType: &ast.Type{
				NonNull: true,
				Elem: &ast.Type{
					Elem: &ast.Type{
						NamedType: "String",
					},
				},
			},
			expected: "String",
		},
		{
			name: "empty name with elem",
			astType: &ast.Type{
				NonNull: true,
				Elem: &ast.Type{
					Elem: &ast.Type{
						NamedType: "Int",
					},
				},
			},
			expected: "Int",
		},
		{
			name: "no name found",
			astType: &ast.Type{
				NonNull: true,
				Elem: &ast.Type{
					Elem: &ast.Type{},
				},
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := GetASTTypeName(tt.astType)
			if result != tt.expected {
				t.Errorf("GetASTTypeName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsASTTypeList(t *testing.T) {
	tests := []struct {
		name     string
		astType  *ast.Type
		expected bool
	}{
		{
			name:     "nil AST type",
			astType:  nil,
			expected: false,
		},
		{
			name: "not a list",
			astType: &ast.Type{
				NamedType: "String",
			},
			expected: false,
		},
		{
			name: "is a list",
			astType: &ast.Type{
				Elem: &ast.Type{
					NamedType: "String",
				},
			},
			expected: true,
		},
		{
			name: "non-null list",
			astType: &ast.Type{
				NonNull: true,
				Elem: &ast.Type{
					Elem: &ast.Type{
						NamedType: "String",
					},
				},
			},
			expected: true,
		},
		{
			name: "non-null non-list",
			astType: &ast.Type{
				NonNull: true,
				Elem: &ast.Type{
					NamedType: "String",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsASTTypeList(tt.astType)
			if result != tt.expected {
				t.Errorf("IsASTTypeList() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestIsASTTypeNonNull(t *testing.T) {
	tests := []struct {
		name     string
		astType  *ast.Type
		expected bool
	}{
		{
			name:     "nil AST type",
			astType:  nil,
			expected: false,
		},
		{
			name: "not non-null",
			astType: &ast.Type{
				NamedType: "String",
			},
			expected: false,
		},
		{
			name: "is non-null",
			astType: &ast.Type{
				NonNull: true,
				Elem: &ast.Type{
					NamedType: "String",
				},
			},
			expected: true,
		},
		{
			name: "list is not non-null",
			astType: &ast.Type{
				Elem: &ast.Type{
					NamedType: "String",
				},
			},
			expected: false,
		},
		{
			name: "non-null list is non-null",
			astType: &ast.Type{
				NonNull: true,
				Elem: &ast.Type{
					Elem: &ast.Type{
						NamedType: "String",
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsASTTypeNonNull(tt.astType)
			if result != tt.expected {
				t.Errorf("IsASTTypeNonNull() = %v, want %v", result, tt.expected)
			}
		})
	}
}
