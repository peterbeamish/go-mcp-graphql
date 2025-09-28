# Go MCP GraphQL Makefile
# This Makefile provides commands to build, run, and manage the GraphQL MCP project

.PHONY: help install generate build run-graphql run-client run-demo clean test lint

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
CLIENT_DIR := example/client
FULL_DEMO_DIR := example/full-demo

# Ports
GRAPHQL_PORT := 8081
MCP_PORT := 8080

# Help target
help: ## Show this help message
	@echo "$(BLUE)Go MCP GraphQL Project$(NC)"
	@echo "============================="
	@echo ""
	@echo "Available commands:"
	@echo ""
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(GREEN)%-20s$(NC) %s\n", $$1, $$2}'
	@echo ""
	@echo "Quick start:"
	@echo "  $(YELLOW)make demo$(NC) - Run the complete demo"
	@echo "  $(YELLOW)make run-graphql$(NC) - Start GraphQL server only"
	@echo "  $(YELLOW)make run-client$(NC) - Start client only"
	@echo ""
	@echo "Development:"
	@echo "  $(YELLOW)make format$(NC) - Format all code"
	@echo "  $(YELLOW)make test$(NC) - Run tests"
	@echo "  $(YELLOW)make clean$(NC) - Clean build artifacts"

# Install dependencies
install: ## Install all dependencies
	@echo "$(GREEN)[INFO]$(NC) Installing dependencies..."
	@cd $(GRAPHQL_SERVER_DIR) && go mod tidy
	@cd $(CLIENT_DIR) && go mod tidy
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
	@cd $(GRAPHQL_SERVER_DIR) && go build -o gqlgen-server .
	@echo "$(GREEN)[INFO]$(NC) Building client..."
	@cd $(CLIENT_DIR) && go build -o client .
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

# Run client
run-client: ## Start the client
	@echo "$(GREEN)[INFO]$(NC) Starting client..."
	@cd $(CLIENT_DIR) && go run .

# Run full demo
run-demo: ## Run the complete demo (full-demo application)
	@echo "$(GREEN)[INFO]$(NC) Starting complete demo..."
	@echo "$(YELLOW)[INFO]$(NC) This will run the full-demo application"
	@echo "$(YELLOW)[INFO]$(NC) Press Ctrl+C to stop"
	@echo ""
	@$(MAKE) run-full-demo

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
	@$(MAKE) run-client &
	@sleep 3
	@$(MAKE) run-full-demo &
	@echo "$(GREEN)[INFO]$(NC) Demo started! Check the logs above for endpoints."
	@echo "$(YELLOW)[INFO]$(NC) To stop: make stop-demo"

# Stop background demo
stop-demo: ## Stop all background demo processes
	@echo "$(GREEN)[INFO]$(NC) Stopping demo processes..."
	@pkill -f "go run ." || true
	@pkill -f "gqlgen-server" || true
	@pkill -f "client" || true
	@pkill -f "full-demo" || true
	@echo "$(GREEN)[INFO]$(NC) Demo stopped!"

# Test all services
test: ## Run tests
	@echo "$(GREEN)[INFO]$(NC) Running tests..."
	@go test ./...

# Lint code
lint: ## Run linter on all Go code
	@echo "$(GREEN)[INFO]$(NC) Running linter..."
	@go vet ./...
	@go fmt ./...
	@echo "$(GREEN)[INFO]$(NC) Linting complete!"

# Format code
format: ## Format all Go code
	@echo "$(GREEN)[INFO]$(NC) Formatting code..."
	@gofmt -w .
	@echo "$(GREEN)[INFO]$(NC) Code formatting complete!"

# Clean build artifacts
clean: ## Clean all build artifacts
	@echo "$(GREEN)[INFO]$(NC) Cleaning build artifacts..."
	@rm -f $(GRAPHQL_SERVER_DIR)/gqlgen-server
	@rm -f $(CLIENT_DIR)/client
	@rm -f $(FULL_DEMO_DIR)/full-demo
	@echo "$(GREEN)[INFO]$(NC) Clean complete!"

# Clean and rebuild
rebuild: clean build ## Clean and rebuild all components
