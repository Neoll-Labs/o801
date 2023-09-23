/*
 license x
*/

package router

import (
	"context"
	"github.com/nelsonstr/o801/internal"
	"github.com/nelsonstr/o801/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRouter_UserEndpoints_GetUser(t *testing.T) {
	router := NewRouter()
	server := NewUserFakeServer(&UserFakeRepository{})
	router.UserEndpoints(server)
	req := httptest.NewRequest("GET", "/123", nil) // Assuming /{id} is the Get route
	rr := httptest.NewRecorder()

	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Get", rr.Body.String())
}

func TestRouter_UserEndpoints_CreateUserEmtpyBody(t *testing.T) {
	router := NewRouter()
	server := NewUserFakeServer(&UserFakeRepository{})
	router.UserEndpoints(server)

	req := httptest.NewRequest("POST", "/", nil)
	rr := httptest.NewRecorder()
	router.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Created", rr.Body.String())
}

func NewUserFakeServer(repo internal.Repository[*models.User]) internal.HandlerFuncAPI {
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
	Repository   internal.Repository[*models.User]
	ErrorHandler internal.ErrorHandler
}

func (f *FakeHandlerAPI) Create(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Created"))
}

func (f *FakeHandlerAPI) Get(w http.ResponseWriter, _ *http.Request) {
	_, _ = w.Write([]byte("Get"))
}

var _ internal.Repository[*models.User] = (*UserFakeRepository)(nil)

type UserFakeRepository struct{}

func (f *UserFakeRepository) Create(_ context.Context, _ string) (*models.User, error) {
	panic("ignore me")
}

func (f *UserFakeRepository) Get(_ context.Context, _ int) (*models.User, error) {
	panic("ignore me")
}
