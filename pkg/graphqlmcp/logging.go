package graphqlmcp

import (
	"log/slog"
	"os"

	"github.com/go-logr/logr"
	"github.com/go-logr/logr/slogr"
)

// ConfigureLogging sets up structured logging for the MCP GraphQL server
func ConfigureLogging(level slog.Level, jsonOutput bool) logr.Logger {
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

	slogger := slog.New(handler)
	slog.SetDefault(slogger)
	return slogr.NewLogr(handler)
}

// ConfigureVerboseLogging sets up verbose logging for debugging
func ConfigureVerboseLogging() logr.Logger {
	return ConfigureLogging(slog.LevelDebug, false)
}

// ConfigureProductionLogging sets up production-appropriate logging
func ConfigureProductionLogging() logr.Logger {
	return ConfigureLogging(slog.LevelInfo, true)
}
