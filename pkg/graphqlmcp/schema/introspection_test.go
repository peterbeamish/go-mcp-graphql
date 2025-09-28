package schema

import (
	"encoding/json"
	"os"
	"testing"
)

func TestParseIntrospectionResponse(t *testing.T) {
	tests := []struct {
		name     string
		input    map[string]interface{}
		expected func(*Schema) bool
		wantErr  bool
	}{
		{
			name: "basic schema with input types",
			input: map[string]interface{}{
				"__schema": map[string]interface{}{
					"queryType": map[string]interface{}{
						"name": "Query",
					},
					"mutationType": map[string]interface{}{
						"name": "Mutation",
					},
					"types": []interface{}{
						// Query type
						map[string]interface{}{
							"name":        "Query",
							"kind":        "OBJECT",
							"description": "Root query type",
							"fields": []interface{}{
								map[string]interface{}{
									"name":        "getUser",
									"description": "Get a user by ID",
									"type": map[string]interface{}{
										"name": "User",
										"kind": "OBJECT",
									},
									"args": []interface{}{
										map[string]interface{}{
											"name":        "id",
											"description": "User ID",
											"type": map[string]interface{}{
												"kind": "NON_NULL",
												"ofType": map[string]interface{}{
													"name": "ID",
													"kind": "SCALAR",
												},
											},
										},
									},
								},
							},
						},
						// Mutation type
						map[string]interface{}{
							"name":        "Mutation",
							"kind":        "OBJECT",
							"description": "Root mutation type",
							"fields": []interface{}{
								map[string]interface{}{
									"name":        "createUser",
									"description": "Create a new user",
									"type": map[string]interface{}{
										"name": "User",
										"kind": "OBJECT",
									},
									"args": []interface{}{
										map[string]interface{}{
											"name":        "input",
											"description": "User input data",
											"type": map[string]interface{}{
												"kind": "NON_NULL",
												"ofType": map[string]interface{}{
													"name": "CreateUserInput",
													"kind": "INPUT_OBJECT",
												},
											},
										},
									},
								},
							},
						},
						// User type
						map[string]interface{}{
							"name":        "User",
							"kind":        "OBJECT",
							"description": "A user in the system",
							"fields": []interface{}{
								map[string]interface{}{
									"name":        "id",
									"description": "User ID",
									"type": map[string]interface{}{
										"kind": "NON_NULL",
										"ofType": map[string]interface{}{
											"name": "ID",
											"kind": "SCALAR",
										},
									},
								},
								map[string]interface{}{
									"name":        "name",
									"description": "User name",
									"type": map[string]interface{}{
										"kind": "NON_NULL",
										"ofType": map[string]interface{}{
											"name": "String",
											"kind": "SCALAR",
										},
									},
								},
							},
						},
						// Input type
						map[string]interface{}{
							"name":        "CreateUserInput",
							"kind":        "INPUT_OBJECT",
							"description": "Input for creating a user",
							"inputFields": []interface{}{
								map[string]interface{}{
									"name":        "name",
									"description": "User name",
									"type": map[string]interface{}{
										"kind": "NON_NULL",
										"ofType": map[string]interface{}{
											"name": "String",
											"kind": "SCALAR",
										},
									},
								},
								map[string]interface{}{
									"name":        "email",
									"description": "User email",
									"type": map[string]interface{}{
										"kind": "NON_NULL",
										"ofType": map[string]interface{}{
											"name": "String",
											"kind": "SCALAR",
										},
									},
								},
								map[string]interface{}{
									"name":        "age",
									"description": "User age",
									"type": map[string]interface{}{
										"name": "Int",
										"kind": "SCALAR",
									},
								},
							},
						},
						// Enum type
						map[string]interface{}{
							"name":        "UserRole",
							"kind":        "ENUM",
							"description": "User roles",
							"enumValues": []interface{}{
								map[string]interface{}{
									"name":        "ADMIN",
									"description": "Administrator role",
								},
								map[string]interface{}{
									"name":        "USER",
									"description": "Regular user role",
								},
							},
						},
						// Scalar types
						map[string]interface{}{
							"name":        "String",
							"kind":        "SCALAR",
							"description": "String scalar type",
						},
						map[string]interface{}{
							"name":        "Int",
							"kind":        "SCALAR",
							"description": "Int scalar type",
						},
						map[string]interface{}{
							"name":        "ID",
							"kind":        "SCALAR",
							"description": "ID scalar type",
						},
					},
				},
			},
			expected: func(schema *Schema) bool {
				if schema == nil {
					return false
				}

				// Check that we have the expected types
				if len(schema.Types) == 0 {
					t.Logf("Expected types but got none")
					return false
				}

				// Check for input type
				var inputType *Type
				for _, typ := range schema.Types {
					if typ.Name == "CreateUserInput" {
						inputType = typ
						break
					}
				}

				if inputType == nil {
					t.Logf("Expected CreateUserInput type but not found")
					return false
				}

				if inputType.Kind != "INPUT_OBJECT" {
					t.Logf("Expected INPUT_OBJECT kind but got %s", inputType.Kind)
					return false
				}

				if len(inputType.Fields) != 3 {
					t.Logf("Expected 3 input fields but got %d", len(inputType.Fields))
					return false
				}

				// Check input fields
				fieldNames := make(map[string]bool)
				for _, field := range inputType.Fields {
					fieldNames[field.Name] = true

					if field.Name == "name" && field.Type.GetTypeName() != "String" {
						t.Logf("Expected name field to have String type but got %s", field.Type.GetTypeName())
						return false
					}

					if field.Name == "email" && field.Type.GetTypeName() != "String" {
						t.Logf("Expected email field to have String type but got %s", field.Type.GetTypeName())
						return false
					}

					if field.Name == "age" && field.Type.GetTypeName() != "Int" {
						t.Logf("Expected age field to have Int type but got %s", field.Type.GetTypeName())
						return false
					}
				}

				expectedFields := []string{"name", "email", "age"}
				for _, expectedField := range expectedFields {
					if !fieldNames[expectedField] {
						t.Logf("Expected field %s not found", expectedField)
						return false
					}
				}

				return true
			},
			wantErr: false,
		},
		{
			name: "invalid schema data",
			input: map[string]interface{}{
				"invalid": "data",
			},
			expected: func(schema *Schema) bool {
				return schema == nil
			},
			wantErr: true,
		},
		{
			name: "empty schema",
			input: map[string]interface{}{
				"__schema": map[string]interface{}{
					"types": []interface{}{},
				},
			},
			expected: func(schema *Schema) bool {
				return schema != nil && len(schema.Types) == 0
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseIntrospectionResponse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseIntrospectionResponse() expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("ParseIntrospectionResponse() unexpected error: %v", err)
				return
			}

			if !tt.expected(result) {
				t.Errorf("ParseIntrospectionResponse() result validation failed")
			}

			// Additional debugging output
			if result != nil {
				t.Logf("Parsed schema has %d types", len(result.Types))
				for _, typ := range result.Types {
					t.Logf("  Type: %s (kind: %s, fields: %d)", typ.Name, typ.Kind, len(typ.Fields))
					for _, field := range typ.Fields {
						t.Logf("    Field: %s (type: %s)", field.Name, field.Type.GetTypeName())
					}
				}
			}
		})
	}
}

