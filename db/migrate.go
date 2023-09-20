/*
 license x
*/

package db

import (
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
)

const dbURLDefault = "postgres://postgres:postgres@127.0.0.1:5432/?sslmode=disable"

func MigrateDB() {
	log.Printf("start migration.")

	m, err := migrate.New("file://db/migrations", DbURL())
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal(err)
	}

	log.Printf("end migration.")
}

func DbURL() string {
	if v, exists := os.LookupEnv("DBURL"); exists {
		return v
	}

	return dbURLDefault
}
