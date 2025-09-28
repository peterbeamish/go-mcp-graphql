package schema

import (
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// convertASTToType converts gqlparser AST Definition to legacy Type for backward compatibility
func convertASTToType(astDef *ast.Definition) *Type {
	if astDef == nil {
		return nil
	}

	typ := &Type{
		Name:          astDef.Name,
		Kind:          convertKindFromAST(astDef.Kind),
		Description:   astDef.Description,
		Interfaces:    astDef.Interfaces,
		PossibleTypes: astDef.Types,
	}

	// Convert fields
	typ.Fields = make([]*Field, 0, len(astDef.Fields))
	for _, astField := range astDef.Fields {
		field := &Field{
			Name:        astField.Name,
			Description: astField.Description,
			Type:        ConvertTypeFromAST(astField.Type),
			ASTType:     astField.Type, // Store the AST type for dynamic query generation
		}

		// Convert arguments
		field.Args = make([]*Argument, 0, len(astField.Arguments))
		for _, astArg := range astField.Arguments {
			arg := &Argument{
				Name:        astArg.Name,
				Description: astArg.Description,
				Type:        ConvertTypeFromAST(astArg.Type),
			}

			// Copy default value if present
			if astArg.DefaultValue != nil {
				arg.DefaultValue = astArg.DefaultValue.Raw
			}

			field.Args = append(field.Args, arg)
		}

		typ.Fields = append(typ.Fields, field)
	}

	return typ
}

// ConvertTypeFromAST converts gqlparser AST Type to legacy TypeRef
func ConvertTypeFromAST(astType *ast.Type) *TypeRef {
	if astType == nil {
		return nil
	}

	typ := &TypeRef{}

	// Handle wrapper types first (NON_NULL or LIST)
	if astType.NonNull {
		typ.Kind = "NON_NULL"
		// For NON_NULL types, we need to create the inner type
		if astType.NamedType != "" {
			// NON_NULL with a named type (e.g., NON_NULL String)
			typ.OfType = &TypeRef{
				Name: astType.NamedType,
				Kind: "SCALAR",
			}
		} else {
			// NON_NULL with a complex type (e.g., NON_NULL LIST)
			typ.OfType = ConvertTypeFromAST(astType.Elem)
		}
	} else if astType.Elem != nil {
		// This is a LIST type
		typ.Kind = "LIST"
		typ.OfType = ConvertTypeFromAST(astType.Elem)
	} else if astType.NamedType != "" {
		// Handle the innermost type (scalar or named type)
		typ.Name = astType.NamedType
		typ.Kind = "SCALAR" // Default to SCALAR for named types
	} else {
		// This shouldn't happen, but handle gracefully
		typ.Kind = "SCALAR"
		typ.Name = "String"
	}

	return typ
}

// convertKindFromAST converts gqlparser AST DefinitionKind to string
func convertKindFromAST(kind ast.DefinitionKind) string {
	return convertKindToString(kind)
}

// convertKindToString converts gqlparser AST DefinitionKind to string
func convertKindToString(kind ast.DefinitionKind) string {
	switch kind {
	case ast.Object:
		return "OBJECT"
	case ast.Interface:
		return "INTERFACE"
	case ast.Union:
		return "UNION"
	case ast.Enum:
		return "ENUM"
	case ast.InputObject:
		return "INPUT_OBJECT"
	case ast.Scalar:
		return "SCALAR"
	default:
		return "OBJECT"
	}
}

// convertStringToKind converts GraphQL kind string to gqlparser AST DefinitionKind
func convertStringToKind(kind string) ast.DefinitionKind {
	switch kind {
	case "OBJECT":
		return ast.Object
	case "INTERFACE":
		return ast.Interface
	case "UNION":
		return ast.Union
	case "ENUM":
		return ast.Enum
	case "INPUT_OBJECT":
		return ast.InputObject
	case "SCALAR":
		return ast.Scalar
	default:
		return ast.Object
	}
}

// IsBuiltinType checks if a type name is a built-in GraphQL type
func IsBuiltinType(name string) bool {
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
	return builtinTypes[name]
}

// isBuiltinType checks if a type name is a built-in GraphQL type (internal helper)
func isBuiltinType(name string) bool {
	return IsBuiltinType(name)
}

// isIntrospectionType checks if a type name is an introspection type
func isIntrospectionType(name string) bool {
	return strings.HasPrefix(name, "__")
}

// isScalarType checks if a type is a scalar using built-in types
func isScalarType(astType *ast.Type) bool {
	if astType == nil {
		return false
	}

	// Unwrap non-null and list wrappers
	currentType := astType
	for currentType != nil {
		if currentType.NamedType != "" {
			scalarTypes := map[string]bool{
				"String":  true,
				"Int":     true,
				"Float":   true,
				"Boolean": true,
				"ID":      true,
			}
			return scalarTypes[currentType.NamedType]
		}
		currentType = currentType.Elem
	}
	return false
}

// isScalarTypeWithSchema checks if a type is a scalar using the schema's type registry
func isScalarTypeWithSchema(astType *ast.Type, schema *Schema) bool {
	if astType == nil || schema == nil || schema.typeRegistry == nil {
		return false
	}

	// Unwrap non-null and list wrappers
	currentType := astType
	for currentType != nil {
		if currentType.NamedType != "" {
			// Check if it's a built-in scalar type
			if isScalarType(currentType) {
				return true
			}

			// Check the type registry to see if it's defined as a scalar
			if typeDef := schema.GetTypeDefinition(currentType.NamedType); typeDef != nil {
				return typeDef.Kind == ast.Scalar
			}

			return false
		}
		currentType = currentType.Elem
	}
	return false
}

// getString safely extracts a string value from a map
func getString(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
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
		// For non-builtin types, default to object (input objects and other complex types)
		// The actual enum values will be handled by the calling function with schema context
		return "object"
	}
}

// ASTTypeToJSONSchemaTypeWithSchema converts AST type to JSON Schema type with schema context
func ASTTypeToJSONSchemaTypeWithSchema(astType *ast.Type, schema *Schema) string {
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

// GetEnumValuesFromAST extracts enum values from an AST type with schema context
func GetEnumValuesFromAST(astType *ast.Type, schema *Schema) []string {
	if astType == nil || schema == nil || schema.typeRegistry == nil {
		return nil
	}

	// Get the base type name
	baseType := GetASTTypeName(astType)
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

// GetASTTypeName extracts the type name from an AST type
func GetASTTypeName(astType *ast.Type) string {
	return extractTypeNameFromAST(astType)
}

// extractTypeNameFromAST is the unified implementation for extracting type names from AST types
func extractTypeNameFromAST(astType *ast.Type) string {
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

	// A list type has Elem set and NamedType empty
	// This works for both nullable and non-nullable lists
	// The key is that the current type has no NamedType (it's a wrapper)
	return astType.Elem != nil && astType.NamedType == ""
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
