package config

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// loadFromEnv loads the configuration from environment variables.
func (c *Config) loadFromEnv() error {
	prefix := "MT_AUTH_"

	// Server
	c.Server.Host = getEnv(prefix+"SERVER_HOST", "0.0.0.0")
	c.Server.Port = getIntEnv(prefix+"SERVER_PORT", 8080)
	c.Server.ReadTimeout = getDurationEnv(prefix+"SERVER_READ_TIMEOUT", 10*time.Second)
	c.Server.WriteTimeout = getDurationEnv(prefix+"SERVER_WRITE_TIMEOUT", 20*time.Second)
	c.Server.RequestTimeout = getDurationEnv(prefix+"SERVER_REQUEST_TIMEOUT", 10*time.Second)

	// Database
	c.Database.DSN = getEnv(prefix+"DATABASE_DSN", "")
	c.Database.MaxOpenConns = getIntEnv(prefix+"DATABASE_MAX_OPEN_CONNS", 15)
	c.Database.MaxIdleConns = getIntEnv(prefix+"DATABASE_MAX_IDLE_CONNS", 5)
	c.Database.ConnMaxLifetime = getDurationEnv(prefix+"DATABASE_CONN_MAX_LIFETIME", time.Hour)

	// Logging
	c.Logging.Level = getEnv(prefix+"LOGGING_LEVEL", "info")
	c.Logging.Format = getEnv(prefix+"LOGGING_FORMAT", "json")
	c.Logging.Output = getEnv(prefix+"LOGGING_OUTPUT", "stdout")
	c.Logging.FilePath = getEnv(prefix+"LOGGING_FILE_PATH", "")

	// CORS
	c.CORS.AllowedOrigins = getArrayEnv(prefix+"CORS_ALLOWED_ORIGINS", []string{"*"})
	c.CORS.AllowedMethods = getArrayEnv(prefix+"CORS_ALLOWED_METHODS", []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"})
	c.CORS.AllowedHeaders = getArrayEnv(prefix+"CORS_ALLOWED_HEADERS", []string{
		"Origin",
		"Content-Type",
		"Authorization",
		"X-Requested-With",
	})
	c.CORS.AllowCredentials = getEnv(prefix+"CORS_ALLOW_CREDENTIALS", "false") == "true"

	return nil
}

// getEnv returns the value of the environment variable with the given key.
// If the variable does not exist, or its value is empty, the defaultValue is returned.
func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists || value == "" {
		return defaultValue
	}
	return value
}

// getIntEnv returns the value of the environment variable with the given key,
// parsed as an integer. If the variable does not exist, or its value
// is empty, the defaultValue is returned. If the value cannot be parsed
// as an integer, the defaultValue is returned.
func getIntEnv(key string, defaultValue int) int {
	valueStr, exists := os.LookupEnv(key)
	if !exists || valueStr == "" {
		return defaultValue
	}
	value, err := strconv.Atoi(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}

// getDurationEnv returns the value of the environment variable with the given key,
// parsed as a time.Duration. If the variable does not exist, or its value
// is empty, the defaultValue is returned. If the value cannot be parsed
// as a time.Duration, the defaultValue is returned.
func getDurationEnv(key string, defaultValue time.Duration) time.Duration {
	valueStr, exists := os.LookupEnv(key)
	if !exists || valueStr == "" {
		return defaultValue
	}
	duration, err := time.ParseDuration(valueStr)
	if err != nil {
		return defaultValue
	}
	return duration
}

// getArrayEnv returns the value of the environment variable with the given key,
// split by commas into a slice of strings. If the variable does not exist,
// or its value is empty, the defaultValue is returned.
func getArrayEnv(key string, defaultValue []string) []string {
	valueStr, exists := os.LookupEnv(key)
	if !exists || valueStr == "" {
		return defaultValue
	}
	var values []string
	for v := range strings.SplitSeq(valueStr, ",") {
		if v == "" {
			continue
		}

		v := strings.TrimLeft(v, " ")

		values = append(values, v)
	}
	return values
}
