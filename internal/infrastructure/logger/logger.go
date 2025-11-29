package logger

import (
	"mytro-backend-auth/internal/infrastructure/config"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// NewLogger creates a new zap logger based on the provided configuration.
//
// The logger is configured with the provided log level, format, and
// output paths. If a file path is provided, the logger will write logs
// to that file in addition to the provided output paths.
func NewLogger(config config.LoggingConfig) (*zap.Logger, error) {
	// Create a new zap logger configuration
	cfg := zap.NewProductionConfig()

	// Set the log level
	cfg.Level = zap.NewAtomicLevelAt(getLogLevel(config.Level))

	// Set the log format (json or console)
	cfg.Encoding = config.Format

	// Set the process output
	cfg.OutputPaths = []string{config.Output}

	// If a file path is provided, add it to the output paths
	if config.FilePath != "" {
		cfg.OutputPaths = append(cfg.OutputPaths, config.FilePath)
	}

	// Build the logger
	return cfg.Build()
}

// getLogLevel takes a string representing a log level and returns the
// corresponding zapcore.Level. Recognized log levels are "debug", "info",
// "warn", and "error". If an unrecognized log level is provided, the
// function returns zap.InfoLevel.
func getLogLevel(level string) zapcore.Level {
	switch level {
	case "debug":
		return zap.DebugLevel
	case "warn":
		return zap.WarnLevel
	case "error":
		return zap.ErrorLevel
	default:
		return zap.InfoLevel
	}
}
