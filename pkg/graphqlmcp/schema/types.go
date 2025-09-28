package schema

import (
	"github.com/vektah/gqlparser/v2/ast"
)

// Schema represents a GraphQL schema
type Schema struct {
	QueryType    *Type   `json:"queryType"`
	MutationType *Type   `json:"mutationType"`
	Types        []*Type `json:"types"`

	// Parsed schema for dynamic introspection using gqlparser
	parsedSchema *ast.Schema
	typeRegistry map[string]*ast.Definition
}

// Type represents a GraphQL type
type Type struct {
	Name        string      `json:"name"`
	Kind        string      `json:"kind"`
	Description string      `json:"description"`
	Fields      []*Field    `json:"fields"`
	Args        []*Argument `json:"args"`

	// Interface-specific fields
	Interfaces    []string `json:"interfaces,omitempty"`    // For objects that implement interfaces
	PossibleTypes []string `json:"possibleTypes,omitempty"` // For interfaces and unions
}

// Field represents a GraphQL field
type Field struct {
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        *TypeRef    `json:"type"`
	Args        []*Argument `json:"args"`

	// AST type information for dynamic query generation
	ASTType *ast.Type
}

// Argument represents a GraphQL argument
type Argument struct {
	Name         string   `json:"name"`
	Description  string   `json:"description"`
	Type         *TypeRef `json:"type"`
	DefaultValue string   `json:"defaultValue"`
}

// TypeRef represents a GraphQL type reference
type TypeRef struct {
	Name   string   `json:"name"`
	Kind   string   `json:"kind"`
	OfType *TypeRef `json:"ofType"`
}
