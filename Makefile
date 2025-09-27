# Industrial Machinery MCP GraphQL Demo Makefile
# This Makefile provides commands to build, run, and manage the GraphQL MCP demo

.PHONY: help install generate build run-graphql run-mcp run-demo clean test lint

# Default target
.DEFAULT_GOAL := help

# Colors for output
RED := \033[0;31m
GREEN := \033[0;32m
YELLOW := \033[1;33m
BLUE := \033[0;34m
NC := \033[0m # No Color

# Project directories
GRAPHQL_SERVER_DIR := example/gqlgen-server
MCP_CLIENT_DIR := example/gqlgen-mcp-client
FULL_DEMO_DIR := example/full-demo

# Ports
GRAPHQL_PORT := 8081
MCP_PORT := 8080

# Help target
help: ## Show this help message
	@echo "$(BLUE)Industrial Machinery MCP GraphQL Demo$(NC)"
	@echo "=========================================="
	@echo ""
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "Quick start:"
	@echo "  $(YELLOW)make demo$(NC) - Run the complete demo"
	@echo "  $(YELLOW)make run-graphql$(NC) - Start GraphQL server only"
	@echo "  $(YELLOW)make run-mcp$(NC) - Start MCP client only"
	@echo ""
	@echo "Development:"
	@echo "  $(YELLOW)make install-tools$(NC) - Install development tools"
	@echo "  $(YELLOW)make format$(NC) - Format all code"
	@echo "  $(YELLOW)make quality-all$(NC) - Run all quality checks"
	@echo "  $(YELLOW)make dev$(NC) - Run in development mode with auto-reload"

# Install dependencies
install: ## Install all dependencies
	@echo "$(GREEN)[INFO]$(NC) Installing dependencies..."
	@cd $(GRAPHQL_SERVER_DIR) && go mod tidy
	@cd $(MCP_CLIENT_DIR) && go mod tidy
	@cd $(FULL_DEMO_DIR) && go mod tidy
	@go mod tidy

# Generate gqlgen code
generate: ## Generate gqlgen code for GraphQL server
	@echo "$(GREEN)[INFO]$(NC) Generating gqlgen code..."
	@cd $(GRAPHQL_SERVER_DIR) && go generate
	@echo "$(GREEN)[INFO]$(NC) Code generation complete!"

# Build all components
build: generate ## Build all components
	@echo "$(GREEN)[INFO]$(NC) Building GraphQL server..."
	@cd $(GRAPHQL_SERVER_DIR) && go build -o graphql-server .
	@echo "$(GREEN)[INFO]$(NC) Building MCP client..."
	@cd $(MCP_CLIENT_DIR) && go build -o mcp-client .
	@echo "$(GREEN)[INFO]$(NC) Building full demo..."
	@cd $(FULL_DEMO_DIR) && go build -o full-demo .
	@echo "$(GREEN)[INFO]$(NC) All components built successfully!"

# Run GraphQL server
run-graphql: generate ## Start the GraphQL server
	@echo "$(GREEN)[INFO]$(NC) Starting GraphQL server on port $(GRAPHQL_PORT)..."
	@echo "$(YELLOW)[INFO]$(NC) GraphQL Playground: http://localhost:$(GRAPHQL_PORT)"
	@echo "$(YELLOW)[INFO]$(NC) GraphQL Endpoint: http://localhost:$(GRAPHQL_PORT)/query"
	@echo "$(YELLOW)[INFO]$(NC) Introspection: http://localhost:$(GRAPHQL_PORT)/graphql"
	@cd $(GRAPHQL_SERVER_DIR) && go run .

# Run MCP client (requires GraphQL server to be running)
run-mcp: ## Start the MCP client
	@echo "$(GREEN)[INFO]$(NC) Starting MCP client on port $(MCP_PORT)..."
	@echo "$(YELLOW)[INFO]$(NC) MCP Endpoint: http://localhost:$(MCP_PORT)/mcp"
	@echo "$(YELLOW)[INFO]$(NC) Health Check: http://localhost:$(MCP_PORT)/health"
	@echo "$(YELLOW)[INFO]$(NC) Tools: http://localhost:$(MCP_PORT)/tools"
	@cd $(MCP_CLIENT_DIR) && go run .

# Run full demo
run-demo: ## Run the complete demo (GraphQL + MCP + Demo)
	@echo "$(GREEN)[INFO]$(NC) Starting complete demo..."
	@echo "$(YELLOW)[INFO]$(NC) This will start:"
	@echo "  - GraphQL server on port $(GRAPHQL_PORT)"
	@echo "  - MCP client on port $(MCP_PORT)"
	@echo "  - Full demo application"
	@echo ""
	@echo "$(YELLOW)[INFO]$(NC) Press Ctrl+C to stop all services"
	@echo ""
	@$(MAKE) -j3 run-graphql run-mcp run-full-demo

# Run full demo application
run-full-demo: ## Run the full demo application
	@echo "$(GREEN)[INFO]$(NC) Starting full demo application..."
	@sleep 3
	@cd $(FULL_DEMO_DIR) && go run .

# Run demo in background (for testing)
demo: ## Run the complete demo in background
	@echo "$(GREEN)[INFO]$(NC) Starting demo in background..."
	@$(MAKE) run-graphql &
	@sleep 3
	@$(MAKE) run-mcp &
	@sleep 3
	@$(MAKE) run-full-demo &
	@echo "$(GREEN)[INFO]$(NC) Demo started! Check the logs above for endpoints."
	@echo "$(YELLOW)[INFO]$(NC) To stop: make stop-demo"

# Stop background demo
stop-demo: ## Stop all background demo processes
	@echo "$(GREEN)[INFO]$(NC) Stopping demo processes..."
	@pkill -f "go run ." || true
	@pkill -f "graphql-server" || true
	@pkill -f "mcp-client" || true
	@pkill -f "full-demo" || true
	@echo "$(GREEN)[INFO]$(NC) Demo stopped!"

# Test GraphQL server
test-graphql: ## Test GraphQL server endpoints
	@echo "$(GREEN)[INFO]$(NC) Testing GraphQL server..."
	@echo "Testing GraphQL endpoint..."
	@curl -s -X POST http://localhost:$(GRAPHQL_PORT)/query \
		-H "Content-Type: application/json" \
		-d '{"query":"{ equipment { id name } }"}' | jq . || echo "GraphQL server not responding"
	@echo ""
	@echo "Testing introspection..."
	@curl -s -X POST http://localhost:$(GRAPHQL_PORT)/graphql \
		-H "Content-Type: application/json" \
		-d '{"query":"{ __schema { queryType { name } } }"}' | jq . || echo "Introspection not responding"

# Test MCP server
test-mcp: ## Test MCP server endpoints
	@echo "$(GREEN)[INFO]$(NC) Testing MCP server..."
	@echo "Testing health endpoint..."
	@curl -s http://localhost:$(MCP_PORT)/health | jq . || echo "MCP server not responding"
	@echo ""
	@echo "Testing tools endpoint..."
	@curl -s http://localhost:$(MCP_PORT)/tools | jq . || echo "Tools endpoint not responding"
	@echo ""
	@echo "Testing MCP endpoint..."
	@curl -s -X POST http://localhost:$(MCP_PORT)/mcp \
		-H "Content-Type: application/json" \
		-d '{"method": "tools/list", "params": {}}' | jq . || echo "MCP endpoint not responding"

# Test all services
test: test-graphql test-mcp ## Test all services

# Lint code
lint: ## Run linter on all Go code
	@echo "$(GREEN)[INFO]$(NC) Running linter..."
	@go vet ./...
	@go fmt ./...
	@echo "$(GREEN)[INFO]$(NC) Linting complete!"

# Run comprehensive linting with golangci-lint
lint-comprehensive: ## Run comprehensive linting with golangci-lint
	@echo "$(GREEN)[INFO]$(NC) Running comprehensive linting..."
	@golangci-lint run ./...
	@echo "$(GREEN)[INFO]$(NC) Comprehensive linting complete!"

# Run security scanning
security: ## Run security scanning with gosec
	@echo "$(GREEN)[INFO]$(NC) Running security scan..."
	@gosec ./...
	@echo "$(GREEN)[INFO]$(NC) Security scan complete!"

# Run code quality checks
quality: ## Run code quality checks
	@echo "$(GREEN)[INFO]$(NC) Running code quality checks..."
	@ineffassign ./...
	@errcheck ./...
	@echo "$(GREEN)[INFO]$(NC) Code quality checks complete!"

# Format code
format: ## Format all Go code
	@echo "$(GREEN)[INFO]$(NC) Formatting code..."
	@goimports -w .
	@gofmt -w .
	@gci write --skip-generated .
	@echo "$(GREEN)[INFO]$(NC) Code formatting complete!"

# Generate mocks
mocks: ## Generate mocks using mockery
	@echo "$(GREEN)[INFO]$(NC) Generating mocks..."
	@mockery --all
	@echo "$(GREEN)[INFO]$(NC) Mock generation complete!"

# Run all quality checks
quality-all: format lint-comprehensive security quality ## Run all quality checks

# Clean build artifacts
clean: ## Clean all build artifacts
	@echo "$(GREEN)[INFO]$(NC) Cleaning build artifacts..."
	@rm -f $(GRAPHQL_SERVER_DIR)/graphql-server
	@rm -f $(MCP_CLIENT_DIR)/mcp-client
	@rm -f $(FULL_DEMO_DIR)/full-demo
	@rm -f $(GRAPHQL_SERVER_DIR)/generated.go
	@rm -f $(GRAPHQL_SERVER_DIR)/models_gen.go
	@rm -rf $(GRAPHQL_SERVER_DIR)/resolver
	@echo "$(GREEN)[INFO]$(NC) Clean complete!"

# Clean and rebuild
rebuild: clean build ## Clean and rebuild all components

# Show status of all services
status: ## Show status of all services
	@echo "$(GREEN)[INFO]$(NC) Checking service status..."
	@echo ""
	@echo "GraphQL Server (port $(GRAPHQL_PORT)):"
	@curl -s http://localhost:$(GRAPHQL_PORT)/query > /dev/null && echo "  $(GREEN)✓ Running$(NC)" || echo "  $(RED)✗ Not running$(NC)"
	@echo ""
	@echo "MCP Server (port $(MCP_PORT)):"
	@curl -s http://localhost:$(MCP_PORT)/health > /dev/null && echo "  $(GREEN)✓ Running$(NC)" || echo "  $(RED)✗ Not running$(NC)"

# Development mode - run with auto-reload (requires air)
dev: ## Run in development mode with auto-reload
	@echo "$(GREEN)[INFO]$(NC) Starting development mode..."
	@echo "$(YELLOW)[INFO]$(NC) This requires 'air' to be installed: go install github.com/cosmtrek/air@latest"
	@$(MAKE) run-graphql &
	@sleep 3
	@$(MAKE) run-mcp &
	@sleep 3
	@cd $(FULL_DEMO_DIR) && air

# Install development tools
install-tools: ## Install development tools
	@echo "$(GREEN)[INFO]$(NC) Installing development tools..."
	@go install github.com/cosmtrek/air@latest
	@go install github.com/99designs/gqlgen@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@go install github.com/vektra/mockery/v2@latest
	@go install github.com/securecodewarrior/gosec/v2@latest
	@go install github.com/gordonklaus/ineffassign@latest
	@go install github.com/kisielk/errcheck@latest
	@go install golang.org/x/tools/cmd/goimports@latest
	@go install golang.org/x/tools/cmd/godoc@latest
	@go install github.com/daixiang0/gci@latest
	@echo "$(GREEN)[INFO]$(NC) Development tools installed!"

# Install tools from go.mod
install-tools-mod: ## Install tools from go.mod dependencies
	@echo "$(GREEN)[INFO]$(NC) Installing tools from go.mod..."
	@go mod download
	@go install $(shell go list -f '{{range .Imports}}{{.}} {{end}}' -tags tools ./example/gqlgen-server)
	@echo "$(GREEN)[INFO]$(NC) Tools from go.mod installed!"

# Show logs
logs: ## Show logs from all services
	@echo "$(GREEN)[INFO]$(NC) Showing logs from all services..."
	@echo "$(YELLOW)[INFO]$(NC) Use 'make logs-graphql' or 'make logs-mcp' for specific services"

# Show GraphQL server logs
logs-graphql: ## Show GraphQL server logs
	@echo "$(GREEN)[INFO]$(NC) GraphQL server logs:"
	@ps aux | grep "go run" | grep -v grep || echo "No GraphQL server running"

# Show MCP server logs
logs-mcp: ## Show MCP server logs
	@echo "$(GREEN)[INFO]$(NC) MCP server logs:"
	@ps aux | grep "mcp-client" | grep -v grep || echo "No MCP server running"

# Quick test - just check if services are running
quick-test: ## Quick test to check if services are running
	@echo "$(GREEN)[INFO]$(NC) Quick service check..."
	@$(MAKE) status

# Full test suite
test-suite: build test ## Run full test suite (build + test)

# Production build
prod-build: ## Build production-ready binaries
	@echo "$(GREEN)[INFO]$(NC) Building production binaries..."
	@cd $(GRAPHQL_SERVER_DIR) && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o graphql-server .
	@cd $(MCP_CLIENT_DIR) && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o mcp-client .
	@cd $(FULL_DEMO_DIR) && CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o full-demo .
	@echo "$(GREEN)[INFO]$(NC) Production binaries built!"

# Docker support (if needed)
docker-build: ## Build Docker images
	@echo "$(GREEN)[INFO]$(NC) Building Docker images..."
	@echo "$(YELLOW)[INFO]$(NC) Docker support not implemented yet"

# Show project information
info: ## Show project information
	@echo "$(BLUE)Industrial Machinery MCP GraphQL Demo$(NC)"
	@echo "=========================================="
	@echo ""
	@echo "Project Structure:"
	@echo "  $(GRAPHQL_SERVER_DIR)/     - GraphQL server (gqlgen)"
	@echo "  $(MCP_CLIENT_DIR)/         - MCP client"
	@echo "  $(FULL_DEMO_DIR)/          - Full demo application"
	@echo ""
	@echo "Ports:"
	@echo "  GraphQL Server: $(GRAPHQL_PORT)"
	@echo "  MCP Server: $(MCP_PORT)"
	@echo ""
	@echo "Endpoints:"
	@echo "  GraphQL Playground: http://localhost:$(GRAPHQL_PORT)"
	@echo "  GraphQL API: http://localhost:$(GRAPHQL_PORT)/query"
	@echo "  MCP API: http://localhost:$(MCP_PORT)/mcp"
	@echo "  MCP Health: http://localhost:$(MCP_PORT)/health"
	@echo "  MCP Tools: http://localhost:$(MCP_PORT)/tools"
