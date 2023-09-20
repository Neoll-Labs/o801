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
	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func main() {

	dbc := o801db.InitDB()
	defer func() { _ = dbc.Close() }()

	router := api.NewRouter()
	v1 := router.ApiVersion(1)
	v1.Prefix("/auth").AuthEndpoints()
	v1.Resources("/users").UserEndpoints()

	v1.Resources("/docs").OpenAPI()

	//server0801 := server2.NewServer(db0801.NewUserStorage(dbc))

	log.Printf("start server.")

	//router := internal.NewRouter(server0801)

	//openApiUI(router)

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

// openApiUI add swagger UI for demo proposal only.
func openApiUI(router *http.ServeMux) {
	router.HandleFunc("/swagger/", httpSwagger.Handler(
		httpSwagger.URL("http://127.0.0.1:8080/docs/openapi.yaml"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	))
}
