package middlewares_test

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmitrymomot/go-server/pkg/middlewares"
)

func TestTestingMiddleware(t *testing.T) {
	// Test case 1: must_err parameter is not set, so next handler should be called
	req, err := http.NewRequest("GET", "/login", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	rr := httptest.NewRecorder()

	nextHandlerCalled := false
	nextHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		nextHandlerCalled = true
	})

	middleware := middlewares.Testing()
	handler := middleware(nextHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: expected %v but got %v",
			http.StatusOK, status)
	}

	if !nextHandlerCalled {
		t.Error("next handler was not called")
	}

	// Test case 2: must_err parameter is set with a valid error code, so error response should be returned
	req, err = http.NewRequest("GET", "/login?must_err=403", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	rr = httptest.NewRecorder()

	middleware = middlewares.Testing()
	handler = middleware(nextHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusForbidden {
		t.Errorf("handler returned wrong status code: expected %v but got %v",
			http.StatusForbidden, status)
	}

	expectedBody := "Forbidden\n"
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("handler returned wrong body: expected %q but got %q",
			expectedBody, body)
	}

	// Test case 3: must_err parameter is set with an invalid error code, so error response should be returned
	req, err = http.NewRequest("GET", "/login?must_err=777", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	rr = httptest.NewRecorder()

	middleware = middlewares.Testing()
	handler = middleware(nextHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusBadRequest {
		t.Errorf("handler returned wrong status code: expected %v but got %v",
			http.StatusBadRequest, status)
	}

	expectedErr := errors.New("Invalid error code\n")
	if body := rr.Body.String(); body != expectedErr.Error() {
		t.Errorf("handler returned wrong body: expected %q but got %q",
			expectedErr.Error(), body)
	}

	// Test case 4: must_err parameter is set with a valid error code but unsupported content type, so error response should be returned
	req, err = http.NewRequest("POST", "/login?must_err=404", nil)
	if err != nil {
		t.Fatalf("Error creating request: %v", err)
	}
	req.Header.Set("Content-Type", "text/plain")
	rr = httptest.NewRecorder()

	middleware = middlewares.Testing()
	handler = middleware(nextHandler)

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: expected %v but got %v",
			http.StatusNotFound, status)
	}

	expectedBody = "Not Found\n"
	if body := rr.Body.String(); body != expectedBody {
		t.Errorf("handler returned wrong body: expected %q but got %q",
			expectedBody, body)
	}
}
