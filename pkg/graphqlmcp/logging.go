package graphqlmcp

import (
	"log/slog"
	"os"
)

// ConfigureLogging sets up structured logging for the MCP GraphQL server
func ConfigureLogging(level slog.Level, jsonOutput bool) *slog.Logger {
	var handler slog.Handler

	if jsonOutput {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	} else {
		handler = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
			Level: level,
		})
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)
	return logger
}

// ConfigureVerboseLogging sets up verbose logging for debugging
func ConfigureVerboseLogging() *slog.Logger {
	return ConfigureLogging(slog.LevelDebug, false)
}

// ConfigureProductionLogging sets up production-appropriate logging
func ConfigureProductionLogging() *slog.Logger {
	return ConfigureLogging(slog.LevelInfo, true)
}
