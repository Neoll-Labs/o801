/*
 license x
*/

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/nelsonstr/o801/api"
	"github.com/nelsonstr/o801/internal"
	"github.com/nelsonstr/o801/internal/router"
	"github.com/nelsonstr/o801/models"
)

// Interface assertions.
var (
	_ http.HandlerFunc = (*UserHandlerAPI)(nil).Get
	_ http.HandlerFunc = (*UserHandlerAPI)(nil).Create
)

func NewUserServer(repo api.Repository[*models.User]) *UserHandlerAPI {
	return &UserHandlerAPI{
		mutex:        sync.Mutex{},
		UserCache:    make(map[int64]models.User),
		Repository:   repo,
		ErrorHandler: internal.DefaultErrorHandler,
	}
}

type UserHandlerAPI struct {
	mutex        sync.Mutex
	UserCache    map[int64]models.User
	Repository   api.Repository[*models.User]
	ErrorHandler internal.ErrorHandler
}

// Get the user from the storage.
func (s *UserHandlerAPI) Get(w http.ResponseWriter, r *http.Request) {
	const (
		expectedParameters = 2
		parameterID        = 1
	)

	p := r.Context().Value(router.ParametersName)

	ps, okParamsKind := p.([]string)
	if !okParamsKind {
		s.ErrorHandler(w, r, &internal.ParsingError{})

		return
	}

	if len(ps) < expectedParameters {
		s.ErrorHandler(w, r, &internal.ParsingError{})

		return
	}

	id, err := strconv.Atoi(ps[parameterID]) // ID
	if err != nil {
		s.ErrorHandler(w, r, &internal.ParsingError{})

		return
	}

	user := s.GetFromCached(id)
	if user != &models.NilUser {
		b, _ := json.Marshal(user)
		_, _ = w.Write(b)

		return
	}

	user, err = s.Repository.Get(r.Context(), id)
	if err != nil {
		s.ErrorHandler(w, r, &internal.StorageError{Err: err})

		return
	}

	if user == &models.NilUser {
		s.ErrorHandler(w, r, &internal.NotFoundError{})

		return
	}

	s.AddToCache(user)

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)
}

func (s *UserHandlerAPI) AddToCache(user *models.User) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	s.UserCache[user.ID] = *user
}

func (s *UserHandlerAPI) GetFromCached(id int) *models.User {
	s.mutex.Lock()

	defer s.mutex.Unlock()

	if u, exist := s.UserCache[int64(id)]; exist {
		return &u
	}

	return &models.NilUser
}

// Create user.
func (s *UserHandlerAPI) Create(w http.ResponseWriter, r *http.Request) {
	createUserReq := struct {
		Name string `json:"name"`
	}{}

	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		s.ErrorHandler(w, r, &internal.ParsingError{})

		return
	}

	user, err := s.Repository.Create(r.Context(), createUserReq.Name)
	if err != nil {
		s.ErrorHandler(w, r, &internal.StorageError{Err: err})

		return
	}

	s.AddToCache(user)

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)
}
