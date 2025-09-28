package graphqlmcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp/schema"
)

// GraphQLClient handles introspection and queries to a GraphQL server
type GraphQLClient struct {
	endpoint   string
	httpClient *http.Client
	headers    map[string]string
	logger     *slog.Logger
}

// NewGraphQLClient creates a new GraphQL client
func NewGraphQLClient(endpoint string) *GraphQLClient {
	return &GraphQLClient{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
		logger:  slog.Default(),
	}
}

// SetHeader sets a custom header for GraphQL requests
func (c *GraphQLClient) SetHeader(key, value string) {
	c.headers[key] = value
}

// SetLogger sets a custom logger for the GraphQL client
func (c *GraphQLClient) SetLogger(logger *slog.Logger) {
	c.logger = logger
}

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string `json:"message"`
	} `json:"errors,omitempty"`
}

// IntrospectionQuery is the standard GraphQL introspection query
const IntrospectionQuery = `
query IntrospectionQuery {
  __schema {
    queryType {
      name
      fields {
        name
        description
        type {
          name
          kind
          ofType {
            name
            kind
            ofType {
              name
              kind
              ofType {
                name
                kind
              }
            }
          }
        }
        args {
          name
          type {
            name
            kind
            ofType {
              name
              kind
            }
          }
          defaultValue
        }
      }
    }
    mutationType {
      name
      fields {
        name
        description
        type {
          name
          kind
          ofType {
            name
            kind
            ofType {
              name
              kind
              ofType {
                name
                kind
              }
            }
          }
        }
        args {
          name
          type {
            name
            kind
            ofType {
              name
              kind
            }
          }
          defaultValue
        }
      }
    }
    types {
      name
      kind
      description
      fields {
        name
        description
        type {
          name
          kind
          ofType {
            name
            kind
            ofType {
              name
              kind
              ofType {
                name
                kind
              }
            }
          }
        }
        args {
          name
          type {
            name
            kind
            ofType {
              name
              kind
            }
          }
          defaultValue
        }
      }
      inputFields {
        name
        description
        type {
          name
          kind
          ofType {
            name
            kind
            ofType {
              name
              kind
              ofType {
                name
                kind
              }
            }
          }
        }
        defaultValue
      }
      enumValues {
        name
        description
      }
    }
  }
}
`

// IntrospectSchema performs GraphQL introspection to get the schema
func (c *GraphQLClient) IntrospectSchema(ctx context.Context) (*schema.Schema, error) {
	requestID := fmt.Sprintf("introspect_%d", time.Now().UnixNano())

	c.logger.Info("Introspecting GraphQL schema",
		"request_id", requestID,
		"endpoint", c.endpoint,
	)

	req := &GraphQLRequest{
		Query: IntrospectionQuery,
	}

	resp, err := c.executeRequest(ctx, req, requestID)
	if err != nil {
		c.logger.Error("Introspection query execution failed",
			"request_id", requestID,
			"endpoint", c.endpoint,
			"error", err,
		)
		return nil, fmt.Errorf("failed to execute introspection query: %w", err)
	}

	if len(resp.Errors) > 0 {
		c.logger.Error("Introspection query returned errors",
			"request_id", requestID,
			"endpoint", c.endpoint,
			"error_count", len(resp.Errors),
			"errors", resp.Errors,
		)
		return nil, fmt.Errorf("introspection query failed: %v", resp.Errors)
	}

	// Parse the introspection response
	schemaData, ok := resp.Data.(map[string]interface{})
	if !ok {
		c.logger.Error("Invalid introspection response format",
			"request_id", requestID,
			"endpoint", c.endpoint,
			"response_data_type", fmt.Sprintf("%T", resp.Data),
		)
		return nil, fmt.Errorf("invalid introspection response format")
	}

	schema, err := schema.ParseIntrospectionResponse(schemaData)
	if err != nil {
		c.logger.Error("Failed to parse introspection response",
			"request_id", requestID,
			"endpoint", c.endpoint,
			"error", err,
		)
		return nil, fmt.Errorf("failed to parse introspection response: %w", err)
	}

	c.logger.Info("Schema introspection completed successfully",
		"request_id", requestID,
		"endpoint", c.endpoint,
		"schema_types", len(schema.Types),
	)

	return schema, nil
}

// ExecuteQuery executes a GraphQL query
func (c *GraphQLClient) ExecuteQuery(ctx context.Context, query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	requestID := fmt.Sprintf("gql_%d", time.Now().UnixNano())

	// Log GraphQL query execution
	c.logger.Info("Executing GraphQL query",
		"request_id", requestID,
		"endpoint", c.endpoint,
		"query_length", len(query),
		"variables_count", len(variables),
	)

	req := &GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	return c.executeRequest(ctx, req, requestID)
}

// executeRequest performs the actual HTTP request
func (c *GraphQLClient) executeRequest(ctx context.Context, req *GraphQLRequest, requestID string) (*GraphQLResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		c.logger.Error("Failed to marshal GraphQL request",
			"request_id", requestID,
			"error", err,
		)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Log the request details at debug level
	c.logger.Debug("GraphQL request details",
		"request_id", requestID,
		"query", req.Query,
		"variables", req.Variables,
		"request_size_bytes", len(jsonData),
	)

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create HTTP request",
			"request_id", requestID,
			"endpoint", c.endpoint,
			"error", err,
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	for key, value := range c.headers {
		httpReq.Header.Set(key, value)
	}

	startTime := time.Now()
	resp, err := c.httpClient.Do(httpReq)
	duration := time.Since(startTime)

	if err != nil {
		c.logger.Error("GraphQL HTTP request failed",
			"request_id", requestID,
			"endpoint", c.endpoint,
			"duration_ms", duration.Milliseconds(),
			"error", err,
		)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read GraphQL response body",
			"request_id", requestID,
			"endpoint", c.endpoint,
			"duration_ms", duration.Milliseconds(),
			"status_code", resp.StatusCode,
			"error", err,
		)
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("GraphQL request failed with non-OK status",
			"request_id", requestID,
			"endpoint", c.endpoint,
			"duration_ms", duration.Milliseconds(),
			"status_code", resp.StatusCode,
			"response_body", string(body),
		)
		return nil, fmt.Errorf("GraphQL request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var graphqlResp GraphQLResponse
	if err := json.Unmarshal(body, &graphqlResp); err != nil {
		c.logger.Error("Failed to unmarshal GraphQL response",
			"request_id", requestID,
			"endpoint", c.endpoint,
			"duration_ms", duration.Milliseconds(),
			"response_size_bytes", len(body),
			"error", err,
		)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Log successful completion
	c.logger.Info("GraphQL query completed successfully",
		"request_id", requestID,
		"endpoint", c.endpoint,
		"duration_ms", duration.Milliseconds(),
		"response_size_bytes", len(body),
		"has_errors", len(graphqlResp.Errors) > 0,
		"error_count", len(graphqlResp.Errors),
	)

	return &graphqlResp, nil
}
