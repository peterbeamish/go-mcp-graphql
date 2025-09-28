# Query Generation Details

Learn how the library converts GraphQL operations into MCP tools and generates optimized queries.

## How Query Generation Works

The library automatically introspects GraphQL schemas and converts each query and mutation into an MCP tool. Here's how it works:

### 1. Schema Introspection

The library performs a complete GraphQL introspection to understand:
- Available queries and mutations
- Input types and arguments
- Return types and fields
- Type relationships and hierarchies

### 2. Tool Creation

Each GraphQL operation becomes an MCP tool:
- **Query tools**: `query_operationName`
- **Mutation tools**: `mutation_operationName`

### 3. Input Schema Generation

GraphQL input types are converted to JSON Schema for MCP tool parameters.

## Union Type Support

The library provides comprehensive support for GraphQL union types:

### Automatic Union Detection

- Introspects union types and their possible member types
- Generates proper inline fragments for each union member
- Handles field conflicts between union member types with automatic aliasing

### Query Generation for Unions

```graphql
query {
  equipmentNotifications {
    __typename
    ... on EquipmentAlert {
      id
      description
      EquipmentAlert_type: type
      EquipmentAlert_severity: severity
    }
    ... on MaintenanceReminder {
      id
      description
      MaintenanceReminder_type: type
      MaintenanceReminder_priority: priority
    }
    ... on StatusUpdate {
      id
      description
      newStatus
      changedAt
    }
    ... on PerformanceAlert {
      id
      description
      metricType
      currentValue
      expectedValue
    }
  }
}
```

### Union-Specific Methods

- `GetUnions()` - Returns all union types in the schema
- `GetUnionPossibleTypes(unionName)` - Gets possible types for a union
- `IsUnionType(typeName)` - Checks if a type is a union
- `GetUnionByName(unionName)` - Gets a specific union by name

## Interface Type Support

### Automatic Interface Handling

- Detects interface types and their implementations
- Generates inline fragments for each implementation
- Preserves interface field inheritance

### Query Generation for Interfaces

```graphql
query {
  personnel {
    __typename
    id
    name
    email
    ... on Manager {
      department
      directReports
      level
    }
    ... on Associate {
      jobTitle
      reportsTo {
        id
        name
      }
    }
  }
}
```

## Complex Type Support

### Nested Relationships

- Handles deeply nested object relationships
- Prevents circular reference issues
- Generates appropriate selection sets

### Type Safety

- Validates field types and relationships
- Ensures proper GraphQL syntax generation
- Handles nullable and non-nullable types correctly

### Schema Introspection

- Complete introspection query with all GraphQL features
- Supports directives, subscriptions, and advanced schema features
- Preserves all metadata and documentation

## Input Type Conversion

### GraphQL to JSON Schema Mapping

| GraphQL Type | JSON Schema Type | Notes |
|--------------|------------------|-------|
| String | string | |
| Int | integer | |
| Float | number | |
| Boolean | boolean | |
| ID | string | |
| Enum | string with enum values | |
| List | array | |
| NonNull | required field | |
| Input Object | object | |
| Custom Scalar | string | |

### Example Input Conversion

**GraphQL Input Type:**
```graphql
input CreateUserInput {
  name: String!
  email: String!
  age: Int
  role: UserRole!
  tags: [String!]
}
```

**Generated JSON Schema:**
```json
{
  "type": "object",
  "properties": {
    "name": {
      "type": "string"
    },
    "email": {
      "type": "string"
    },
    "age": {
      "type": "integer"
    },
    "role": {
      "type": "string",
      "enum": ["ADMIN", "USER", "MODERATOR"]
    },
    "tags": {
      "type": "array",
      "items": {
        "type": "string"
      }
    }
  },
  "required": ["name", "email", "role"]
}
```

## Query Optimization

### Selection Set Generation

The library generates optimized selection sets by:
- Including all scalar fields by default
- Adding `__typename` for union and interface types
- Handling nested object relationships
- Avoiding circular references

### Field Aliasing

When field names conflict between union members or interface implementations, the library automatically creates aliases:

```graphql
query {
  searchResults {
    __typename
    ... on User {
      id
      name
      User_email: email
    }
    ... on Post {
      id
      name
      Post_content: content
    }
  }
}
```

## Mutation Support

### Input Validation

Mutations are converted to MCP tools with proper input validation:
- Required fields are marked as required
- Type validation based on GraphQL schema
- Default values are preserved

### Example Mutation Tool

**GraphQL Mutation:**
```graphql
mutation CreatePost($input: CreatePostInput!) {
  createPost(input: $input) {
    id
    title
    content
    author {
      id
      name
    }
  }
}
```

**Generated MCP Tool:**
- Name: `mutation_createPost`
- Input: JSON Schema for `CreatePostInput`
- Description: Includes input field information

## Error Handling

### GraphQL Errors

The library handles GraphQL errors gracefully:
- Captures and reports GraphQL execution errors
- Preserves error messages and locations
- Returns structured error responses

### Validation Errors

Input validation errors are caught before GraphQL execution:
- Type validation
- Required field validation
- Enum value validation

## Performance Considerations

### Query Complexity

The library generates queries that:
- Include necessary fields only
- Avoid over-fetching data
- Handle large result sets efficiently

### Caching

- Schema introspection results are cached
- Tool definitions are generated once
- HTTP client connections are reused

## Customization

### Custom Field Selection

While the library automatically generates selection sets, you can customize them by:
- Modifying the schema introspection
- Creating custom query generators
- Using GraphQL fragments

### Custom Input Handling

For complex input types, you can:
- Provide custom JSON Schema definitions
- Add input validation logic
- Handle special input formats

## Debugging

### Query Logging

Enable debug logging to see generated queries:

```go
logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
    Level: slog.LevelDebug,
}))

server, err := graphqlmcp.NewMCPGraphQLServer("https://api.example.com/graphql",
    graphqlmcp.WithLogger(logger),
)
```

### Schema Inspection

View the introspected schema:

```bash
curl http://localhost:8080/schema
```

### Tool Inspection

List all available tools:

```bash
curl http://localhost:8080/tools
```