func TestParseIntrospectionResponse_ComplexInputTypes(t *testing.T) {
	// Test with more complex input types that have nested structures
	input := map[string]interface{}{
		"__schema": map[string]interface{}{
			"queryType": map[string]interface{}{
				"name": "Query",
			},
			"mutationType": map[string]interface{}{
				"name": "Mutation",
			},
			"types": []interface{}{
				// Nested input type
				map[string]interface{}{
					"name":        "CreateUserInput",
					"kind":        "INPUT_OBJECT",
					"description": "Input for creating a user",
					"inputFields": []interface{}{
						map[string]interface{}{
							"name":        "name",
							"description": "User name",
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"ofType": map[string]interface{}{
									"name": "String",
									"kind": "SCALAR",
								},
							},
						},
						map[string]interface{}{
							"name":        "profile",
							"description": "User profile",
							"type": map[string]interface{}{
								"kind": "NON_NULL",
								"ofType": map[string]interface{}{
									"name": "UserProfileInput",
									"kind": "INPUT_OBJECT",
								},
							},
						},
						map[string]interface{}{
							"name":        "tags",
							"description": "User tags",
							"type": map[string]interface{}{
								"kind": "LIST",
								"ofType": map[string]interface{}{
									"kind": "NON_NULL",
									"ofType": map[string]interface{}{
										"name": "String",
										"kind": "SCALAR",
									},
								},
							},
						},
					},
				},
				// Nested input type
				map[string]interface{}{
					"name":        "UserProfileInput",
					"kind":        "INPUT_OBJECT",
					"description": "User profile input",
					"inputFields": []interface{}{
						map[string]interface{}{
							"name":        "bio",
							"description": "User biography",
							"type": map[string]interface{}{
								"name": "String",
								"kind": "SCALAR",
							},
						},
						map[string]interface{}{
							"name":        "age",
							"description": "User age",
							"type": map[string]interface{}{
								"name": "Int",
								"kind": "SCALAR",
							},
						},
					},
				},
				// Scalar types
				map[string]interface{}{
					"name":        "String",
					"kind":        "SCALAR",
					"description": "String scalar type",
				},
				map[string]interface{}{
					"name":        "Int",
					"kind":        "SCALAR",
					"description": "Int scalar type",
				},
			},
		},
	}

	result, err := ParseIntrospectionResponse(input)
	if err != nil {
		t.Fatalf("ParseIntrospectionResponse() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("ParseIntrospectionResponse() returned nil result")
	}

	// Check that we have the expected types
	if len(result.Types) == 0 {
		t.Fatal("Expected types but got none")
	}

	// Check for input type
	var inputType *Type
	for _, typ := range result.Types {
		if typ.Name == "CreateUserInput" {
			inputType = typ
			break
		}
	}

	if inputType == nil {
		t.Fatal("Expected CreateUserInput type but not found")
	}

	if inputType.Kind != "INPUT_OBJECT" {
		t.Errorf("Expected INPUT_OBJECT kind but got %s", inputType.Kind)
	}

	if len(inputType.Fields) != 3 {
		t.Errorf("Expected 3 input fields but got %d", len(inputType.Fields))
	}

	// Check specific fields
	fieldMap := make(map[string]*Field)
	for _, field := range inputType.Fields {
		fieldMap[field.Name] = field
	}

	// Check name field
	if nameField, ok := fieldMap["name"]; ok {
		if nameField.Type.GetTypeName() != "String" {
			t.Errorf("Expected name field to have String type but got %s", nameField.Type.GetTypeName())
		}
		if !nameField.Type.IsNonNull() {
			t.Errorf("Expected name field to be non-null")
		}
	} else {
		t.Error("Expected name field not found")
	}

	// Check profile field
	if profileField, ok := fieldMap["profile"]; ok {
		if profileField.Type.GetTypeName() != "UserProfileInput" {
			t.Errorf("Expected profile field to have UserProfileInput type but got %s", profileField.Type.GetTypeName())
		}
		if !profileField.Type.IsNonNull() {
			t.Errorf("Expected profile field to be non-null")
		}
	} else {
		t.Error("Expected profile field not found")
	}

	// Check tags field
	if tagsField, ok := fieldMap["tags"]; ok {
		if tagsField.Type.GetTypeName() != "String" {
			t.Errorf("Expected tags field to have String type but got %s", tagsField.Type.GetTypeName())
		}
		if !tagsField.Type.IsList() {
			t.Errorf("Expected tags field to be a list")
		}
	} else {
		t.Error("Expected tags field not found")
	}

	// Debug output
	t.Logf("Parsed schema has %d types", len(result.Types))
	for _, typ := range result.Types {
		t.Logf("  Type: %s (kind: %s, fields: %d)", typ.Name, typ.Kind, len(typ.Fields))
		for _, field := range typ.Fields {
			t.Logf("    Field: %s (type: %s, nonNull: %v, isList: %v)",
				field.Name, field.Type.GetTypeName(), field.Type.IsNonNull(), field.Type.IsList())
		}
	}
}

