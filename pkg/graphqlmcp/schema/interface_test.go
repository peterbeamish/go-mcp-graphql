package schema

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestSchema_GetInterfaces(t *testing.T) {
	tests := []struct {
		name     string
		schema   *Schema
		expected []string
	}{
		{
			name: "no interfaces",
			schema: &Schema{
				parsedSchema: &ast.Schema{
					Types: map[string]*ast.Definition{
						"User": {
							Name: "User",
							Kind: ast.Object,
						},
					},
				},
			},
			expected: []string{},
		},
		{
			name: "with interfaces",
			schema: &Schema{
				parsedSchema: &ast.Schema{
					Types: map[string]*ast.Definition{
						"Personnel": {
							Name: "Personnel",
							Kind: ast.Interface,
							Fields: []*ast.FieldDefinition{
								{Name: "id", Type: ast.NamedType("ID", nil)},
								{Name: "name", Type: ast.NamedType("String", nil)},
							},
						},
						"Manager": {
							Name:       "Manager",
							Kind:       ast.Object,
							Interfaces: []string{"Personnel"},
						},
						"User": {
							Name: "User",
							Kind: ast.Object,
						},
					},
				},
			},
			expected: []string{"Personnel"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.GetInterfaces()
			if len(result) != len(tt.expected) {
				t.Errorf("GetInterfaces() returned %d interfaces, want %d", len(result), len(tt.expected))
				return
			}
			for i, iface := range result {
				if iface.Name != tt.expected[i] {
					t.Errorf("GetInterfaces()[%d].Name = %v, want %v", i, iface.Name, tt.expected[i])
				}
			}
		})
	}
}

func TestSchema_GetImplementations(t *testing.T) {
	tests := []struct {
		name          string
		schema        *Schema
		interfaceName string
		expected      []string
	}{
		{
			name: "no implementations",
			schema: &Schema{
				parsedSchema: &ast.Schema{
					Types: map[string]*ast.Definition{
						"Personnel": {
							Name: "Personnel",
							Kind: ast.Interface,
						},
					},
				},
			},
			interfaceName: "Personnel",
			expected:      []string{},
		},
		{
			name: "with implementations",
			schema: &Schema{
				parsedSchema: &ast.Schema{
					Types: map[string]*ast.Definition{
						"Personnel": {
							Name: "Personnel",
							Kind: ast.Interface,
						},
						"Manager": {
							Name:       "Manager",
							Kind:       ast.Object,
							Interfaces: []string{"Personnel"},
						},
						"Associate": {
							Name:       "Associate",
							Kind:       ast.Object,
							Interfaces: []string{"Personnel"},
						},
						"User": {
							Name: "User",
							Kind: ast.Object,
						},
					},
				},
			},
			interfaceName: "Personnel",
			expected:      []string{"Manager", "Associate"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.GetImplementations(tt.interfaceName)
			if len(result) != len(tt.expected) {
				t.Errorf("GetImplementations(%s) returned %d implementations, want %d", tt.interfaceName, len(result), len(tt.expected))
				return
			}

			// Check that all expected implementations are present (order doesn't matter)
			found := make(map[string]bool)
			for _, impl := range result {
				found[impl.Name] = true
			}
			for _, expectedName := range tt.expected {
				if !found[expectedName] {
					t.Errorf("GetImplementations(%s) missing expected implementation: %s", tt.interfaceName, expectedName)
				}
			}
		})
	}
}

