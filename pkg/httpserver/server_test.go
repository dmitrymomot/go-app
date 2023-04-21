package httpserver_test

import (
	"context"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/dmitrymomot/go-app/pkg/httpserver"
	"github.com/stretchr/testify/assert"
)

func TestServer(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Hello, World!") })
	server := httpserver.NewServer("localhost:9999", handler)

	// Run the server in a separate goroutine
	go func() {
		if err := server.Run(context.TODO()); err != nil {
			t.Errorf("Server returned an error: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Perform an HTTP request to the server
	resp, err := http.Get(fmt.Sprintf("http://%s", server.Server.Addr))
	assert.NoError(t, err, "Unexpected error in GET request")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code")

	// Shutdown the server
	server.Shutdown()

	// Wait for the server to shut down
	time.Sleep(1 * time.Second)

	// Perform an HTTP request to the server after it has shut down
	_, err = http.Get(fmt.Sprintf("http://%s", server.Server.Addr))
	assert.Error(t, err, "Expected error after server shutdown")
}

func TestServer_Shutdown(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Hello, World!") })
	server := httpserver.NewServer("localhost:9999", handler)

	// Start the server
	go func() {
		if err := server.Run(context.TODO()); err != nil {
			t.Errorf("Server returned an error: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Perform an HTTP request to the server
	resp, err := http.Get(fmt.Sprintf("http://%s", server.Server.Addr))
	assert.NoError(t, err, "Unexpected error in GET request")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code")

	// Shutdown the server
	server.Shutdown()

	// Wait for the server to shut down
	time.Sleep(1 * time.Second)

	// Call Shutdown again and ensure no panic occurs
	assert.NotPanics(t, func() {
		server.Shutdown()
	}, "Calling Shutdown multiple times should not panic")
}

func TestServer_ContextCancellation(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) { fmt.Fprintln(w, "Hello, World!") })
	server := httpserver.NewServer("localhost:9999", handler, httpserver.WithShutdownTimeout(time.Millisecond))

	// Create a context with a timeout and set the Server's context
	ctx, cancel := context.WithTimeout(context.TODO(), 3*time.Second)
	defer cancel()

	// Run the server in a separate goroutine
	go func() {
		if err := server.Run(ctx); err != nil {
			t.Errorf("Server returned an error: %v", err)
		}
	}()

	// Wait for the server to start
	time.Sleep(1 * time.Second)

	// Perform an HTTP request to the server
	resp, err := http.Get(fmt.Sprintf("http://%s", server.Server.Addr))
	assert.NoError(t, err, "Unexpected error in GET request")
	assert.Equal(t, http.StatusOK, resp.StatusCode, "Unexpected status code")

	// Wait for the context to be canceled
	<-ctx.Done()

	// Wait for the server to shut down
	time.Sleep(1 * time.Second)

	// Perform an HTTP request to the server after the context has been canceled
	_, err = http.Get(fmt.Sprintf("http://%s", server.Server.Addr))
	assert.Error(t, err, "Expected error after server shutdown due to context cancellation")
}
