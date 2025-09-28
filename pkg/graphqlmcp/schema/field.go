package schema

import (
	"fmt"
	"strings"

	"github.com/vektah/gqlparser/v2/ast"
)

// GenerateQueryStringWithSchema generates a GraphQL query string for a field with schema context
func (f *Field) GenerateQueryStringWithSchema(schema *Schema) (string, error) {
	return f.generateOperationString(schema, "query")
}

// GenerateMutationStringWithSchema generates a GraphQL mutation string for a field
func (f *Field) GenerateMutationStringWithSchema(schema *Schema) (string, error) {
	return f.generateOperationString(schema, "mutation")
}

// generateOperationString generates a GraphQL operation string (query or mutation)
func (f *Field) generateOperationString(schema *Schema, operationType string) (string, error) {
	if schema == nil {
		return "", fmt.Errorf("schema is nil")
	}

	var operation strings.Builder

	// Start with the operation keyword and variable declarations
	operation.WriteString(operationType)
	if len(f.Args) > 0 {
		operation.WriteString("(")
		for i, arg := range f.Args {
			if i > 0 {
				operation.WriteString(", ")
			}
			operation.WriteString("$")
			operation.WriteString(arg.Name)
			operation.WriteString(": ")
			operation.WriteString(arg.Type.GetTypeName())
			if arg.Type.IsNonNull() {
				operation.WriteString("!")
			}
		}
		operation.WriteString(")")
	}
	operation.WriteString(" {\n  ")
	operation.WriteString(f.Name)

	// Add field arguments (using variables)
	if len(f.Args) > 0 {
		operation.WriteString("(")
		for i, arg := range f.Args {
			if i > 0 {
				operation.WriteString(", ")
			}
			operation.WriteString(arg.Name)
			operation.WriteString(": $")
			operation.WriteString(arg.Name)
		}
		operation.WriteString(")")
	}

	// Add selection set based on the return type
	operation.WriteString(" {\n    ")

	selectionSet, err := f.generateSelectionSetFromAST(schema)
	if err != nil {
		return "", fmt.Errorf("failed to generate selection set: %w", err)
	}
	operation.WriteString(selectionSet)

	operation.WriteString("\n  }\n}")

	return operation.String(), nil
}

// generateSelectionSetFromAST generates a selection set using the parsed schema
func (f *Field) generateSelectionSetFromAST(schema *Schema) (string, error) {
	if schema == nil || schema.typeRegistry == nil {
		return "", fmt.Errorf("schema or type registry is nil")
	}

	// Get the return type name from the field's type
	typeName := f.getReturnTypeNameFromAST()
	if typeName == "" {
		return "", fmt.Errorf("type name is empty")
	}

	// Get the type definition
	typeDef := schema.GetTypeDefinition(typeName)
	if typeDef == nil {
		return "", fmt.Errorf("type definition is nil for type: %s", typeName)
	}

	// Generate selection set based on the actual type definition
	// Use a visited map to prevent circular references and track the original type
	visited := make(map[string]bool)
	originalType := typeName // Track the original type to avoid self-referencing fields
	return f.generateSelectionSetForTypeWithVisited(typeDef, schema, 0, visited, originalType)
}

// getReturnTypeName extracts the type name from the field's return type
func (f *Field) getReturnTypeName() string {
	if f.Type == nil {
		return ""
	}

	// Handle non-null and list wrappers
	currentType := f.Type
	for currentType != nil {
		if currentType.Name != "" {
			return currentType.Name
		}
		currentType = currentType.OfType
	}
	return ""
}

// getReturnTypeNameFromAST extracts the type name from the field's AST type
func (f *Field) getReturnTypeNameFromAST() string {
	if f.ASTType == nil {
		// Fallback to legacy method if AST type is not available
		return f.getReturnTypeName()
	}

	// Use the unified AST type name extraction
	return GetASTTypeName(f.ASTType)
}