func TestSchema_GetInterfaceFields(t *testing.T) {
	tests := []struct {
		name           string
		schema         *Schema
		interfaceName  string
		expectedFields []string
	}{
		{
			name: "interface not found",
			schema: &Schema{
				parsedSchema: &ast.Schema{
					Types: map[string]*ast.Definition{},
				},
			},
			interfaceName:  "Personnel",
			expectedFields: []string{},
		},
		{
			name: "not an interface",
			schema: &Schema{
				parsedSchema: &ast.Schema{
					Types: map[string]*ast.Definition{
						"User": {
							Name: "User",
							Kind: ast.Object,
						},
					},
				},
			},
			interfaceName:  "User",
			expectedFields: []string{},
		},
		{
			name: "interface with fields",
			schema: &Schema{
				parsedSchema: &ast.Schema{
					Types: map[string]*ast.Definition{
						"Personnel": {
							Name: "Personnel",
							Kind: ast.Interface,
							Fields: []*ast.FieldDefinition{
								{Name: "id", Type: ast.NamedType("ID", nil)},
								{Name: "name", Type: ast.NamedType("String", nil)},
								{Name: "email", Type: ast.NamedType("String", nil)},
							},
						},
					},
				},
				typeRegistry: map[string]*ast.Definition{
					"Personnel": {
						Name: "Personnel",
						Kind: ast.Interface,
						Fields: []*ast.FieldDefinition{
							{Name: "id", Type: ast.NamedType("ID", nil)},
							{Name: "name", Type: ast.NamedType("String", nil)},
							{Name: "email", Type: ast.NamedType("String", nil)},
						},
					},
				},
			},
			interfaceName:  "Personnel",
			expectedFields: []string{"id", "name", "email"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.schema.GetInterfaceFields(tt.interfaceName)
			if len(result) != len(tt.expectedFields) {
				t.Errorf("GetInterfaceFields(%s) returned %d fields, want %d", tt.interfaceName, len(result), len(tt.expectedFields))
				return
			}
			for i, field := range result {
				if field.Name != tt.expectedFields[i] {
					t.Errorf("GetInterfaceFields(%s)[%d].Name = %v, want %v", tt.interfaceName, i, field.Name, tt.expectedFields[i])
				}
			}
		})
	}
}

