package db

import (
	"database/sql"
	"errors"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

var beginError = errors.New("begin error")

func TestExecuteTablesScriptsDBBeginError(t *testing.T) {
	db := &mockDB{BeginFunc: func() (*sql.Tx, error) {
		return nil, beginError
	}}

	m := &migrate{db: db}

	err := m.executeTablesScripts()

	assert.Equal(t, beginError, err)
}

func TestMigrateDB(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("^CREATE TABLE IF NOT EXISTS users \\(ID serial primary key , Name varchar \\);$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	MigrateDB(db)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestMigrateSteps(t *testing.T) {
	type dummystrut struct {
		ID   int    `sql:"type:serial,primary key"`
		Name string `sql:"type:varchar(255)"`
	}

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectBegin()
	mock.ExpectExec("^CREATE TABLE IF NOT EXISTS dummystruts \\(ID serial primary key , Name varchar\\(255\\) \\);$").
		WillReturnResult(sqlmock.NewResult(1, 1))

	mock.ExpectCommit()

	m := &migrate{db: db}
	m.createTableScript(dummystrut{})
	err = m.executeTablesScripts()
	if err != nil {
		t.Fatalf("%s", err)
	}

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestExecuteTablesScriptsCreateTableError(t *testing.T) {
	db, _, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	ddl := "CREATE TABLE - with invalid sql"

	m := &migrate{db: db}
	m.queries = []string{ddl}

	err = m.executeTablesScripts()
	assert.NotNilf(t, err, "expected an error, got %v", err)
}

func TestProcessStruct(t *testing.T) {
	// dummy test model
	type dummytable struct {
		ID        int    `sql:"type:serial,primary key"`
		Name      string `sql:"type:varchar(255)"`
		EmptyType string `sql:"type:"`
		NoType    string
	}

	// Call processStruct with the test model
	result := processStruct(dummytable{})

	// Add assertions or checks based on the expected behavior of processStruct.
	// For example, you can check if the result contains the expected table name and columns.
	assert.Equal(t, result.Name, "dummytables")
	assert.Equal(t, len(result.Columns), 3)
	assert.Equal(t, len(result.Values), 3)
	assert.Equal(t, result.Columns[0].Name, "ID")
	assert.Equal(t, result.Columns[0].Types, "serial primary key")
	assert.Equal(t, result.Columns[1].Name, "Name")
	assert.Equal(t, result.Columns[1].Types, "varchar(255)")
}

func TestGetTableName(t *testing.T) {
	type TestStruct struct{}
	testCases := []struct {
		name     string
		input    reflect.Type
		expected string
	}{
		{
			name:     "test TestStruct",
			input:    reflect.TypeOf(TestStruct{}),
			expected: "teststructs",
		},
		{
			name:     "test nil input",
			input:    nil,
			expected: "",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := getTableName(tc.input)
			if result != tc.expected {
				t.Errorf("Expected %s, but got %s", tc.expected, result)
			}
		})
	}
}
