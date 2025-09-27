package schema

import (
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// GetQueries returns all query fields from the schema
func (s *Schema) GetQueries() []*Field {
	if s.QueryType == nil {
		return nil
	}
	return s.QueryType.Fields
}

// GetMutations returns all mutation fields from the schema
func (s *Schema) GetMutations() []*Field {
	if s.MutationType == nil {
		return nil
	}
	return s.MutationType.Fields
}

// GetTypeDefinition returns the AST definition for a type name
func (s *Schema) GetTypeDefinition(typeName string) *ast.Definition {
	if s.typeRegistry == nil {
		return nil
	}
	return s.typeRegistry[typeName]
}

// GetSchemaSDL returns the GraphQL schema in SDL format
func (s *Schema) GetSchemaSDL() string {
	if s.parsedSchema == nil {
		return ""
	}

	var sdl strings.Builder

	// Add Query type
	if s.parsedSchema.Query != nil {
		sdl.WriteString(s.generateTypeSDL(s.parsedSchema.Query))
		sdl.WriteString("\n\n")
	}

	// Add Mutation type
	if s.parsedSchema.Mutation != nil {
		sdl.WriteString(s.generateTypeSDL(s.parsedSchema.Mutation))
		sdl.WriteString("\n\n")
	}

	// Add all other types (excluding built-in types and introspection types)
	for _, typeDef := range s.parsedSchema.Types {
		if typeDef != nil && !isBuiltinType(typeDef.Name) && !isIntrospectionType(typeDef.Name) {
			sdl.WriteString(s.generateTypeSDL(typeDef))
			sdl.WriteString("\n\n")
		}
	}

	return strings.TrimSpace(sdl.String())
}
