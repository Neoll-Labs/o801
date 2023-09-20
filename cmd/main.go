/*
 license x
*/

package main

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/nelsonstr/o801/config"
	db0801 "github.com/nelsonstr/o801/db"
	"github.com/nelsonstr/o801/internal"
	server2 "github.com/nelsonstr/o801/internal/server"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {

	dbc, err := sql.Open("postgres", config.DbURL())
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer func() { _ = dbc.Close() }()

	db0801.MigrateDB()

	server0801 := server2.NewServer(db0801.NewUserStorage(dbc))

	log.Printf("start server.")

	mux := internal.NewRouter(server0801)

	// publishing the openapi.yaml.
	mux.Handle("/docs/", http.StripPrefix("/docs/", http.FileServer(http.Dir("api"))))

	openApiUI(mux)

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		// Error starting or closing listener.
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}
}

// openApiUI add swagger UI for demo proposal only.
func openApiUI(router *http.ServeMux) {
	router.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://127.0.0.1:8080/docs/openapi.yaml"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))
}
