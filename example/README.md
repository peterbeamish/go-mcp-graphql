# GraphQL MCP Examples

This directory contains examples demonstrating how to use the GraphQL MCP library.

## Quick Start

Run the complete demo from the project root:

```bash
make run-demo
```

This runs a standalone application that demonstrates the full workflow.

## Available Examples

### 1. Gqlgen GraphQL Server (`gqlgen-server/`)
A complete industrial machinery management system built with [gqlgen](https://github.com/99designs/gqlgen).

**Run:**
```bash
make run-graphql
```

### 2. MCP Client (`client/`)
An MCP client that connects to GraphQL servers and creates MCP tools.

**Run:**
```bash
make run-mcp
```

### 3. Full Demo (`full-demo/`)
A complete demonstration showing the entire workflow in one application.

**Run:**
```bash
make run-demo
```

## Makefile Commands

From the project root:

```bash
make help          # Show all available commands
make run-demo      # Run complete demo
make run-graphql   # Start GraphQL server only
make run-mcp       # Start MCP client only
make test          # Test all services
make clean         # Clean build artifacts
```

## Detailed Documentation

For comprehensive examples, configuration options, and real-world usage patterns, see:

- **[Examples & Use Cases](../docs/examples.md)** - Real-world examples and API reference
- **[Configuration Guide](../docs/config.md)** - Authentication, logging, and filtering
- **[Query Generation](../docs/query-generation.md)** - How GraphQL operations become MCP tools
- **[Quick Start Guide](../docs/quickstart.md)** - Getting started with the library

## Testing the Demo

Once running, test the services:

- **GraphQL Playground:** http://localhost:8081
- **MCP Health Check:** http://localhost:8080/health
- **MCP Tools List:** http://localhost:8080/tools

For detailed testing examples and troubleshooting, see the [Examples documentation](../docs/examples.md).