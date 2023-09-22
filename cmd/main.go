/*
 license x
*/

package main

import (
	"errors"
	s "github.com/nelsonstr/o801/internal/handlers"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"

	o801db "github.com/nelsonstr/o801/db"
	"github.com/nelsonstr/o801/internal/api"
)

func main() {

	dbc, err := o801db.InitDB()
	if err != nil {
		log.Fatalf("db connection error: %v", err)
	}
	defer func() { _ = dbc.Close() }()

	o801db.MigrateDB(dbc)

	router := api.NewRouter()

	router.Metrics()

	v1 := router.Version(1)

	v1.Resource("/users").UserEndpoints(s.NewServer(o801db.NewUserStorage(dbc)))

	log.Printf("start server.")

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Error starting or closing listener.
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}
