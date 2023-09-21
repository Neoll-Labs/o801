/*
 license x
*/

package server

import (
	"encoding/json"
	"net/http"
	"strconv"
	"sync"

	"github.com/nelsonstr/o801/db"
	"github.com/nelsonstr/o801/internal"
	"github.com/nelsonstr/o801/models"
)

// Interface assertions.
var (
	_ http.HandlerFunc = (*Server)(nil).GetUser
	_ http.HandlerFunc = (*Server)(nil).CreateUser
)

func NewServer(service db.CRService[*models.User]) *Server {
	return &Server{
		Mutex:     sync.Mutex{},
		userCache: make(map[int64]models.User),
		service:   service,
	}
}

type Server struct {
	sync.Mutex
	userCache    map[int64]models.User
	service      db.CRService[*models.User]
	errorHandler internal.ErrorHandler
}

// GetUser gets the user from the storage.
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	ps := r.Context().Value("params").([]string)

	if len(ps) < 2 {
		s.errorHandler(w, r, &internal.BadRequestError{}, nil)
		return
	}
	id, err := strconv.Atoi(ps[1]) // ID
	if err != nil {
		s.errorHandler(w, r, &internal.BadRequestError{}, nil)
		return
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if u := s.userCache[int64(id)]; u != models.NilUser {
		b, _ := json.Marshal(u)
		_, _ = w.Write(b)

		return
	}

	user, err := s.service.Get(r.Context(), id)
	if err != nil {
		s.errorHandler(w, r, &internal.StorageError{Err: err}, nil)
		return
	}

	s.userCache[user.ID] = *user

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)
}

// CreateUser create user
func (s *Server) CreateUser(w http.ResponseWriter, r *http.Request) {
	createUserReq := struct {
		Name string
	}{}

	if err := json.NewDecoder(r.Body).Decode(&createUserReq); err != nil {
		s.errorHandler(w, r, &internal.BadRequestError{}, nil)
		return
	}

	user, err := s.service.Create(r.Context(), createUserReq.Name)
	if err != nil {
		s.errorHandler(w, r, &internal.StorageError{Err: err}, nil)
		return
	}

	s.userCache[user.ID] = *user

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)

	return
}
