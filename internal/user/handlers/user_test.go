/*
 license x
*/

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/nelsonstr/o801/internal"
	userModel "github.com/nelsonstr/o801/internal/model"
	"github.com/nelsonstr/o801/internal/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var mockUser = &userModel.UserView{ID: 14, Name: "nelson"}

func TestUserHandlerAPI_Get_Success(t *testing.T) {
	t.Parallel()
	handler := NewUserHandler(&MockService{
		user: mockUser,
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{"/api/user", "1"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseUser userModel.UserView
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)

	assert.NoError(t, err)
	assert.Equal(t, mockUser, &responseUser)
}

func TestUserHandlerAPI_Get_InvalidParameters(t *testing.T) {
	t.Parallel()
	handler := NewUserHandler(&MockService{
		user: mockUser,
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandlerAPI_Get_InvalidParameterKind(t *testing.T) {
	t.Parallel()
	handler := NewUserHandler(&MockService{
		user: mockUser,
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, "")
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandlerAPI_Get_InvalidID(t *testing.T) {
	t.Parallel()
	handler := NewUserHandler(&MockService{
		user: mockUser,
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{"/api/user", "a"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandlerAPI_Get_FromCache(t *testing.T) {
	t.Parallel()
	handler := NewUserHandler(&MockService{
		user: mockUser,
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{"/api/user", "14"})
	handler.Get(rr, req.WithContext(ctx))

	rr = httptest.NewRecorder()
	handler.Get(rr, req.WithContext(ctx))
	assert.Equal(t, http.StatusOK, rr.Code)

	var responseUser userModel.UserView
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)

	assert.NoError(t, err)
	assert.Equal(t, mockUser, &responseUser)
}

func TestUserHandlerAPI_Get_DBError(t *testing.T) {
	t.Parallel()
	handler := NewUserHandler(&MockService{
		error: &internal.StorageError{},
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{"/api/user", "14"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadGateway, rr.Code)
}

func TestUserHandlerAPI_Get_NotFound(t *testing.T) {
	t.Parallel()
	handler := NewUserHandler(&MockService{
		user:  &userModel.NilUserView,
		error: &internal.NotFoundError{},
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{"/api/user", "14"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUserHandlerAPI_Create(t *testing.T) {
	t.Parallel()
	handler := NewUserHandler(&MockService{
		user: mockUser,
	})

	createUserReq := struct {
		Name string
	}{Name: "nelson"}

	reqBody, err := json.Marshal(createUserReq)
	assert.NoError(t, err)

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(reqBody))

	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseUser userModel.UserView
	err = json.Unmarshal(rr.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, &responseUser)
}

func TestUserHandlerAPI_Create_BadRequest(t *testing.T) {
	t.Parallel()
	handler := NewUserHandler(&MockService{})

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer([]byte("invalid")))

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandlerAPI_Create_Invalid(t *testing.T) {
	t.Parallel()
	handler := NewUserHandler(&MockService{
		error: &internal.StorageError{},
	})

	createUserReq := struct {
		Name string
	}{Name: "Alice"}

	reqBody, err := json.Marshal(createUserReq)
	assert.NoError(t, err)
	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer(reqBody))

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadGateway, rr.Code)
}

// MockService mocks to run unit test.
type MockService struct {
	user  *userModel.UserView
	error error
}

func (s *MockService) Create(_ context.Context, _ *userModel.UserView) (*userModel.UserView, error) {
	return s.user, s.error
}

func (s *MockService) Get(_ context.Context, _ *userModel.UserView) (*userModel.UserView, error) {
	return s.user, s.error
}
