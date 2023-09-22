package db

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestInitDBPingError(t *testing.T) {
	// given
	env := "DB_URL"
	os.Setenv(env, "postgres://postgres:postgres@127.0.0.1:9779/?sslmode=disable")

	// when
	_, err := InitDB()

	// then
	assert.Error(t, err)
	assert.Equal(t, "dial tcp 127.0.0.1:9779: connectex: No connection could be made because the target machine actively refused it.", err.Error())

	os.Unsetenv(env)
}

func TestInitDBInvalidURLError(t *testing.T) {
	// given
	env := "DB_DRIVER"
	os.Setenv(env, "fake:fake")

	// when
	_, err := InitDB()

	// then
	assert.Error(t, err)
	assert.Equal(t, "sql: unknown driver \"fake:fake\" (forgotten import?)", err.Error())

	os.Unsetenv(env)
}
