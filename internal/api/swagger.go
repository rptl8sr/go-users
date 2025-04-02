package api

import (
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
)

// OpenAPISpecPath defines the default file path to the OpenAPI specification in YAML format.
var (
	OpenAPISpecPath = "openapi/openapi.yaml"
)

// SetOpenAPISpecPath sets the file path to the OpenAPI specification, overwriting the default or current path.
func SetOpenAPISpecPath(path string) {
	OpenAPISpecPath = path
}

// GetOpenAPISpec handles HTTP requests to serve the OpenAPI specification file in YAML format.
func GetOpenAPISpec(w http.ResponseWriter, _ *http.Request) {
	spec, err := os.ReadFile(OpenAPISpecPath)
	if err != nil {
		http.Error(w, "Failed to read OpenAPI specification", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/yaml")
	w.Write(spec)
}

// SwaggerUI returns an HTTP handler function that serves a Swagger UI for exploring the OpenAPI specification.
func SwaggerUI() http.HandlerFunc {
	return httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.yaml"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	)
}

// RegisterSwaggerRoutes registers routes for serving the OpenAPI specification and Swagger UI in the provided router.
func RegisterSwaggerRoutes(r chi.Router) {
	r.Get("/swagger/doc.yaml", GetOpenAPISpec)
	r.Get("/swagger/*", SwaggerUI())
}
