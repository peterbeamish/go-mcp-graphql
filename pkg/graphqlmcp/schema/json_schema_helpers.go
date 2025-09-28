package schema

import (
	"fmt"
	"strconv"

	"github.com/vektah/gqlparser/v2/ast"
)

// CreateInputObjectSchema creates a detailed JSON schema for an input object type
func (s *Schema) CreateInputObjectSchema(typeName string) map[string]interface{} {
	return s.createInputObjectSchemaWithDepth(typeName, make(map[string]bool), 0)
}

// createInputObjectSchemaWithDepth creates a JSON schema with recursion depth tracking
func (s *Schema) createInputObjectSchemaWithDepth(typeName string, visited map[string]bool, depth int) map[string]interface{} {
	// Check recursion depth to prevent stack overflow
	const maxDepth = 10
	if depth > maxDepth {
		return map[string]interface{}{
			"type":        "object",
			"description": fmt.Sprintf("Recursive type %s (max depth %d exceeded)", typeName, maxDepth),
			"properties":  map[string]interface{}{},
		}
	}

	// Check for circular references
	if visited[typeName] {
		return map[string]interface{}{
			"type":        "object",
			"description": fmt.Sprintf("Circular reference to %s", typeName),
			"properties":  map[string]interface{}{},
		}
	}

	// Get the type definition from the schema
	typeDef := s.GetTypeDefinition(typeName)
	if typeDef == nil {
		return nil
	}

	// Only handle input object types
	if typeDef.Kind != "INPUT_OBJECT" {
		return nil
	}

	// Mark this type as visited
	visited[typeName] = true
	defer delete(visited, typeName) // Clean up when done

	properties := make(map[string]interface{})
	required := []string{}

	// Process each field in the input object
	for _, field := range typeDef.Fields {
		fieldSchema := s.createInputFieldSchemaFromASTWithDepth(field, visited, depth+1)
		properties[field.Name] = fieldSchema

		// Add to required if it's non-null and has no default value
		if IsASTTypeNonNull(field.Type) && (field.DefaultValue == nil || field.DefaultValue.Raw == "") {
			required = append(required, field.Name)
		}
	}

	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	// Add description if available
	if typeDef.Description != "" {
		schema["description"] = typeDef.Description
	}

	return schema
}

// CreateInputFieldSchemaFromAST creates a JSON schema for an input field from AST
func (s *Schema) CreateInputFieldSchemaFromAST(field *ast.FieldDefinition) map[string]interface{} {
	return s.createBaseSchemaFromAST(field.Type, field.Description, field.DefaultValue)
}

// createInputFieldSchemaFromASTWithDepth creates a JSON schema for an input field with depth tracking
func (s *Schema) createInputFieldSchemaFromASTWithDepth(field *ast.FieldDefinition, visited map[string]bool, depth int) map[string]interface{} {
	return s.createBaseSchemaFromASTWithDepth(field.Type, field.Description, field.DefaultValue, visited, depth)
}

// createBaseSchemaFromAST creates a base JSON schema from AST type with common logic
func (s *Schema) createBaseSchemaFromAST(astType *ast.Type, description string, defaultValue *ast.Value) map[string]interface{} {
	return s.createBaseSchemaFromASTWithDepth(astType, description, defaultValue, make(map[string]bool), 0)
}

// createBaseSchemaFromASTWithDepth creates a base JSON schema with depth tracking
func (s *Schema) createBaseSchemaFromASTWithDepth(astType *ast.Type, description string, defaultValue *ast.Value, visited map[string]bool, depth int) map[string]interface{} {
	schema := map[string]interface{}{
		"type": ASTTypeToJSONSchemaTypeWithSchema(astType, s),
	}

	// Add description if available
	if description != "" {
		schema["description"] = description
	}

	// Handle list types
	if IsASTTypeList(astType) {
		schema["type"] = "array"
		// For list items, use a different function that doesn't add array wrappers
		itemSchema := s.createItemSchemaFromASTWithDepth(astType.Elem, defaultValue, visited, depth+1)
		schema["items"] = itemSchema
		return schema
	}

	// Handle enum types
	typeName := GetASTTypeName(astType)
	if typeName != "" && !isBuiltinType(typeName) {
		// Check if it's an enum type
		if typeDef := s.GetTypeDefinition(typeName); typeDef != nil {
			if typeDef.Kind == ast.Enum {
				// Add enum values to the schema
				enumValues := GetEnumValuesFromAST(astType, s)
				if len(enumValues) > 0 {
					schema["enum"] = enumValues
				}
			}
		}

		// Check for input object types
		if inputObjectSchema := s.createInputObjectSchemaWithDepth(typeName, visited, depth+1); inputObjectSchema != nil {
			// Handle nested input object types
			// Merge the input object schema with the current schema
			for key, value := range inputObjectSchema {
				schema[key] = value
			}
		}
	}

	// Add default value if available
	if defaultValue != nil {
		schema["default"] = convertDefaultValue(defaultValue, astType)
	}

	return schema
}

