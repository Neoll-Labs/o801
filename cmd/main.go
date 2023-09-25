/*
 license x
*/

package main

import (
	"errors"
	"github.com/nelsonstr/o801/internal/user/routes"
	"log"
	"net/http"
	"time"

	_ "github.com/lib/pq"
	"github.com/nelsonstr/o801/internal/repository"
	"github.com/nelsonstr/o801/internal/router"
)

func main() {

	dbc, err := repository.InitDB()
	if err != nil {
		log.Fatalf("db connection error: %s", err)
	}

	defer func() { _ = dbc.Close() }()

	if err := repository.MigrateDB(dbc); err != nil {
		log.Fatalf("db migration error: %s", err)
	}

	r := router.NewRouter()
	v1 := r.Version(1)

	routes.InitUserRoutes(dbc, v1)

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
