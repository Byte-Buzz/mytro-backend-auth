package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres - creates a new database connection
func NewPostgres(databaseURL string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(databaseURL))
	if err != nil {
		return nil, err
	}
	
	return db, nil
}
