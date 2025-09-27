package main

import (
	"log"
	"net/http"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server/graphql"
	"github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server/resolver"
)

// StartHTTPServer starts the GraphQL HTTP server
func StartHTTPServer(addr string) error {
	// Create a new GraphQL server
	srv := handler.NewDefaultServer(graphql.NewExecutableSchema(graphql.Config{Resolvers: resolver.NewResolver()}))

	// Add playground for testing
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	// Add introspection endpoint
	http.Handle("/graphql", srv)

	log.Printf("GraphQL server starting on %s", addr)
	log.Println("Playground available at: http://localhost" + addr)
	log.Println("GraphQL endpoint: http://localhost" + addr + "/query")
	log.Println("Introspection endpoint: http://localhost" + addr + "/graphql")

	return http.ListenAndServe(addr, nil)
}

func main() {
	log.Fatal(StartHTTPServer(":8081"))
}
