package main

import (
	"errors"
	"github.com/nelsonstr/o801/db"
	"log"
	"net/http"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	openapi "github.com/nelsonstr/o801/pkg/go"
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {

	db.MigrateDB()

	log.Printf("server started.")

	HealthCheckAPIService := openapi.NewHealthCheckAPIService()
	HealthCheckAPIController := openapi.NewHealthCheckAPIController(HealthCheckAPIService)

	MonitoringAPIService := openapi.NewMonitoringAPIService()
	MonitoringAPIController := openapi.NewMonitoringAPIController(MonitoringAPIService)

	ServicesAPIService := openapi.NewServicesAPIService()
	ServicesAPIController := openapi.NewServicesAPIController(ServicesAPIService)

	mux := openapi.NewRouter(HealthCheckAPIController, MonitoringAPIController, ServicesAPIController)

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
