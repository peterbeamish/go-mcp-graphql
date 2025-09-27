package schema

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// generateTypeSDL generates SDL for a specific type definition
func (s *Schema) generateTypeSDL(typeDef *ast.Definition) string {
	var sdl strings.Builder

	// Add description if present
	if typeDef.Description != "" {
		sdl.WriteString(fmt.Sprintf("\"%s\"\n", typeDef.Description))
	}

	// Generate type definition based on kind
	switch typeDef.Kind {
	case ast.Object:
		sdl.WriteString(fmt.Sprintf("type %s", typeDef.Name))
		if len(typeDef.Interfaces) > 0 {
			sdl.WriteString(fmt.Sprintf(" implements %s", strings.Join(typeDef.Interfaces, " & ")))
		}
		sdl.WriteString(" {\n")
		for _, field := range typeDef.Fields {
			sdl.WriteString(s.generateFieldSDL(field))
		}
		sdl.WriteString("}")

	case ast.Interface:
		sdl.WriteString(fmt.Sprintf("interface %s", typeDef.Name))
		if len(typeDef.Interfaces) > 0 {
			sdl.WriteString(fmt.Sprintf(" implements %s", strings.Join(typeDef.Interfaces, " & ")))
		}
		sdl.WriteString(" {\n")
		for _, field := range typeDef.Fields {
			sdl.WriteString(s.generateFieldSDL(field))
		}
		sdl.WriteString("}")

	case ast.Union:
		sdl.WriteString(fmt.Sprintf("union %s = %s", typeDef.Name, strings.Join(typeDef.Types, " | ")))

	case ast.Enum:
		sdl.WriteString(fmt.Sprintf("enum %s {\n", typeDef.Name))
		for _, enumValue := range typeDef.EnumValues {
			if enumValue.Description != "" {
				sdl.WriteString(fmt.Sprintf("  \"%s\"\n", enumValue.Description))
			}
			sdl.WriteString(fmt.Sprintf("  %s\n", enumValue.Name))
		}
		sdl.WriteString("}")

	case ast.InputObject:
		sdl.WriteString(fmt.Sprintf("input %s {\n", typeDef.Name))
		for _, field := range typeDef.Fields {
			sdl.WriteString(s.generateFieldSDL(field))
		}
		sdl.WriteString("}")

	case ast.Scalar:
		sdl.WriteString(fmt.Sprintf("scalar %s", typeDef.Name))
	}

	return sdl.String()
}

// generateFieldSDL generates SDL for a field definition
func (s *Schema) generateFieldSDL(field *ast.FieldDefinition) string {
	var sdl strings.Builder

	// Add description if present
	if field.Description != "" {
		sdl.WriteString(fmt.Sprintf("  \"%s\"\n", field.Description))
	}

	// Add field name
	sdl.WriteString(fmt.Sprintf("  %s", field.Name))

	// Add arguments if present
	if len(field.Arguments) > 0 {
		sdl.WriteString("(")
		for i, arg := range field.Arguments {
			if i > 0 {
				sdl.WriteString(", ")
			}
			sdl.WriteString(s.generateArgumentSDL(arg))
		}
		sdl.WriteString(")")
	}

	// Add return type
	sdl.WriteString(fmt.Sprintf(": %s\n", s.generateTypeRefSDL(field.Type)))

	return sdl.String()
}

// generateArgumentSDL generates SDL for an argument definition
func (s *Schema) generateArgumentSDL(arg *ast.ArgumentDefinition) string {
	var sdl strings.Builder

	// Add description if present
	if arg.Description != "" {
		sdl.WriteString(fmt.Sprintf("\"%s\" ", arg.Description))
	}

	// Add argument name and type
	sdl.WriteString(fmt.Sprintf("%s: %s", arg.Name, s.generateTypeRefSDL(arg.Type)))

	// Add default value if present
	if arg.DefaultValue != nil {
		sdl.WriteString(fmt.Sprintf(" = %s", arg.DefaultValue.Raw))
	}

	return sdl.String()
}

// generateTypeRefSDL generates SDL for a type reference
func (s *Schema) generateTypeRefSDL(astType *ast.Type) string {
	if astType == nil {
		return "String"
	}

	if astType.NonNull {
		return s.generateTypeRefSDL(astType.Elem) + "!"
	}

	if astType.Elem != nil {
		return "[" + s.generateTypeRefSDL(astType.Elem) + "]"
	}

	return astType.NamedType
}
