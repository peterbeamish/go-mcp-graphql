module github.com/peterbeamish/go-mcp-graphql/example/full-demo

go 1.24.0

require (
	github.com/99designs/gqlgen v0.17.81
	github.com/peterbeamish/go-mcp-graphql v0.0.0
	github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server v0.0.0
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/google/jsonschema-go v0.3.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/modelcontextprotocol/go-sdk v0.8.0 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/vektah/gqlparser/v2 v2.5.30 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
)

replace github.com/peterbeamish/go-mcp-graphql => ../../

replace github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server => ../gqlgen-server
