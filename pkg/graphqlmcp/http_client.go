package graphqlmcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/go-logr/logr"
)

// MCPRequest represents an MCP request over HTTP
type MCPRequest struct {
	Method string      `json:"method"`
	Params interface{} `json:"params,omitempty"`
	ID     interface{} `json:"id,omitempty"`
}

// MCPResponse represents an MCP response over HTTP
type MCPResponse struct {
	Result interface{} `json:"result,omitempty"`
	Error  *MCPError   `json:"error,omitempty"`
	ID     interface{} `json:"id,omitempty"`
}

// MCPError represents an MCP error
type MCPError struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// HTTPMCPClient is an HTTP client for MCP communication
type HTTPMCPClient struct {
	baseURL string
	client  *http.Client
	logger  logr.Logger
}

// CreateHTTPClient creates an HTTP client for communicating with the MCP server
func CreateHTTPClient(baseURL string) *HTTPMCPClient {
	return &HTTPMCPClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
		logger:  logr.Discard(),
	}
}

// SetLogger sets a custom logger for the HTTP client
func (c *HTTPMCPClient) SetLogger(logger logr.Logger) {
	c.logger = logger
}

// CallTool calls an MCP tool via HTTP
func (c *HTTPMCPClient) CallTool(ctx context.Context, toolName string, arguments map[string]interface{}) (*MCPResponse, error) {
	requestID := fmt.Sprintf("http_req_%d", time.Now().UnixNano())

	// Log HTTP tool call initiation
	c.logger.Info("HTTP tool call initiated",
		"request_id", requestID,
		"tool_name", toolName,
		"arguments", arguments,
		"base_url", c.baseURL,
	)

	request := MCPRequest{
		Method: "tools/call",
		Params: map[string]interface{}{
			"name":      toolName,
			"arguments": arguments,
		},
		ID: requestID,
	}

	jsonData, err := json.Marshal(request)
	if err != nil {
		c.logger.Error(err, "Failed to marshal HTTP request",
			"request_id", requestID,
			"tool_name", toolName,
		)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/mcp", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error(err, "Failed to create HTTP request",
			"request_id", requestID,
			"tool_name", toolName,
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Log the HTTP request details
	c.logger.V(1).Info("Sending HTTP request",
		"request_id", requestID,
		"tool_name", toolName,
		"url", req.URL.String(),
		"method", req.Method,
		"request_size_bytes", len(jsonData),
	)

	startTime := time.Now()
	resp, err := c.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		c.logger.Error(err, "HTTP request execution failed",
			"request_id", requestID,
			"tool_name", toolName,
			"duration_ms", duration.Milliseconds(),
		)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err, "Failed to read HTTP response body",
			"request_id", requestID,
			"tool_name", toolName,
			"duration_ms", duration.Milliseconds(),
			"status_code", resp.StatusCode,
		)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Info("HTTP request failed with non-OK status",
			"request_id", requestID,
			"tool_name", toolName,
			"duration_ms", duration.Milliseconds(),
			"status_code", resp.StatusCode,
			"response_body", string(body),
		)
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var mcpResp MCPResponse
	if err := json.Unmarshal(body, &mcpResp); err != nil {
		c.logger.Error(err, "Failed to unmarshal HTTP response",
			"request_id", requestID,
			"tool_name", toolName,
			"duration_ms", duration.Milliseconds(),
			"response_size_bytes", len(body),
		)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Log successful completion
	c.logger.Info("HTTP tool call completed successfully",
		"request_id", requestID,
		"tool_name", toolName,
		"duration_ms", duration.Milliseconds(),
		"response_size_bytes", len(body),
		"has_error", mcpResp.Error != nil,
	)

	return &mcpResp, nil
}

// ListTools lists available MCP tools
func (c *HTTPMCPClient) ListTools(ctx context.Context) ([]map[string]interface{}, error) {
	requestID := fmt.Sprintf("list_tools_%d", time.Now().UnixNano())

	c.logger.Info("Listing MCP tools",
		"request_id", requestID,
		"base_url", c.baseURL,
	)

	req, err := http.NewRequestWithContext(ctx, "GET", c.baseURL+"/tools", nil)
	if err != nil {
		c.logger.Error(err, "Failed to create tools list request",
			"request_id", requestID,
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	startTime := time.Now()
	resp, err := c.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		c.logger.Error(err, "Failed to execute tools list request",
			"request_id", requestID,
			"duration_ms", duration.Milliseconds(),
		)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error(err, "Failed to read tools list response",
			"request_id", requestID,
			"duration_ms", duration.Milliseconds(),
			"status_code", resp.StatusCode,
		)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Info("Tools list request failed with non-OK status",
			"request_id", requestID,
			"duration_ms", duration.Milliseconds(),
			"status_code", resp.StatusCode,
			"response_body", string(body),
		)
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.logger.Error(err, "Failed to unmarshal tools list response",
			"request_id", requestID,
			"duration_ms", duration.Milliseconds(),
			"response_size_bytes", len(body),
		)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	tools, ok := result["tools"].([]interface{})
	if !ok {
		c.logger.Info("Invalid tools response format",
			"request_id", requestID,
			"duration_ms", duration.Milliseconds(),
			"response", result,
		)
		return nil, fmt.Errorf("invalid tools response format")
	}

	toolMaps := make([]map[string]interface{}, len(tools))
	for i, tool := range tools {
		if toolMap, ok := tool.(map[string]interface{}); ok {
			toolMaps[i] = toolMap
		}
	}

	c.logger.Info("Successfully listed MCP tools",
		"request_id", requestID,
		"duration_ms", duration.Milliseconds(),
		"tool_count", len(toolMaps),
	)

	return toolMaps, nil
}
