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

	db0801 "github.com/nelsonstr/o801/db"
	openapi "github.com/nelsonstr/o801/internal/go"
	server2 "github.com/nelsonstr/o801/internal/server"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {

	dbc, err := sql.Open("postgres", db0801.DbURL())
	if err != nil {
		log.Fatalf("database error: %v", err)
	}
	defer func() { _ = dbc.Close() }()

	db0801.MigrateDB()

	server0801 := server2.NewServer(db0801.NewUserStorage(dbc))

	log.Printf("start server.")

	HealthCheckAPIService := openapi.NewHealthCheckAPIService(dbc)
	HealthCheckAPIController := openapi.NewHealthCheckAPIController(HealthCheckAPIService)

	MonitoringAPIController := openapi.NewMonitoringAPIController()

	mux := openapi.NewRouter(server0801, HealthCheckAPIController, MonitoringAPIController)

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
