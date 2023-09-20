package db

import (
	"github.com/golang-migrate/migrate/v4"
	"log"
)

func MigrateDB() {
	log.Printf("start migration")
	m, err := migrate.New(
		"file://db/migrations",
		"postgres://postgres:postgres@localhost:5432/?sslmode=disable") // TODO create configuration
	if err != nil {
		log.Fatal(err)
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
	log.Printf("end migration")
}
