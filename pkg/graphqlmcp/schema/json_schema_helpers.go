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
	return s.createBaseSchemaFromAST(field.Type, field.Description, field.DefaultValue)
}

// createBaseSchemaFromAST creates a base JSON schema from AST type with common logic
func (s *Schema) createBaseSchemaFromAST(astType *ast.Type, description string, defaultValue *ast.Value) map[string]interface{} {
	schema := map[string]interface{}{
		"type": ASTTypeToJSONSchemaType(astType),
	}

	// Add description if available
	if description != "" {
		schema["description"] = description
	}

	// Handle list types
	if IsASTTypeList(astType) {
		schema["type"] = "array"
		itemSchema := s.createBaseSchemaFromAST(astType.Elem, description, defaultValue)
		schema["items"] = itemSchema
		return schema
	}

	// Handle nested input object types
	typeName := GetASTTypeName(astType)
	if typeName != "" && !isBuiltinType(typeName) {
		if inputObjectSchema := s.CreateInputObjectSchema(typeName); inputObjectSchema != nil {
			// Merge the input object schema with the current schema
			for key, value := range inputObjectSchema {
				schema[key] = value
			}
		}
	}

	// Add default value if available
	if defaultValue != nil {
		schema["default"] = defaultValue.Raw
	}

	return schema
}

// CreateTypeRefSchema creates a JSON schema for a TypeRef
func (s *Schema) CreateTypeRefSchema(typeRef *TypeRef, description string) map[string]interface{} {
	return s.createBaseSchemaFromTypeRef(typeRef, description, "")
}

// createBaseSchemaFromTypeRef creates a base JSON schema from TypeRef with common logic
func (s *Schema) createBaseSchemaFromTypeRef(typeRef *TypeRef, description, defaultValue string) map[string]interface{} {
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
		itemSchema := s.createBaseSchemaFromTypeRef(typeRef.OfType, description, defaultValue)
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

	// Add default value if available
	if defaultValue != "" {
		schema["default"] = defaultValue
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