// generateSelectionSetForTypeWithVisited generates a selection set with circular reference protection
func (f *Field) generateSelectionSetForTypeWithVisited(typeDef *ast.Definition, schema *Schema, depth int, visited map[string]bool, originalType string) (string, error) {
	if depth > 5 { // Reduced depth limit
		return "", fmt.Errorf("depth is greater than 5")
	}

	// Check for circular references
	if visited[typeDef.Name] {
		return "", nil // Skip this type to prevent circular reference
	}

	// Mark this type as visited
	visited[typeDef.Name] = true
	defer func() {
		// Unmark when we're done with this type
		delete(visited, typeDef.Name)
	}()

	var fields []string

	switch typeDef.Kind {
	case ast.Object, ast.Interface:
		// For objects and interfaces, select their fields
		for _, field := range typeDef.Fields {
			if f.shouldIncludeField(field) {
				// Skip fields that return the same type as the current type being processed to avoid self-referencing
				// This prevents infinite recursion while allowing legitimate cross-references
				fieldTypeName := GetASTTypeName(field.Type)
				if fieldTypeName == typeDef.Name {
					continue // Skip this field to avoid self-referencing
				}

				// Skip fields that would create a circular reference based on the visited map
				if visited[fieldTypeName] {
					continue // Skip this field to avoid circular reference
				}

				fieldSelection, err := f.generateFieldSelectionWithVisited(field, schema, depth+1, visited, originalType)
				if err != nil {
					return "", fmt.Errorf("failed to generate field selection for %s: %w", field.Name, err)
				}
				if fieldSelection != "" {
					fields = append(fields, fieldSelection)
				}
			}
		}
	case ast.Union:
		// For unions, select common fields from all possible types
		for _, possibleType := range typeDef.Types {
			if possibleTypeDef := schema.GetTypeDefinition(possibleType); possibleTypeDef != nil {
				unionFields, err := f.generateSelectionSetForTypeWithVisited(possibleTypeDef, schema, depth+1, visited, originalType)
				if err != nil {
					return "", fmt.Errorf("failed to generate union fields for %s: %w", possibleType, err)
				}
				if unionFields != "" {
					fields = append(fields, unionFields)
				}
			}
		}
	case ast.Enum:
		// For enums, just return empty (no selection set needed)
		return "", nil
	case ast.Scalar:
		// For scalars, just return empty (no selection set needed)
		return "", nil
	}

	return strings.Join(fields, "\n    "), nil
}

// shouldIncludeField determines if a field should be included in the selection set
func (f *Field) shouldIncludeField(field *ast.FieldDefinition) bool {
	// Skip fields that are likely to cause circular references or are too complex
	skipFields := map[string]bool{
		"__typename": true,
		"__schema":   true,
		"__type":     true,
	}

	if skipFields[field.Name] {
		return false
	}

	// Skip fields with complex arguments (to keep queries simple)
	if len(field.Arguments) > 2 {
		return false
	}

	// Skip fields that might cause circular references
	// These are common patterns that can lead to infinite recursion
	circularPatterns := []string{
		"parent", "children", "related", "linked", "associated",
		"owner", "owned", "creator", "created", "updater", "updated",
		"source", "target", "from", "to", "next", "previous",
	}

	for _, pattern := range circularPatterns {
		if strings.Contains(strings.ToLower(field.Name), pattern) {
			return false
		}
	}

	return true
}

// generateFieldSelectionWithVisited generates a selection for a specific field with circular reference protection
func (f *Field) generateFieldSelectionWithVisited(field *ast.FieldDefinition, schema *Schema, depth int, visited map[string]bool, originalType string) (string, error) {
	// Check if this is a scalar type using the schema
	if isScalarTypeWithSchema(field.Type, schema) {
		return field.Name, nil
	}

	// For complex types, generate nested selection
	typeName := GetASTTypeName(field.Type)
	if typeName == "" {
		return field.Name, nil
	}

	typeDef := schema.GetTypeDefinition(typeName)
	if typeDef == nil {
		return field.Name, nil
	}

	// Generate nested selection
	nestedSelection, err := f.generateSelectionSetForTypeWithVisited(typeDef, schema, depth, visited, originalType)
	if err != nil {
		return "", fmt.Errorf("failed to generate nested selection for field %s: %w", field.Name, err)
	}
	if nestedSelection == "" {
		return field.Name, nil
	}

	return fmt.Sprintf("%s {\n      %s\n    }", field.Name, strings.ReplaceAll(nestedSelection, "\n    ", "\n      ")), nil
}
