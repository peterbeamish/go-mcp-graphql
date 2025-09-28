package schema

import (
	"fmt"

	"github.com/vektah/gqlparser/v2/ast"
)

// ParseIntrospectionResponse parses the introspection response and builds gqlparser AST
// This function manually parses the introspection JSON response since gqlparser's LoadSchema
// is designed for SDL sources, not introspection responses
func ParseIntrospectionResponse(data map[string]interface{}) (*Schema, error) {
	schemaData, ok := data["__schema"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid schema data")
	}

	// Build the gqlparser AST schema directly
	astSchema := &ast.Schema{
		Types: make(map[string]*ast.Definition),
	}

	// Parse all types first (needed for references)
	if typesData, ok := schemaData["types"].([]interface{}); ok {
		for _, typeData := range typesData {
			if typeMap, ok := typeData.(map[string]interface{}); ok {
				astDef, err := parseTypeToAST(typeMap)
				if err != nil {
					return nil, fmt.Errorf("failed to parse type to AST: %w", err)
				}
				if astDef != nil {
					astSchema.Types[astDef.Name] = astDef
				}
			}
		}
	}

	// Set query and mutation types
	if queryTypeData, ok := schemaData["queryType"].(map[string]interface{}); ok {
		if queryTypeName, ok := queryTypeData["name"].(string); ok {
			astSchema.Query = astSchema.Types[queryTypeName]
		}
	}

	if mutationTypeData, ok := schemaData["mutationType"].(map[string]interface{}); ok {
		if mutationTypeName, ok := mutationTypeData["name"].(string); ok {
			astSchema.Mutation = astSchema.Types[mutationTypeName]
		}
	}

	// Create the schema with the parsed AST
	schema := &Schema{
		parsedSchema: astSchema,
		typeRegistry: astSchema.Types,
	}

	// Convert to legacy types for backward compatibility
	schema.QueryType = convertASTToType(astSchema.Query)
	schema.MutationType = convertASTToType(astSchema.Mutation)

	// Convert all types for backward compatibility
	schema.Types = make([]*Type, 0, len(astSchema.Types))
	for _, astDef := range astSchema.Types {
		if astDef != nil && !isBuiltinType(astDef.Name) {
			schema.Types = append(schema.Types, convertASTToType(astDef))
		}
	}

	return schema, nil
}

// parseTypeToAST converts introspection data to gqlparser AST Definition
func parseTypeToAST(data map[string]interface{}) (*ast.Definition, error) {
	name, ok := data["name"].(string)
	if !ok {
		return nil, nil // Skip types without names
	}

	kind, ok := data["kind"].(string)
	if !ok {
		return nil, fmt.Errorf("type %s missing kind", name)
	}

	astDef := &ast.Definition{
		Name:        name,
		Kind:        parseKindToAST(kind),
		Description: getString(data, "description"),
	}

	// Parse fields (for objects, interfaces, etc.)
	if fieldsData, ok := data["fields"].([]interface{}); ok {
		astDef.Fields = make([]*ast.FieldDefinition, 0, len(fieldsData))
		for _, fieldData := range fieldsData {
			if fieldMap, ok := fieldData.(map[string]interface{}); ok {
				field, err := parseFieldToAST(fieldMap)
				if err != nil {
					return nil, fmt.Errorf("failed to parse field: %w", err)
				}
				astDef.Fields = append(astDef.Fields, field)
			}
		}
	}

	// Parse input fields (for input objects)
	if inputFieldsData, ok := data["inputFields"].([]interface{}); ok && len(inputFieldsData) > 0 {
		astDef.Fields = make([]*ast.FieldDefinition, 0, len(inputFieldsData))
		for _, inputFieldData := range inputFieldsData {
			if inputFieldMap, ok := inputFieldData.(map[string]interface{}); ok {
				// Input fields are parsed the same way as regular fields
				// but they don't have arguments, so we use parseInputFieldToAST
				inputField, err := parseInputFieldToAST(inputFieldMap)
				if err != nil {
					return nil, fmt.Errorf("failed to parse input field: %w", err)
				}
				astDef.Fields = append(astDef.Fields, inputField)
			}
		}
	}

	// Parse enum values
	if enumValuesData, ok := data["enumValues"].([]interface{}); ok {
		astDef.EnumValues = make([]*ast.EnumValueDefinition, 0, len(enumValuesData))
		for _, enumValueData := range enumValuesData {
			if enumValueMap, ok := enumValueData.(map[string]interface{}); ok {
				enumValue := &ast.EnumValueDefinition{
					Name:        getString(enumValueMap, "name"),
					Description: getString(enumValueMap, "description"),
				}
				astDef.EnumValues = append(astDef.EnumValues, enumValue)
			}
		}
	}

	// Parse union types
	if possibleTypesData, ok := data["possibleTypes"].([]interface{}); ok {
		astDef.Types = make([]string, 0, len(possibleTypesData))
		for _, possibleTypeData := range possibleTypesData {
			if possibleTypeMap, ok := possibleTypeData.(map[string]interface{}); ok {
				if typeName, ok := possibleTypeMap["name"].(string); ok {
					astDef.Types = append(astDef.Types, typeName)
				}
			}
		}
	}

	// Parse interfaces
	if interfacesData, ok := data["interfaces"].([]interface{}); ok {
		astDef.Interfaces = make([]string, 0, len(interfacesData))
		for _, interfaceData := range interfacesData {
			if interfaceMap, ok := interfaceData.(map[string]interface{}); ok {
				if interfaceName, ok := interfaceMap["name"].(string); ok {
					astDef.Interfaces = append(astDef.Interfaces, interfaceName)
				}
			}
		}
	}

	return astDef, nil
}

