package repsoitory

import (
	"database/sql"
	"github.com/nelsonstr/o801/internal/config"
	"time"
)

func InitDB() (*sql.DB, error) {
	dbc, err := sql.Open(config.DbDriver(), config.DbURL())
	if err != nil {
		return nil, err
	}

	dbc.SetMaxIdleConns(0)
	dbc.SetMaxOpenConns(100)
	dbc.SetConnMaxLifetime(60 * time.Second)

	if err := dbc.Ping(); err != nil {
		return nil, err
	}

	return dbc, nil
}
