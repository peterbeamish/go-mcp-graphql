package graphqlmcp

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// SSETransportManager manages SSE transports for MCP servers
type SSETransportManager struct {
	transports map[string]*mcp.Server
	mu         sync.RWMutex
}

// NewSSETransportManager creates a new transport manager
func NewSSETransportManager() *SSETransportManager {
	return &SSETransportManager{
		transports: make(map[string]*mcp.Server),
	}
}

// GetOrCreateServer gets an existing server or creates a new one for the session
func (tm *SSETransportManager) GetOrCreateServer(sessionID string, createServer func() (*mcp.Server, error)) (*mcp.Server, error) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	if server, exists := tm.transports[sessionID]; exists {
		return server, nil
	}

	server, err := createServer()
	if err != nil {
		return nil, err
	}

	tm.transports[sessionID] = server
	return server, nil
}

// RemoveServer removes a server from the manager
func (tm *SSETransportManager) RemoveServer(sessionID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()
	delete(tm.transports, sessionID)
}

// HTTPServer wraps an MCP server to be served over HTTP with SSE support
type HTTPServer struct {
	transportManager *SSETransportManager
	createServer     func() (*mcp.Server, error)
	timeout          time.Duration
}

// NewHTTPServer creates a new HTTP server wrapper for MCP
func NewHTTPServer(createServer func() (*mcp.Server, error)) *HTTPServer {
	return &HTTPServer{
		transportManager: NewSSETransportManager(),
		createServer:     createServer,
		timeout:          30 * time.Second,
	}
}

// SetTimeout sets the request timeout
func (h *HTTPServer) SetTimeout(timeout time.Duration) {
	h.timeout = timeout
}

// ServeHTTP implements the http.Handler interface
func (h *HTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Handle CORS
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type,Mcp-Session-Id")

	if r.Method == "OPTIONS" {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Use the streamable handler for all requests
	// This handles both GET (SSE) and POST (MCP messages) automatically
	handler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			// Get or create server for this request
			sessionID := r.URL.Query().Get("sessionId")
			if sessionID == "" {
				sessionID = fmt.Sprintf("session_%d", time.Now().UnixNano())
			}

			server, err := h.transportManager.GetOrCreateServer(sessionID, h.createServer)
			if err != nil {
				log.Printf("Failed to create MCP server for session %s: %v", sessionID, err)
				return nil
			}
			return server
		},
		nil,
	)

	// Delegate to the MCP handler
	handler.ServeHTTP(w, r)
}

// Note: We now use the MCP SDK's built-in StreamableHTTPHandler
// which handles all the HTTP transport details including SSE support

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

// StartHTTPServer starts an HTTP server with the MCP GraphQL server
func GetMux(server *MCPGraphQLServer) *http.ServeMux {
	// Create a mux for routing
	mux := http.NewServeMux()

	// MCP endpoint using the streamable handler
	mcpHandler := mcp.NewStreamableHTTPHandler(
		func(r *http.Request) *mcp.Server {
			return server.GetMCPServer()
		},
		nil,
	)
	mux.Handle("/mcp", mcpHandler)

	// Add a health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "healthy",
			"service": "graphql-mcp-server",
		})
	})

	// Add a schema endpoint to view the GraphQL schema
	mux.HandleFunc("/schema", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		schema := server.GetSchema()
		response := map[string]interface{}{
			"schema": schema,
			"sdl":    schema.GetSchemaSDL(),
		}
		json.NewEncoder(w).Encode(response)
	})

	// Add a tools endpoint to list available MCP tools
	mux.HandleFunc("/tools", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		// Get available tools from the schema
		queries := server.GetSchema().GetQueries()
		mutations := server.GetSchema().GetMutations()

		tools := make([]map[string]interface{}, 0, len(queries)+len(mutations))

		// Add query tools
		for _, query := range queries {
			inputSchema := server.createInputSchema(query)
			tools = append(tools, map[string]interface{}{
				"name":        "query_" + query.Name,
				"description": query.Description,
				"type":        "query",
				"inputSchema": inputSchema,
			})
		}

		// Add mutation tools
		for _, mutation := range mutations {
			inputSchema := server.createInputSchema(mutation)
			tools = append(tools, map[string]interface{}{
				"name":        "mutation_" + mutation.Name,
				"description": mutation.Description,
				"type":        "mutation",
				"inputSchema": inputSchema,
			})
		}

		json.NewEncoder(w).Encode(map[string]interface{}{
			"tools": tools,
			"count": len(tools),
		})
	})

	return mux
}

// CreateHTTPClient creates an HTTP client for communicating with the MCP server
func CreateHTTPClient(baseURL string) *HTTPMCPClient {
	return &HTTPMCPClient{
		baseURL: baseURL,
		client:  &http.Client{Timeout: 30 * time.Second},
		logger:  slog.Default(),
	}
}

// HTTPMCPClient is an HTTP client for MCP communication
type HTTPMCPClient struct {
	baseURL string
	client  *http.Client
	logger  *slog.Logger
}

// SetLogger sets a custom logger for the HTTP client
func (c *HTTPMCPClient) SetLogger(logger *slog.Logger) {
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
		c.logger.Error("Failed to marshal HTTP request",
			"request_id", requestID,
			"tool_name", toolName,
			"error", err,
		)
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.baseURL+"/mcp", bytes.NewBuffer(jsonData))
	if err != nil {
		c.logger.Error("Failed to create HTTP request",
			"request_id", requestID,
			"tool_name", toolName,
			"error", err,
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	// Log the HTTP request details
	c.logger.Debug("Sending HTTP request",
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
		c.logger.Error("HTTP request execution failed",
			"request_id", requestID,
			"tool_name", toolName,
			"duration_ms", duration.Milliseconds(),
			"error", err,
		)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read HTTP response body",
			"request_id", requestID,
			"tool_name", toolName,
			"duration_ms", duration.Milliseconds(),
			"status_code", resp.StatusCode,
			"error", err,
		)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("HTTP request failed with non-OK status",
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
		c.logger.Error("Failed to unmarshal HTTP response",
			"request_id", requestID,
			"tool_name", toolName,
			"duration_ms", duration.Milliseconds(),
			"response_size_bytes", len(body),
			"error", err,
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
		c.logger.Error("Failed to create tools list request",
			"request_id", requestID,
			"error", err,
		)
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	startTime := time.Now()
	resp, err := c.client.Do(req)
	duration := time.Since(startTime)

	if err != nil {
		c.logger.Error("Failed to execute tools list request",
			"request_id", requestID,
			"duration_ms", duration.Milliseconds(),
			"error", err,
		)
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.logger.Error("Failed to read tools list response",
			"request_id", requestID,
			"duration_ms", duration.Milliseconds(),
			"status_code", resp.StatusCode,
			"error", err,
		)
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		c.logger.Error("Tools list request failed with non-OK status",
			"request_id", requestID,
			"duration_ms", duration.Milliseconds(),
			"status_code", resp.StatusCode,
			"response_body", string(body),
		)
		return nil, fmt.Errorf("HTTP request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		c.logger.Error("Failed to unmarshal tools list response",
			"request_id", requestID,
			"duration_ms", duration.Milliseconds(),
			"response_size_bytes", len(body),
			"error", err,
		)
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	tools, ok := result["tools"].([]interface{})
	if !ok {
		c.logger.Error("Invalid tools response format",
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
