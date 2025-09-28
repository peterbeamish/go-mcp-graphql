package schema

import (
	"github.com/vektah/gqlparser/v2/ast"
)

// IsBuiltinType checks if a type name is a built-in GraphQL type
func IsBuiltinType(typeName string) bool {
	builtinTypes := map[string]bool{
		"String":       true,
		"Int":          true,
		"Float":        true,
		"Boolean":      true,
		"ID":           true,
		"__Schema":     true,
		"__Type":       true,
		"__Field":      true,
		"__InputValue": true,
		"__EnumValue":  true,
		"__Directive":  true,
	}
	return builtinTypes[typeName]
}

// ASTTypeToJSONSchemaType converts AST type to JSON Schema type
func ASTTypeToJSONSchemaType(astType *ast.Type) string {
	if astType == nil {
		return "string"
	}

	// Get the base type name
	baseType := GetASTTypeName(astType)

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
		return "object"
	}
}

// GetASTTypeName extracts the type name from an AST type
func GetASTTypeName(astType *ast.Type) string {
	if astType == nil {
		return ""
	}

	currentType := astType
	for currentType != nil {
		if currentType.NamedType != "" {
			return currentType.NamedType
		}
		currentType = currentType.Elem
	}
	return ""
}

// IsASTTypeList checks if an AST type is a list
func IsASTTypeList(astType *ast.Type) bool {
	if astType == nil {
		return false
	}

	// Check if this is a list type
	if astType.Elem != nil && !astType.NonNull {
		return true
	}

	// Check if wrapped in non-null
	if astType.NonNull && astType.Elem != nil {
		return IsASTTypeList(astType.Elem)
	}

	return false
}

// IsASTTypeNonNull checks if an AST type is non-null
func IsASTTypeNonNull(astType *ast.Type) bool {
	if astType == nil {
		return false
	}

	// Check if this is a non-null type
	if astType.NonNull {
		return true
	}

	// For other types (including LIST), they are nullable unless wrapped in NON_NULL
	return false
}