// createItemSchemaFromAST creates a JSON schema for list items without adding array wrappers
func (s *Schema) createItemSchemaFromAST(astType *ast.Type, defaultValue *ast.Value) map[string]interface{} {
	return s.createItemSchemaFromASTWithDepth(astType, defaultValue, make(map[string]bool), 0)
}

// createItemSchemaFromASTWithDepth creates a JSON schema for list items with depth tracking
func (s *Schema) createItemSchemaFromASTWithDepth(astType *ast.Type, defaultValue *ast.Value, visited map[string]bool, depth int) map[string]interface{} {
	if astType == nil {
		return map[string]interface{}{"type": "string"}
	}

	// Handle the innermost type (scalar or named type)
	if astType.NamedType != "" {
		schema := map[string]interface{}{
			"type": ASTTypeToJSONSchemaTypeWithSchema(astType, s),
		}

		// Handle enum types for list items
		typeName := GetASTTypeName(astType)
		if typeName != "" && !isBuiltinType(typeName) {
			if typeDef := s.GetTypeDefinition(typeName); typeDef != nil && typeDef.Kind == "ENUM" {
				// Add enum values to the schema
				enumValues := GetEnumValuesFromAST(astType, s)
				if len(enumValues) > 0 {
					schema["enum"] = enumValues
				}
			}
		}

		return schema
	}

	// Handle wrapper types (NON_NULL)
	if astType.NonNull {
		return s.createItemSchemaFromASTWithDepth(astType.Elem, defaultValue, visited, depth+1)
	}

	// Handle LIST types - this shouldn't happen in normal cases, but handle it
	if astType.Elem != nil {
		return s.createItemSchemaFromASTWithDepth(astType.Elem, defaultValue, visited, depth+1)
	}

	return map[string]interface{}{"type": "string"}
}

// CreateTypeRefSchema creates a JSON schema for a TypeRef
func (s *Schema) CreateTypeRefSchema(typeRef *TypeRef, description string) map[string]interface{} {
	return s.createBaseSchemaFromTypeRef(typeRef, description, "")
}

// createBaseSchemaFromTypeRef creates a base JSON schema from TypeRef with common logic
func (s *Schema) createBaseSchemaFromTypeRef(typeRef *TypeRef, description, defaultValue string) map[string]interface{} {
	schema := map[string]interface{}{
		"type": typeRef.ToJSONSchemaTypeWithSchema(s),
	}

	// Add description if available
	if description != "" {
		schema["description"] = description
	}

	// Handle list types
	if typeRef.IsList() {
		schema["type"] = "array"
		itemSchema := s.createBaseSchemaFromTypeRef(typeRef.OfType, description, defaultValue)
		schema["items"] = itemSchema
		return schema
	}

	// Handle enum types and input object types
	if typeRef.GetTypeName() != "" && !isBuiltinType(typeRef.GetTypeName()) {
		// Check if it's an enum type
		if typeDef := s.GetTypeDefinition(typeRef.GetTypeName()); typeDef != nil {
			if typeDef.Kind == ast.Enum {
				// Add enum values to the schema
				enumValues := GetEnumValuesFromTypeRef(typeRef, s)
				if len(enumValues) > 0 {
					schema["enum"] = enumValues
				}
			}
		}

		// Check for input object types
		if inputObjectSchema := s.CreateInputObjectSchema(typeRef.GetTypeName()); inputObjectSchema != nil {
			// Handle input object types - but don't flatten them when creating argument schemas
			// For input object types, we want to reference the schema, not flatten it
			// This preserves the proper nesting structure for arguments
			schema["properties"] = inputObjectSchema["properties"]
			if required, ok := inputObjectSchema["required"].([]string); ok && len(required) > 0 {
				schema["required"] = required
			}
		}
	}

	// Add default value if available
	if defaultValue != "" {
		schema["default"] = convertDefaultValueFromString(defaultValue, typeRef)
	}

	return schema
}

