/*
 license x
*/

package repository

import (
	"database/sql"
	"fmt"
	"time"
)

const (
	maxConnetions   = 100
	maxIdleConns    = 0
	connMaxLifetime = 60 * time.Second
)

func InitDB() (*sql.DB, error) {
	dbc, err := sql.Open(DBDriver(), DBURL())
	if err != nil {
		return nil, fmt.Errorf("failed to open connection: %w", err)
	}

	dbc.SetMaxIdleConns(maxIdleConns)

	dbc.SetMaxOpenConns(maxConnetions)
	dbc.SetConnMaxLifetime(connMaxLifetime)

	if err := dbc.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping db: %w", err)
	}

	return dbc, nil
}
