package router

import (
	"context"
	"github.com/nelsonstr/o801/internal"
	"github.com/nelsonstr/o801/internal/repsoitory"
	"github.com/nelsonstr/o801/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
)

func TestRouter_UserEndpoints_GetUser(t *testing.T) {
	// Create a new router
	router := NewRouter()

	// Create a test server
	server := NewUserFakeServer(&UserFakeRepository{})
	router.UserEndpoints(server)

	// Define a test request for Get
	req := httptest.NewRequest("GET", "/123", nil) // Assuming /{id} is the Get route

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Check the response status code (assert as needed)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "Get", rr.Body.String())
}

func TestRouter_UserEndpoints_CreateUserEmtpyBody(t *testing.T) {
	// Create a new router
	router := NewRouter()

	// Create a test server
	server := NewUserFakeServer(&UserFakeRepository{})
	router.UserEndpoints(server)

	// Define a test request for Create (POST)
	req := httptest.NewRequest("POST", "/", nil)

	// Create a response recorder to capture the response
	rr := httptest.NewRecorder()

	// Serve the request using the router
	router.ServeHTTP(rr, req)

	// Check the response status code (assert as needed)
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
	w.Write([]byte("Created"))
}

func (f *FakeHandlerAPI) Get(w http.ResponseWriter, _ *http.Request) {
	w.Write([]byte("Get"))
}

var _ internal.Repository[*models.User] = (*UserFakeRepository)(nil)

func NewUserStorage(db repsoitory.DBInterface) *UserFakeRepository {
	return &UserFakeRepository{}
}

type UserFakeRepository struct {
}

func (f *UserFakeRepository) Create(ctx context.Context, s string) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}

func (f *UserFakeRepository) Get(ctx context.Context, i int) (*models.User, error) {
	//TODO implement me
	panic("implement me")
}
