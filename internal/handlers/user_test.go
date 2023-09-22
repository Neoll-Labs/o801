package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/nelsonstr/o801/models"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

var mockUser = &models.User{ID: 14, Name: "nelson"}

func TestUserHandlerAPI_Get_Success(t *testing.T) {
	handler := NewUserServer(&UserFakeRepository{
		user: mockUser,
	})
	req := httptest.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "params", []string{"/api/user", "1"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseUser models.User
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)

	assert.NoError(t, err)
	assert.Equal(t, mockUser, &responseUser)
}

func TestUserHandlerAPI_Get_InvalidParameters(t *testing.T) {
	handler := NewUserServer(&UserFakeRepository{
		user: mockUser,
	})
	req := httptest.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "params", []string{})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandlerAPI_Get_InvalidID(t *testing.T) {
	handler := NewUserServer(&UserFakeRepository{
		user: mockUser,
	})
	req := httptest.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "params", []string{"/api/user", "a"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandlerAPI_Get_FromCache(t *testing.T) {
	handler := NewUserServer(&UserFakeRepository{
		user: mockUser,
	})
	req := httptest.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "params", []string{"/api/user", "14"})
	handler.Get(rr, req.WithContext(ctx))

	rr = httptest.NewRecorder()
	handler.Get(rr, req.WithContext(ctx))
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseUser models.User
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)

	assert.NoError(t, err)
	assert.Equal(t, mockUser, &responseUser)
}

func TestUserHandlerAPI_Get_DBError(t *testing.T) {
	handler := NewUserServer(&UserFakeRepository{
		error: errors.New("DB error"),
	})
	req := httptest.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "params", []string{"/api/user", "14"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadGateway, rr.Code)
}

func TestUserHandlerAPI_Get_NotFound(t *testing.T) {
	handler := NewUserServer(&UserFakeRepository{
		user: &models.NilUser,
	})
	req := httptest.NewRequest("GET", "/user/1", nil)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), "params", []string{"/api/user", "14"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

//func TestUserHandlerAPI_Create(t *testing.T) {
//	// Create a test repository with a mock Create method
//	mockRepo := &internal.MockRepository{}
//	mockUser := &models.user{ID: 1, Name: "John"}
//	mockRepo.On("Create").Return(mockUser, nil)
//
//	// Create a new UserHandlerAPI instance with the mock repository
//	handler := NewUserServer(mockRepo)
//
//	// Create a test request for CreateUser (POST)
//	createUserReq := struct {
//		Name string
//	}{Name: "Alice"}
//
//	// Encode the request body
//	reqBody, err := json.Marshal(createUserReq)
//	assert.NoError(t, err)
//
//	// Create a test request with the encoded body
//	req := httptest.NewRequest("POST", "/user", bytes.NewBuffer(reqBody))
//
//	// Create a response recorder to capture the response
//	rr := httptest.NewRecorder()
//
//	// Serve the request using the handler
//	handler.Create(rr, req)
//
//	// Check the response status code (assert as needed)
//	assert.Equal(t, http.StatusOK, rr.Code)
//
//	// Parse the response JSON and check its contents
//	var responseUser models.user
//	err = json.Unmarshal(rr.Body.Bytes(), &responseUser)
//	assert.NoError(t, err)
//	assert.Equal(t, mockUser, &responseUser)
//
//	// Assert that the mock repository's Create method was called
//	//mockRepo.AssertCalled(t, "Create")
//}

// You can add more test cases to cover error scenarios, invalid requests, and edge cases.

type UserFakeRepository struct {
	user  *models.User
	error error
}

func (f *UserFakeRepository) Create(ctx context.Context, s string) (*models.User, error) {
	return f.user, f.error
}

func (f *UserFakeRepository) Get(ctx context.Context, i int) (*models.User, error) {
	return f.user, f.error
}
