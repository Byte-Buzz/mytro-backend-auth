package main

import (
	"flag"
	"log"
	"mytro-backend-auth/internal/infrastructure/config"
	"mytro-backend-auth/internal/infrastructure/database"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main() {
	direction := flag.String("direction", "up", "Migration direction: up or down")
	flag.Parse()

	cfg, err := config.LoadFromEnv()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.NewPostgres(cfg.Database)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}

	defer sqlDB.Close()

	m, err := migrate.New(
		"file://migrations",
		cfg.Database.DSN,
	)
	if err != nil {
		log.Fatal(err)
	}

	switch *direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	default:
		log.Fatal("Invalid direction")
	}

	log.Println("Migration completed")
}
