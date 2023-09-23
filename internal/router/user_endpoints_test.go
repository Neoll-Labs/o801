/*
 license x
*/

package router

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/nelsonstr/o801/api"
	"github.com/nelsonstr/o801/internal"
	"github.com/nelsonstr/o801/models"
	"github.com/stretchr/testify/assert"
)

func TestRouter_UserEndpoints_GetUser(t *testing.T) {
	t.Parallel()
	router := NewRouter()
	server := NewUserFakeServer(&UserFakeRepository{})
	router.UserEndpoints(server)
	req := httptest.NewRequest(http.MethodGet, "/123", http.NoBody) // Assuming /{id} is the Get route
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, http.MethodGet, rr.Body.String())
}

func TestRouter_UserEndpoints_CreateUserEmtpyBody(t *testing.T) {
	t.Parallel()
	router := NewRouter()
	server := NewUserFakeServer(&UserFakeRepository{})
	router.UserEndpoints(server)

	req := httptest.NewRequest(http.MethodPost, "/", http.NoBody)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Created", rr.Body.String())
}

func NewUserFakeServer(repo api.Repository[*models.User]) api.HandlerFuncAPI {
	return &FakeHandlerAPI{
		Mutex:        sync.Mutex{},
		UserCache:    make(map[int64]models.User),
		Repository:   repo,
		ErrorHandler: internal.DefaultErrorHandler,
	}
}

type FakeHandlerAPI struct {
	sync.Mutex
	UserCache    map[int64]models.User
	Repository   api.Repository[*models.User]
	ErrorHandler internal.ErrorHandler
}

func (f *FakeHandlerAPI) Create(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Created"))
}

func (f *FakeHandlerAPI) Get(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte(http.MethodGet))
}

var _ api.Repository[*models.User] = (*UserFakeRepository)(nil)

type UserFakeRepository struct{}

func (f *UserFakeRepository) Create(_ context.Context, _ string) (*models.User, error) {
	panic("ignore me")
}

func (f *UserFakeRepository) Get(_ context.Context, _ int) (*models.User, error) {
	panic("ignore me")
}
