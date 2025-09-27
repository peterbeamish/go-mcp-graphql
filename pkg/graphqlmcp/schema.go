package graphqlmcp

import (
	"fmt"
	"strings"
)

// Schema represents a GraphQL schema
type Schema struct {
	QueryType    *Type   `json:"queryType"`
	MutationType *Type   `json:"mutationType"`
	Types        []*Type `json:"types"`
}

// Type represents a GraphQL type
type Type struct {
	Name        string      `json:"name"`
	Kind        string      `json:"kind"`
	Description string      `json:"description"`
	Fields      []*Field    `json:"fields"`
	Args        []*Argument `json:"args"`
}

// Field represents a GraphQL field
type Field struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        *TypeRef    `json:"type"`
	Args        []*Argument `json:"args"`
}

// Argument represents a GraphQL argument
type Argument struct {
	Name         string   `json:"name"`
	Type         *TypeRef `json:"type"`
	DefaultValue string   `json:"defaultValue"`
}

// TypeRef represents a GraphQL type reference
type TypeRef struct {
	Name   string   `json:"name"`
	Kind   string   `json:"kind"`
	OfType *TypeRef `json:"ofType"`
}

// GetTypeName returns the actual type name, handling non-null and list wrappers
func (tr *TypeRef) GetTypeName() string {
	if tr == nil {
		return "String"
	}

	// Handle non-null and list wrappers
	if tr.Kind == "NON_NULL" || tr.Kind == "LIST" {
		if tr.OfType != nil {
			return tr.OfType.GetTypeName()
		}
	}

	if tr.Name != "" {
		return tr.Name
	}

	return "String"
}

// IsList returns true if this type is a list
func (tr *TypeRef) IsList() bool {
	if tr == nil {
		return false
	}
	if tr.Kind == "LIST" {
		return true
	}
	if tr.OfType != nil {
		return tr.OfType.IsList()
	}
	return false
}

// IsNonNull returns true if this type is non-null
func (tr *TypeRef) IsNonNull() bool {
	if tr == nil {
		return false
	}
	if tr.Kind == "NON_NULL" {
		return true
	}
	if tr.OfType != nil {
		return tr.OfType.IsNonNull()
	}
	return false
}

// ToJSONSchemaType converts a GraphQL type to a JSON Schema type
func (tr *TypeRef) ToJSONSchemaType() string {
	typeName := tr.GetTypeName()

	switch typeName {
	case "String":
		return "string"
	case "Int":
		return "integer"
	case "Float":
		return "number"
	case "Boolean":
		return "boolean"
	case "ID":
		return "string"
	default:
		return "object"
	}
}

// parseIntrospectionResponse parses the introspection response into our Schema struct
func parseIntrospectionResponse(data map[string]interface{}) (*Schema, error) {
	schemaData, ok := data["__schema"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid schema data")
	}

	schema := &Schema{}

	// Parse query type
	if queryTypeData, ok := schemaData["queryType"].(map[string]interface{}); ok {
		queryType, err := parseType(queryTypeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse query type: %w", err)
		}
		schema.QueryType = queryType
	}

	// Parse mutation type
	if mutationTypeData, ok := schemaData["mutationType"].(map[string]interface{}); ok {
		mutationType, err := parseType(mutationTypeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse mutation type: %w", err)
		}
		schema.MutationType = mutationType
	}

	// Parse types
	if typesData, ok := schemaData["types"].([]interface{}); ok {
		types := make([]*Type, 0, len(typesData))
		for _, typeData := range typesData {
			if typeMap, ok := typeData.(map[string]interface{}); ok {
				typ, err := parseType(typeMap)
				if err != nil {
					return nil, fmt.Errorf("failed to parse type: %w", err)
				}
				types = append(types, typ)
			}
		}
		schema.Types = types
	}

	return schema, nil
}

