/*
 license x
*/

package main

import (
	"errors"
	s "github.com/nelsonstr/o801/internal/handlers"
	"github.com/nelsonstr/o801/internal/repsoitory"
	"github.com/nelsonstr/o801/internal/router"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
)

func main() {

	dbc, err := repsoitory.InitDB()
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer func() { _ = dbc.Close() }()

	if err := repsoitory.MigrateDB(dbc); err != nil {
		log.Fatalf("db migration error: %v", err)
	}

	r := router.NewRouter()

	v1 := r.Version(1)
	v1.Resource("/users").UserEndpoints(s.NewUserServer(repsoitory.NewUserRepo(dbc)))

	log.Printf("start server.")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Error starting or closing listener.
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
