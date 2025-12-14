package main

import (
	"context"
	"mytro-backend-auth/internal/app"
	"mytro-backend-auth/internal/infrastructure/config"
	"mytro-backend-auth/internal/infrastructure/database"
	"mytro-backend-auth/internal/infrastructure/logger"
	httpTransport "mytro-backend-auth/internal/transport/http"
	"syscall"
	"time"

	"net/http"
	"os"
	"os/signal"
	"strconv"

	"go.uber.org/zap"
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

	// Create a new application instance with the configuration, database, and logger
	app := app.NewApp(config, db, logger)

	// Create a new HTTP router
	router := httpTransport.NewRouter(app)

	// Create a new HTTP server
	srv := http.Server{
		Addr:         config.Server.Host + ":" + strconv.Itoa(config.Server.Port),
		Handler:      router,
		ReadTimeout:  config.Server.ReadTimeout,
		WriteTimeout: config.Server.WriteTimeout,
	}

	// Start the HTTP server
	go func() {
		logger.Info("Starting server", zap.String("address", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Failed to start server", zap.Error(err))
		}
	}()

	// Wait for interrupt signal to gracefully shut down the server
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stop

	logger.Info("Shutting down server...")

	// Shutdown the server gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown", zap.Error(err))
	} else {
		logger.Info("Server exited properly")
	}
}
