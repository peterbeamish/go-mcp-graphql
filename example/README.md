# GraphQL MCP Examples

This directory contains comprehensive examples demonstrating how to use the GraphQL MCP library with various GraphQL servers.

## Quick Start

The easiest way to run the complete demo is using the Makefile from the project root:

```bash
# From the project root directory
make demo
```

This will start:
- GraphQL server on port 8081
- MCP client on port 8080  
- Full demo application

## Available Examples

### 1. Gqlgen GraphQL Server (`gqlgen-server/`)
A complete industrial machinery management system built with [gqlgen](https://github.com/99designs/gqlgen).

**Features:**
- Equipment management (CNC machines, robots, conveyors, etc.)
- Facility management with location tracking
- Maintenance scheduling and tracking
- Operational metrics and performance monitoring
- Comprehensive GraphQL schema with full documentation

**Run:**
```bash
make run-graphql
# or
cd gqlgen-server
go generate  # Generate GraphQL code
go run .     # Start the server
```

### 2. MCP Client (`gqlgen-mcp-client/`)
An MCP client that introspects the gqlgen GraphQL server and creates MCP tools.

**Features:**
- Automatic GraphQL schema introspection
- MCP tool generation for all queries and mutations
- HTTP server for MCP protocol
- Health check and tools listing endpoints

**Run:**
```bash
make run-mcp
# or
cd gqlgen-mcp-client
go run .
```

### 3. Full Demo (`full-demo/`)
A complete demonstration showing the entire workflow.

**Features:**
- Starts GraphQL server
- Introspects schema and creates MCP tools
- Tests MCP tools via HTTP client
- Demonstrates end-to-end functionality

**Run:**
```bash
make run-full-demo
# or
cd full-demo
go run .
```

### 4. Simple Examples
- `main.go` - HTTP server with GitHub GraphQL API
- `simple/main.go` - Stdio server with countries API
- `client/main.go` - HTTP client example

## Makefile Commands

From the project root, you can use these commands:

```bash
make help          # Show all available commands
make install       # Install all dependencies
make generate      # Generate gqlgen code
make build         # Build all components
make demo          # Run complete demo
make run-graphql   # Start GraphQL server only
make run-mcp       # Start MCP client only
make test          # Test all services
make status        # Check service status
make clean         # Clean build artifacts
make logs          # Show logs from all services
```

## Testing the Demo

Once the demo is running, you can test the services:

### GraphQL Server
- **Playground:** http://localhost:8081
- **API:** http://localhost:8081/query
- **Introspection:** http://localhost:8081/graphql

### MCP Server
- **Health Check:** http://localhost:8080/health
- **Tools List:** http://localhost:8080/tools
- **MCP Protocol:** http://localhost:8080/mcp

### Example GraphQL Queries

```graphql
# Get all equipment
query {
  equipment {
    id
    name
    type
    status
    efficiency
    facility {
      name
      location {
        latitude
        longitude
      }
    }
  }
}

# Get all facilities
query {
  facilities {
    id
    name
    address
    capacity
    utilization
    equipment {
      id
      name
      type
    }
  }
}

# Create new equipment
mutation {
  createEquipment(input: {
    name: "CNC Mill #1"
    description: "High-precision CNC milling machine"
    manufacturer: "Haas Automation"
    model: "VF-2"
    serialNumber: "HM123456"
    type: CNC_MILL
    facilityId: "1"
    specifications: {
      powerConsumption: 15.5
      maxSpeed: 12000
      operatingTemperature: { min: 10, max: 40 }
      weight: 2500
      dimensions: { length: 2.5, width: 1.5, height: 2.0 }
      electricalSpecs: {
        voltage: 380
        current: 25
        powerFactor: 0.9
        frequency: 50
      }
      environmentalRequirements: ["Clean environment", "Stable temperature"]
      certifications: ["CE", "ISO 9001"]
    }
    installedAt: "2024-01-15T00:00:00Z"
  }) {
    id
    name
    type
    status
  }
}
```

### Example MCP Tool Calls

```bash
# List available tools
curl http://localhost:8080/tools

# Call a query tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "query_equipment", "arguments": {}}}'

# Call a mutation tool
curl -X POST http://localhost:8080/mcp \
  -H "Content-Type: application/json" \
  -d '{"method": "tools/call", "params": {"name": "mutation_createEquipment", "arguments": {"input": {"name": "Test Equipment", "description": "Test", "manufacturer": "Test Corp", "model": "T1", "serialNumber": "TEST123", "type": "CNC_MILL", "facilityId": "1", "specifications": {"powerConsumption": 10, "maxSpeed": 5000, "operatingTemperature": {"min": 0, "max": 50}, "weight": 1000, "dimensions": {"length": 1, "width": 1, "height": 1}, "electricalSpecs": {"voltage": 220, "current": 10, "powerFactor": 0.8, "frequency": 60}, "environmentalRequirements": ["Clean"], "certifications": ["CE"]}, "installedAt": "2024-01-01T00:00:00Z"}}}}'
```

## Troubleshooting

### Services Not Starting
- Check if ports 8080 and 8081 are available
- Ensure all dependencies are installed: `make install`
- Check service status: `make status`

### GraphQL Code Generation Issues
- Run `make generate` to regenerate gqlgen code
- Or use `go generate` from the gqlgen-server directory
- Check if gqlgen is installed: `go install github.com/99designs/gqlgen@latest`

### MCP Connection Issues
- Ensure GraphQL server is running before starting MCP client
- Check GraphQL server is accessible at http://localhost:8081/graphql
- Verify MCP server is running: `curl http://localhost:8080/health`

### Clean Restart
```bash
make clean
make install
make demo
```
