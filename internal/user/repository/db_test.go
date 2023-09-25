/*
 license x
*/

package repository

import (
	"context"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	userModel "github.com/nelsonstr/o801/internal/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateUser(t *testing.T) {
	t.Parallel()
	db := &MockDB{}
	repo := NewUserRepository(db)

	assert.Equal(t, db, repo.db)
}

func TestCreateUserSuccess(t *testing.T) {
	t.Parallel()
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}

	defer func() { _ = db.Close() }()

	mock.ExpectBegin()

	expectedSQL := "^INSERT INTO users \\(name\\) VALUES \\(\\$1\\)  RETURNING id,name$"
	name := "new user"
	mock.ExpectPrepare(expectedSQL)

	mock.ExpectQuery(expectedSQL).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, name))

	mock.ExpectCommit()

	// when
	userStorage := &UserRepository{db: db}
	user, err := userStorage.Create(context.Background(), &userModel.User{Name: name})

	// then
	assert.NoError(t, err)
	assert.Equal(t, int64(1), user.ID)
	assert.Equal(t, name, user.Name)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestCreateUserBeginError(t *testing.T) {
	t.Parallel()
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}

	defer func() { _ = db.Close() }()

	mock.ExpectBegin().WillReturnError(errors.New("error"))

	// when
	userStorage := &UserRepository{db: db}
	user, err := userStorage.Create(context.Background(), &userModel.User{Name: "name"})

	// then
	assert.Equal(t, errors.Unwrap(err), errors.New("error"))
	assert.Equal(t, user, &userModel.NilUser)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestCreateUserPrepareError(t *testing.T) {
	t.Parallel()
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}

	defer func() { _ = db.Close() }()

	mock.ExpectBegin()

	expectedSQL := "^INSERT INTO users \\(name\\) VALUES \\(\\$1\\)  RETURNING id,name$"
	name := "new user"
	mock.ExpectPrepare(expectedSQL).WillReturnError(errors.New("error"))

	mock.ExpectRollback()

	// when
	userStorage := &UserRepository{db: db}
	user, err := userStorage.Create(context.Background(), &userModel.User{Name: name})

	// then
	assert.Equal(t, errors.Unwrap(err), errors.New("error"))
	assert.Equal(t, user, &userModel.NilUser)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestCreateUserErrorInsertError(t *testing.T) {
	t.Parallel()
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}

	defer func() { _ = db.Close() }()

	mock.ExpectBegin()

	expectedSQL := "^INSERT INTO users \\(name\\) VALUES \\(\\$1\\)  RETURNING id,name$"
	name := "new user"
	mock.ExpectPrepare(expectedSQL)

	mock.ExpectQuery(expectedSQL).WillReturnError(errors.New("error"))

	mock.ExpectRollback()

	// when
	userStorage := &UserRepository{db: db}
	user, err := userStorage.Create(context.Background(), &userModel.User{Name: name})

	// then
	assert.Equal(t, errors.Unwrap(err), errors.New("error"))
	assert.Equal(t, user, &userModel.NilUser)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestCreateUsesCommitError(t *testing.T) {
	t.Parallel()
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}

	defer func() { _ = db.Close() }()

	mock.ExpectBegin()

	expectedSQL := "^INSERT INTO users \\(name\\) VALUES \\(\\$1\\)  RETURNING id,name$"
	name := "new user"
	mock.ExpectPrepare(expectedSQL)

	mock.ExpectQuery(expectedSQL).WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, name))

	mock.ExpectCommit().WillReturnError(errors.New("error"))

	// when
	userStorage := &UserRepository{db: db}
	user, err := userStorage.Create(context.Background(), &userModel.User{Name: name})

	// then
	assert.Error(t, err)
	assert.Equal(t, errors.Unwrap(err), errors.New("error"))
	assert.Equal(t, user, &userModel.NilUser)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestGetUserSuccess(t *testing.T) {
	t.Parallel()
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}

	defer func() { _ = db.Close() }()

	row := sqlmock.NewRows([]string{"id", "name"}).
		AddRow(2, "name2")
	mock.ExpectQuery("^select id, name from users where id = \\$1$").
		WithArgs(2).
		WillReturnRows(row)

	// when
	userStorage := &UserRepository{db: db}
	user, err := userStorage.Get(context.Background(), &userModel.User{ID: 2})

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
	t.Parallel()
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}

	defer func() { _ = db.Close() }()

	row := sqlmock.NewRows([]string{"id", "name"})
	mock.ExpectQuery("^select id, name from users where id = \\$1$").
		WithArgs(2).
		WillReturnRows(row).WillReturnError(errors.New("sql: no rows in result set"))

	// when
	userStorage := &UserRepository{db: db}
	user, err := userStorage.Get(context.Background(), &userModel.User{ID: 2})

	// then
	assert.NoError(t, err)
	assert.Equal(t, &userModel.NilUser, user)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestGetUserError(t *testing.T) {
	t.Parallel()
	// given
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error'%s' was not expected when opening a stub database connection", err)
	}

	defer func() { _ = db.Close() }()

	mock.ExpectQuery("^select id, name from users where id = \\$1$").
		WithArgs(2).WillReturnError(errors.New("sql error"))

	// when
	userStorage := &UserRepository{db: db}
	user, err := userStorage.Get(context.Background(), &userModel.User{ID: 2})

	// then
	assert.Error(t, err)
	assert.Equal(t, &userModel.NilUser, user)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}
