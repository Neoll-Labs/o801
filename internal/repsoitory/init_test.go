package repsoitory

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"

	_ "github.com/lib/pq"
)

func TestInitDB_PingError(t *testing.T) {
	// given
	env := "DB_URL"
	_ = os.Setenv(env, "postgres://postgres:postgres@127.0.0.1:9779/?sslmode=disable")

	// when
	_, err := InitDB()

	// then
	assert.Error(t, err)
	assert.Equal(t, "dial tcp 127.0.0.1:9779: connect: connection refused", err.Error())

	_ = os.Unsetenv(env)
}

func TestInitDB_InvalidURLError(t *testing.T) {
	// given
	env := "DB_DRIVER"
	_ = os.Setenv(env, "fake:fake")

	// when
	_, err := InitDB()

	// then
	assert.Error(t, err)
	assert.Equal(t, "sql: unknown driver \"fake:fake\" (forgotten import?)", err.Error())

	_ = os.Unsetenv(env)
}
