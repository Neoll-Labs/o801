/*
 license x
*/

package repository

import (
	"os"
	"strings"
	"testing"

	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
)

//nolint:paraleltest
func TestInitDB_PingError(t *testing.T) {

	// given
	env := "DB_URL"
	t.Setenv(env, "postgres://postgres:postgres@127.0.0.1:9779/?sslmode=disable")

	// when
	_, err := InitDB()

	// then
	assert.Error(t, err)
	// running in container: "dial tcp 127.0.0.1:9779: connect: connection refused"
	// running from command line: "dial tcp 127.0.0.1:9779: connectex: No connection could be made because the target machine actively refused it."
	assert.True(t, strings.Contains(err.Error(), "dial tcp 127.0.0.1:9779: connect"))
}

//nolint:paraleltest
func TestInitDB_InvalidURLError(t *testing.T) {

	// given
	env := "DB_DRIVER"
	t.Setenv(env, "fake:fake")

	// when
	_, err := InitDB()

	// then
	assert.Error(t, err)
	assert.Equal(t, "failed to open connection: sql: unknown driver \"fake:fake\" (forgotten import?)", err.Error())

	_ = os.Unsetenv(env)
}
