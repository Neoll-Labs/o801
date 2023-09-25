/*
 license x
*/

package routes

import (
	"bytes"
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nelsonstr/o801/internal/interfaces"
	"github.com/nelsonstr/o801/internal/model"
	router "github.com/nelsonstr/o801/internal/router"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/nelsonstr/o801/internal"
	"github.com/stretchr/testify/assert"
)

func TestRouter_UserEndpoints_GetUser(t *testing.T) {
	t.Parallel()
	r := router.NewRouter()
	server := NewMockService(&MockRepository{})
	UserEndpoints(r, server)
	req := httptest.NewRequest(http.MethodGet, "/123", http.NoBody) // Assuming /{id} is the Read route
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, http.MethodGet, rr.Body.String())
}

func TestRouter_UserEndpoints_CreateUserEmtpyBody(t *testing.T) {
	t.Parallel()
	r := router.NewRouter()
	server := NewMockService(&MockRepository{})
	UserEndpoints(r, server)

	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Created", rr.Body.String())
}

func TestInitUserRoutes_Get(t *testing.T) {
	t.Parallel()
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

	r := router.NewRouter()
	InitUserRoutes(db, r)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest(http.MethodGet, "/users/2", http.NoBody))

	assert.Equal(t, http.StatusOK, rr.Code)

	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

func TestInitUserRoutes_Post(t *testing.T) {
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

	r := router.NewRouter()
	InitUserRoutes(db, r)

	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, httptest.NewRequest(http.MethodPost, "/users/", bytes.NewBuffer([]byte("{\"name\": \"nelson\"}"))))

	assert.Equal(t, http.StatusOK, rr.Code)
	err = mock.ExpectationsWereMet()
	if err != nil {
		t.Fatalf("%s", err)
	}
}

// MOCKS

func NewMockService(repo interfaces.Repository[*model.User]) *MockHandlerAPI {
	return &MockHandlerAPI{
		Mutex:        sync.Mutex{},
		UserCache:    make(map[int64]model.User),
		Repository:   repo,
		ErrorHandler: internal.DefaultErrorHandler,
	}
}

type MockHandlerAPI struct {
	sync.Mutex
	UserCache    map[int64]model.User
	Repository   interfaces.Repository[*model.User]
	ErrorHandler internal.ErrorHandler
}

func (f *MockHandlerAPI) Create(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Created"))
}

func (f *MockHandlerAPI) Get(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(http.MethodGet))
}

// mocks to run unit test
type MockRepository struct {
	user  *model.User
	error error
}

func (f *MockRepository) Create(_ context.Context, _ *model.User) (*model.User, error) {
	return f.user, f.error
}

func (f *MockRepository) Get(_ context.Context, _ *model.User) (*model.User, error) {
	return f.user, f.error
}
