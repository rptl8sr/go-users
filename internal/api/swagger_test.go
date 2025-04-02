package api

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"

	"github.com/stretchr/testify/assert"
)

func TestGetOpenAPISpec(t *testing.T) {
	t.Run("Successful read of OpenAPI spec", func(t *testing.T) {
		specContent := []byte("openapi: '3.0.0'")
		tmpFile, err := os.CreateTemp("", "openapi*.yaml")
		assert.NoError(t, err)
		defer os.Remove(tmpFile.Name())

		_, err = tmpFile.Write(specContent)
		assert.NoError(t, err)
		tmpFile.Close()

		SetOpenAPISpecPath(tmpFile.Name())

		req := httptest.NewRequest(http.MethodGet, "/swagger/doc.yaml", nil)
		rr := httptest.NewRecorder()

		GetOpenAPISpec(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.Equal(t, "application/yaml", rr.Header().Get("Content-Type"))
		assert.Equal(t, string(specContent), rr.Body.String())
	})

	t.Run("OpenAPI spec file not found", func(t *testing.T) {
		SetOpenAPISpecPath("nonexistent.yaml")

		req := httptest.NewRequest(http.MethodGet, "/swagger/doc.yaml", nil)
		rr := httptest.NewRecorder()

		GetOpenAPISpec(rr, req)

		assert.Equal(t, http.StatusInternalServerError, rr.Code)
		assert.Contains(t, rr.Body.String(), "Failed to read OpenAPI specification")
	})
}

func TestSwaggerUI(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	rr := httptest.NewRecorder()

	handler := SwaggerUI()
	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "<title>Swagger UI</title>")
}

func TestRegisterSwaggerRoutes(t *testing.T) {
	r := chi.NewRouter()
	RegisterSwaggerRoutes(r)

	req := httptest.NewRequest(http.MethodGet, "/swagger/doc.yaml", nil)
	rr := httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusInternalServerError, rr.Code)

	req = httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	rr = httptest.NewRecorder()

	r.ServeHTTP(rr, req)
	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Contains(t, rr.Body.String(), "<title>Swagger UI</title>")
}