// parseFieldToAST converts field introspection data to gqlparser AST FieldDefinition
func parseFieldToAST(data map[string]interface{}) (*ast.FieldDefinition, error) {
	return parseFieldDefinition(data, true)
}

// parseInputFieldToAST converts input field introspection data to gqlparser AST FieldDefinition
func parseInputFieldToAST(data map[string]interface{}) (*ast.FieldDefinition, error) {
	return parseFieldDefinition(data, false)
}

// parseFieldDefinition is the unified implementation for parsing field definitions
func parseFieldDefinition(data map[string]interface{}, includeArgs bool) (*ast.FieldDefinition, error) {
	name, ok := data["name"].(string)
	if !ok {
		return nil, fmt.Errorf("field missing name")
	}

	field := &ast.FieldDefinition{
		Name:        name,
		Description: getString(data, "description"),
	}

	// Parse type
	if typeData, ok := data["type"].(map[string]interface{}); ok {
		fieldType, err := parseTypeRefToAST(typeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse field type: %w", err)
		}
		field.Type = fieldType
	}

	// Parse arguments (only for regular fields, not input fields)
	if includeArgs {
		if argsData, ok := data["args"].([]interface{}); ok {
			field.Arguments = make([]*ast.ArgumentDefinition, 0, len(argsData))
			for _, argData := range argsData {
				if argMap, ok := argData.(map[string]interface{}); ok {
					arg, err := parseArgumentToAST(argMap)
					if err != nil {
						return nil, fmt.Errorf("failed to parse argument: %w", err)
					}
					field.Arguments = append(field.Arguments, arg)
				}
			}
		}
	}

	// Parse default value (for input fields)
	if !includeArgs {
		if defaultValue, ok := data["defaultValue"]; ok && defaultValue != nil {
			field.DefaultValue = &ast.Value{
				Kind: ast.StringValue,
				Raw:  fmt.Sprintf("%v", defaultValue),
			}
		}
	}

	return field, nil
}

// parseArgumentToAST converts argument introspection data to gqlparser AST ArgumentDefinition
func parseArgumentToAST(data map[string]interface{}) (*ast.ArgumentDefinition, error) {
	name, ok := data["name"].(string)
	if !ok {
		return nil, fmt.Errorf("argument missing name")
	}

	arg := &ast.ArgumentDefinition{
		Name:        name,
		Description: getString(data, "description"),
	}

	// Parse type
	if typeData, ok := data["type"].(map[string]interface{}); ok {
		argType, err := parseTypeRefToAST(typeData)
		if err != nil {
			return nil, fmt.Errorf("failed to parse argument type: %w", err)
		}
		arg.Type = argType
	}

	// Parse default value
	if defaultValue, ok := data["defaultValue"]; ok && defaultValue != nil {
		arg.DefaultValue = &ast.Value{
			Kind: ast.StringValue,
			Raw:  fmt.Sprintf("%v", defaultValue),
		}
	}

	return arg, nil
}

// parseTypeRefToAST converts type reference introspection data to gqlparser AST Type
func parseTypeRefToAST(data map[string]interface{}) (*ast.Type, error) {
	kind, ok := data["kind"].(string)
	if !ok {
		return nil, fmt.Errorf("type reference missing kind")
	}

	switch kind {
	case "NON_NULL":
		if ofType, ok := data["ofType"].(map[string]interface{}); ok {
			innerType, err := parseTypeRefToAST(ofType)
			if err != nil {
				return nil, err
			}
			// For NON_NULL, we need to check if the inner type is a named type
			// If it is, we should set NamedType directly instead of Elem
			if innerType.NamedType != "" {
				return &ast.Type{
					NonNull:   true,
					NamedType: innerType.NamedType,
				}, nil
			}
			// For complex types (lists, etc.), use Elem
			return &ast.Type{
				NonNull: true,
				Elem:    innerType,
			}, nil
		}
		return nil, fmt.Errorf("NON_NULL type missing ofType")
	case "LIST":
		if ofType, ok := data["ofType"].(map[string]interface{}); ok {
			innerType, err := parseTypeRefToAST(ofType)
			if err != nil {
				return nil, err
			}
			return ast.ListType(innerType, nil), nil
		}
		return nil, fmt.Errorf("LIST type missing ofType")
	default:
		// For scalar types and other named types
		if name, ok := data["name"].(string); ok && name != "" {
			return ast.NamedType(name, nil), nil
		}
		return nil, fmt.Errorf("type missing name")
	}
}

// parseKindToAST converts GraphQL kind string to gqlparser AST DefinitionKind
func parseKindToAST(kind string) ast.DefinitionKind {
	return convertStringToKind(kind)
}
