/*
 license x
*/

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/nelsonstr/o801/internal"
	"github.com/nelsonstr/o801/internal/interfaces"
	userModel "github.com/nelsonstr/o801/internal/model"
	"github.com/nelsonstr/o801/internal/router"
)

// Interface assertions.
var (
	_ http.HandlerFunc = (*UserHandler)(nil).Get
	_ http.HandlerFunc = (*UserHandler)(nil).Create
)

type UserHandler struct {
	service      interfaces.ServiceAPI[*userModel.UserView]
	ErrorHandler internal.ErrorHandler
}

func NewUserHandler(service interfaces.ServiceAPI[*userModel.UserView]) *UserHandler {
	return &UserHandler{
		service:      service,
		ErrorHandler: internal.DefaultErrorHandler,
	}
}

// Get the user from the storage.
func (h *UserHandler) Get(w http.ResponseWriter, r *http.Request) {
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

	user, err := h.service.Get(r.Context(), &userModel.UserView{ID: int64(id)})
	if err != nil {
		h.ErrorHandler(w, r, err)

		return
	}

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)
}

// Create user.
func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	createUserReq := struct {
		Name string `json:"name"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		h.ErrorHandler(w, r, &internal.ParsingError{})

		return
	}

	user, err := h.service.Create(r.Context(), &userModel.UserView{Name: createUserReq.Name})
	if err != nil {
		h.ErrorHandler(w, r, &internal.StorageError{})

		return
	}

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)
}
