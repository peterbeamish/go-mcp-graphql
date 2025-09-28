package schema

import (
	"github.com/vektah/gqlparser/v2/ast"
)

// CreateInputObjectSchema creates a detailed JSON schema for an input object type
func (s *Schema) CreateInputObjectSchema(typeName string) map[string]interface{} {
	// Get the type definition from the schema
	typeDef := s.GetTypeDefinition(typeName)
	if typeDef == nil {
		return nil
	}

	// Only handle input object types
	if typeDef.Kind != "INPUT_OBJECT" {
		return nil
	}

	properties := make(map[string]interface{})
	required := []string{}

	// Process each field in the input object
	for _, field := range typeDef.Fields {
		fieldSchema := s.CreateInputFieldSchemaFromAST(field)
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
	schema := map[string]interface{}{
		"type": ASTTypeToJSONSchemaType(field.Type),
	}

	// Add description if available
	if field.Description != "" {
		schema["description"] = field.Description
	}

	// Handle list types
	if IsASTTypeList(field.Type) {
		schema["type"] = "array"
		itemSchema := s.CreateInputFieldSchemaFromAST(&ast.FieldDefinition{
			Name:        field.Name,
			Description: field.Description,
			Type:        field.Type.Elem,
		})
		schema["items"] = itemSchema
		return schema
	}

	// Handle nested input object types
	typeName := GetASTTypeName(field.Type)
	if typeName != "" && !isBuiltinType(typeName) {
		if inputObjectSchema := s.CreateInputObjectSchema(typeName); inputObjectSchema != nil {
			// Merge the input object schema with the current schema
			for key, value := range inputObjectSchema {
				schema[key] = value
			}
		}
	}

	// Add default value if available
	if field.DefaultValue != nil {
		schema["default"] = field.DefaultValue.Raw
	}

	return schema
}

// CreateTypeRefSchema creates a JSON schema for a TypeRef
func (s *Schema) CreateTypeRefSchema(typeRef *TypeRef, description string) map[string]interface{} {
	schema := map[string]interface{}{
		"type": typeRef.ToJSONSchemaType(),
	}

	// Add description if available
	if description != "" {
		schema["description"] = description
	}

	// Handle list types
	if typeRef.IsList() {
		schema["type"] = "array"
		itemSchema := s.CreateTypeRefSchema(typeRef.OfType, description)
		schema["items"] = itemSchema
		return schema
	}

	// Handle input object types
	if typeRef.GetTypeName() != "" && !isBuiltinType(typeRef.GetTypeName()) {
		if inputObjectSchema := s.CreateInputObjectSchema(typeRef.GetTypeName()); inputObjectSchema != nil {
			// Merge the input object schema with the current schema
			for key, value := range inputObjectSchema {
				schema[key] = value
			}
		}
	}

	return schema
}

// CreateArgumentSchema creates a JSON schema for a GraphQL argument
func (s *Schema) CreateArgumentSchema(arg *Argument) map[string]interface{} {
	schema := map[string]interface{}{
		"type": arg.Type.ToJSONSchemaType(),
	}

	// Add description if available
	if arg.Description != "" {
		schema["description"] = arg.Description
	}

	// Handle list types
	if arg.Type.IsList() {
		schema["type"] = "array"
		// Create schema for list items
		itemSchema := s.CreateTypeRefSchema(arg.Type.OfType, arg.Description)
		schema["items"] = itemSchema
		return schema
	}

	// Handle input object types - resolve the actual input object definition
	if arg.Type.GetTypeName() != "" && !isBuiltinType(arg.Type.GetTypeName()) {
		if inputObjectSchema := s.CreateInputObjectSchema(arg.Type.GetTypeName()); inputObjectSchema != nil {
			// Merge the input object schema with the current schema
			for key, value := range inputObjectSchema {
				schema[key] = value
			}
		}
	}

	// Add default value if available
	if arg.DefaultValue != "" {
		schema["default"] = arg.DefaultValue
	}

	return schema
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
