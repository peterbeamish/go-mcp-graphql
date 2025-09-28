package graphqlmcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp/schema"
)

// GraphQLClient handles introspection and queries to a GraphQL server
type GraphQLClient struct {
	endpoint   string
	httpClient *http.Client
	headers    map[string]string
}

// NewGraphQLClient creates a new GraphQL client
func NewGraphQLClient(endpoint string) *GraphQLClient {
	return &GraphQLClient{
		endpoint: endpoint,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		headers: make(map[string]string),
	}
}

// SetHeader sets a custom header for GraphQL requests
func (c *GraphQLClient) SetHeader(key, value string) {
	c.headers[key] = value
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
    }
  }
}
`

// IntrospectSchema performs GraphQL introspection to get the schema
func (c *GraphQLClient) IntrospectSchema(ctx context.Context) (*schema.Schema, error) {
	req := &GraphQLRequest{
		Query: IntrospectionQuery,
	}

	resp, err := c.executeRequest(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute introspection query: %w", err)
	}

	if len(resp.Errors) > 0 {
		return nil, fmt.Errorf("introspection query failed: %v", resp.Errors)
	}

	// Parse the introspection response
	schemaData, ok := resp.Data.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid introspection response format")
	}

	schema, err := schema.ParseIntrospectionResponse(schemaData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse introspection response: %w", err)
	}

	return schema, nil
}

// ExecuteQuery executes a GraphQL query
func (c *GraphQLClient) ExecuteQuery(ctx context.Context, query string, variables map[string]interface{}) (*GraphQLResponse, error) {
	req := &GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	return c.executeRequest(ctx, req)
}

// executeRequest performs the actual HTTP request
func (c *GraphQLClient) executeRequest(ctx context.Context, req *GraphQLRequest) (*GraphQLResponse, error) {
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	log.Printf("Executing request: %s", string(jsonData))

	httpReq, err := http.NewRequestWithContext(ctx, "POST", c.endpoint, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	for key, value := range c.headers {
		httpReq.Header.Set(key, value)
	}

	resp, err := c.httpClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GraphQL request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var graphqlResp GraphQLResponse
	if err := json.Unmarshal(body, &graphqlResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &graphqlResp, nil
}
