/*
 license x
*/

package config

import "os"

const (
	dbURLDefault    = "postgres://postgres:postgres@127.0.0.1:5432/?sslmode=disable"
	dbDriverDefault = "postgres"
)

func DBURL() string {
	if v, exists := os.LookupEnv("DB_URL"); exists {
		return v
	}

	return dbURLDefault
}

func DBDriver() string {
	if v, exists := os.LookupEnv("DB_DRIVER"); exists {
		return v
	}

	return dbDriverDefault
}
