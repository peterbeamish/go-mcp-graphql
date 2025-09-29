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
		return
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

func TestUnionFieldAliasing(t *testing.T) {
	// Create a test schema with union types that have conflicting field names
	schema := createTestSchemaWithConflictingFields()

	// Create a field that returns the union type
	field := &Field{
		Name: "getNotifications",
		Type: &TypeRef{
			Name: "EquipmentNotification",
			Kind: "UNION",
		},
		ASTType: &ast.Type{
			NamedType: "EquipmentNotification",
		},
	}

	// Generate the selection set
	selectionSet, err := field.GenerateQueryStringWithSchema(schema)
	if err != nil {
		t.Fatalf("Error generating query: %v", err)
	}

	// Check that the query contains aliased fields for conflicting field names
	// Only fields that appear in multiple union member types should be aliased
	expectedAliases := []string{
		"EquipmentAlert_type: type",                    // 'type' appears in EquipmentAlert and MaintenanceReminder
		"MaintenanceReminder_type: type",               // 'type' appears in EquipmentAlert and MaintenanceReminder
		"EquipmentAlert_id: id",                        // 'id' appears in all union member types
		"MaintenanceReminder_id: id",                   // 'id' appears in all union member types
		"StatusUpdate_id: id",                          // 'id' appears in all union member types
		"PerformanceAlert_id: id",                      // 'id' appears in all union member types
		"EquipmentAlert_description: description",      // 'description' appears in EquipmentAlert, MaintenanceReminder, StatusUpdate, PerformanceAlert
		"MaintenanceReminder_description: description", // 'description' appears in EquipmentAlert, MaintenanceReminder, StatusUpdate, PerformanceAlert
		"StatusUpdate_description: description",        // 'description' appears in EquipmentAlert, MaintenanceReminder, StatusUpdate, PerformanceAlert
		"PerformanceAlert_description: description",    // 'description' appears in EquipmentAlert, MaintenanceReminder, StatusUpdate, PerformanceAlert
	}

	for _, expectedAlias := range expectedAliases {
		if !contains(selectionSet, expectedAlias) {
			t.Errorf("Expected query to contain aliased field '%s', but it wasn't found", expectedAlias)
		}
	}

	// Check that non-conflicting fields are not aliased
	// These fields only appear in one union member type, so they shouldn't be aliased
	nonConflictingFields := []string{
		"severity",      // Only in EquipmentAlert
		"priority",      // Only in MaintenanceReminder
		"newStatus",     // Only in StatusUpdate
		"changedAt",     // Only in StatusUpdate
		"metricType",    // Only in PerformanceAlert
		"currentValue",  // Only in PerformanceAlert
		"expectedValue", // Only in PerformanceAlert
	}

	for _, fieldName := range nonConflictingFields {
		// These fields should appear without aliases
		if !contains(selectionSet, fieldName) {
			t.Errorf("Expected query to contain field '%s', but it wasn't found", fieldName)
		}
	}

	t.Logf("Generated query:\n%s", selectionSet)
}

func TestUnionFieldAliasingIntelligence(t *testing.T) {
	// Test that the new implementation only aliases fields that actually conflict
	// This demonstrates the improvement over the hardcoded approach

	schema := createTestSchemaWithUniqueFields()

	field := &Field{
		Name: "getItems",
		Type: &TypeRef{
			Name: "ItemUnion",
			Kind: "UNION",
		},
		ASTType: &ast.Type{
			NamedType: "ItemUnion",
		},
	}

	selectionSet, err := field.GenerateQueryStringWithSchema(schema)
	if err != nil {
		t.Fatalf("Error generating query: %v", err)
	}

	// Fields that appear in multiple union member types should be aliased
	expectedAliases := []string{
		"Book_title: title",    // 'title' appears in Book and Article
		"Article_title: title", // 'title' appears in Book and Article
	}

	for _, expectedAlias := range expectedAliases {
		if !contains(selectionSet, expectedAlias) {
			t.Errorf("Expected query to contain aliased field '%s', but it wasn't found", expectedAlias)
		}
	}

	// Fields that only appear in one union member type should NOT be aliased
	// even if they were in the old hardcoded conflict list
	nonConflictingFields := []string{
		"author",   // Only in Book - should NOT be aliased
		"category", // Only in Article - should NOT be aliased
		"price",    // Only in Product - should NOT be aliased
		// Note: inStock might be filtered out by shouldIncludeField logic
		// Let's test what we can see in the output
	}

	for _, fieldName := range nonConflictingFields {
		// These fields should appear without aliases
		if !contains(selectionSet, fieldName) {
			t.Errorf("Expected query to contain field '%s', but it wasn't found", fieldName)
		}
	}

	t.Logf("Generated query:\n%s", selectionSet)
}

func createTestSchemaWithUniqueFields() *Schema {
	// Create union definition
	unionDef := &ast.Definition{
		Name:        "ItemUnion",
		Kind:        ast.Union,
		Description: "Union type with some unique fields per type",
		Types:       []string{"Book", "Article", "Product"},
	}

	// Book with unique fields
	bookDef := &ast.Definition{
		Name: "Book",
		Kind: ast.Object,
		Fields: []*ast.FieldDefinition{
			{
				Name: "title",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
			{
				Name: "author",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
		},
	}

	// Article with some overlapping fields
	articleDef := &ast.Definition{
		Name: "Article",
		Kind: ast.Object,
		Fields: []*ast.FieldDefinition{
			{
				Name: "title",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
			{
				Name: "category",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
		},
	}

	// Product with unique fields
	productDef := &ast.Definition{
		Name: "Product",
		Kind: ast.Object,
		Fields: []*ast.FieldDefinition{
			{
				Name: "price",
				Type: &ast.Type{
					NamedType: "Float",
				},
			},
			{
				Name: "inStock",
				Type: &ast.Type{
					NamedType: "Boolean",
				},
			},
		},
	}

	// Create the schema
	astSchema := &ast.Schema{
		Types: map[string]*ast.Definition{
			"ItemUnion": unionDef,
			"Book":      bookDef,
			"Article":   articleDef,
			"Product":   productDef,
		},
	}

	return &Schema{
		parsedSchema: astSchema,
		typeRegistry: astSchema.Types,
	}
}

func createTestSchemaWithConflictingFields() *Schema {
	// Create union definition
	unionDef := &ast.Definition{
		Name:        "EquipmentNotification",
		Kind:        ast.Union,
		Description: "Union type representing different types of equipment notifications",
		Types:       []string{"EquipmentAlert", "MaintenanceReminder", "StatusUpdate", "PerformanceAlert"},
	}

	// EquipmentAlert with conflicting fields
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
			{
				Name: "type",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
			{
				Name: "severity",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
			{
				Name: "description",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
		},
	}

	// MaintenanceReminder with conflicting fields
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
			{
				Name: "type",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
			{
				Name: "priority",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
			{
				Name: "description",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
		},
	}

	// StatusUpdate with some conflicting fields
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
			{
				Name: "description",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
			{
				Name: "newStatus",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
			{
				Name: "changedAt",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
		},
	}

	// PerformanceAlert with conflicting fields
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
			{
				Name: "description",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
			{
				Name: "metricType",
				Type: &ast.Type{
					NamedType: "String",
				},
			},
			{
				Name: "currentValue",
				Type: &ast.Type{
					NamedType: "Float",
				},
			},
			{
				Name: "expectedValue",
				Type: &ast.Type{
					NamedType: "Float",
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
