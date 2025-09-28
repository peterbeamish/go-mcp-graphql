package graphqlmcp

import (
	"regexp"

	"github.com/go-logr/logr"
)

// MCPGraphQLServerOptions holds configuration options for the MCP GraphQL server
type MCPGraphQLServerOptions struct {
	Logger          logr.Logger
	Mask            *MaskConfig
	PassthruHeaders []string
}

// MaskConfig defines how to filter queries and mutations
type MaskConfig struct {
	// AllowList defines patterns that are allowed (if non-empty, only these patterns are allowed)
	AllowList []string
	// BlockList defines patterns that are blocked (these patterns are never allowed)
	BlockList []string
	// CompiledAllowList contains compiled regex patterns for AllowList
	CompiledAllowList []*regexp.Regexp
	// CompiledBlockList contains compiled regex patterns for BlockList
	CompiledBlockList []*regexp.Regexp
}

// MCPGraphQLServerOption is a function that configures MCPGraphQLServerOptions
type MCPGraphQLServerOption func(*MCPGraphQLServerOptions)

// WithLogger sets a custom logger for the MCP GraphQL server
func WithLogger(logger logr.Logger) MCPGraphQLServerOption {
	return func(opts *MCPGraphQLServerOptions) {
		opts.Logger = logger
	}
}

// WithMask configures query/mutation filtering by name or pattern
// allowList: patterns that are allowed (if non-empty, only these patterns are allowed)
// blockList: patterns that are blocked (these patterns are never allowed)
func WithMask(allowList, blockList []string) MCPGraphQLServerOption {
	return func(opts *MCPGraphQLServerOptions) {
		opts.Mask = &MaskConfig{
			AllowList: allowList,
			BlockList: blockList,
		}

		// Compile regex patterns for efficient matching
		opts.Mask.CompiledAllowList = make([]*regexp.Regexp, 0, len(allowList))
		for _, pattern := range allowList {
			if compiled, err := regexp.Compile(pattern); err == nil {
				opts.Mask.CompiledAllowList = append(opts.Mask.CompiledAllowList, compiled)
			}
		}

		opts.Mask.CompiledBlockList = make([]*regexp.Regexp, 0, len(blockList))
		for _, pattern := range blockList {
			if compiled, err := regexp.Compile(pattern); err == nil {
				opts.Mask.CompiledBlockList = append(opts.Mask.CompiledBlockList, compiled)
			}
		}
	}
}

// isOperationAllowed checks if a query or mutation operation is allowed based on the mask configuration
func (opts *MCPGraphQLServerOptions) isOperationAllowed(operationName string) bool {
	if opts.Mask == nil {
		return true // No masking configured, allow all
	}

	// Check block list first - if it matches, deny regardless of allow list
	for _, pattern := range opts.Mask.CompiledBlockList {
		if pattern.MatchString(operationName) {
			return false
		}
	}

	// If allow list is empty, allow all (except those in block list)
	if len(opts.Mask.CompiledAllowList) == 0 {
		return true
	}

	// Check allow list - must match at least one pattern
	for _, pattern := range opts.Mask.CompiledAllowList {
		if pattern.MatchString(operationName) {
			return true
		}
	}

	return false
}

// WithPassthruHeaders configures which headers to pass through from MCP requests to GraphQL requests
func WithPassthruHeaders(headers []string) MCPGraphQLServerOption {
	return func(opts *MCPGraphQLServerOptions) {
		opts.PassthruHeaders = headers
	}
}

// NewMCPGraphQLServerOptions creates a new options struct with default values
func NewMCPGraphQLServerOptions() *MCPGraphQLServerOptions {
	return &MCPGraphQLServerOptions{
		Logger:          logr.Discard(), // Will use default logger if not set
		Mask:            nil,            // No masking by default
		PassthruHeaders: nil,            // No passthru headers by default
	}
}
