package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"sync"

	"log"

	"net/http"
)

type User struct {
	ID   int64
	Name string
}

var nilUser User

func NewServer(db *sql.DB) Server {
	return Server{db: db}
}

type Server struct {
	sync.Mutex
	db        *sql.DB
	userCache map[int64]User
}

// gets the user from the database
func (s *Server) GetUser(w http.ResponseWriter, r *http.Request) {
	idReq := struct {
		ID int
	}{}

	if err := json.NewDecoder(r.Body).Decode(&idReq); err != nil {
		log.Fatalf("error parsing json id request: %w", err)
		return
	}

	s.Mutex.Lock()
	defer s.Mutex.Unlock()

	if u := s.userCache[int64(idReq.ID)]; u != nilUser {
		b, _ := json.Marshal(u)
		w.Write(b)
	}

	row := s.db.QueryRow("select * from users where id=" + fmt.Sprint(idReq.ID))
	u := &User{}
	err := row.Scan(u.ID, u.Name)
	if err != nil {
		log.Fatalf("scan user: %w", err)
		return
	}

	b, _ := json.Marshal(u)
	w.Write(b)
}

// CreateUser creates a user in the database
func (s *Server) CreateUser(writer http.ResponseWriter, request *http.Request) {
	createUserReq := struct {
		Name string
	}{}

	if err := json.NewDecoder(request.Body).Decode(&createUserReq); err != nil {
		log.Fatalf("error parsing json id request: %w", err)
		return
	}

	_, err := s.db.ExecContext(context.Background(), "insert into users set name="+createUserReq.Name)
	if err != nil {
		log.Fatalf("insert user: %w", err)
		return
	}

	return
}

// Interface assertions.
var (
	_ http.HandlerFunc = (*Server)(nil).GetUser
	_ http.HandlerFunc = (*Server)(nil).CreateUser
)
