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

// GetInterfaces returns all interface types in the schema
func (s *Schema) GetInterfaces() []*Type {
	if s.parsedSchema == nil {
		return nil
	}

	var interfaces []*Type
	for _, typeDef := range s.parsedSchema.Types {
		if typeDef != nil && typeDef.Kind == ast.Interface && !isBuiltinType(typeDef.Name) && !isIntrospectionType(typeDef.Name) {
			interfaces = append(interfaces, convertASTToType(typeDef))
		}
	}
	return interfaces
}

// GetImplementations returns all types that implement the given interface
func (s *Schema) GetImplementations(interfaceName string) []*Type {
	if s.parsedSchema == nil {
		return nil
	}

	var implementations []*Type
	for _, typeDef := range s.parsedSchema.Types {
		if typeDef != nil && typeDef.Kind == ast.Object {
			for _, iface := range typeDef.Interfaces {
				if iface == interfaceName {
					implementations = append(implementations, convertASTToType(typeDef))
					break
				}
			}
		}
	}
	return implementations
}

// GetInterfaceFields returns all fields defined by an interface
func (s *Schema) GetInterfaceFields(interfaceName string) []*Field {
	if s.parsedSchema == nil {
		return nil
	}

	typeDef := s.GetTypeDefinition(interfaceName)
	if typeDef == nil || typeDef.Kind != ast.Interface {
		return nil
	}

	// Convert AST fields to Field objects
	fields := make([]*Field, 0, len(typeDef.Fields))
	for _, astField := range typeDef.Fields {
		field := &Field{
			Name:        astField.Name,
			Description: astField.Description,
			Type:        ConvertTypeFromAST(astField.Type),
			ASTType:     astField.Type,
		}
		fields = append(fields, field)
	}

	return fields
}

// GetUnions returns all union types in the schema
func (s *Schema) GetUnions() []*Type {
	if s.parsedSchema == nil {
		return nil
	}

	var unions []*Type
	for _, typeDef := range s.parsedSchema.Types {
		if typeDef != nil && typeDef.Kind == ast.Union && !isBuiltinType(typeDef.Name) && !isIntrospectionType(typeDef.Name) {
			unions = append(unions, convertASTToType(typeDef))
		}
	}
	return unions
}

// GetUnionPossibleTypes returns all possible types for a given union
func (s *Schema) GetUnionPossibleTypes(unionName string) []*Type {
	if s.parsedSchema == nil {
		return nil
	}

	typeDef := s.GetTypeDefinition(unionName)
	if typeDef == nil || typeDef.Kind != ast.Union {
		return nil
	}

	var possibleTypes []*Type
	for _, possibleTypeName := range typeDef.Types {
		if possibleTypeDef := s.GetTypeDefinition(possibleTypeName); possibleTypeDef != nil {
			possibleTypes = append(possibleTypes, convertASTToType(possibleTypeDef))
		}
	}
	return possibleTypes
}

// IsUnionType checks if a type is a union
func (s *Schema) IsUnionType(typeName string) bool {
	if s.parsedSchema == nil {
		return false
	}

	typeDef := s.GetTypeDefinition(typeName)
	return typeDef != nil && typeDef.Kind == ast.Union
}

// GetUnionByName returns a union type by name
func (s *Schema) GetUnionByName(unionName string) *Type {
	if s.parsedSchema == nil {
		return nil
	}

	typeDef := s.GetTypeDefinition(unionName)
	if typeDef == nil || typeDef.Kind != ast.Union {
		return nil
	}

	return convertASTToType(typeDef)
}

// GetMaxDepth returns the maximum depth for query generation
func (s *Schema) GetMaxDepth() int {
	if s.MaxDepth <= 0 {
		return 5 // Default value if not set
	}
	return s.MaxDepth
}