func TestField_GenerateInterfaceQueryString(t *testing.T) {
	// Create a test schema with interface and implementations
	schema := &Schema{
		parsedSchema: &ast.Schema{
			Types: map[string]*ast.Definition{
				"Personnel": {
					Name: "Personnel",
					Kind: ast.Interface,
					Fields: []*ast.FieldDefinition{
						{Name: "id", Type: ast.NamedType("ID", nil)},
						{Name: "name", Type: ast.NamedType("String", nil)},
						{Name: "email", Type: ast.NamedType("String", nil)},
					},
				},
				"Manager": {
					Name:       "Manager",
					Kind:       ast.Object,
					Interfaces: []string{"Personnel"},
					Fields: []*ast.FieldDefinition{
						{Name: "id", Type: ast.NamedType("ID", nil)},
						{Name: "name", Type: ast.NamedType("String", nil)},
						{Name: "email", Type: ast.NamedType("String", nil)},
						{Name: "department", Type: ast.NamedType("String", nil)},
						{Name: "level", Type: ast.NamedType("Int", nil)},
					},
				},
				"Associate": {
					Name:       "Associate",
					Kind:       ast.Object,
					Interfaces: []string{"Personnel"},
					Fields: []*ast.FieldDefinition{
						{Name: "id", Type: ast.NamedType("ID", nil)},
						{Name: "name", Type: ast.NamedType("String", nil)},
						{Name: "email", Type: ast.NamedType("String", nil)},
						{Name: "jobTitle", Type: ast.NamedType("String", nil)},
						{Name: "reportsTo", Type: ast.NamedType("Manager", nil)},
					},
				},
			},
		},
		typeRegistry: map[string]*ast.Definition{
			"Personnel": {
				Name: "Personnel",
				Kind: ast.Interface,
				Fields: []*ast.FieldDefinition{
					{Name: "id", Type: ast.NamedType("ID", nil)},
					{Name: "name", Type: ast.NamedType("String", nil)},
					{Name: "email", Type: ast.NamedType("String", nil)},
				},
			},
			"Manager": {
				Name:       "Manager",
				Kind:       ast.Object,
				Interfaces: []string{"Personnel"},
				Fields: []*ast.FieldDefinition{
					{Name: "id", Type: ast.NamedType("ID", nil)},
					{Name: "name", Type: ast.NamedType("String", nil)},
					{Name: "email", Type: ast.NamedType("String", nil)},
					{Name: "department", Type: ast.NamedType("String", nil)},
					{Name: "level", Type: ast.NamedType("Int", nil)},
				},
			},
			"Associate": {
				Name:       "Associate",
				Kind:       ast.Object,
				Interfaces: []string{"Personnel"},
				Fields: []*ast.FieldDefinition{
					{Name: "id", Type: ast.NamedType("ID", nil)},
					{Name: "name", Type: ast.NamedType("String", nil)},
					{Name: "email", Type: ast.NamedType("String", nil)},
					{Name: "jobTitle", Type: ast.NamedType("String", nil)},
					{Name: "reportsTo", Type: ast.NamedType("Manager", nil)},
				},
			},
		},
	}

	// Create a field that returns Personnel interface
	field := &Field{
		Name: "personnel",
		Type: &TypeRef{
			Name: "Personnel",
			Kind: "INTERFACE",
		},
		ASTType: ast.NamedType("Personnel", nil),
	}

	result, err := field.GenerateInterfaceQueryString(schema)
	if err != nil {
		t.Fatalf("GenerateInterfaceQueryString() unexpected error: %v", err)
	}

	// Check that the result contains interface fields
	expectedInterfaceFields := []string{"id", "name", "email", "__typename"}
	for _, expectedField := range expectedInterfaceFields {
		if !containsString(result, expectedField) {
			t.Errorf("Expected query to contain interface field '%s'", expectedField)
		}
	}

	// Check that the result contains inline fragments
	expectedFragments := []string{"... on Manager", "... on Associate"}
	for _, expectedFragment := range expectedFragments {
		if !containsString(result, expectedFragment) {
			t.Errorf("Expected query to contain inline fragment '%s'", expectedFragment)
		}
	}

	// Check that Manager-specific fields are included
	expectedManagerFields := []string{"department", "level"}
	for _, expectedField := range expectedManagerFields {
		if !containsString(result, expectedField) {
			t.Errorf("Expected query to contain Manager field '%s'", expectedField)
		}
	}

	// Check that Associate-specific fields are included
	expectedAssociateFields := []string{"jobTitle", "reportsTo"}
	for _, expectedField := range expectedAssociateFields {
		if !containsString(result, expectedField) {
			t.Errorf("Expected query to contain Associate field '%s'", expectedField)
		}
	}

	t.Logf("Generated interface query:\n%s", result)
}

func TestField_GenerateInterfaceQueryString_NonInterface(t *testing.T) {
	// Create a test schema with a regular object type
	schema := &Schema{
		parsedSchema: &ast.Schema{
			Types: map[string]*ast.Definition{
				"User": {
					Name: "User",
					Kind: ast.Object,
					Fields: []*ast.FieldDefinition{
						{Name: "id", Type: ast.NamedType("ID", nil)},
						{Name: "name", Type: ast.NamedType("String", nil)},
					},
				},
			},
		},
		typeRegistry: map[string]*ast.Definition{
			"User": {
				Name: "User",
				Kind: ast.Object,
				Fields: []*ast.FieldDefinition{
					{Name: "id", Type: ast.NamedType("ID", nil)},
					{Name: "name", Type: ast.NamedType("String", nil)},
				},
			},
		},
	}

	// Create a field that returns User object (not interface)
	field := &Field{
		Name: "user",
		Type: &TypeRef{
			Name: "User",
			Kind: "OBJECT",
		},
		ASTType: ast.NamedType("User", nil),
	}

	result, err := field.GenerateInterfaceQueryString(schema)
	if err != nil {
		t.Fatalf("GenerateInterfaceQueryString() unexpected error: %v", err)
	}

	// Should fall back to regular query generation
	if !containsString(result, "user") {
		t.Errorf("Expected query to contain field name 'user'")
	}
}
