package http

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestSwaggerDocRoute(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{}, &fakeTokenManager{})

	request := httptest.NewRequest(http.MethodGet, "/swagger/doc.json", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	if got := response.Header().Get("Content-Type"); !strings.Contains(got, "application/json") {
		t.Fatalf("content type = %q, want application/json", got)
	}

	if !strings.Contains(response.Body.String(), "\"openapi\": \"3.0.3\"") {
		t.Fatalf("response body does not contain openapi version")
	}
}

func TestSwaggerUIRoute(t *testing.T) {
	t.Parallel()

	router := newTestRouter(t, &fakeUserRegistrationUseCase{}, &fakeTokenManager{})

	request := httptest.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", response.Code, http.StatusOK)
	}

	if !strings.Contains(response.Body.String(), "SwaggerUIBundle") {
		t.Fatalf("response body does not look like swagger ui")
	}
}
