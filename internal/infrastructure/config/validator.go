package config

import (
	"errors"
	"strings"
)

// Validate checks that the configuration is valid and returns an error if it is not.
// This function checks that all configuration options are valid and returns
// an error if any of them are not.
func (c *Config) Validate() error {
	// Check that the server port is valid
	if c.Server.Port <= 0 || c.Server.Port > 65535 {
		return errors.New("server.port must be a number between 1 and 65535")
	}
	// Check that the database DSN is set
	if c.Database.DSN == "" {
		return errors.New("database.dsn is required")
	}
	// Check that the logging level is valid
	switch c.Logging.Level {
	case "debug", "info", "warn", "error":
		// do nothing
	default:
		return errors.New("logging.level must be one of debug, info, warn, error")
	}
	// Check that the logging format is valid
	switch c.Logging.Format {
	case "json", "console":
		// do nothing
	default:
		return errors.New("logging.format must be one of json, console")
	}
	// Check that the logging output is valid
	switch c.Logging.Output {
	case "stdout", "stderr":
		// do nothing
	default:
		return errors.New("logging.output must be one of stdout, stderr")
	}

	// Check that the keys are set
	c.Keys.privateKey = strings.TrimSpace(c.Keys.privateKey)
	if c.Keys.privateKey == "" {
		return errors.New("keys.private_key is required")
	}

	c.Keys.publicKey = strings.TrimSpace(c.Keys.publicKey)
	if c.Keys.publicKey == "" {
		return errors.New("keys.public_key is required")
	}

	return nil
}
