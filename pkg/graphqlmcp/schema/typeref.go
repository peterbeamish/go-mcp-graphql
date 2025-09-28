package schema

import (
	"github.com/vektah/gqlparser/v2/ast"
)

// GetTypeName returns the actual type name, handling non-null and list wrappers
func (tr *TypeRef) GetTypeName() string {
	if tr == nil {
		return "String"
	}

	// Handle non-null and list wrappers
	currentType := tr
	for currentType != nil {
		if currentType.Name != "" {
			return currentType.Name
		}
		currentType = currentType.OfType
	}
	return "String"
}

// IsList checks if the type is a list
func (tr *TypeRef) IsList() bool {
	if tr == nil {
		return false
	}

	// Check if this is a list type
	if tr.Kind == "LIST" {
		return true
	}

	// Check if wrapped in non-null
	if tr.Kind == "NON_NULL" && tr.OfType != nil {
		return tr.OfType.IsList()
	}

	return false
}

// IsNonNull checks if the type is non-null
func (tr *TypeRef) IsNonNull() bool {
	if tr == nil {
		return false
	}

	// Check if this is a non-null type
	if tr.Kind == "NON_NULL" {
		return true
	}

	// For other types (including LIST), they are nullable unless wrapped in NON_NULL
	return false
}

// ToJSONSchemaType converts GraphQL type to JSON Schema type
func (tr *TypeRef) ToJSONSchemaType() string {
	if tr == nil {
		return "string"
	}

	// Get the base type name
	baseType := tr.GetTypeName()

	// Convert GraphQL types to JSON Schema types
	switch baseType {
	case "String", "ID":
		return "string"
	case "Int":
		return "integer"
	case "Float":
		return "number"
	case "Boolean":
		return "boolean"
	default:
		// For non-builtin types, default to object (input objects and other complex types)
		// The actual enum values will be handled by the calling function with schema context
		return "object"
	}
}

// ToJSONSchemaTypeWithSchema converts GraphQL type to JSON Schema type with schema context
func (tr *TypeRef) ToJSONSchemaTypeWithSchema(schema *Schema) string {
	if tr == nil {
		return "string"
	}

	// Get the base type name
	baseType := tr.GetTypeName()

	// Convert GraphQL types to JSON Schema types
	switch baseType {
	case "String", "ID":
		return "string"
	case "Int":
		return "integer"
	case "Float":
		return "number"
	case "Boolean":
		return "boolean"
	default:
		// Check if it's an enum type
		if schema != nil && schema.typeRegistry != nil {
			if typeDef := schema.GetTypeDefinition(baseType); typeDef != nil {
				switch typeDef.Kind {
				case ast.Enum:
					return "string" // Enums are represented as strings in JSON Schema
				case ast.InputObject:
					return "object" // Input objects are represented as objects in JSON Schema
				case ast.Object:
					return "object" // Objects are represented as objects in JSON Schema
				case ast.Interface:
					return "object" // Interfaces are represented as objects in JSON Schema
				case ast.Union:
					return "object" // Unions are represented as objects in JSON Schema
				case ast.Scalar:
					return "string" // Custom scalars are represented as strings in JSON Schema
				default:
					return "object" // Default to object for unknown complex types
				}
			}
		}
		return "object" // Default to object for unknown types
	}
}

// GetEnumValuesFromTypeRef extracts enum values from a TypeRef with schema context
func GetEnumValuesFromTypeRef(typeRef *TypeRef, schema *Schema) []string {
	if typeRef == nil || schema == nil || schema.typeRegistry == nil {
		return nil
	}

	// Get the base type name
	baseType := typeRef.GetTypeName()
	if baseType == "" {
		return nil
	}

	// Get the type definition
	typeDef := schema.GetTypeDefinition(baseType)
	if typeDef == nil || typeDef.Kind != ast.Enum {
		return nil
	}

	// Extract enum values
	enumValues := make([]string, 0, len(typeDef.EnumValues))
	for _, enumValue := range typeDef.EnumValues {
		enumValues = append(enumValues, enumValue.Name)
	}

	return enumValues
}
