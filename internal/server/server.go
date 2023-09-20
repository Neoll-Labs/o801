/*
 license x
*/

package server

import (
	"encoding/json"
	"github.com/nelsonstr/o801/internal"
	"io"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/nelsonstr/o801/db"
	"github.com/nelsonstr/o801/models"
)

// Interface assertions.
var (
	_ http.HandlerFunc = (*Server)(nil).GetUser
	_ http.HandlerFunc = (*Server)(nil).GetOrCreateUser
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

func (s *Server) Routes() internal.Routes {
	return internal.Routes{
		"GetOrCreateUser": internal.Route{
			"/users",
			s.GetOrCreateUser,
		},
		"GetUser": internal.Route{
			"/users/",
			s.GetUser,
		},
	}
}

// GetUser gets the user from the storage.
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	idReq := struct {
		ID int
	}{}

	err := json.NewDecoder(r.Body).Decode(&idReq)
	if err != nil && err != io.EOF {
		s.errorHandler(w, r, &internal.BadRequestError{}, nil)
		return
	}

	// empty body, check if is path parameter
	if err == io.EOF {
		p := strings.Split(r.URL.String(), "/")
		if len(p) < 3 {
			s.errorHandler(w, r, &internal.BadRequestError{}, nil)
			return
		}
		idReq.ID, _ = strconv.Atoi(p[2])
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if u := s.userCache[int64(idReq.ID)]; u != models.NilUser {
		b, _ := json.Marshal(u)
		_, _ = w.Write(b)

		return
	}

	user, err := s.service.Get(r.Context(), idReq.ID)
	if err != nil {
		s.errorHandler(w, r, &internal.StorageError{Err: err}, nil)
		return
	}

	s.userCache[int64(idReq.ID)] = *user

	b, _ := json.Marshal(user)
	_, _ = w.Write(b)
}

// GetOrCreateUser creates or get an user from storage
func (s *Server) GetOrCreateUser(writer http.ResponseWriter, request *http.Request) {
	if request.Method == "POST" {
		s.CreateUser(writer, request)
		return
	} else if request.Method == "GET" {
		s.GetUser(writer, request)
		return
	} else {
		s.errorHandler(writer, request, &internal.MethodNotAllowedError{}, nil)
	}
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