func TestParseIntrospectionResponse_RealData(t *testing.T) {
	// Load the real introspection response from file
	file, err := os.Open("testdata/real_introspection_response.json")
	if err != nil {
		t.Skipf("Skipping real data test - file not found: %v", err)
		return
	}
	defer file.Close()

	var response map[string]interface{}
	if err := json.NewDecoder(file).Decode(&response); err != nil {
		t.Fatalf("Failed to decode real introspection response: %v", err)
	}

	// Extract the schema data from the GraphQL response wrapper
	schemaData, ok := response["data"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected 'data' field in response")
	}

	// Parse the introspection response
	result, err := ParseIntrospectionResponse(schemaData)
	if err != nil {
		t.Fatalf("ParseIntrospectionResponse() unexpected error: %v", err)
	}

	if result == nil {
		t.Fatal("ParseIntrospectionResponse() returned nil result")
	}

	// Verify we have a reasonable number of types
	if len(result.Types) < 10 {
		t.Errorf("Expected at least 10 types but got %d", len(result.Types))
	}

	// Check for specific input types that should be present
	expectedInputTypes := []string{
		"CreateEquipmentInput",
		"UpdateEquipmentInput",
		"CreateFacilityInput",
		"ContactInfoInput",
		"DimensionsInput",
		"ElectricalSpecsInput",
		"EquipmentSpecificationsInput",
		"LocationInput",
		"TemperatureRangeInput",
		"ScheduleMaintenanceInput",
		"CompleteMaintenanceInput",
		"RecordOperationalMetricInput",
		"UpdateMaintenanceRecordInput",
		"UpdateFacilityInput",
	}

	typeMap := make(map[string]*Type)
	for _, typ := range result.Types {
		typeMap[typ.Name] = typ
	}

	for _, expectedType := range expectedInputTypes {
		if typ, exists := typeMap[expectedType]; !exists {
			t.Errorf("Expected input type %s not found", expectedType)
		} else if typ.Kind != "INPUT_OBJECT" {
			t.Errorf("Expected %s to be INPUT_OBJECT but got %s", expectedType, typ.Kind)
		} else if len(typ.Fields) == 0 {
			t.Errorf("Expected %s to have input fields but got none", expectedType)
		}
	}

	// Test a specific input type with complex nested structure
	createEquipmentInput, exists := typeMap["CreateEquipmentInput"]
	if !exists {
		t.Fatal("CreateEquipmentInput not found")
	}

	// Check that it has the expected fields
	fieldMap := make(map[string]*Field)
	for _, field := range createEquipmentInput.Fields {
		fieldMap[field.Name] = field
	}

	expectedFields := []string{
		"name", "description", "manufacturer", "model", "serialNumber",
		"type", "facilityId", "specifications", "installedAt",
	}

	for _, expectedField := range expectedFields {
		if field, exists := fieldMap[expectedField]; !exists {
			t.Errorf("Expected field %s not found in CreateEquipmentInput", expectedField)
		} else {
			// Check that required fields are non-null
			if expectedField == "name" && !field.Type.IsNonNull() {
				t.Errorf("Expected name field to be non-null")
			}
			if expectedField == "specifications" && field.Type.GetTypeName() != "EquipmentSpecificationsInput" {
				t.Errorf("Expected specifications field to have type EquipmentSpecificationsInput but got %s", field.Type.GetTypeName())
			}
		}
	}

	// Test nested input type
	equipmentSpecsInput, exists := typeMap["EquipmentSpecificationsInput"]
	if !exists {
		t.Fatal("EquipmentSpecificationsInput not found")
	}

	// Check that it has nested input types
	specsFieldMap := make(map[string]*Field)
	for _, field := range equipmentSpecsInput.Fields {
		specsFieldMap[field.Name] = field
	}

	if operatingTempField, exists := specsFieldMap["operatingTemperature"]; exists {
		if operatingTempField.Type.GetTypeName() != "TemperatureRangeInput" {
			t.Errorf("Expected operatingTemperature field to have type TemperatureRangeInput but got %s", operatingTempField.Type.GetTypeName())
		}
	}

	// Test list types
	if envReqsField, exists := specsFieldMap["environmentalRequirements"]; exists {
		if !envReqsField.Type.IsList() {
			t.Errorf("Expected environmentalRequirements field to be a list")
		}
		if envReqsField.Type.GetTypeName() != "String" {
			t.Errorf("Expected environmentalRequirements field to have String type but got %s", envReqsField.Type.GetTypeName())
		}
	}

	// Test Query type
	if result.QueryType == nil {
		t.Fatal("Expected QueryType to be set")
	}
	if result.QueryType.Name != "Query" {
		t.Errorf("Expected QueryType name to be 'Query' but got '%s'", result.QueryType.Name)
	}
	if result.QueryType.Kind != "OBJECT" {
		t.Errorf("Expected QueryType kind to be 'OBJECT' but got '%s'", result.QueryType.Kind)
	}
	if len(result.QueryType.Fields) == 0 {
		t.Error("Expected QueryType to have fields")
	}

	// Check specific query fields
	queryFieldMap := make(map[string]*Field)
	for _, field := range result.QueryType.Fields {
		queryFieldMap[field.Name] = field
	}

	expectedQueryFields := []string{
		"equipment", "equipmentById", "facilities", "facilityById",
		"maintenanceRecords", "operationalMetrics", "facilityStatus",
	}

	for _, expectedField := range expectedQueryFields {
		if field, exists := queryFieldMap[expectedField]; !exists {
			t.Errorf("Expected query field '%s' not found", expectedField)
		} else {
			// Check that equipment field returns a list
			if expectedField == "equipment" && !field.Type.IsList() {
				t.Errorf("Expected 'equipment' field to be a list")
			}
			// Check that equipmentById has an id argument
			if expectedField == "equipmentById" && len(field.Args) == 0 {
				t.Errorf("Expected 'equipmentById' field to have arguments")
			}
		}
	}

	// Test Mutation type
	if result.MutationType == nil {
		t.Fatal("Expected MutationType to be set")
	}
	if result.MutationType.Name != "Mutation" {
		t.Errorf("Expected MutationType name to be 'Mutation' but got '%s'", result.MutationType.Name)
	}
	if result.MutationType.Kind != "OBJECT" {
		t.Errorf("Expected MutationType kind to be 'OBJECT' but got '%s'", result.MutationType.Kind)
	}
	if len(result.MutationType.Fields) == 0 {
		t.Error("Expected MutationType to have fields")
	}

	// Check specific mutation fields
	mutationFieldMap := make(map[string]*Field)
	for _, field := range result.MutationType.Fields {
		mutationFieldMap[field.Name] = field
	}

	expectedMutationFields := []string{
		"createEquipment", "updateEquipment", "deleteEquipment",
		"createFacility", "updateFacility", "deleteFacility",
		"scheduleMaintenance", "recordOperationalMetric",
	}

	for _, expectedField := range expectedMutationFields {
		if field, exists := mutationFieldMap[expectedField]; !exists {
			t.Errorf("Expected mutation field '%s' not found", expectedField)
		} else {
			// Check that createEquipment has an input argument
			if expectedField == "createEquipment" {
				if len(field.Args) == 0 {
					t.Errorf("Expected 'createEquipment' field to have arguments")
				} else {
					// Check that the input argument is of the correct type
					inputArg := field.Args[0]
					if inputArg.Name != "input" {
						t.Errorf("Expected first argument to be named 'input' but got '%s'", inputArg.Name)
					}
					if inputArg.Type.GetTypeName() != "CreateEquipmentInput" {
						t.Errorf("Expected input argument type to be 'CreateEquipmentInput' but got '%s'", inputArg.Type.GetTypeName())
					}
					if !inputArg.Type.IsNonNull() {
						t.Errorf("Expected input argument to be non-null")
					}
				}
			}
		}
	}

	// Debug output for verification
	t.Logf("Successfully parsed real introspection response:")
	t.Logf("  Total types: %d", len(result.Types))
	t.Logf("  Input types found: %d", countInputTypes(result.Types))
	t.Logf("  Query type: %s (%d fields)", result.QueryType.Name, len(result.QueryType.Fields))
	t.Logf("  Mutation type: %s (%d fields)", result.MutationType.Name, len(result.MutationType.Fields))
}

// Helper function to count input types
func countInputTypes(types []*Type) int {
	count := 0
	for _, typ := range types {
		if typ.Kind == "INPUT_OBJECT" {
			count++
		}
	}
	return count
}
