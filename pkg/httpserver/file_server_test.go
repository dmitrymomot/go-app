package httpserver_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dmitrymomot/go-server/pkg/httpserver"
	"github.com/stretchr/testify/require"
)

func TestFileServer(t *testing.T) {
	// create a mock HTTP request
	req, err := http.NewRequest("GET", "/static/file.html", nil)
	require.NoError(t, err)

	// create a mock HTTP response recorder
	rr := httptest.NewRecorder()

	// define the test server directory and prefix
	dir := "./testdata"
	prefix := "/static/"

	// create a handler using the test directory and prefix
	handler := httpserver.FileServer(dir, prefix)

	// call the handler's ServeHTTP method with the mock request and response recorder
	handler.ServeHTTP(rr, req)

	// check that the response status code is 200 OK
	require.Equal(t, http.StatusOK, rr.Code)

	// check that the Content-Type header is set correctly
	expectedContentType := "text/html; charset=utf-8"
	require.Equal(t, expectedContentType, rr.Header().Get("Content-Type"))

	// check that the response body contains the expected content
	expectedBody := "This is a test file."
	require.Contains(t, rr.Body.String(), expectedBody)
}
