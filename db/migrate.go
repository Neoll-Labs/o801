/*
 license x
*/

package db

import (
	"errors"
	"github.com/golang-migrate/migrate/v4"
	"github.com/nelsonstr/o801/config"
	"log"
)

func MigrateDB() {
	log.Printf("start migration.")

	m, err := migrate.New("file://db/migrations", config.DbURL())
	if err != nil {
		log.Fatal(err)
	}

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal(err)
	}

	log.Printf("end migration.")
}
