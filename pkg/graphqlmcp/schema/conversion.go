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
		Name:        astDef.Name,
		Kind:        convertKindFromAST(astDef.Kind),
		Description: astDef.Description,
	}

	// Convert fields
	typ.Fields = make([]*Field, 0, len(astDef.Fields))
	for _, astField := range astDef.Fields {
		field := &Field{
			Name:        astField.Name,
			Description: astField.Description,
			Type:        convertTypeFromAST(astField.Type),
			ASTType:     astField.Type, // Store the AST type for dynamic query generation
		}

		// Convert arguments
		field.Args = make([]*Argument, 0, len(astField.Arguments))
		for _, astArg := range astField.Arguments {
			arg := &Argument{
				Name:        astArg.Name,
				Description: astArg.Description,
				Type:        convertTypeFromAST(astArg.Type),
			}
			field.Args = append(field.Args, arg)
		}

		typ.Fields = append(typ.Fields, field)
	}

	return typ
}

// convertTypeFromAST converts gqlparser AST Type to legacy TypeRef
func convertTypeFromAST(astType *ast.Type) *TypeRef {
	if astType == nil {
		return nil
	}

	typ := &TypeRef{}

	if astType.NamedType != "" {
		typ.Name = astType.NamedType
		// Don't set a default kind here - let the actual type lookup determine it
	}

	if astType.NonNull {
		typ.Kind = "NON_NULL"
		typ.OfType = convertTypeFromAST(astType.Elem)
		// For NON_NULL types, don't set the name - it should only be on the innermost type
	} else if astType.Elem != nil {
		typ.Kind = "LIST"
		typ.OfType = convertTypeFromAST(astType.Elem)
		// For LIST types, don't set the name - it should only be on the innermost type
	}

	return typ
}

// convertKindFromAST converts gqlparser AST DefinitionKind to string
func convertKindFromAST(kind ast.DefinitionKind) string {
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

// isBuiltinType checks if a type name is a built-in GraphQL type
func isBuiltinType(name string) bool {
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

// isIntrospectionType checks if a type name is an introspection type
func isIntrospectionType(name string) bool {
	return strings.HasPrefix(name, "__")
}

// getString safely extracts a string value from a map
func getString(data map[string]interface{}, key string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return ""
}
