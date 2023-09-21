package db

import (
	"database/sql"
	"github.com/nelsonstr/o801/config"
	"log"
)

func InitDB() *sql.DB {

	dbc, err := sql.Open("postgres", config.DbURL())

	if err != nil {
		log.Fatalf("database error: %v", err)
	}

	MigrateDB(dbc)

	return dbc
}
