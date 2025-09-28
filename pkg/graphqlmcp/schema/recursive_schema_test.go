package schema

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestRecursiveInputSchemaGeneration(t *testing.T) {
	// Create a test schema with the recursive AddOrgChainInput type
	schema := createTestSchemaWithRecursiveInputForSchema()

	tests := []struct {
		name        string
		typeName    string
		expectError bool
	}{
		{
			name:        "AddManagerInput should generate without recursion",
			typeName:    "AddManagerInput",
			expectError: false,
		},
		{
			name:        "AddAssociateInput should generate without recursion",
			typeName:    "AddAssociateInput",
			expectError: false,
		},
		{
			name:        "AddOrgChainInput should generate with recursion protection",
			typeName:    "AddOrgChainInput",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// This should not panic or cause stack overflow
			schemaResult := schema.CreateInputObjectSchema(tt.typeName)

			if tt.expectError {
				if schemaResult != nil {
					t.Errorf("expected error but got schema: %v", schemaResult)
				}
			} else {
				if schemaResult == nil {
					t.Errorf("expected schema but got nil")
					return
				}

				// Verify it's a valid object schema
				if schemaResult["type"] != "object" {
					t.Errorf("expected object type, got %v", schemaResult["type"])
				}

				// For recursive types, check that it has proper recursion protection
				if tt.typeName == "AddOrgChainInput" {
					properties, ok := schemaResult["properties"].(map[string]interface{})
					if !ok {
						t.Errorf("expected properties map, got %T", schemaResult["properties"])
						return
					}

					// Check that nextLevel field exists and has recursion protection
					nextLevel, exists := properties["nextLevel"]
					if !exists {
						t.Errorf("expected nextLevel field in AddOrgChainInput")
						return
					}

					nextLevelMap, ok := nextLevel.(map[string]interface{})
					if !ok {
						t.Errorf("expected nextLevel to be a map, got %T", nextLevel)
						return
					}

					// The nextLevel should either be a proper object schema or have recursion protection
					if nextLevelMap["type"] != "object" {
						t.Errorf("expected nextLevel to be object type, got %v", nextLevelMap["type"])
					}
				}
			}
		})
	}
}

func TestRecursiveInputSchemaGeneration_DeepRecursion(t *testing.T) {
	// Create a schema with a deeply recursive type
	schema := createTestSchemaWithRecursiveInputForSchema()

	// This should not cause stack overflow even with deep recursion
	schemaResult := schema.CreateInputObjectSchema("AddOrgChainInput")

	if schemaResult == nil {
		t.Fatalf("expected schema but got nil")
	}

	// Verify the schema has proper structure
	if schemaResult["type"] != "object" {
		t.Errorf("expected object type, got %v", schemaResult["type"])
	}

	properties, ok := schemaResult["properties"].(map[string]interface{})
	if !ok {
		t.Errorf("expected properties map, got %T", schemaResult["properties"])
		return
	}

	// Check that all expected fields are present
	expectedFields := []string{"manager", "associate", "nextLevel"}
	for _, field := range expectedFields {
		if _, exists := properties[field]; !exists {
			t.Errorf("expected field %s not found in schema", field)
		}
	}
}

func TestRecursiveInputSchemaGeneration_CircularReference(t *testing.T) {
	// Create a schema with circular reference
	schema := createTestSchemaWithCircularReference()

	// This should not cause stack overflow
	schemaResult := schema.CreateInputObjectSchema("CircularInput")

	if schemaResult == nil {
		t.Fatalf("expected schema but got nil")
	}

	// Should have recursion protection
	if schemaResult["type"] != "object" {
		t.Errorf("expected object type, got %v", schemaResult["type"])
	}
}

// Helper functions

