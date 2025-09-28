package schema

import (
	"testing"

	"github.com/vektah/gqlparser/v2/ast"
)

func TestUnionFunctionality(t *testing.T) {
	// Create a test schema with union types
	schema := createTestSchemaWithUnions()

	// Test GetUnions
	unions := schema.GetUnions()
	if len(unions) != 1 {
		t.Errorf("Expected 1 union, got %d", len(unions))
	}

	union := unions[0]
	if union.Name != "EquipmentNotification" {
		t.Errorf("Expected union name 'EquipmentNotification', got '%s'", union.Name)
	}

	if union.Kind != "UNION" {
		t.Errorf("Expected union kind 'UNION', got '%s'", union.Kind)
	}

	// Test GetUnionByName
	unionByName := schema.GetUnionByName("EquipmentNotification")
	if unionByName == nil {
		t.Error("Expected to find union by name")
	}

	if unionByName.Name != "EquipmentNotification" {
		t.Errorf("Expected union name 'EquipmentNotification', got '%s'", unionByName.Name)
	}

	// Test IsUnionType
	if !schema.IsUnionType("EquipmentNotification") {
		t.Error("Expected EquipmentNotification to be a union type")
	}

	if schema.IsUnionType("Equipment") {
		t.Error("Expected Equipment to not be a union type")
	}

	// Test GetUnionPossibleTypes
	possibleTypes := schema.GetUnionPossibleTypes("EquipmentNotification")
	if len(possibleTypes) != 4 {
		t.Errorf("Expected 4 possible types, got %d", len(possibleTypes))
	}

	// Check that all expected types are present
	expectedTypes := []string{"EquipmentAlert", "MaintenanceReminder", "StatusUpdate", "PerformanceAlert"}
	for _, expectedType := range expectedTypes {
		found := false
		for _, possibleType := range possibleTypes {
			if possibleType.Name == expectedType {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected to find possible type '%s'", expectedType)
		}
	}
}

func TestUnionSDLGeneration(t *testing.T) {
	schema := createTestSchemaWithUnions()
	sdl := schema.GetSchemaSDL()

	// Check that the union is included in the SDL
	if !contains(sdl, "union EquipmentNotification") {
		t.Error("Expected SDL to contain union definition")
	}

	if !contains(sdl, "EquipmentAlert | MaintenanceReminder | StatusUpdate | PerformanceAlert") {
		t.Error("Expected SDL to contain union member types")
	}
}

func TestUnionIntrospection(t *testing.T) {
	// Test that union types are properly parsed from introspection
	introspectionData := map[string]interface{}{
		"__schema": map[string]interface{}{
			"queryType": map[string]interface{}{
				"name": "Query",
			},
			"types": []interface{}{
				map[string]interface{}{
					"name":        "EquipmentNotification",
					"kind":        "UNION",
					"description": "Union type representing different types of equipment notifications",
					"possibleTypes": []interface{}{
						map[string]interface{}{
							"name": "EquipmentAlert",
						},
						map[string]interface{}{
							"name": "MaintenanceReminder",
						},
						map[string]interface{}{
							"name": "StatusUpdate",
						},
						map[string]interface{}{
							"name": "PerformanceAlert",
						},
					},
				},
				map[string]interface{}{
					"name": "EquipmentAlert",
					"kind": "OBJECT",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "id",
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"ofType": map[string]interface{}{
									"kind": "SCALAR",
									"name": "ID",
								},
							},
						},
					},
				},
				map[string]interface{}{
					"name": "MaintenanceReminder",
					"kind": "OBJECT",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "id",
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"ofType": map[string]interface{}{
									"kind": "SCALAR",
									"name": "ID",
								},
							},
						},
					},
				},
				map[string]interface{}{
					"name": "StatusUpdate",
					"kind": "OBJECT",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "id",
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"ofType": map[string]interface{}{
									"kind": "SCALAR",
									"name": "ID",
								},
							},
						},
					},
				},
				map[string]interface{}{
					"name": "PerformanceAlert",
					"kind": "OBJECT",
					"fields": []interface{}{
						map[string]interface{}{
							"name": "id",
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"ofType": map[string]interface{}{
									"kind": "SCALAR",
									"name": "ID",
								},
							},
						},
					},
				},
			},
		},
	}

	schema, err := ParseIntrospectionResponse(introspectionData)
	if err != nil {
		t.Fatalf("Failed to parse introspection response: %v", err)
	}

	// Test union functionality
	unions := schema.GetUnions()
	if len(unions) != 1 {
		t.Errorf("Expected 1 union, got %d", len(unions))
	}

	union := unions[0]
	if union.Name != "EquipmentNotification" {
		t.Errorf("Expected union name 'EquipmentNotification', got '%s'", union.Name)
	}

	possibleTypes := schema.GetUnionPossibleTypes("EquipmentNotification")
	if len(possibleTypes) != 4 {
		t.Errorf("Expected 4 possible types, got %d", len(possibleTypes))
	}
}

func createTestSchemaWithUnions() *Schema {
	// Create AST definitions for union and its member types
	unionDef := &ast.Definition{
		Name:        "EquipmentNotification",
		Kind:        ast.Union,
		Description: "Union type representing different types of equipment notifications",
		Types:       []string{"EquipmentAlert", "MaintenanceReminder", "StatusUpdate", "PerformanceAlert"},
	}

	equipmentAlertDef := &ast.Definition{
		Name: "EquipmentAlert",
		Kind: ast.Object,
		Fields: []*ast.FieldDefinition{
			{
				Name: "id",
				Type: &ast.Type{
					NonNull:   true,
					NamedType: "ID",
				},
			},
		},
	}

	maintenanceReminderDef := &ast.Definition{
		Name: "MaintenanceReminder",
		Kind: ast.Object,
		Fields: []*ast.FieldDefinition{
			{
				Name: "id",
				Type: &ast.Type{
					NonNull:   true,
					NamedType: "ID",
				},
			},
		},
	}

	statusUpdateDef := &ast.Definition{
		Name: "StatusUpdate",
		Kind: ast.Object,
		Fields: []*ast.FieldDefinition{
			{
				Name: "id",
				Type: &ast.Type{
					NonNull:   true,
					NamedType: "ID",
				},
			},
		},
	}

	performanceAlertDef := &ast.Definition{
		Name: "PerformanceAlert",
		Kind: ast.Object,
		Fields: []*ast.FieldDefinition{
			{
				Name: "id",
				Type: &ast.Type{
					NonNull:   true,
					NamedType: "ID",
				},
			},
		},
	}

	// Create the schema
	astSchema := &ast.Schema{
		Types: map[string]*ast.Definition{
			"EquipmentNotification": unionDef,
			"EquipmentAlert":        equipmentAlertDef,
			"MaintenanceReminder":   maintenanceReminderDef,
			"StatusUpdate":          statusUpdateDef,
			"PerformanceAlert":      performanceAlertDef,
		},
	}

	return &Schema{
		parsedSchema: astSchema,
		typeRegistry: astSchema.Types,
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(substr) == 0 ||
		(len(s) > len(substr) && (s[:len(substr)] == substr ||
			s[len(s)-len(substr):] == substr ||
			containsSubstring(s, substr))))
}
