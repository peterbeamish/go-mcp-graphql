module github.com/peterbeamish/go-mcp-graphql/example/full-demo

go 1.25

require (
	github.com/99designs/gqlgen v0.17.87
	github.com/go-logr/logr v1.4.3
	github.com/peterbeamish/go-mcp-graphql v0.0.0
	github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server v0.0.0
)

require (
	github.com/agnivade/levenshtein v1.2.1 // indirect
	github.com/go-viper/mapstructure/v2 v2.5.0 // indirect
	github.com/google/jsonschema-go v0.4.2 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/modelcontextprotocol/go-sdk v1.4.0 // indirect
	github.com/segmentio/asm v1.1.3 // indirect
	github.com/segmentio/encoding v0.5.3 // indirect
	github.com/sosodev/duration v1.3.1 // indirect
	github.com/vektah/gqlparser/v2 v2.5.32 // indirect
	github.com/yosida95/uritemplate/v3 v3.0.2 // indirect
	golang.org/x/oauth2 v0.34.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.41.0 // indirect
)

replace github.com/peterbeamish/go-mcp-graphql => ../../

replace github.com/peterbeamish/go-mcp-graphql/example/gqlgen-server => ../gqlgen-server
