/*
 license x
*/

package repository

import (
	"database/sql"
)

// MockDB is a mock implementation of the database connection for testing.
type MockDB struct {
	BeginFunc    func() (*sql.Tx, error)
	PrepareFunc  func(query string) (*sql.Stmt, error)
	ExecFunc     func(query string, args ...interface{}) (sql.Result, error)
	QueryFunc    func(query string, args ...interface{}) (*sql.Rows, error)
	QueryRowFunc func(query string, args ...interface{}) *sql.Row
	CloseFunc    func() error
	CommitFunc   func() error
	RollbackFunc func() error
}

func (m *MockDB) Prepare(query string) (*sql.Stmt, error) {
	if m.PrepareFunc != nil {
		return m.PrepareFunc(query)
	}
	return nil, nil
}

func (m *MockDB) Query(query string, args ...interface{}) (*sql.Rows, error) {
	if m.QueryFunc != nil {
		return m.QueryFunc(query, args...)
	}
	return nil, nil
}

func (m *MockDB) QueryRow(query string, args ...interface{}) *sql.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(query, args...)
	}
	return nil
}

func (m *MockDB) Close() error {
	if m.CloseFunc != nil {
		return m.CloseFunc()
	}
	return nil
}

func (m *MockDB) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(query, args...)
	}
	return nil, nil
}

func (m *MockDB) Begin() (*sql.Tx, error) {
	if m.BeginFunc != nil {
		return m.BeginFunc()
	}
	return nil, nil
}

func (m *MockDB) Commit() error {
	if m.CommitFunc != nil {
		return m.CommitFunc()
	}
	return nil
}

func (m *MockDB) Rollback() error {
	if m.RollbackFunc != nil {
		return m.RollbackFunc()
	}
	return nil
}
