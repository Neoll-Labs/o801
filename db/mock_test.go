package db

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

/////////////////////////////////////////////////////

// mockTx is a mock implementation of sql.Tx for testing.
type mockTx struct {
	ExecFunc     func(query string, args ...interface{}) (sql.Result, error)
	CommitFunc   func() error
	RollbackFunc func() error
}

func (m *mockTx) Exec(query string, args ...interface{}) (sql.Result, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(query, args...)
	}
	return nil, nil
}

func (m *mockTx) Commit() error {
	if m.CommitFunc != nil {
		return m.CommitFunc()
	}
	return nil
}

func (m *mockTx) Rolback() error {
	if m.RollbackFunc != nil {
		return m.RollbackFunc()
	}
	return nil
}

/////////////////////////////////////

// mockStmt is a mock implementation of sql.Stmt for testing.
type mockStmt struct {
	QueryRowFunc func(args ...interface{}) *sql.Row
	ExecFunc     func(args ...interface{}) (sql.Result, error)
	ScanFunc     func(dest ...interface{}) error
}

// MockStmtWrapper is a wrapper type that embeds both *mockStmt and *sql.Stmt.
type MockStmtWrapper struct {
	*mockStmt
	*sql.Stmt
}

func (m *MockStmtWrapper) QueryRow(args ...interface{}) *sql.Row {
	if m.QueryRowFunc != nil {
		return m.QueryRowFunc(args...)
	}
	return &sql.Row{} // Return an empty row by default
}

func (m *MockStmtWrapper) Exec(args ...interface{}) (sql.Result, error) {
	if m.ExecFunc != nil {
		return m.ExecFunc(args...)
	}
	return nil, nil // Return nil result and error by default
}

func (m *MockStmtWrapper) Scan(dest ...interface{}) error {
	if m.ScanFunc != nil {
		return m.ScanFunc(dest...)
	}
	return nil // Return no error by default
}
