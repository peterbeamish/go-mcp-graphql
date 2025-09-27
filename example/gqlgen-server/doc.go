// Package main provides an industrial machinery management GraphQL server.
//
// This package implements a comprehensive GraphQL API for managing industrial equipment,
// manufacturing facilities, maintenance records, and operational metrics. The server is
// built using gqlgen for type-safe GraphQL operations.
//
// The API provides the following main capabilities:
//   - Equipment Management: Register, update, and monitor industrial machinery
//   - Facility Management: Manage manufacturing facilities and their locations
//   - Maintenance Scheduling: Plan and track equipment maintenance activities
//   - Performance Monitoring: Record and analyze operational metrics
//   - Alert Management: Handle equipment alerts and issues
//
// GraphQL Schema Features:
//   - Comprehensive equipment types (CNC machines, robots, conveyors, etc.)
//   - Detailed specifications and technical parameters
//   - Real-time status monitoring and alerting
//   - Maintenance scheduling and tracking
//   - Performance metrics and analytics
//   - Geographic location tracking for facilities
//
// Code Generation:
//
//	This package uses gqlgen for automatic code generation. To regenerate the
//	GraphQL code after schema changes, run:
//
//	  go generate
//
//	Or use the Makefile from the project root:
//
//	  make generate
//
//	The generated code includes:
//	- Type definitions for all GraphQL types
//	- Resolver interfaces and implementations
//	- Input validation and serialization
//	- Query and mutation handlers
//
// Example Usage:
//
//	Start the GraphQL server:
//
//	  go run .
//
//	Or use the Makefile:
//
//	  make run-graphql
//
//	The server will be available at:
//	- GraphQL Playground: http://localhost:8081
//	- GraphQL API: http://localhost:8081/query
//	- Introspection: http://localhost:8081/graphql
//
// Integration with MCP:
//
//	This GraphQL server is designed to work with the MCP (Model Context Protocol)
//	library. The MCP client can introspect this server and automatically generate
//	tools for all available queries and mutations.
//
//	To test the MCP integration:
//
//	  make demo
//
//	This will start both the GraphQL server and the MCP client that introspects it.
//
// Schema Documentation:
//
//	The GraphQL schema includes comprehensive documentation for all types, fields,
//	and operations. This documentation is automatically included in the GraphQL
//	introspection and can be viewed in the GraphQL Playground.
//
//	Key schema types:
//	- Equipment: Industrial machinery with specifications and status
//	- Facility: Manufacturing facilities with location and capacity info
//	- MaintenanceRecord: Maintenance activities and scheduling
//	- OperationalMetric: Performance data and KPIs
//	- EquipmentAlert: Equipment issues and warnings
//
//	The schema supports complex queries and mutations for:
//	- Equipment lifecycle management
//	- Maintenance planning and execution
//	- Performance monitoring and analytics
//	- Facility operations and management
//	- Alert handling and resolution
//
// Development:
//
//	When modifying the GraphQL schema (schema.graphql), you must regenerate
//	the Go code using gqlgen. The generated code will be placed in:
//	- generated.go: Main GraphQL server code
//	- models_gen.go: Generated type definitions
//	- resolver/: Resolver implementations
//
//	After regeneration, you may need to update the resolver implementations
//	in the resolver/ directory to match the new schema.
//
// Testing:
//
//	The server can be tested using:
//	- GraphQL Playground (http://localhost:8081)
//	- MCP client integration
//	- Direct HTTP requests to the GraphQL endpoint
//
//	Example test queries are available in the project documentation.
package main

//go:generate go run github.com/99designs/gqlgen generate
