/*
 license x
*/

package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/nelsonstr/o801/internal"
	"github.com/nelsonstr/o801/models"
)

// Interface assertions.
var (
	_ http.HandlerFunc = (*UserHandlerAPI)(nil).Get
	_ http.HandlerFunc = (*UserHandlerAPI)(nil).Create
)

func NewUserServer(repo internal.Repository[*models.User]) internal.Handlers {
	return &UserHandlerAPI{
		Mutex:        sync.Mutex{},
		UserCache:    make(map[int64]models.User),
		Repository:   repo,
		ErrorHandler: internal.DefaultErrorHandler,
	}
}

type UserHandlerAPI struct {
	sync.Mutex
	UserCache    map[int64]models.User
	Repository   internal.Repository[*models.User]
	ErrorHandler internal.ErrorHandler
}

// Get the user from the storage.
func (s *UserHandlerAPI) Get(w http.ResponseWriter, r *http.Request) {
	ps := r.Context().Value("params").([]string)

	if len(ps) < 2 {
		s.ErrorHandler(w, r, &internal.BadRequestError{})
		return
	}

	id, err := strconv.Atoi(ps[1]) // ID
	if err != nil {
		s.ErrorHandler(w, r, &internal.BadRequestError{})
		return
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if u := s.UserCache[int64(id)]; u != models.NilUser {
		b, _ := json.Marshal(u)
		_, _ = w.Write(b)

		return
	}

	user, err := s.Repository.Get(r.Context(), id)
	if err != nil {
		s.ErrorHandler(w, r, &internal.StorageError{Err: err})
		return
	}
	if user == &models.NilUser {
		s.ErrorHandler(w, r, &internal.NotFoundError{})
		return
	}

	s.UserCache[user.ID] = *user

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)
}

// Create user.
func (s *UserHandlerAPI) Create(w http.ResponseWriter, r *http.Request) {
	createUserReq := struct {
		Name string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		s.ErrorHandler(w, r, &internal.BadRequestError{})
		return
	}

	user, err := s.Repository.Create(r.Context(), createUserReq.Name)
	if err != nil {
		s.ErrorHandler(w, r, &internal.StorageError{Err: err})
		return
	}

	s.UserCache[user.ID] = *user

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)

	return
}