func createTestSchemaWithRecursiveInputForSchema() *Schema {
	schema := &Schema{
		Types:        []*Type{},
		typeRegistry: make(map[string]*ast.Definition),
	}

	// Add AddManagerInput type
	managerInputType := &Type{
		Name: "AddManagerInput",
		Kind: "INPUT_OBJECT",
		Fields: []*Field{
			{Name: "name", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
			{Name: "email", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
			{Name: "phone", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
			{Name: "department", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
			{Name: "level", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "Int", Kind: "SCALAR"}}},
			{Name: "joinedAt", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
		},
	}

	// Add AddAssociateInput type
	associateInputType := &Type{
		Name: "AddAssociateInput",
		Kind: "INPUT_OBJECT",
		Fields: []*Field{
			{Name: "name", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
			{Name: "email", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
			{Name: "phone", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
			{Name: "jobTitle", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
			{Name: "department", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
			{Name: "reportsToId", Type: &TypeRef{Name: "String", Kind: "SCALAR"}}, // Optional field
			{Name: "joinedAt", Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}}},
		},
	}

	// Add AddOrgChainInput type (recursive)
	orgChainInputType := &Type{
		Name: "AddOrgChainInput",
		Kind: "INPUT_OBJECT",
		Fields: []*Field{
			{
				Name: "manager",
				Type: &TypeRef{
					Kind: "LIST",
					OfType: &TypeRef{
						Name: "AddManagerInput",
						Kind: "INPUT_OBJECT",
					},
				},
			},
			{
				Name: "associate",
				Type: &TypeRef{
					Kind: "LIST",
					OfType: &TypeRef{
						Name: "AddAssociateInput",
						Kind: "INPUT_OBJECT",
					},
				},
			},
			{
				Name: "nextLevel",
				Type: &TypeRef{
					Name: "AddOrgChainInput",
					Kind: "INPUT_OBJECT",
				},
			},
		},
	}

	schema.Types = append(schema.Types, managerInputType, associateInputType, orgChainInputType)

	// Create AST definitions for the type registry
	managerASTDef := &ast.Definition{
		Name: "AddManagerInput",
		Kind: ast.InputObject,
		Fields: []*ast.FieldDefinition{
			{Name: "name", Type: &ast.Type{NamedType: "String", NonNull: true}},
			{Name: "email", Type: &ast.Type{NamedType: "String", NonNull: true}},
			{Name: "phone", Type: &ast.Type{NamedType: "String", NonNull: true}},
			{Name: "department", Type: &ast.Type{NamedType: "String", NonNull: true}},
			{Name: "level", Type: &ast.Type{NamedType: "Int", NonNull: true}},
			{Name: "joinedAt", Type: &ast.Type{NamedType: "String", NonNull: true}},
		},
	}

	associateASTDef := &ast.Definition{
		Name: "AddAssociateInput",
		Kind: ast.InputObject,
		Fields: []*ast.FieldDefinition{
			{Name: "name", Type: &ast.Type{NamedType: "String", NonNull: true}},
			{Name: "email", Type: &ast.Type{NamedType: "String", NonNull: true}},
			{Name: "phone", Type: &ast.Type{NamedType: "String", NonNull: true}},
			{Name: "jobTitle", Type: &ast.Type{NamedType: "String", NonNull: true}},
			{Name: "department", Type: &ast.Type{NamedType: "String", NonNull: true}},
			{Name: "reportsToId", Type: &ast.Type{NamedType: "String", NonNull: false}},
			{Name: "joinedAt", Type: &ast.Type{NamedType: "String", NonNull: true}},
		},
	}

	orgChainASTDef := &ast.Definition{
		Name: "AddOrgChainInput",
		Kind: ast.InputObject,
		Fields: []*ast.FieldDefinition{
			{
				Name: "manager",
				Type: &ast.Type{
					Elem: &ast.Type{
						NamedType: "AddManagerInput",
					},
				},
			},
			{
				Name: "associate",
				Type: &ast.Type{
					Elem: &ast.Type{
						NamedType: "AddAssociateInput",
					},
				},
			},
			{
				Name: "nextLevel",
				Type: &ast.Type{
					NamedType: "AddOrgChainInput",
				},
			},
		},
	}

	// Add to type registry
	schema.typeRegistry["AddManagerInput"] = managerASTDef
	schema.typeRegistry["AddAssociateInput"] = associateASTDef
	schema.typeRegistry["AddOrgChainInput"] = orgChainASTDef

	return schema
}

func createTestSchemaWithCircularReference() *Schema {
	schema := &Schema{
		Types:        []*Type{},
		typeRegistry: make(map[string]*ast.Definition),
	}

	// Create a circular reference: A -> A
	typeA := &Type{
		Name: "CircularInput",
		Kind: "INPUT_OBJECT",
		Fields: []*Field{
			{
				Name: "name",
				Type: &TypeRef{Kind: "NON_NULL", OfType: &TypeRef{Name: "String", Kind: "SCALAR"}},
			},
			{
				Name: "child",
				Type: &TypeRef{
					Name: "CircularInput",
					Kind: "INPUT_OBJECT",
				},
			},
		},
	}

	schema.Types = append(schema.Types, typeA)

	// Create AST definition for circular reference
	circularASTDef := &ast.Definition{
		Name: "CircularInput",
		Kind: ast.InputObject,
		Fields: []*ast.FieldDefinition{
			{Name: "name", Type: &ast.Type{NamedType: "String", NonNull: true}},
			{Name: "child", Type: &ast.Type{NamedType: "CircularInput", NonNull: false}},
		},
	}

	schema.typeRegistry["CircularInput"] = circularASTDef

	return schema
}
