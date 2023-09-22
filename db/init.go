package db

import (
	"database/sql"
	"github.com/nelsonstr/o801/config"
)

func InitDB() (*sql.DB, error) {
	dbc, err := sql.Open(config.DbDriver(), config.DbURL())
	if err != nil {
		return nil, err
	}

	if err := dbc.Ping(); err != nil {
		return nil, err
	}

	return dbc, nil
}
