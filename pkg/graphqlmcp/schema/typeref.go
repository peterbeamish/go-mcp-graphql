package schema

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
		return "object"
	}
}
