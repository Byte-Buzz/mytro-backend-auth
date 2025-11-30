package main

import (
	"mytro-backend-auth/internal/infrastructure/app"
	"mytro-backend-auth/internal/infrastructure/config"
	"mytro-backend-auth/internal/infrastructure/database"
	"mytro-backend-auth/internal/infrastructure/logger"
)

// main is the entry point for the application.
func main() {
	// Load the configuration from environment variables
	config, err := config.LoadFromEnv()
	if err != nil {
		// If there is an error loading the configuration, panic with the error
		panic(err)
	}

	// Create a new logger based on the configuration
	logger, err := logger.NewLogger(config.Logging)
	if err != nil {
		// If there is an error creating the logger, panic with the error
		panic(err)
	}

	// Create a new database connection based on the configuration
	db, err := database.NewPostgres(config.Database)
	if err != nil {
		// If there is an error creating the database connection, panic with the error
		panic(err)
	}

	app := app.NewApp(config, db, logger)

}
