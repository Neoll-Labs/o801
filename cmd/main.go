/*
 license x
*/

package main

import (
	"errors"
	"log"
	"net/http"
	"time"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	o801db "github.com/nelsonstr/o801/db"
	"github.com/nelsonstr/o801/internal/api"
)

func main() {

	dbc := o801db.InitDB()
	defer func() { _ = dbc.Close() }()

	router := api.NewRouter()
	router.Metrics()

	v1 := router.ApiVersion(1)
	v1.SwaggerUI()
	v1.Prefix("/docs").OpenAPI()
	v1.Prefix("/auth").AuthEndpoints()
	v1.Resources("/users").UserEndpoints()

	//server0801 := server2.NewServer(db0801.NewUserStorage(dbc))

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
