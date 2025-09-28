package schema

import (
	"testing"
)

func TestTypeRef_GetTypeName(t *testing.T) {
	tests := []struct {
		name     string
		typeref  *TypeRef
		expected string
	}{
		{
			name:     "nil TypeRef",
			typeref:  nil,
			expected: "String",
		},
		{
			name: "simple type",
			typeref: &TypeRef{
				Name: "User",
				Kind: "OBJECT",
			},
			expected: "User",
		},
		{
			name: "non-null wrapper",
			typeref: &TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Name: "String",
					Kind: "SCALAR",
				},
			},
			expected: "String",
		},
		{
			name: "list wrapper",
			typeref: &TypeRef{
				Kind: "LIST",
				OfType: &TypeRef{
					Name: "String",
					Kind: "SCALAR",
				},
			},
			expected: "String",
		},
		{
			name: "non-null list wrapper",
			typeref: &TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Kind: "LIST",
					OfType: &TypeRef{
						Name: "String",
						Kind: "SCALAR",
					},
				},
			},
			expected: "String",
		},
		{
			name: "empty name with ofType",
			typeref: &TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					OfType: &TypeRef{
						Name: "Int",
						Kind: "SCALAR",
					},
				},
			},
			expected: "Int",
		},
		{
			name: "no name found",
			typeref: &TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Kind: "LIST",
					OfType: &TypeRef{
						Kind: "SCALAR",
					},
				},
			},
			expected: "String",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.typeref.GetTypeName()
			if result != tt.expected {
				t.Errorf("GetTypeName() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTypeRef_IsList(t *testing.T) {
	tests := []struct {
		name     string
		typeref  *TypeRef
		expected bool
	}{
		{
			name:     "nil TypeRef",
			typeref:  nil,
			expected: false,
		},
		{
			name: "not a list",
			typeref: &TypeRef{
				Name: "String",
				Kind: "SCALAR",
			},
			expected: false,
		},
		{
			name: "is a list",
			typeref: &TypeRef{
				Kind: "LIST",
				OfType: &TypeRef{
					Name: "String",
					Kind: "SCALAR",
				},
			},
			expected: true,
		},
		{
			name: "non-null list",
			typeref: &TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Kind: "LIST",
					OfType: &TypeRef{
						Name: "String",
						Kind: "SCALAR",
					},
				},
			},
			expected: true,
		},
		{
			name: "non-null non-list",
			typeref: &TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Name: "String",
					Kind: "SCALAR",
				},
			},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.typeref.IsList()
			if result != tt.expected {
				t.Errorf("IsList() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTypeRef_IsNonNull(t *testing.T) {
	tests := []struct {
		name     string
		typeref  *TypeRef
		expected bool
	}{
		{
			name:     "nil TypeRef",
			typeref:  nil,
			expected: false,
		},
		{
			name: "not non-null",
			typeref: &TypeRef{
				Name: "String",
				Kind: "SCALAR",
			},
			expected: false,
		},
		{
			name: "is non-null",
			typeref: &TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Name: "String",
					Kind: "SCALAR",
				},
			},
			expected: true,
		},
		{
			name: "list is not non-null",
			typeref: &TypeRef{
				Kind: "LIST",
				OfType: &TypeRef{
					Name: "String",
					Kind: "SCALAR",
				},
			},
			expected: false,
		},
		{
			name: "non-null list is non-null",
			typeref: &TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Kind: "LIST",
					OfType: &TypeRef{
						Name: "String",
						Kind: "SCALAR",
					},
				},
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.typeref.IsNonNull()
			if result != tt.expected {
				t.Errorf("IsNonNull() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestTypeRef_ToJSONSchemaType(t *testing.T) {
	tests := []struct {
		name     string
		typeref  *TypeRef
		expected string
	}{
		{
			name:     "nil TypeRef",
			typeref:  nil,
			expected: "string",
		},
		{
			name: "String type",
			typeref: &TypeRef{
				Name: "String",
				Kind: "SCALAR",
			},
			expected: "string",
		},
		{
			name: "ID type",
			typeref: &TypeRef{
				Name: "ID",
				Kind: "SCALAR",
			},
			expected: "string",
		},
		{
			name: "Int type",
			typeref: &TypeRef{
				Name: "Int",
				Kind: "SCALAR",
			},
			expected: "integer",
		},
		{
			name: "Float type",
			typeref: &TypeRef{
				Name: "Float",
				Kind: "SCALAR",
			},
			expected: "number",
		},
		{
			name: "Boolean type",
			typeref: &TypeRef{
				Name: "Boolean",
				Kind: "SCALAR",
			},
			expected: "boolean",
		},
		{
			name: "Object type",
			typeref: &TypeRef{
				Name: "User",
				Kind: "OBJECT",
			},
			expected: "object",
		},
		{
			name: "non-null String",
			typeref: &TypeRef{
				Kind: "NON_NULL",
				OfType: &TypeRef{
					Name: "String",
					Kind: "SCALAR",
				},
			},
			expected: "string",
		},
		{
			name: "list of String",
			typeref: &TypeRef{
				Kind: "LIST",
				OfType: &TypeRef{
					Name: "String",
					Kind: "SCALAR",
				},
			},
			expected: "string",
		},
		{
			name: "unknown type",
			typeref: &TypeRef{
				Name: "CustomType",
				Kind: "OBJECT",
			},
			expected: "object",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.typeref.ToJSONSchemaType()
			if result != tt.expected {
				t.Errorf("ToJSONSchemaType() = %v, want %v", result, tt.expected)
			}
		})
	}
}
