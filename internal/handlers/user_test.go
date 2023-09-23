/*
 license x
*/

package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/nelsonstr/o801/internal/router"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nelsonstr/o801/models"
	"github.com/stretchr/testify/assert"
)

var mockUser = &models.User{ID: 14, Name: "nelson"}

func TestUserHandlerAPI_Get_Success(t *testing.T) {
	t.Parallel()
	handler := NewUserServer(&UserFakeRepository{
		user: mockUser,
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{"/api/user", "1"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusOK, rr.Code)

	var responseUser models.User
	err := json.Unmarshal(rr.Body.Bytes(), &responseUser)

	assert.NoError(t, err)
	assert.Equal(t, mockUser, &responseUser)
}

func TestUserHandlerAPI_Get_InvalidParameters(t *testing.T) {
	t.Parallel()
	handler := NewUserServer(&UserFakeRepository{
		user: mockUser,
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandlerAPI_Get_InvalidID(t *testing.T) {
	t.Parallel()
	handler := NewUserServer(&UserFakeRepository{
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
	handler := NewUserServer(&UserFakeRepository{
		user: mockUser,
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{"/api/user", "14"})
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
	t.Parallel()
	handler := NewUserServer(&UserFakeRepository{
		error: errors.New("DB error"),
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{"/api/user", "14"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusBadGateway, rr.Code)
}

func TestUserHandlerAPI_Get_NotFound(t *testing.T) {
	t.Parallel()
	handler := NewUserServer(&UserFakeRepository{
		user: &models.NilUser,
	})
	req := httptest.NewRequest(http.MethodGet, "/user/1", http.NoBody)
	rr := httptest.NewRecorder()

	ctx := context.WithValue(req.Context(), router.ParametersName, []string{"/api/user", "14"})
	handler.Get(rr, req.WithContext(ctx))

	assert.Equal(t, http.StatusNotFound, rr.Code)
}

func TestUserHandlerAPI_Create(t *testing.T) {
	t.Parallel()
	handler := NewUserServer(&UserFakeRepository{
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

	var responseUser models.User
	err = json.Unmarshal(rr.Body.Bytes(), &responseUser)
	assert.NoError(t, err)
	assert.Equal(t, mockUser, &responseUser)
}

func TestUserHandlerAPI_Create_BadRequest(t *testing.T) {
	t.Parallel()
	handler := NewUserServer(&UserFakeRepository{})

	req := httptest.NewRequest(http.MethodPost, "/user", bytes.NewBuffer([]byte("invalid")))

	rr := httptest.NewRecorder()
	handler.Create(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
}

func TestUserHandlerAPI_Create_Invalid(t *testing.T) {
	t.Parallel()
	handler := NewUserServer(&UserFakeRepository{
		error: errors.New("DB error"),
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

// mocks to run unit test

type UserFakeRepository struct {
	user  *models.User
	error error
}

func (f *UserFakeRepository) Create(_ context.Context, _ string) (*models.User, error) {
	return f.user, f.error
}

func (f *UserFakeRepository) Get(_ context.Context, _ int) (*models.User, error) {
	return f.user, f.error
}
