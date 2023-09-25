/*
 license x
*/

package handlers

import (
	"encoding/json"
	"github.com/nelsonstr/o801/api"
	"github.com/nelsonstr/o801/internal"
	"github.com/nelsonstr/o801/internal/router"
	"github.com/nelsonstr/o801/internal/user/service"
	"net/http"
	"strconv"
)

// Interface assertions.
var (
	_ http.HandlerFunc = (*userHandler)(nil).Get
	_ http.HandlerFunc = (*userHandler)(nil).Create
)

type userHandler struct {
	service      api.ServiceAPI[*service.User]
	ErrorHandler internal.ErrorHandler
}

func NewUserHandler(service api.ServiceAPI[*service.User]) *userHandler {
	return &userHandler{
		service:      service,
		ErrorHandler: internal.DefaultErrorHandler,
	}
}

// Get the user from the storage.
func (h *userHandler) Get(w http.ResponseWriter, r *http.Request) {
	const (
		expectedParameters = 2
		parameterID        = 1
	)

	p := r.Context().Value(router.ParametersName)

	ps, okParamsKind := p.([]string)
	if !okParamsKind {
		h.ErrorHandler(w, r, &internal.ParsingError{})

		return
	}

	if len(ps) < expectedParameters {
		h.ErrorHandler(w, r, &internal.ParsingError{})

		return
	}

	id, err := strconv.Atoi(ps[parameterID])
	if err != nil {
		h.ErrorHandler(w, r, &internal.ParsingError{})

		return
	}

	user, err := h.service.Get(r.Context(), &service.User{ID: int64(id)})
	if err != nil {
		h.ErrorHandler(w, r, err)

		return
	}

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)
}

// Create user.
func (h *userHandler) Create(w http.ResponseWriter, r *http.Request) {
	createUserReq := struct {
		Name string `json:"name"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		h.ErrorHandler(w, r, &internal.ParsingError{})

		return
	}

	user, err := h.service.Create(r.Context(), &service.User{Name: createUserReq.Name})
	if err != nil {
		h.ErrorHandler(w, r, &internal.StorageError{})

		return
	}

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)
}
