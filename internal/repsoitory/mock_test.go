package repsoitory

import (
	"database/sql"
)

// mockDB is a mock implementation of the database connection for testing.
type mockDB struct {
	BeginFunc    func() (*sql.Tx, error)
	PrepareFunc  func(query string) (*sql.Stmt, error)
	ExecFunc     func(query string, args ...interface{}) (sql.Result, error)
	QueryFunc    func(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowFunc func(query string, args ...interface{}) *sql.Row
	CloseFunc    func() error
	CommitFunc   func() error
	RollbackFunc func() error
}

func (m *mockDB) Prepare(query string) (*sql.Stmt, error) {
	if m.PrepareFunc != nil {
		return m.PrepareFunc(query)
	}
	return nil, nil
}

func (m *mockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(query, args...)
	}
	return nil, nil
}

func (m *mockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(query, args...)
	}
	return nil
}

func (m *mockDB) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *mockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(query, args...)
	}
	return nil, nil
}

func (m *mockDB) Begin() (*sql.Tx, error) {
	if m.BeginFunc != nil {
		return m.BeginFunc()
	}
	return nil, nil
}

func (m *mockDB) Commit() error {
	if m.CommitFunc != nil {
		return m.CommitFunc()
	}
	return nil
}

func (m *mockDB) Rollback() error {
	if m.RollbackFunc != nil {
		return m.RollbackFunc()
	}
	return nil
}
