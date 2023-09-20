package config

import "os"

const dbURLDefault = "postgres://postgres:postgres@127.0.0.1:5432/?sslmode=disable"

func DbURL() string {
	if v, exists := os.LookupEnv("DBURL"); exists {
		return v
	}

	return dbURLDefault
}
