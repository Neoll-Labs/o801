package api

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger/v2"
)

func (r *router) SwaggerUI() {
	r.Handle("/swagger/", http.StripPrefix("/swagger",
		httpSwagger.Handler(
			httpSwagger.URL("http://127.0.0.1:8080/api/v1/docs/openapi.yaml"), //The url pointing to API definition
			httpSwagger.DeepLinking(true),
			httpSwagger.DocExpansion("none"),
			httpSwagger.DomID("swagger-ui"),
		)))
}
