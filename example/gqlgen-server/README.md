# Gqlgen GraphQL Server Example

This example demonstrates how to create a GraphQL server using [gqlgen](https://github.com/99designs/gqlgen) and then generate MCP tools from it.

## What this example includes

- A complete GraphQL schema for a blog API with users and posts
- gqlgen-generated resolvers and models
- In-memory data store for demonstration
- Full CRUD operations for both users and posts

## GraphQL Schema

The server provides the following operations:

### Queries
- `posts`: Get all posts
- `post(id: ID!)`: Get a specific post by ID
- `users`: Get all users
- `user(id: ID!)`: Get a specific user by ID

### Mutations
- `createPost(input: CreatePostInput!)`: Create a new post
- `updatePost(id: ID!, input: UpdatePostInput!)`: Update an existing post
- `deletePost(id: ID!)`: Delete a post
- `createUser(input: CreateUserInput!)`: Create a new user
- `updateUser(id: ID!, input: UpdateUserInput!)`: Update an existing user
- `deleteUser(id: ID!)`: Delete a user

## Running the example

1. **Generate the gqlgen code:**
   ```bash
   cd example/gqlgen-server
   ./generate.sh
   ```

2. **Start the GraphQL server:**
   ```bash
   go run .
   ```

3. **Test the GraphQL server:**
   - Open http://localhost:8081 in your browser for the GraphQL playground
   - The GraphQL endpoint is available at http://localhost:8081/query
   - Introspection is available at http://localhost:8081/graphql

## Testing with GraphQL Playground

You can test the following queries in the playground:

```graphql
# Get all posts
query {
  posts {
    id
    title
    content
    author {
      id
      name
      email
    }
    publishedAt
    tags
  }
}

# Get all users
query {
  users {
    id
    name
    email
    createdAt
  }
}

# Create a new post
mutation {
  createPost(input: {
    title: "My New Post"
    content: "This is the content of my new post"
    authorId: "1"
    tags: ["demo", "test"]
  }) {
    id
    title
    content
    author {
      name
    }
  }
}

# Create a new user
mutation {
  createUser(input: {
    name: "Alice Johnson"
    email: "alice@example.com"
  }) {
    id
    name
    email
    createdAt
  }
}
```

## Next Steps

Once the GraphQL server is running, you can:

1. Use the MCP client example to introspect this server and generate MCP tools
2. Test the generated MCP tools via HTTP
3. Integrate the MCP tools with other applications

See the parent directory for examples of how to use this GraphQL server with the MCP library.
