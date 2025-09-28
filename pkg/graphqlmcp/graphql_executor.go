package graphqlmcp

import (
	"context"

	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp/schema"
)

// GraphQLExecutor defines the interface for executing GraphQL operations
// This interface allows for easy mocking in tests
type GraphQLExecutor interface {
	ExecuteQuery(ctx context.Context, query string, variables map[string]interface{}) (*GraphQLResponse, error)
	IntrospectSchema(ctx context.Context) (*schema.Schema, error)
}

// GraphQLClient implements GraphQLExecutor
var _ GraphQLExecutor = (*GraphQLClient)(nil)