// CreateArgumentSchema creates a JSON schema for a GraphQL argument
func (s *Schema) CreateArgumentSchema(arg *Argument) map[string]interface{} {
	return s.createBaseSchemaFromTypeRef(arg.Type, arg.Description, arg.DefaultValue)
}

// CreateInputSchema creates a JSON schema for the tool input
func (s *Schema) CreateInputSchema(field *Field) map[string]interface{} {
	properties := make(map[string]interface{})
	required := []string{}

	// Add arguments as properties
	for _, arg := range field.Args {
		argSchema := s.CreateArgumentSchema(arg)
		properties[arg.Name] = argSchema

		// Add to required if it's non-null and has no default value
		if arg.Type.IsNonNull() && arg.DefaultValue == "" {
			required = append(required, arg.Name)
		}
	}

	schema := map[string]interface{}{
		"type":       "object",
		"properties": properties,
	}

	if len(required) > 0 {
		schema["required"] = required
	}

	return schema
}

// convertDefaultValue converts a GraphQL default value to the appropriate JSON type
func convertDefaultValue(defaultValue *ast.Value, astType *ast.Type) interface{} {
	if defaultValue == nil {
		return nil
	}

	// Get the base type name
	baseType := GetASTTypeName(astType)

	// Convert based on the GraphQL type
	switch baseType {
	case "Boolean":
		if defaultValue.Raw == "true" {
			return true
		} else if defaultValue.Raw == "false" {
			return false
		}
		// Fallback to string if not a valid boolean
		return defaultValue.Raw
	case "Int":
		if intVal, err := strconv.Atoi(defaultValue.Raw); err == nil {
			return intVal
		}
		// Fallback to string if not a valid int
		return defaultValue.Raw
	case "Float":
		if floatVal, err := strconv.ParseFloat(defaultValue.Raw, 64); err == nil {
			return floatVal
		}
		// Fallback to string if not a valid float
		return defaultValue.Raw
	case "String", "ID":
		// Remove quotes if present
		raw := defaultValue.Raw
		if len(raw) >= 2 && raw[0] == '"' && raw[len(raw)-1] == '"' {
			return raw[1 : len(raw)-1]
		}
		return raw
	default:
		// For other types (enums, objects, etc.), return as string
		return defaultValue.Raw
	}
}

// convertDefaultValueFromString converts a string default value to the appropriate JSON type based on TypeRef
func convertDefaultValueFromString(defaultValue string, typeRef *TypeRef) interface{} {
	if defaultValue == "" {
		return nil
	}

	// Get the base type name
	baseType := typeRef.GetTypeName()

	// Convert based on the GraphQL type
	switch baseType {
	case "Boolean":
		if defaultValue == "true" {
			return true
		} else if defaultValue == "false" {
			return false
		}
		// Fallback to string if not a valid boolean
		return defaultValue
	case "Int":
		if intVal, err := strconv.Atoi(defaultValue); err == nil {
			return intVal
		}
		// Fallback to string if not a valid int
		return defaultValue
	case "Float":
		if floatVal, err := strconv.ParseFloat(defaultValue, 64); err == nil {
			return floatVal
		}
		// Fallback to string if not a valid float
		return defaultValue
	case "String", "ID":
		// Remove quotes if present
		if len(defaultValue) >= 2 && defaultValue[0] == '"' && defaultValue[len(defaultValue)-1] == '"' {
			return defaultValue[1 : len(defaultValue)-1]
		}
		return defaultValue
	default:
		// For other types (enums, objects, etc.), return as string
		return defaultValue
	}
}