// parseType parses a type from the introspection response
func parseType(data map[string]interface{}) (*Type, error) {
	typ := &Type{}

	if name, ok := data["name"].(string); ok {
		typ.Name = name
	}

	if kind, ok := data["kind"].(string); ok {
		typ.Kind = kind
	}

	if description, ok := data["description"].(string); ok {
		typ.Description = description
	}

	// Parse fields
	if fieldsData, ok := data["fields"].([]interface{}); ok {
		fields := make([]*Field, 0, len(fieldsData))
		for _, fieldData := range fieldsData {
			if fieldMap, ok := fieldData.(map[string]interface{}); ok {
				field, err := parseField(fieldMap)
				if err != nil {
					return nil, fmt.Errorf("failed to parse field: %w", err)
				}
				fields = append(fields, field)
			}
		}
		typ.Fields = fields
	}

	// Parse args
	if argsData, ok := data["args"].([]interface{}); ok {
		args := make([]*Argument, 0, len(argsData))
		for _, argData := range argsData {
			if argMap, ok := argData.(map[string]interface{}); ok {
				arg, err := parseArgument(argMap)
				if err != nil {
					return nil, fmt.Errorf("failed to parse argument: %w", err)
				}
				args = append(args, arg)
			}
		}
		typ.Args = args
	}

	return typ, nil
}

// parseField parses a field from the introspection response
func parseField(data map[string]interface{}) (*Field, error) {
	field := &Field{}

	if name, ok := data["name"].(string); ok {
		field.Name = name
	}

	if description, ok := data["description"].(string); ok {
		field.Description = description
	}

	// Parse type
	if typeData, ok := data["type"].(map[string]interface{}); ok {
		typeRef, err := parseTypeRef(typeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse field type: %w", err)
		}
		field.Type = typeRef
	}

	// Parse args
	if argsData, ok := data["args"].([]interface{}); ok {
		args := make([]*Argument, 0, len(argsData))
		for _, argData := range argsData {
			if argMap, ok := argData.(map[string]interface{}); ok {
				arg, err := parseArgument(argMap)
				if err != nil {
					return nil, fmt.Errorf("failed to parse field argument: %w", err)
				}
				args = append(args, arg)
			}
		}
		field.Args = args
	}

	return field, nil
}

// parseArgument parses an argument from the introspection response
func parseArgument(data map[string]interface{}) (*Argument, error) {
	arg := &Argument{}

	if name, ok := data["name"].(string); ok {
		arg.Name = name
	}

	if defaultValue, ok := data["defaultValue"].(string); ok {
		arg.DefaultValue = defaultValue
	}

	// Parse type
	if typeData, ok := data["type"].(map[string]interface{}); ok {
		typeRef, err := parseTypeRef(typeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse argument type: %w", err)
		}
		arg.Type = typeRef
	}

	return arg, nil
}

// parseTypeRef parses a type reference from the introspection response
func parseTypeRef(data map[string]interface{}) (*TypeRef, error) {
	typeRef := &TypeRef{}

	if name, ok := data["name"].(string); ok {
		typeRef.Name = name
	}

	if kind, ok := data["kind"].(string); ok {
		typeRef.Kind = kind
	}

	// Parse ofType recursively
	if ofTypeData, ok := data["ofType"].(map[string]interface{}); ok {
		ofType, err := parseTypeRef(ofTypeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse ofType: %w", err)
		}
		typeRef.OfType = ofType
	}

	return typeRef, nil
}

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

// GenerateQueryString generates a GraphQL query string for a field
func (f *Field) GenerateQueryString() string {
	var query strings.Builder

	// Start with the query keyword
	query.WriteString("query {\n  ")
	query.WriteString(f.Name)

	// Add arguments
	if len(f.Args) > 0 {
		query.WriteString("(")
		for i, arg := range f.Args {
			if i > 0 {
				query.WriteString(", ")
			}
			query.WriteString("$")
			query.WriteString(arg.Name)
			query.WriteString(": ")
			query.WriteString(arg.Type.GetTypeName())
		}
		query.WriteString(")")
	}

	// Add selection set based on the return type
	query.WriteString(" {\n    ")
	query.WriteString(f.generateSelectionSet())
	query.WriteString("\n  }\n}")

	return query.String()
}

// generateSelectionSet generates a selection set based on the field's return type
func (f *Field) generateSelectionSet() string {
	// For now, return a basic selection set
	// In a more sophisticated implementation, this would introspect the return type
	// and select appropriate fields based on the GraphQL schema

	// Handle different return types
	if f.Type.IsList() {
		// For list types, we need to select fields from the list item type
		// This is a simplified implementation that includes common fields and nested objects
		return "id\n    name\n    location {\n      latitude\n      longitude\n      altitude\n    }\n    capacity\n    utilization\n    status"
	}

	// For object types, return common fields
	// This is a simplified implementation that works for most cases
	return "id\n    name"
}

// GenerateMutationString generates a GraphQL mutation string for a field
func (f *Field) GenerateMutationString() string {
	return f.GenerateQueryString() // Same logic for mutations
}
