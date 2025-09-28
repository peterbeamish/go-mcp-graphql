package graphqlmcp

import (
	"context"

	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp/schema"
)

// Context key for passthru headers
type passthruHeadersKey struct{}

// AddPassthruHeadersToContext adds passthru headers to the context
func AddPassthruHeadersToContext(ctx context.Context, headers map[string]string) context.Context {
	return context.WithValue(ctx, passthruHeadersKey{}, headers)
}

// GetPassthruHeaders retrieves passthru headers from the context
func GetPassthruHeaders(ctx context.Context) map[string]string {
	if headers, ok := ctx.Value(passthruHeadersKey{}).(map[string]string); ok {
		return headers
	}
	return nil
}

// GraphQLExecutor defines the interface for executing GraphQL operations
// This interface allows for easy mocking in tests
type GraphQLExecutor interface {
	ExecuteQuery(ctx context.Context, query string, variables map[string]interface{}) (*GraphQLResponse, error)
	IntrospectSchema(ctx context.Context) (*schema.Schema, error)
}

// GraphQLClient implements GraphQLExecutor
var _ GraphQLExecutor = (*GraphQLClient)(nil)
