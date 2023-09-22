package db

import (
	"context"
	"database/sql"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nelsonstr/o801/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateUserPrepareStmtError(t *testing.T) {
	// given
	db := &mockDB{
		PrepareFunc: func(query string) (*sql.Stmt, error) {
			return nil, errors.New("Prepare error")
		},
	}

	// when
	userStorage := &UserStorage{db: db}

	// then
	_, err := userStorage.Create(context.Background(), "new name")
	assert.Error(t, err)
}

func TestCreateUserSuccess(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	expectedSQL := "^INSERT INTO users \\(name\\) VALUES \\(\\$1\\)  RETURNING id$"
	name := "new user"
	mock.ExpectPrepare(expectedSQL).ExpectExec().WithArgs(name).WillReturnResult(sqlmock.NewResult(1, 1))

	//mock.ExpectBegin()

	rs := sqlmock.NewRows([]string{"id", "title"}).FromCSVString("5,hello world")

	mock.ExpectQuery(expectedSQL).
		WithArgs(5).
		WillReturnRows(rs)

	mock.ExpectCommit()

	// when
	userStorage := &UserStorage{db: db}
	user, err := userStorage.Create(context.Background(), name)

	// then
	assert.NoError(t, err)
	assert.Equal(t, int64(2), user.ID)
	assert.Equal(t, "name2", user.Name)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestGetUserSuccess(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(2, "name2")
	mock.ExpectQuery("^select id, name from users where id = \\$1$").
		WithArgs(2).
		WillReturnRows(row)

	// when
	userStorage := &UserStorage{db: db}
	user, err := userStorage.Get(context.Background(), 2)

	// then
	assert.NoError(t, err)
	assert.Equal(t, int64(2), user.ID)
	assert.Equal(t, "name2", user.Name)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestGetUserNotFound(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("^select id, name from users where id = \\$1$").
		WithArgs(2).
		WillReturnRows(row).WillReturnError(errors.New("sql: no rows in result set"))

	// when
	userStorage := &UserStorage{db: db}
	user, err := userStorage.Get(context.Background(), 2)

	// then
	assert.NoError(t, err)
	assert.Equal(t, &models.NilUser, user)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestGetUserError(t *testing.T) {
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectQuery("^select id, name from users where id = \\$1$").
		WithArgs(2).WillReturnError(errors.New("sql error"))

	// when
	userStorage := &UserStorage{db: db}
	user, err := userStorage.Get(context.Background(), 2)

	// then
	assert.Error(t, err)
	assert.Equal(t, &models.NilUser, user)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}
