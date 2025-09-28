package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server/graphql"
	"github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server/resolver"
	"github.com/peterbeamish/go-mcp-graphql/pkg/graphqlmcp"
)

var (
	mcpServerOnce sync.Once
	mcpServer     *graphqlmcp.MCPGraphQLServer
	mcpServerErr  error
)

func main() {
	// Configure structured logging
	logger := graphqlmcp.ConfigureVerboseLogging()
	logger.Info("Starting MCP GraphQL server with verbose logging")

	// Create context with cancellation tied to OS signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Create WaitGroup to wait for both servers to exit
	var wg sync.WaitGroup

	// Set up signal handling
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	// Start signal handler in background
	go func() {
		<-sigChan
		log.Println("Received shutdown signal, canceling context...")
		cancel()
	}()

	// Start GraphQL server
	wg.Add(1)
	go func() {
		defer wg.Done()
		startGraphQLServer(ctx)
	}()

	// Wait a moment for GraphQL server to start
	time.Sleep(2 * time.Second)

	// Start MCP server
	wg.Add(1)
	go func() {
		defer wg.Done()
		startMCPServer(ctx)
	}()

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Shutting down servers...")

	// Wait for both servers to exit
	wg.Wait()
	log.Println("Server shutdown complete")
}

// startGraphQLServer starts the GraphQL server with context cancellation
func startGraphQLServer(ctx context.Context) {
	// Create GraphQL server
	graphqlResolver := resolver.NewResolver()
	graphqlServer := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: graphqlResolver}))

	// Create HTTP server
	graphqlMux := http.NewServeMux()
	graphqlMux.Handle("/", playground.Handler("GraphQL playground", "/query"))
	graphqlMux.Handle("/query", graphqlServer)
	graphqlMux.Handle("/graphql", graphqlServer)

	server := &http.Server{
		Addr:    ":8080",
		Handler: graphqlMux,
	}

	// Start server in background
	go func() {
		log.Println("Starting GraphQL server on :8080...")
		log.Println("📊 GraphQL Playground: http://localhost:8080")
		log.Println("🔍 GraphQL Endpoint: http://localhost:8080/query")
		log.Println("📋 Introspection: http://localhost:8080/graphql")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("GraphQL server failed: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	log.Println("Shutting down GraphQL server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("GraphQL server shutdown error: %v", err)
	} else {
		log.Println("GraphQL server shutdown complete")
	}
}

// startMCPServer starts the MCP server with context cancellation
func startMCPServer(ctx context.Context) {
	// Create MCP server
	server, err := getOrCreateMCPServer()
	if err != nil {
		log.Fatalf("Failed to create MCP server: %v", err)
	}

	// Create MCP HTTP server with individual handlers
	mcpMux := http.NewServeMux()
	mcpMux.Handle("/mcp", graphqlmcp.GetMCPHandler(server))
	mcpMux.HandleFunc("/health", graphqlmcp.GetHealthHandler())
	mcpMux.HandleFunc("/schema", graphqlmcp.GetSchemaHandler(server))
	mcpMux.HandleFunc("/tools", graphqlmcp.GetToolsHandler(server))

	httpServer := &http.Server{
		Addr:    ":8081",
		Handler: mcpMux,
	}
	// Start server in background
	go func() {
		log.Println("🚀 Starting MCP server on :8081")
		log.Println("🤖 MCP Server: http://localhost:8081/mcp")
		log.Println("🛠️  Tools List: http://localhost:8081/tools")
		log.Println("📋 Schema Info: http://localhost:8081/schema")
		log.Println("❤️  Health Check: http://localhost:8081/health")
		log.Println("")
		log.Println("Available MCP tools:")
		log.Println("  - MCP tools are available at http://localhost:8081/mcp")
		log.Println("  - Uses Server-Sent Events (SSE) for streaming")
		log.Println("  - Uses http_server.go StartHTTPServer pattern")
		log.Println("")

		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Printf("MCP server failed: %v", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()

	log.Println("Shutting down MCP server...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Printf("GraphQL server shutdown error: %v", err)
	} else {
		log.Println("GraphQL server shutdown complete")
	}
}

func getOrCreateMCPServer() (*graphqlmcp.MCPGraphQLServer, error) {
	mcpServerOnce.Do(func() {
		logger := slog.Default()
		logger.Info("Creating MCP server...")
		graphqlURL := "http://localhost:8080/query"
		mcpServer, mcpServerErr = graphqlmcp.NewMCPGraphQLServer(graphqlURL)
		if mcpServerErr != nil {
			logger.Error("Failed to create MCP server", "error", mcpServerErr)
		} else {
			logger.Info("MCP server created successfully")
			// Set the logger on the server for tool call logging
			mcpServer.SetLogger(logger)
		}
	})
	return mcpServer, mcpServerErr
}
