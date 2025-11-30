package database

import (
	"mytro-backend-auth/internal/infrastructure/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres returns a new GORM database object connected to a Postgres database.
// The databaseURL parameter should be a valid PostgreSQL connection string.
// If an error occurs during the connection, it is returned as the second result.
func NewPostgres(config config.DatabaseConfig) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(config.MaxOpenConns)
	sqlDB.SetMaxIdleConns(config.MaxIdleConns)
	sqlDB.SetConnMaxLifetime(config.ConnMaxLifetime)

	return db, nil
}
