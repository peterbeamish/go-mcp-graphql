package main

import (
	"context"
	"log"
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
)

func main() {
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

	// Wait for context cancellation// Wait for context cancellation
	<-ctx.Done()
	log.Println("Shutting down servers...")

	// Wait for both servers to exit
	wg.Wait()
	log.Println("Server shutdown complete")
}

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
